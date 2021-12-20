//go:build integration && !unit
// +build integration,!unit

package integration

import (
	"context"
	"github.com/nakiner/faceit/internal/database"
	"github.com/nakiner/faceit/internal/repository/user"
	"log"
	"testing"

	"github.com/nakiner/faceit/configs"
	"github.com/stretchr/testify/require"
)

func TestDatabaseUserServiceCreateUser(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := configs.NewConfig()
	err := cfg.Read()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	require.NoError(t, err)

	repo := user.NewRepository(db)
	_, err = repo.Create(ctx, &user.User{Nickname: "sample"})
	require.NoError(t, err)
}

func TestDatabaseUserServiceUpdateUser(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := configs.NewConfig()
	err := cfg.Read()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	require.NoError(t, err)

	repo := user.NewRepository(db)
	id, err := repo.Create(ctx, &user.User{Nickname: "sample"})
	require.NoError(t, err)

	err = repo.Update(ctx, &user.User{
		ID:       id,
		Nickname: "sample",
	})
	require.NoError(t, err)
}

func TestDatabaseUserServiceDeleteUser(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := configs.NewConfig()
	err := cfg.Read()
	if err != nil {
		log.Fatal(err)
	}

	db, err := database.Connect(ctx, cfg)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	require.NoError(t, err)

	repo := user.NewRepository(db)
	id, err := repo.Create(ctx, &user.User{Nickname: "sample"})
	require.NoError(t, err)

	err = repo.Delete(ctx, id)
	require.NoError(t, err)
}
