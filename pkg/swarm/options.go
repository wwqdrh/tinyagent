package swarm

import (
	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/volume"
)

const (
	SupportedDockerAPIVersion = "1.30"
)

type (
	DockerEngineStatus int

	// DockerNodeRole represent the role of a Docker swarm node
	DockerNodeRole int
)

const (
	_ DockerNodeRole = iota
	// NodeRoleManager represent a Docker swarm manager node role
	NodeRoleManager
	// NodeRoleWorker represent a Docker swarm worker node role
	NodeRoleWorker
)

const (
	_ DockerEngineStatus = iota
	// EngineStatusStandalone represent a standalone Docker environment
	EngineStatusStandalone
	// EngineStatusSwarm represent a Docker swarm environment
	EngineStatusSwarm
)

type RuntimeConfiguration struct {
	AgentPort           string
	EdgeKeySet          bool
	NodeName            string
	DockerConfiguration DockerRuntimeConfiguration
}

type DockerRuntimeConfiguration struct {
	EngineStatus DockerEngineStatus
	Leader       bool
	NodeRole     DockerNodeRole
}

type DockerSnapshot struct {
	Time                    int64             `json:"Time"`
	DockerVersion           string            `json:"DockerVersion"`
	Swarm                   bool              `json:"Swarm"`
	TotalCPU                int               `json:"TotalCPU"`
	TotalMemory             int64             `json:"TotalMemory"`
	RunningContainerCount   int               `json:"RunningContainerCount"`
	StoppedContainerCount   int               `json:"StoppedContainerCount"`
	HealthyContainerCount   int               `json:"HealthyContainerCount"`
	UnhealthyContainerCount int               `json:"UnhealthyContainerCount"`
	VolumeCount             int               `json:"VolumeCount"`
	ImageCount              int               `json:"ImageCount"`
	ServiceCount            int               `json:"ServiceCount"`
	StackCount              int               `json:"StackCount"`
	SnapshotRaw             DockerSnapshotRaw `json:"DockerSnapshotRaw"`
	NodeCount               int               `json:"NodeCount"`
	GpuUseAll               bool              `json:"GpuUseAll"`
	GpuUseList              []string          `json:"GpuUseList"`
}

type DockerSnapshotRaw struct {
	Containers []types.Container       `json:"Containers" swaggerignore:"true"`
	Volumes    volume.VolumeListOKBody `json:"Volumes" swaggerignore:"true"`
	Networks   []types.NetworkResource `json:"Networks" swaggerignore:"true"`
	Images     []types.ImageSummary    `json:"Images" swaggerignore:"true"`
	Info       types.Info              `json:"Info" swaggerignore:"true"`
	Version    types.Version           `json:"Version" swaggerignore:"true"`
}
