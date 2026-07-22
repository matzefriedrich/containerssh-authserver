package services

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh"
)

func Test_PostgresUserProfileService_Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test in short mode")
	}

	ctx := context.Background()

	pgContainer, err := postgres.Run(ctx,
		"postgres:16-alpine",
		postgres.WithDatabase("authserver"),
		postgres.WithUsername("user"),
		postgres.WithPassword("password"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(30*time.Second)),
	)
	if err != nil {
		t.Fatalf("failed to start container: %s", err)
	}

	defer func() {
		if err := pgContainer.Terminate(ctx); err != nil {
			t.Fatalf("failed to terminate container: %s", err)
		}
	}()

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	assert.NoError(t, err)

	pool, err := pgxpool.New(ctx, connStr)
	assert.NoError(t, err)
	defer pool.Close()

	const tableName = "authserver.user_profiles"
	_, err = pool.Exec(ctx, fmt.Sprintf(`
		CREATE SCHEMA IF NOT EXISTS authserver;
		CREATE TABLE %s (
			username TEXT PRIMARY KEY,
			binds TEXT[],
			image TEXT,
			networks TEXT[],
			public_keys TEXT[],
			shell_command TEXT[],
			secret TEXT
		)
	`, tableName))
	assert.NoError(t, err)

	const username = "johndoe"
	const password = "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	publicKeyStr := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOMY6S2Mv8Xp/D8QY68Jz3J+5S/9W8y9G+5p6f8W7Xo/ user@host"

	_, err = pool.Exec(ctx, fmt.Sprintf(`
		INSERT INTO %s (username, binds, image, networks, public_keys, shell_command, secret)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`, tableName),
		username,
		[]string{"/a:/b"},
		"some-image",
		[]string{"net1"},
		[]string{publicKeyStr},
		[]string{"/bin/bash"},
		string(hashedPassword),
	)
	assert.NoError(t, err)

	sut := NewPostgresUserProfileService(pool, tableName)

	t.Run("GetProfile returns correct data", func(t *testing.T) {
		profile, err := sut.GetProfile(username)
		assert.NoError(t, err)
		assert.Equal(t, "some-image", profile.Image)
		assert.Equal(t, []string{"/a:/b"}, profile.Binds)
		assert.Equal(t, []string{"net1"}, profile.Networks)
		assert.Equal(t, []string{publicKeyStr}, profile.PublicKeys)
		assert.Equal(t, []string{"/bin/bash"}, profile.ShellCommand)
		assert.Equal(t, string(hashedPassword), profile.Secret)
	})

	t.Run("VerifySecret works", func(t *testing.T) {
		passwordBase64 := base64.StdEncoding.EncodeToString([]byte(password))
		success, err := sut.VerifySecret(username, passwordBase64)
		assert.NoError(t, err)
		assert.True(t, success)
	})

	t.Run("VerifyPublicKey works", func(t *testing.T) {
		parsedKey, _, _, _, _ := ssh.ParseAuthorizedKey([]byte(publicKeyStr))
		success, err := sut.VerifyPublicKey(username, parsedKey)
		assert.NoError(t, err)
		assert.True(t, success)
	})

	t.Run("GetProfile returns error for nonexistent user", func(t *testing.T) {
		_, err := sut.GetProfile("nonexistent")
		assert.Error(t, err)
		assert.Equal(t, ErrorUserProfileNotFound, err.Error())
	})
}
