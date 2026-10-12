package docker

import (
	"context"
	"io"
	"strings"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/pkg/stdcopy"
)

func (d *docker) ContainersList(ctx context.Context, filters filters.Args) ([]*entities.Container, error) {
	raw, err := d.cli.ContainerList(ctx, container.ListOptions{
		All:     true,
		Filters: filters,
	})
	if err != nil {
		return nil, err
	}
	containers := make([]*entities.Container, 0, len(raw))
	for _, c := range raw {
		name := ""
		if len(c.Names) > 0 {
			name = c.Names[0]
		}
		containers = append(containers, &entities.Container{
			ID:     c.ID,
			Name:   cleanName(name),
			Image:  c.Image,
			Status: c.Status,
			State:  entities.ContainerState(c.State),
		})
	}
	return containers, nil
}

func (d *docker) GetContainer(ctx context.Context, containerID string) (*entities.Container, error) {
	c, err := d.cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, err
	}
	cont := &entities.Container{
		ID:   c.ID,
		Name: cleanName(c.Name),
		// Inspect's top-level Image is the image digest; Config.Image is the
		// reference the container was created from, as the list shows it.
		Image: c.Image,
	}
	if c.Config != nil && c.Config.Image != "" {
		cont.Image = c.Config.Image
	}
	if c.State != nil {
		cont.Status = c.State.Status
		cont.State = entities.ContainerState(c.State.Status)
	}
	return cont, nil
}

// cleanName strips the leading "/" the Docker API puts in front of names.
func cleanName(name string) string {
	return strings.TrimPrefix(name, "/")
}

func (d *docker) ContainerStats(ctx context.Context, containerID string) (io.ReadCloser, error) {
	stats, err := d.cli.ContainerStats(ctx, containerID, true)
	if err != nil {
		return nil, err
	}
	return stats.Body, nil
}

func (d *docker) ContainerLogs(ctx context.Context, containerID string) (io.ReadCloser, error) {
	info, err := d.cli.ContainerInspect(ctx, containerID)
	if err != nil {
		return nil, err
	}
	// TODO: Interesting options about logs are available, you can get them from user settings
	stream, err := d.cli.ContainerLogs(ctx, containerID, container.LogsOptions{
		ShowStdout: true,
		ShowStderr: true,
		Follow:     true,
		Tail:       "50",
	})
	if err != nil {
		return nil, err
	}
	if info.Config != nil && info.Config.Tty {
		// TTY containers stream raw text.
		return stream, nil
	}
	// Without a TTY, stdout and stderr are multiplexed with 8-byte frame
	// headers that must be stripped, or they show up as garbage.
	pr, pw := io.Pipe()
	go func() {
		_, err := stdcopy.StdCopy(pw, pw, stream)
		stream.Close()
		pw.CloseWithError(err)
	}()
	return pr, nil
}

func (d *docker) ContainerStart(ctx context.Context, containerID string) error {
	return d.cli.ContainerStart(ctx, containerID, container.StartOptions{})
}

func (d *docker) ContainerStop(ctx context.Context, containerID string) error {
	return d.cli.ContainerStop(ctx, containerID, container.StopOptions{})
}

func (d *docker) ContainerRemove(ctx context.Context, containerID string, removeForm *entities.ContainerRemoveForm) error {
	opts := container.RemoveOptions{}
	if removeForm != nil {
		opts.RemoveVolumes = removeForm.RemoveVolumes
		opts.Force = removeForm.Force
	}
	return d.cli.ContainerRemove(ctx, containerID, opts)
}

func (d *docker) ContainerRename(ctx context.Context, containerID, newName string) error {
	return d.cli.ContainerRename(ctx, containerID, newName)
}

func (d *docker) ContainerPause(ctx context.Context, containerID string) error {
	return d.cli.ContainerPause(ctx, containerID)
}

func (d *docker) ContainerUnpause(ctx context.Context, containerID string) error {
	return d.cli.ContainerUnpause(ctx, containerID)
}
