package entities

import (
	"sort"
	"strings"

	"github.com/containers/libpod/cmd/podman/shared"
	"github.com/containers/libpod/libpod"
	"github.com/cri-o/ocicni/pkg/ocicni"
	"github.com/pkg/errors"
)

// Listcontainer describes a container suitable for listing
type ListContainer struct {
	// Container command
	Command []string
	// Container creation time
	Created int64
	// If container has exited/stopped
	Exited bool
	// Time container exited
	ExitedAt int64
	// If container has exited, the return code from the command
	ExitCode int32
	// The unique identifier for the container
	ID string `json:"Id"`
	// Container image
	Image string
	// If this container is a Pod infra container
	IsInfra bool
	// Labels for container
	Labels map[string]string
	// User volume mounts
	Mounts []string
	// The names assigned to the container
	Names []string
	// Namespaces the container belongs to.  Requires the
	// namespace boolean to be true
	Namespaces ListContainerNamespaces
	// The process id of the container
	Pid int
	// If the container is part of Pod, the Pod ID. Requires the pod
	// boolean to be set
	Pod string
	// If the container is part of Pod, the Pod name. Requires the pod
	// boolean to be set
	PodName string
	// Port mappings
	Ports []ocicni.PortMapping
	// Size of the container rootfs.  Requires the size boolean to be true
	Size *shared.ContainerSize
	// Time when container started
	StartedAt int64
	// State of container
	State string
}

// ListContainer Namespaces contains the identifiers of the container's Linux namespaces
type ListContainerNamespaces struct {
	// Mount namespace
	MNT string `json:"Mnt,omitempty"`
	// Cgroup namespace
	Cgroup string `json:"Cgroup,omitempty"`
	// IPC namespace
	IPC string `json:"Ipc,omitempty"`
	// Network namespace
	NET string `json:"Net,omitempty"`
	// PID namespace
	PIDNS string `json:"Pidns,omitempty"`
	// UTS namespace
	UTS string `json:"Uts,omitempty"`
	// User namespace
	User string `json:"User,omitempty"`
}

// SortContainers helps us set-up ability to sort by createTime
type SortContainers []*libpod.Container

func (a SortContainers) Len() int      { return len(a) }
func (a SortContainers) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type SortCreateTime struct{ SortContainers }

func (a SortCreateTime) Less(i, j int) bool {
	return a.SortContainers[i].CreatedTime().Before(a.SortContainers[j].CreatedTime())
}

type SortListContainers []ListContainer

func (a SortListContainers) Len() int      { return len(a) }
func (a SortListContainers) Swap(i, j int) { a[i], a[j] = a[j], a[i] }

type psSortedCommand struct{ SortListContainers }

func (a psSortedCommand) Less(i, j int) bool {
	return strings.Join(a.SortListContainers[i].Command, " ") < strings.Join(a.SortListContainers[j].Command, " ")
}

type psSortedId struct{ SortListContainers }

func (a psSortedId) Less(i, j int) bool {
	return a.SortListContainers[i].ID < a.SortListContainers[j].ID
}

type psSortedImage struct{ SortListContainers }

func (a psSortedImage) Less(i, j int) bool {
	return a.SortListContainers[i].Image < a.SortListContainers[j].Image
}

type psSortedNames struct{ SortListContainers }

func (a psSortedNames) Less(i, j int) bool {
	return a.SortListContainers[i].Names[0] < a.SortListContainers[j].Names[0]
}

type psSortedPod struct{ SortListContainers }

func (a psSortedPod) Less(i, j int) bool {
	return a.SortListContainers[i].Pod < a.SortListContainers[j].Pod
}

type psSortedRunningFor struct{ SortListContainers }

func (a psSortedRunningFor) Less(i, j int) bool {
	return a.SortListContainers[i].StartedAt < a.SortListContainers[j].StartedAt
}

type psSortedStatus struct{ SortListContainers }

func (a psSortedStatus) Less(i, j int) bool {
	return a.SortListContainers[i].State < a.SortListContainers[j].State
}

type psSortedSize struct{ SortListContainers }

func (a psSortedSize) Less(i, j int) bool {
	if a.SortListContainers[i].Size == nil || a.SortListContainers[j].Size == nil {
		return false
	}
	return a.SortListContainers[i].Size.RootFsSize < a.SortListContainers[j].Size.RootFsSize
}

type PsSortedCreateTime struct{ SortListContainers }

func (a PsSortedCreateTime) Less(i, j int) bool {
	return a.SortListContainers[i].Created < a.SortListContainers[j].Created
}

func SortPsOutput(sortBy string, psOutput SortListContainers) (SortListContainers, error) {
	switch sortBy {
	case "id":
		sort.Sort(psSortedId{psOutput})
	case "image":
		sort.Sort(psSortedImage{psOutput})
	case "command":
		sort.Sort(psSortedCommand{psOutput})
	case "runningfor":
		sort.Sort(psSortedRunningFor{psOutput})
	case "status":
		sort.Sort(psSortedStatus{psOutput})
	case "size":
		sort.Sort(psSortedSize{psOutput})
	case "names":
		sort.Sort(psSortedNames{psOutput})
	case "created":
		sort.Sort(PsSortedCreateTime{psOutput})
	case "pod":
		sort.Sort(psSortedPod{psOutput})
	default:
		return nil, errors.Errorf("invalid option for --sort, options are: command, created, id, image, names, runningfor, size, or status")
	}
	return psOutput, nil
}
