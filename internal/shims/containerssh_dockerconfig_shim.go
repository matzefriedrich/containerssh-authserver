package shims

import (
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/network"
	specs "github.com/opencontainers/image-spec/specs-go/v1"
	"time"
)

type DockerConfigShim struct {
	Connection DockerConnectionConfigShim `json:"connection" `
	// Execution This field uses the DockerExecutionConfigShim type to make the configuration work with the latest Docker SDK.
	Execution DockerExecutionConfigShim `json:"execution" `
	Timeouts  DockerTimeoutConfigShim   `json:"timeouts" `
}

type DockerConnectionConfigShim struct {
	Host   string `json:"host"  default:"unix:///var/run/docker.sock"`
	CaCert string `json:"cacert" `
	Cert   string `json:"cert" `
	Key    string `json:"key" `
}

type DockerTimeoutConfigShim struct {
	ContainerStart time.Duration `json:"containerStart"  default:"60s"`
	ContainerStop  time.Duration `json:"containerStop"  default:"60s"`
	CommandStart   time.Duration `json:"commandStart"  default:"60s"`
	Signal         time.Duration `json:"signal"  default:"60s"`
	Window         time.Duration `json:"window"  default:"60s"`
	HTTP           time.Duration `json:"http"  default:"15s"`
}

type DockerLaunchConfigShim struct {
	ContainerConfig *container.Config         `json:"container"  comment:"DockerConfig configuration." default:"{\"image\":\"containerssh/containerssh-guest-image\"}"`
	HostConfig      *container.HostConfig     `json:"host"  comment:"Host configuration"`
	NetworkConfig   *network.NetworkingConfig `json:"network"  comment:"Network configuration"`
	Platform        *specs.Platform           `json:"platform"  comment:"Platform specification"`
	ContainerName   string                    `json:"containername"  comment:"Name for the container to be launched"`
}

type DockerExecutionMode string

const (
	DockerExecutionModeConnection DockerExecutionMode = "connection"
	DockerExecutionModeSession    DockerExecutionMode = "session"
)

type DockerImagePullPolicy string

const (
	ImagePullPolicyAlways       DockerImagePullPolicy = "Always"
	ImagePullPolicyIfNotPresent DockerImagePullPolicy = "IfNotPresent"
	ImagePullPolicyNever        DockerImagePullPolicy = "Never"
)

type DockerExecutionConfigShim struct {
	DockerLaunchConfigShim `json:",inline" yaml:",inline"`
	// Auth This field uses the AuthConfigShim type to make the configuration work with the latest Docker SDK.
	Auth                    *AuthConfigShim       `json:"auth" `
	Mode                    DockerExecutionMode   `json:"mode"  default:"connection"`
	IdleCommand             []string              `json:"idleCommand"  default:"[\"/usr/bin/containerssh-agent\", \"wait-signal\", \"--signal\", \"INT\", \"--signal\", \"TERM\"]"`
	ShellCommand            []string              `json:"shellCommand"  default:"[\"/bin/bash\"]"`
	AgentPath               string                `json:"agentPath"  default:"/usr/bin/containerssh-agent"`
	DisableAgent            bool                  `json:"disableAgent" `
	Subsystems              map[string]string     `json:"subsystems"  default:"{\"sftp\":\"/usr/lib/openssh/sftp-server\"}"`
	ImagePullPolicy         DockerImagePullPolicy `json:"imagePullPolicy" default:"IfNotPresent"`
	ExposeAuthMetadataAsEnv bool                  `json:"exposeAuthMetadataAsEnv"`
}

// AuthConfigShim This shim recreates the missing AuthConfig type from the Docker SDK.
type AuthConfigShim struct {
	Username      string `json:"username,omitempty"`
	Password      string `json:"password,omitempty"`
	Auth          string `json:"auth,omitempty"`
	Email         string `json:"email,omitempty"`
	ServerAddress string `json:"serveraddress,omitempty"`
	IdentityToken string `json:"identitytoken,omitempty"`
	RegistryToken string `json:"registrytoken,omitempty"`
}
