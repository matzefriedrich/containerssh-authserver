package services

import (
	"fmt"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/matzefriedrich/containerssh-authserver/internal/shims"
	"github.com/rs/zerolog"
)

type dockerAppConfigService struct {
	profileService UserProfileService
	logger         *zerolog.Logger
}

// CreateApplicationConfigFor creates and returns an application configuration for the specified authenticated username.
// It retrieves the user profile, constructs a Docker container configuration, and handles necessary setup based on user profile information.
func (d *dockerAppConfigService) CreateApplicationConfigFor(authenticatedUsername string) (shims.AppConfigShim, error) {

	d.logger.Info().Msgf("Created application config for user %s", authenticatedUsername)

	profile, profileErr := d.profileService.GetProfile(authenticatedUsername)
	if profileErr != nil {
		d.logger.Error().Msgf("Cannot load user profile: %s", profileErr)
		return InvalidAppConfig, fmt.Errorf("cannot load user profile: %s", profileErr)
	}

	endpointsConfig := map[string]*network.EndpointSettings{}
	for _, networkName := range profile.Networks {
		endpointsConfig[networkName] = &network.EndpointSettings{}
	}

	cfg := shims.AppConfigShim{
		Backend: shims.BackendDocker,
		Docker: shims.DockerConfigShim{
			Connection: shims.DockerConnectionConfigShim{
				Host: "unix:///var/run/docker.sock",
			},
			Execution: shims.DockerExecutionConfigShim{
				ImagePullPolicy: shims.ImagePullPolicyIfNotPresent,
				DisableAgent:    true,
				Mode:            shims.DockerExecutionModeSession,
				ShellCommand:    profile.ShellCommand,
				DockerLaunchConfigShim: shims.DockerLaunchConfigShim{
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
