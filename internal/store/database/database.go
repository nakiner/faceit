package database

import (
	"context"
	"fmt"
	"github.com/nakiner/faceit/configs"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/url"
)

// Connection for different types of operations could be used read/read-write connections to perform actions
type Connection struct {
	Master  *gorm.DB
	Replica *gorm.DB
}

// Connect connects to master and replica and returns database descriptor to access connections
func Connect(ctx context.Context, cfg *configs.Config) (*Connection, error) {
	var res Connection
	var err error

	res.Master, err = connectPool(ctx, cfg.Postgres.Master)
	if err != nil {
		return nil, errors.Wrap(err, "Master DB connect")
	}

	res.Replica, err = connectPool(ctx, cfg.Postgres.Replica)
	if err != nil {
		return nil, errors.Wrap(err, "Replica DB connect")
	}

	return &res, nil
}

// connectPool performs a new connection with given configs.Database config
func connectPool(ctx context.Context, db configs.Database) (conn *gorm.DB, err error) {
	dsn := url.URL{
		User:     url.UserPassword(db.User, db.Password),
		Scheme:   "postgres",
		Host:     fmt.Sprintf("%s:%d", db.Host, db.Port),
		Path:     db.DatabaseName,
		RawQuery: (&url.Values{"sslmode": []string{db.Secure}}).Encode(),
	}

	conn, err = gorm.Open(postgres.Open(dsn.String()), &gorm.Config{})

	return conn, err
}

// Close closes all available idle connections to master/replica sets
func (c *Connection) Close() error {
	db, err := c.Master.DB()
	if err != nil {
		return errors.Wrap(err, "err get conn.Master")
	}
	err = db.Close()
	if err != nil {
		return errors.Wrap(err, "err close conn.Master")
	}

	db, err = c.Replica.DB()
	if err != nil {
		return errors.Wrap(err, "err get conn.Replica")
	}
	err = db.Close()
	if err != nil {
		return errors.Wrap(err, "err close conn.Replica")
	}

	return nil
}

// GetReplicaConn allows applying passed context to active operation and takes active Replica connection
func (c *Connection) GetReplicaConn(ctx context.Context) *gorm.DB {
	return c.Replica.WithContext(ctx)
}

// GetMasterConn allows applying passed context to active operation and takes active Master connection
func (c *Connection) GetMasterConn(ctx context.Context) *gorm.DB {
	return c.Master.WithContext(ctx)
}

// CheckConn performs Ping() operation on connections and returns error if database went down
// In common used in repository to perform Readiness checks and restart service if needed.
func (c *Connection) CheckConn() error {
	db, err := c.Master.DB()
	if err != nil {
		return errors.Wrap(err, "err get conn.Master")
	}
	if err := db.Ping(); err != nil {
		return err
	}
	db, err = c.Replica.DB()
	if err != nil {
		return errors.Wrap(err, "err get conn.Replica")
	}
	if err := db.Ping(); err != nil {
		return err
	}

	return nil
}
