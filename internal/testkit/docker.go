package testkit

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"sync"

	"github.com/arshamalh/dockeroller/docker"
	"github.com/arshamalh/dockeroller/entities"
	"github.com/docker/docker/api/types/filters"
)

// Docker is an in-memory docker.Docker. Containers and images are plain
// slices the test can set up; every call is recorded, and Fail makes the
// next call to a method return an error.
type Docker struct {
	mu         sync.Mutex
	containers []*entities.Container
	images     []*entities.Image
	calls      []string
	failNext   map[string]error

	// Logs, when set, produces the log stream of a container. The stream
	// ends when the returned reader does; the fake closes it on cancel.
	Logs func(ctx context.Context, id string) io.ReadCloser
	// Stats, when set, produces the stats stream (one JSON document per
	// line).
	Stats func(ctx context.Context, id string) io.ReadCloser
}

var _ docker.Docker = (*Docker)(nil)

// NewDocker returns a fake with the given containers and images.
func NewDocker(containers []*entities.Container, images []*entities.Image) *Docker {
	return &Docker{containers: containers, images: images, failNext: map[string]error{}}
}

// Fail makes the next call to method return err.
func (d *Docker) Fail(method string, err error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	d.failNext[method] = err
}

// Calls returns the calls made so far, as "Method arg1 arg2".
func (d *Docker) Calls() []string {
	d.mu.Lock()
	defer d.mu.Unlock()
	return append([]string(nil), d.calls...)
}

// Containers returns a copy of the current containers.
func (d *Docker) Containers() []entities.Container {
	d.mu.Lock()
	defer d.mu.Unlock()
	out := make([]entities.Container, 0, len(d.containers))
	for _, c := range d.containers {
		out = append(out, *c)
	}
	return out
}

// Images returns a copy of the current images.
func (d *Docker) Images() []entities.Image {
	d.mu.Lock()
	defer d.mu.Unlock()
	out := make([]entities.Image, 0, len(d.images))
	for _, i := range d.images {
		out = append(out, *i)
	}
	return out
}

// record logs the call and returns the injected error, if any. d.mu must be
// held.
func (d *Docker) record(method string, args ...any) error {
	parts := []string{method}
	for _, a := range args {
		parts = append(parts, fmt.Sprint(a))
	}
	d.calls = append(d.calls, strings.Join(parts, " "))
	if err, ok := d.failNext[method]; ok {
		delete(d.failNext, method)
		return err
	}
	return nil
}

func (d *Docker) find(id string) (*entities.Container, error) {
	for _, c := range d.containers {
		if c.ID == id {
			return c, nil
		}
	}
	return nil, fmt.Errorf("Error response from daemon: No such container: %s", id)
}

func (d *Docker) ContainersList(_ context.Context, f filters.Args) ([]*entities.Container, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.record("ContainersList"); err != nil {
		return nil, err
	}
	ids := f.Get("id")
	out := []*entities.Container{}
	for _, c := range d.containers {
		if len(ids) > 0 && ids[0] != c.ID {
			continue
		}
		cp := *c
		out = append(out, &cp)
	}
	return out, nil
}

func (d *Docker) GetContainer(_ context.Context, id string) (*entities.Container, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.record("GetContainer", id); err != nil {
		return nil, err
	}
	c, err := d.find(id)
	if err != nil {
		return nil, err
	}
	cp := *c
	return &cp, nil
}

func (d *Docker) setState(method, id string, from []entities.ContainerState, to entities.ContainerState, status string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.record(method, id); err != nil {
		return err
	}
	c, err := d.find(id)
	if err != nil {
		return err
	}
	allowed := false
	for _, s := range from {
		allowed = allowed || c.State == s
	}
	if !allowed {
		return fmt.Errorf("Error response from daemon: container %s is %s", c.Name, c.State)
	}
	c.State, c.Status = to, status
	return nil
}

func (d *Docker) ContainerStart(_ context.Context, id string) error {
	return d.setState("ContainerStart", id,
		[]entities.ContainerState{entities.ContainerStateCreated, entities.ContainerStateExited},
		entities.ContainerStateRunning, "Up Less than a second")
}

func (d *Docker) ContainerStop(_ context.Context, id string) error {
	return d.setState("ContainerStop", id,
		[]entities.ContainerState{entities.ContainerStateRunning, entities.ContainerStatePaused},
		entities.ContainerStateExited, "Exited (0) Less than a second ago")
}

func (d *Docker) ContainerPause(_ context.Context, id string) error {
	return d.setState("ContainerPause", id,
		[]entities.ContainerState{entities.ContainerStateRunning},
		entities.ContainerStatePaused, "Up 1 minute (Paused)")
}

func (d *Docker) ContainerUnpause(_ context.Context, id string) error {
	return d.setState("ContainerUnpause", id,
		[]entities.ContainerState{entities.ContainerStatePaused},
		entities.ContainerStateRunning, "Up 1 minute")
}

func (d *Docker) ContainerRemove(_ context.Context, id string, form *entities.ContainerRemoveForm) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	force, volumes := false, false
	if form != nil {
		force, volumes = form.Force, form.RemoveVolumes
	}
	if err := d.record("ContainerRemove", id, fmt.Sprintf("force=%v", force), fmt.Sprintf("volumes=%v", volumes)); err != nil {
		return err
	}
	for i, c := range d.containers {
		if c.ID != id {
			continue
		}
		if c.State == entities.ContainerStateRunning && !force {
			return fmt.Errorf("Error response from daemon: cannot remove container %q: container is running: stop the container before removing or force remove", "/"+c.Name)
		}
		d.containers = append(d.containers[:i], d.containers[i+1:]...)
		return nil
	}
	return fmt.Errorf("Error response from daemon: No such container: %s", id)
}

func (d *Docker) ContainerRename(_ context.Context, id, name string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.record("ContainerRename", id, name); err != nil {
		return err
	}
	if name == "" || strings.ContainsAny(name, " !/") {
		return fmt.Errorf("Error response from daemon: Error when allocating new name: Invalid container name (/%s), only [a-zA-Z0-9][a-zA-Z0-9_.-] are allowed", name)
	}
	c, err := d.find(id)
	if err != nil {
		return err
	}
	c.Name = name
	return nil
}

func (d *Docker) stream(method, id string, ctx context.Context, src func(context.Context, string) io.ReadCloser) (io.ReadCloser, error) {
	d.mu.Lock()
	err := d.record(method, id)
	if err == nil {
		_, err = d.find(id)
	}
	d.mu.Unlock()
	if err != nil {
		return nil, err
	}
	if src == nil {
		return io.NopCloser(strings.NewReader("")), nil
	}
	rc := src(ctx, id)
	go func() { // the real client stops the stream when ctx is cancelled
		<-ctx.Done()
		rc.Close()
	}()
	return rc, nil
}

func (d *Docker) ContainerLogs(ctx context.Context, id string) (io.ReadCloser, error) {
	return d.stream("ContainerLogs", id, ctx, d.Logs)
}

func (d *Docker) ContainerStats(ctx context.Context, id string) (io.ReadCloser, error) {
	return d.stream("ContainerStats", id, ctx, d.Stats)
}

func (d *Docker) ImagesList(context.Context) ([]*entities.Image, error) {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.record("ImagesList"); err != nil {
		return nil, err
	}
	out := make([]*entities.Image, 0, len(d.images))
	for _, img := range d.images {
		cp := *img
		cp.Tags = append([]string(nil), img.Tags...)
		out = append(out, &cp)
	}
	return out, nil
}

func (d *Docker) ImageTag(_ context.Context, id, tag string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.record("ImageTag", id, tag); err != nil {
		return err
	}
	if !strings.Contains(tag, ":") || strings.ContainsAny(tag, " !") {
		return fmt.Errorf("Error response from daemon: invalid reference format: %q", tag)
	}
	for _, img := range d.images {
		if img.ID == id {
			img.Tags = append(img.Tags, tag)
			return nil
		}
	}
	return errors.New("Error response from daemon: No such image: " + id)
}

func (d *Docker) ImageRemove(_ context.Context, id string, force, pruneChildren bool) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	if err := d.record("ImageRemove", id, fmt.Sprintf("force=%v", force), fmt.Sprintf("prune=%v", pruneChildren)); err != nil {
		return err
	}
	for i, img := range d.images {
		if img.ID != id {
			continue
		}
		if img.Status == entities.ImageStatusInUse && !force {
			return fmt.Errorf("Error response from daemon: conflict: unable to delete %s (must be forced) - image is being used by a container", id[7:19])
		}
		d.images = append(d.images[:i], d.images[i+1:]...)
		return nil
	}
	return errors.New("Error response from daemon: No such image: " + id)
}
