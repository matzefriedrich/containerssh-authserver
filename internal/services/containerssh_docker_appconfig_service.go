package services

import (
	"fmt"

	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	"github.com/matzefriedrich/containerssh-authserver/internal/types"
	shims2 "github.com/matzefriedrich/containerssh-authserver/internal/types/shims"
	"github.com/rs/zerolog"
)

type dockerAppConfigService struct {
	profileService UserProfileService
	logger         *zerolog.Logger
}

// CreateApplicationConfigFor creates and returns an application configuration for the specified authenticated username.
// It retrieves the user profile, constructs a Docker container configuration, and handles necessary setup based on user profile information.
func (d *dockerAppConfigService) CreateApplicationConfigFor(authenticatedUsername string) (shims2.AppConfigShim, error) {

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

	defaultContainerStorageSize := types.Gb(5)
	storageOpt := types.NewStorageOptions().Size(defaultContainerStorageSize).AddOrUpdate(profile.StorageOptions)

	cfg := shims2.AppConfigShim{
		Backend: shims2.BackendDocker,
		Docker: shims2.DockerConfigShim{
			Connection: shims2.DockerConnectionConfigShim{
				Host: "unix:///var/run/docker.sock",
			},
			Execution: shims2.DockerExecutionConfigShim{
				ImagePullPolicy: shims2.ImagePullPolicyIfNotPresent,
				DisableAgent:    true,
				Mode:            shims2.DockerExecutionModeSession,
				ShellCommand:    profile.ShellCommand,
				DockerLaunchConfigShim: shims2.DockerLaunchConfigShim{
					ContainerName: containerName,
					HostConfig: &container.HostConfig{
						Privileged: false,
						AutoRemove: true,
						Binds:      profile.Binds,
						StorageOpt: storageOpt.AsMap(),
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
