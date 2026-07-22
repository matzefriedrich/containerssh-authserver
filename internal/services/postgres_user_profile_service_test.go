package services

import (
	"encoding/base64"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/pashagolub/pgxmock/v3"
	"github.com/stretchr/testify/assert"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/ssh"
)

func Test_PostgresUserProfileService_GetProfile_returns_profile(t *testing.T) {
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

func Test_PostgresUserProfileService_GetProfile_returns_error_if_not_found(t *testing.T) {
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

func Test_PostgresUserProfileService_VerifySecret_returns_true_if_password_matches(t *testing.T) {
	// Arrange
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	const username = "johndoe"
	password := "password123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	passwordBase64 := base64.StdEncoding.EncodeToString([]byte(password))

	mock.ExpectQuery("SELECT binds, image, networks, public_keys, shell_command, secret FROM authserver.user_profiles WHERE username = \\$1").
		WithArgs(username).
		WillReturnRows(pgxmock.NewRows([]string{"binds", "image", "networks", "public_keys", "shell_command", "secret"}).
			AddRow([]string{}, "image", []string{}, []string{}, []string{}, string(hashedPassword)))

	sut := NewPostgresUserProfileService(mock)

	// Act
	success, err := sut.VerifySecret(username, passwordBase64)

	// Assert
	assert.NoError(t, err)
	assert.True(t, success)
}

func Test_PostgresUserProfileService_VerifyPublicKey_returns_true_if_key_matches(t *testing.T) {
	// Arrange
	mock, err := pgxmock.NewPool()
	assert.NoError(t, err)
	defer mock.Close()

	const username = "johndoe"
	publicKeyStr := "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIOMY6S2Mv8Xp/D8QY68Jz3J+5S/9W8y9G+5p6f8W7Xo/ user@host"
	parsedKey, _, _, _, _ := ssh.ParseAuthorizedKey([]byte(publicKeyStr))

	mock.ExpectQuery("SELECT binds, image, networks, public_keys, shell_command, secret FROM authserver.user_profiles WHERE username = \\$1").
		WithArgs(username).
		WillReturnRows(pgxmock.NewRows([]string{"binds", "image", "networks", "public_keys", "shell_command", "secret"}).
			AddRow([]string{}, "image", []string{}, []string{publicKeyStr}, []string{}, "secret"))

	sut := NewPostgresUserProfileService(mock)

	// Act
	success, err := sut.VerifyPublicKey(username, parsedKey)

	// Assert
	assert.NoError(t, err)
	assert.True(t, success)
}
