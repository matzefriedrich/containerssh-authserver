package services

import (
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/matzefriedrich/containerssh-authserver/internal/types"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
)

func Test_PostgresUserProfileBackend_GetProfile_returns_profile(t *testing.T) {
	// Arrange
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	const username = "johndoe"
	expectedBinds := []string{"/home/johndoe:/home/johndoe"}
	expectedImage := "containerssh/containerssh-guest-image"
	expectedNetworks := []string{"containerssh"}
	expectedPublicKeys := []string{"ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCv..."}
	expectedShellCommand := []string{"/bin/bash"}
	expectedSecret := "secret"

	mock.ExpectQuery("SELECT binds, image, networks, public_keys, shell_command, secret FROM authserver.user_profiles WHERE username = \\$1").
		WithArgs(username).
		WillReturnRows(pgxmock.NewRows([]string{"binds", "image", "networks", "public_keys", "shell_command", "secret"}).
			AddRow(expectedBinds, expectedImage, expectedNetworks, expectedPublicKeys, expectedShellCommand, expectedSecret))

	sut := &postgresUserProfileService{pool: mock}

	// Act
	profile, err := sut.GetProfile(t.Context(), username)

	// Assert
	assert.NoError(t, err)
	assert.Equal(t, expectedImage, profile.Image)
	assert.Equal(t, expectedBinds, profile.Binds)
	assert.Equal(t, expectedNetworks, profile.Networks)
	assert.Equal(t, expectedPublicKeys, profile.PublicKeys)
	assert.Equal(t, expectedShellCommand, profile.ShellCommand)
	assert.Equal(t, expectedSecret, profile.Secret)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func Test_PostgresUserProfileBackend_GetProfile_returns_error_if_not_found(t *testing.T) {
	// Arrange
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	const username = "nonexistent"

	mock.ExpectQuery("SELECT binds, image, networks, public_keys, shell_command, secret FROM authserver.user_profiles WHERE username = \\$1").
		WithArgs(username).
		WillReturnError(pgx.ErrNoRows)

	sut := &postgresUserProfileService{pool: mock}

	// Act
	_, err = sut.GetProfile(t.Context(), username)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), types.ErrorUserProfileNotFound)
}

func Test_PostgresUserProfileBackend_RegisterUser_executes_upsert(t *testing.T) {
	// Arrange
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	const username = "johndoe"
	profile := types.UserProfile{
		Binds:        []string{"/a:/b"},
		Image:        "alpine",
		Networks:     []string{"net1"},
		PublicKeys:   []string{"key1"},
		ShellCommand: []string{"/bin/sh"},
		Secret:       "hash",
	}

	mock.ExpectExec("MERGE INTO authserver.user_profiles").
		WithArgs(username, profile.Binds, profile.Image, profile.Networks, profile.PublicKeys, profile.ShellCommand, profile.Secret).
		WillReturnResult(pgxmock.NewResult("MERGE", 1))

	sut := &postgresUserProfileService{pool: mock}

	// Act
	err = sut.RegisterUser(t.Context(), username, profile)

	// Assert
	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}
