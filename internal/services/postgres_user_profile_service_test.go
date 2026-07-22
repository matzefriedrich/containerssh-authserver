package services

import (
	"testing"

	"github.com/jackc/pgx/v5"
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

	sut := NewPostgresUserProfileService(mock)

	// Act
	profile, err := sut.GetProfile(username)

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

	sut := NewPostgresUserProfileService(mock)

	// Act
	_, err = sut.GetProfile(username)

	// Assert
	assert.Error(t, err)
	assert.Contains(t, err.Error(), ErrorUserProfileNotFound)
}
