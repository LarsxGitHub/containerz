package docker

import (
	"context"
	"time"

	"github.com/docker/docker/api/types/container"
	"github.com/openconfig/containerz/containers"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"k8s.io/klog/v2"
)

// ContainerStop stops a container. If the Force option is set and a timeout
// is specified in the context, the contains if forcefully terminated after that timeout.
// If the Force option is set but no timeout is provided the container's StopTimeout
// value is used, if set, otherwise the engine default.
// If the Force option is not set, no forceful termination is performed.
func (m *Manager) ContainerStop(ctx context.Context, instance string, opts ...options.Option) error {
	optionz := options.ApplyOptions(opts...)

	cnts, err := m.client.ContainerList(ctx, container.ListOptions{
		// TODO(alshabib): consider filtering for the image we care about
	})
	if err != nil {
		return err
	}

	// check if the container exists.
	if err := checkExistingInstanceAndPorts(instance, nil, cnts); err == nil {
		return status.Errorf(codes.NotFound, "container %s was not found", instance)
	}

	// a negative timeout indicates to docker that no forceful termination should
	// occur.
	duration := -1
	if optionz.Force {
		// compute duration from context deadline
		duration = 0
		timeoutTime, ok := ctx.Deadline()
		if ok {
			duration = int(timeoutTime.Sub(time.Now()).Seconds())
		}
	}

	pDuration := &duration
	if duration == 0 {
		pDuration = nil
	}

	if err := m.client.ContainerStop(ctx, instance, container.StopOptions{Timeout: pDuration}); err != nil {
		klog.Warningf("container %s was not running", instance)
	}

	if err := m.client.ContainerRemove(ctx, instance, container.RemoveOptions{
		Force: optionz.Force,
	}); err != nil {
		return err
	}

	return nil
}
