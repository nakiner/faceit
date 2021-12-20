package user

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/nakiner/faceit/internal/store/database"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"regexp"
	"testing"
	"time"
)

func TestUserDBRepository_Create(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	DB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	dbpool := database.Connection{
		Master:  DB,
		Replica: DB,
	}

	repo := NewRepository(&dbpool)
	u := User{
		FirstName: "firstname",
		LastName:  "lastname",
		Nickname:  "nickname",
		Password:  "password",
		Email:     "email",
		Country:   "country",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO "users" ("id","first_name","last_name","nickname","password","email","country","created_at","updated_at") 
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9)`)).
		WithArgs(sqlmock.AnyArg(), u.FirstName, u.LastName, u.Nickname, u.Password, u.Email, u.Country, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	res, err := repo.Create(context.Background(), &u)
	require.NoError(t, err)

	assert.Equal(t, u.ID, res)
}

func TestUserDBRepository_Delete(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	DB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	dbpool := database.Connection{
		Master:  DB,
		Replica: DB,
	}

	id := "testid"

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM "users" WHERE id = $1`)).
		WithArgs(id).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	repo := NewRepository(&dbpool)
	err = repo.Delete(context.Background(), id)
	require.NoError(t, err)
}

func TestUserDBRepository_Update(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	DB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	dbpool := database.Connection{
		Master:  DB,
		Replica: DB,
	}

	repo := NewRepository(&dbpool)

	u := User{
		ID:       "testid",
		Nickname: "test",
		Country:  "test",
	}

	mock.ExpectBegin()
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE "users" SET "nickname"=$1,"country"=$2,"updated_at"=$3 WHERE "id" = $4`)).
		WithArgs(u.Nickname, u.Country, sqlmock.AnyArg(), u.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))
	mock.ExpectCommit()

	err = repo.Update(context.Background(), &u)
	require.NoError(t, err)
}

func TestUserDBRepository_Get(t *testing.T) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)

	DB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{
		SkipDefaultTransaction: true,
	})
	require.NoError(t, err)

	dbpool := database.Connection{
		Master:  DB,
		Replica: DB,
	}

	repo := NewRepository(&dbpool)

	tm := time.Now()

	u := User{
		ID:        "test",
		FirstName: "firstname",
		LastName:  "lastname",
		Nickname:  "nickname",
		Password:  "password",
		Email:     "email",
		Country:   "country",
		CreatedAt: tm,
		UpdatedAt: tm,
	}

	rows := sqlmock.
		NewRows([]string{"id", "first_name", "last_name", "nickname", "password", "email", "country", "created_at", "updated_at"}).
		AddRow(u.ID, u.FirstName, u.LastName, u.Nickname, u.Password, u.Email, u.Country, u.CreatedAt, u.UpdatedAt)

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE id = $1 AND country = $2 LIMIT 50`)).
		WithArgs(u.ID, u.Country).
		WillReturnRows(rows)

	conds := Conditions{}
	conds["id"] = u.ID
	conds["country"] = u.Country

	res, err := repo.Get(context.Background(), conds, 50, 1)
	require.NoError(t, err)

	var exp []*User
	exp = append(exp, &u)
	assert.Equal(t, exp, res)
}
