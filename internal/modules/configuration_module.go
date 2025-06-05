package modules

import (
	"github.com/matzefriedrich/containerssh-authserver/internal"
	"github.com/matzefriedrich/containerssh-authserver/internal/utils"
	"os"

	"github.com/matzefriedrich/containerssh-authserver/internal/configuration"

	"github.com/matzefriedrich/parsley/pkg/registration"
	"github.com/matzefriedrich/parsley/pkg/types"
)

const (
	AuthServerConfigPathVariableName = "AUTHSERVER_CONFIG_PATH"
)

// ApplicationConfigurationModule configures and initializes the application settings.
func ApplicationConfigurationModule(registry types.ServiceRegistry) error {

	configurationFilesPath := os.Getenv(AuthServerConfigPathVariableName)
	if _, err := os.Stat(configurationFilesPath); os.IsNotExist(err) {
		configurationFilesPath = "/var/run/authserver"
	}

	const configurationType = "yaml"
	configuration.ConfigureApplication(
		configuration.PathOption(configurationFilesPath),
		configuration.ConfigTypesOption(configurationType))

	settings, configurationErr := configuration.LoadApplicationSettings(
		configuration.EnsureUnprivilegedPortOrDefault(5000),
		configuration.EnsureUserProfilesNotEmpty(),
		configuration.EnsurePublicKeyFormat())

	if configurationErr != nil {
		panic(configurationErr)
	}

	info := getApplicationInfo(configurationFilesPath, configurationType)
	registration.RegisterInstance(registry, info)

	return registration.RegisterInstance(registry, settings)
}

func getApplicationInfo(configurationFilesPath string, configurationType string) *configuration.ApplicationInfo {
	info := &configuration.ApplicationInfo{
		Version:           internal.Version,
		Commit:            utils.IIf(len(internal.CommitSha) > 0, internal.CommitSha, "not set"),
		ReleaseName:       utils.IIf(len(internal.ReleaseName) > 0, internal.ReleaseName, "authserver"),
		ConfigurationPath: configurationFilesPath,
		ConfigurationType: configurationType,
	}
	return info
}
