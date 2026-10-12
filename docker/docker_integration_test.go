package docker

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/arshamalh/dockeroller/entities"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/api/types/image"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// These tests talk to a real Docker daemon (found through DOCKER_HOST like
// the docker CLI) and create, change and remove a throwaway container. They
// only run when asked:
//
//	DOCKEROLLER_DOCKER_TESTS=1 go test ./docker/
//
// DOCKEROLLER_TEST_IMAGE picks the image (default busybox:1.36, pulled if
// missing). With a custom image, its own entrypoint must keep printing
// lines; with busybox a shell loop does.
func integration(t *testing.T) (*docker, string) {
	t.Helper()
	if os.Getenv("DOCKEROLLER_DOCKER_TESTS") == "" {
		t.Skip("set DOCKEROLLER_DOCKER_TESTS=1 to run against the local Docker daemon")
	}
	d, err := New()
	require.NoError(t, err)
	t.Cleanup(func() { d.Close() })

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	require.NoError(t, d.Ping(ctx), "Docker daemon not reachable")

	img := os.Getenv("DOCKEROLLER_TEST_IMAGE")
	if img == "" {
		img = "busybox:1.36"
	}
	if _, err := d.cli.ImageInspect(ctx, img); err != nil {
		pull, err := d.cli.ImagePull(context.Background(), img, image.PullOptions{})
		require.NoError(t, err, "pull %s", img)
		_, _ = io.Copy(io.Discard, pull)
		pull.Close()
	}
	return d, img
}

// newContainer creates (but doesn't start) a container whose name contains
// MarkdownV2-hostile characters, and removes it when the test ends.
func newContainer(t *testing.T, d *docker, img string) (id, name string) {
	t.Helper()
	name = fmt.Sprintf("dclr-it_%d.test", rand.Intn(1_000_000))
	cfg := &container.Config{Image: img}
	if strings.HasPrefix(img, "busybox") {
		cfg.Cmd = []string{"sh", "-c", `i=0; while true; do i=$((i+1)); echo "line $i <a_b>"; echo "err $i" >&2; sleep 0.2; done`}
	}
	created, err := d.cli.ContainerCreate(context.Background(), cfg, nil, nil, nil, name)
	require.NoError(t, err)
	t.Cleanup(func() {
		_ = d.cli.ContainerRemove(context.Background(), created.ID, container.RemoveOptions{Force: true})
	})
	return created.ID, name
}

func TestIntegrationContainerLifecycle(t *testing.T) {
	d, img := integration(t)
	id, name := newContainer(t, d, img)
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	list, err := d.ContainersList(ctx, filters.NewArgs(filters.Arg("id", id)))
	require.NoError(t, err)
	require.Len(t, list, 1)
	assert.Equal(t, name, list[0].Name, "the leading / of the API is stripped")
	assert.Equal(t, img, list[0].Image)
	assert.Equal(t, entities.ContainerStateCreated, list[0].State)

	require.NoError(t, d.ContainerStart(ctx, id))
	got, err := d.GetContainer(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, entities.ContainerStateRunning, got.State)
	assert.Equal(t, img, got.Image, "the image reference, not its digest")
	assert.Equal(t, name, got.Name)

	require.NoError(t, d.ContainerPause(ctx, id))
	got, _ = d.GetContainer(ctx, id)
	assert.Equal(t, entities.ContainerStatePaused, got.State)
	require.NoError(t, d.ContainerUnpause(ctx, id))

	renamed := name + "-renamed"
	require.NoError(t, d.ContainerRename(ctx, id, renamed))
	got, _ = d.GetContainer(ctx, id)
	assert.Equal(t, renamed, got.Name)
	assert.Error(t, d.ContainerRename(ctx, id, "bad name!"), "daemon errors are returned")

	assert.Error(t, d.ContainerRemove(ctx, id, &entities.ContainerRemoveForm{}), "a running container needs force")
	require.NoError(t, d.ContainerStop(ctx, id))
	got, _ = d.GetContainer(ctx, id)
	assert.Equal(t, entities.ContainerStateExited, got.State)
	require.NoError(t, d.ContainerRemove(ctx, id, &entities.ContainerRemoveForm{RemoveVolumes: true}))

	_, err = d.GetContainer(ctx, id)
	assert.Error(t, err)
	assert.Error(t, d.ContainerStart(ctx, id), "unknown containers are an error, not a silent no-op")
}

func TestIntegrationLogsAreDemultiplexed(t *testing.T) {
	d, img := integration(t)
	id, _ := newContainer(t, d, img)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	require.NoError(t, d.ContainerStart(ctx, id))

	logsCtx, stop := context.WithCancel(ctx)
	rc, err := d.ContainerLogs(logsCtx, id)
	require.NoError(t, err)
	defer rc.Close()

	sc := bufio.NewScanner(rc)
	var got []string
	for len(got) < 4 && sc.Scan() {
		got = append(got, sc.Text())
	}
	stop()
	require.Len(t, got, 4, "lines: %q, err: %v", got, sc.Err())
	for _, l := range got {
		assert.NotContains(t, l, "\x00", "no 8-byte stream headers in %q", l)
		assert.NotContains(t, l, "\x01", "no 8-byte stream headers in %q", l)
		assert.NotContains(t, l, "\x02", "no 8-byte stream headers in %q", l)
	}
	if strings.HasPrefix(img, "busybox") {
		joined := strings.Join(got, "\n")
		assert.Contains(t, joined, "line 1 <a_b>", "stdout")
		assert.Contains(t, joined, "err 1", "stderr")
	}

	// Cancelling the context ends the stream.
	done := make(chan struct{})
	go func() { io.Copy(io.Discard, rc); close(done) }()
	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("logs stream did not end after cancel")
	}
}

func TestIntegrationStats(t *testing.T) {
	d, img := integration(t)
	id, _ := newContainer(t, d, img)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	require.NoError(t, d.ContainerStart(ctx, id))

	statsCtx, stop := context.WithCancel(ctx)
	defer stop()
	rc, err := d.ContainerStats(statsCtx, id)
	require.NoError(t, err)
	defer rc.Close()

	dec := json.NewDecoder(rc)
	var first, second entities.Stats
	require.NoError(t, dec.Decode(&first))
	require.NoError(t, dec.Decode(&second))
	assert.NotZero(t, second.Memory.Usage, "memory is reported")
	assert.NotZero(t, second.PreCPU.SystemCPUUsage, "the second sample carries the previous CPU counters (precpu_stats)")
	assert.NotZero(t, second.CPU.OnlineCPUs)
	cpu := second.CPUUsage()
	assert.GreaterOrEqual(t, cpu, 0.0)
	assert.Less(t, cpu, float64(second.CPU.OnlineCPUs)*100+1, "a current load, not a cumulative figure")
}

func TestIntegrationImages(t *testing.T) {
	d, img := integration(t)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	info, err := d.cli.ImageInspect(ctx, img)
	require.NoError(t, err)

	tag := fmt.Sprintf("dclr-it:t%d", rand.Intn(1_000_000))
	require.NoError(t, d.ImageTag(ctx, info.ID, tag))
	t.Cleanup(func() { _, _ = d.cli.ImageRemove(context.Background(), tag, image.RemoveOptions{}) })

	images, err := d.ImagesList(ctx)
	require.NoError(t, err)
	var found *entities.Image
	for _, i := range images {
		for _, tg := range i.Tags {
			if tg == tag {
				found = i
			}
		}
	}
	require.NotNil(t, found, "the new tag is listed")
	assert.Positive(t, found.Size)
	assert.NotEmpty(t, found.CreatedAt)

	// Removing by tag only untags; the image itself stays.
	require.NoError(t, d.ImageRemove(ctx, tag, false, false))
	_, err = d.cli.ImageInspect(ctx, img)
	assert.NoError(t, err)
	assert.Error(t, d.ImageTag(ctx, info.ID, "Invalid Tag!"), "daemon errors are returned")
}

func TestIntegrationErrorsAreReturned(t *testing.T) {
	d, _ := integration(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := d.ContainersList(ctx, filters.Args{})
	assert.Error(t, err, "a failed list is an error, not an empty list")
	_, err = d.ImagesList(ctx)
	assert.Error(t, err)
}
