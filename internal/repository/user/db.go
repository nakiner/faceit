package user

import (
	"context"
	"github.com/google/uuid"
	"github.com/nakiner/faceit/internal/store/database"
	"github.com/pkg/errors"
	"gorm.io/gorm"
	"time"
)

var (
	ErrRowsAffectedEmpty = errors.New("result.RowsAffected is empty")
)

type userDBRepository struct {
	db    *database.Connection
	ready bool
}

// NewRepository creates new interfaced repository with CRUD to User entity
func NewRepository(db *database.Connection) Repository {
	rep := userDBRepository{db: db, ready: true}
	go func() {
		tic := time.Tick(time.Minute * 1)
		for range tic {
			if rep.db.CheckConn() != nil {
				rep.ready = false
			}
		}
	}()
	return &rep
}

// IsReady used in Readiness checks to check availability of CRUD operations
func (r *userDBRepository) IsReady() bool {
	return r.ready
}

// Create creates a new User entity in database
func (r *userDBRepository) Create(ctx context.Context, data *User) (string, error) {
	conn := r.db.GetMasterConn(ctx)

	id, err := uuid.NewUUID()
	if err != nil {
		return "", errors.Wrap(err, "userDBRepository generate uuid err")
	}

	data.ID = id.String()

	if err := conn.Create(data).Error; err != nil {
		return "", errors.Wrap(err, "userDBRepository Create err")
	}

	return data.ID, nil
}

// Delete deletes a User entity with given id represented by UUID standard
func (r *userDBRepository) Delete(ctx context.Context, id string) error {
	conn := r.db.GetMasterConn(ctx)

	result := conn.Delete(&User{}, "id = ?", id)

	if err := result.Error; err != nil {
		return errors.Wrap(err, "userDBRepository Delete err")
	}

	if count := result.RowsAffected; count < 1 {
		return errors.Wrap(ErrRowsAffectedEmpty, "userDBRepository Delete err")
	}

	return nil
}

// Update updates a User entity with given id and User payload
func (r *userDBRepository) Update(ctx context.Context, data *User) error {
	conn := r.db.GetMasterConn(ctx)

	result := conn.Model(data).Updates(data)
	if err := result.Error; err != nil {
		return errors.Wrap(err, "userDBRepository Update err")
	}

	if count := result.RowsAffected; count < 1 {
		return errors.Wrap(ErrRowsAffectedEmpty, "userDBRepository Update err")
	}

	return nil
}

// Get performs select from database with set of conditions to fetch User collection
// Conditions parsed into database prepared conditions to take only required data
// Empty rowset does not handled to evade 404 behavior on transport layer.
func (r *userDBRepository) Get(ctx context.Context, conditions Conditions, limit uint32, offset uint32) ([]*User, error) {
	conn := r.db.GetReplicaConn(ctx)

	var users []*User
	var result *gorm.DB

	if len(conditions) > 0 {
		var values map[string]interface{}
		values = conditions
		result = conn.Where(conditions.GetPreparedStatement(), values).Limit(int(limit)).Offset(int(offset - 1)).Find(&users)
	} else {
		result = conn.Limit(int(limit)).Offset(int(offset - 1)).Find(&users)
	}

	if err := result.Error; err != nil {
		return nil, errors.Wrap(err, "userDBRepository Get err")
	}

	if err := result.Error; err != nil {
		return nil, err
	}

	return users, nil
}
