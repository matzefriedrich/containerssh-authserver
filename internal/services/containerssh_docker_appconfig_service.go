package services

import (
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/rs/zerolog"
	"go.containerssh.io/containerssh/config"
)

type dockerAppConfigService struct {
	profileService UserProfileService
	logger         *zerolog.Logger
}

// CreateApplicationConfigFor creates and returns an application configuration for the specified authenticated username.
// It retrieves the user profile, constructs a Docker container configuration, and handles necessary setup based on user profile information.
func (d *dockerAppConfigService) CreateApplicationConfigFor(authenticatedUsername string) (config.AppConfig, error) {

	d.logger.Info().Msgf("Created application config for user %s", authenticatedUsername)

	profile, profileErr := d.profileService.GetProfile(authenticatedUsername)
	if profileErr != nil {
		d.logger.Error().Msgf("Cannot load user profile: %s", profileErr)
		return InvalidAppConfig, fmt.Errorf("cannot load user profile: %s", profileErr)
	}

	containerName := fmt.Sprintf("containerssh-session-%s", authenticatedUsername)

	endpointsConfig := map[string]*network.EndpointSettings{}
	for _, networkName := range profile.Networks {
		endpointsConfig[networkName] = &network.EndpointSettings{}
	}

	cfg := config.AppConfig{
		Backend: config.BackendDocker,
		Docker: config.DockerConfig{
			Connection: config.DockerConnectionConfig{
				Host: "unix:///var/run/docker.sock",
			},
			Execution: config.DockerExecutionConfig{
				ImagePullPolicy: config.ImagePullPolicyIfNotPresent,
				DisableAgent:    true,
				Mode:            config.DockerExecutionModeSession,
				ShellCommand:    profile.ShellCommand,
				DockerLaunchConfig: config.DockerLaunchConfig{
					ContainerName: containerName,
					HostConfig: &container.HostConfig{
						Privileged: false,
						AutoRemove: true,
						Binds:      profile.Binds,
					},
					NetworkConfig: &network.NetworkingConfig{
						EndpointsConfig: endpointsConfig,
					},
					ContainerConfig: &container.Config{
						Image: profile.Image,
					},
				},
			},
		},
	}

	return cfg, nil
}

var _ ContainerAppConfigService = (*dockerAppConfigService)(nil)

func NewDockerAppConfigService(profileService UserProfileService, logger *zerolog.Logger) ContainerAppConfigService {
	return &dockerAppConfigService{
		profileService: profileService,
		logger:         logger,
	}
}
