// Package docker implements a container manager for docker orchestration.
package docker

import (
	"context"
	"io"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/registry"
	"github.com/docker/docker/api/types/volume"

	ocispec "github.com/opencontainers/image-spec/specs-go/v1"
)

type docker interface {
	Close() error
	ContainerCreate(ctx context.Context, config *container.Config, hostConfig *container.HostConfig, networkingConfig *network.NetworkingConfig, platform *ocispec.Platform, containerName string) (container.CreateResponse, error)
	ContainerList(ctx context.Context, options container.ListOptions) ([]types.Container, error)
	ContainerLogs(ctx context.Context, container string, options container.LogsOptions) (io.ReadCloser, error)
	ContainerRemove(ctx context.Context, container string, options container.RemoveOptions) error
	ContainerStart(ctx context.Context, container string, options container.StartOptions) error
	ContainerStop(ctx context.Context, container string, options container.StopOptions) error
	ImageList(ctx context.Context, options image.ListOptions) ([]image.Summary, error)
	ImageLoad(ctx context.Context, input io.Reader, quiet bool) (types.ImageLoadResponse, error)
	ImagePull(ctx context.Context, ref string, options image.PullOptions) (io.ReadCloser, error)
	ImageRemove(ctx context.Context, imageID string, options image.RemoveOptions) ([]image.DeleteResponse, error)
	ImageTag(ctx context.Context, source, target string) error
	RegistryLogin(ctx context.Context, auth registry.AuthConfig) (registry.AuthenticateOKBody, error)
	VolumeCreate(ctx context.Context, options volume.CreateOptions) (volume.Volume, error)
	VolumeList(ctx context.Context, options volume.ListOptions) (volume.ListResponse, error)
	VolumeRemove(ctx context.Context, volumeID string, force bool) error

	ContainersPrune(ctx context.Context, args filters.Args) (types.ContainersPruneReport, error)
	ImagesPrune(ctx context.Context, args filters.Args) (types.ImagesPruneReport, error)
}

// Manager is a docker container orchestration manager.
type Manager struct {
	client  docker
	janitor *Vacuum
}

// New builds a new docker manager given a docker client.
func New(cli docker) *Manager {
	return &Manager{
		client:  cli,
		janitor: NewJanitor(cli),
	}
}

// Start starts a docker session to the host
func (m *Manager) Start(ctx context.Context) error {
	m.janitor.Start(ctx)
	return nil
}

// Stop closes the connection to the docker server.
func (m *Manager) Stop(ctx context.Context) error {
	m.janitor.Stop(ctx)
	return m.client.Close()
}
