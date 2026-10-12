package telegram_test

import (
	"context"
	"errors"
	"io"
	"strings"
	"testing"
	"time"

	"github.com/arshamalh/dockeroller/internal/testkit"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestStartShowsMainMenu(t *testing.T) {
	h := newHarness(t, fixtures())
	h.send("/start")

	p := h.panel()
	assert.Contains(t, p.Text, `Hello Ars\_ham`, "the user's name is escaped")
	assert.Equal(t, [][]string{{"Images 🌄", "Containers 📦"}}, p.Labels())
}

func TestWhitelistDeniesStrangers(t *testing.T) {
	d := fixtures()
	h := newHarness(t, d, 42) // the owner is not whitelisted

	h.send("/containers")
	assert.Contains(t, h.panel().Text, "Your chat ID is `7`")
	assert.Contains(t, h.panel().Text, "`\\-w 7`")

	mark := h.tg.Mark()
	h.press(1, "\fcontStop|0")
	calls := h.tg.CallsSince(mark)
	require.Len(t, calls, 1, "a stranger's button is only dismissed")
	assert.Equal(t, "answerCallbackQuery", calls[0].Method)
	assert.Empty(t, d.Calls(), "a stranger never reaches Docker")
}

func TestContainersNavigationAndStateButtons(t *testing.T) {
	h := newHarness(t, fixtures())
	h.send("/containers")

	p := h.panel()
	assert.True(t, containsAll(p.Text, `web\-1`, `nginx:1\.27`, "running", "Up 2 hours"), p.Text)
	assert.Equal(t, []string{"Stop ⏹", "Pause ⏸"}, h.row(1), "running: stop or pause")

	panelID := p.ID
	h.click("Next ➡")
	p = h.panel()
	assert.Equal(t, panelID, p.ID, "navigation edits the same message")
	assert.Contains(t, p.Text, `api\_server`)
	assert.Equal(t, []string{"Unpause ⏯", "Stop ⏹"}, h.row(1), "paused: unpause or stop")

	h.click("Next ➡")
	assert.Contains(t, h.panel().Text, `stopped\.one`)
	assert.Equal(t, []string{"Start ▶️"}, h.row(1), "exited: start")

	h.click("Next ➡")
	assert.Contains(t, h.panel().Text, `web\-1`, "Next wraps around to the first")
	h.click("⬅ Prev")
	assert.Contains(t, h.panel().Text, `stopped\.one`, "Prev wraps around to the last")
}

func TestContainerLifecycle(t *testing.T) {
	d := fixtures()
	h := newHarness(t, d)
	h.send("/containers")

	mark := h.tg.Mark()
	h.click("Pause ⏸")
	assert.Equal(t, []string{"Container paused"}, h.tg.Toasts(mark))
	assert.Contains(t, h.panel().Text, "paused")
	assert.Equal(t, []string{"Unpause ⏯", "Stop ⏹"}, h.row(1), "buttons follow the new state")

	h.click("Unpause ⏯")
	h.click("Stop ⏹")
	assert.Contains(t, h.panel().Text, "exited")
	h.click("Start ▶️")
	assert.Contains(t, h.panel().Text, "running")

	assert.Equal(t, []string{"ContainerPause c1", "ContainerUnpause c1", "ContainerStop c1", "ContainerStart c1"},
		filter(d.Calls(), "ContainerPause", "ContainerUnpause", "ContainerStop", "ContainerStart"))
}

func TestContainerActionFailureIsAToast(t *testing.T) {
	d := fixtures()
	h := newHarness(t, d)
	h.send("/containers")
	before := h.panel().Text

	long := "Error response from daemon: " + strings.Repeat("x", 300)
	d.Fail("ContainerStop", errors.New(long))
	mark := h.tg.Mark()
	h.click("Stop ⏹")

	toasts := h.tg.Toasts(mark)
	require.Len(t, toasts, 1)
	assert.True(t, strings.HasPrefix(toasts[0], "Failed: Error response from daemon: xxx"), toasts[0])
	assert.LessOrEqual(t, len([]rune(toasts[0])), 200, "Telegram caps toasts at 200 characters")
	assert.Equal(t, before, h.panel().Text, "the view stays as it was")
}

func TestBackShowsCurrentState(t *testing.T) {
	d := fixtures()
	h := newHarness(t, d)
	h.send("/containers")
	h.click("Rename ✏️")

	// Meanwhile someone stops the container from another client.
	require.NoError(t, d.ContainerStop(context.Background(), "c1"))
	h.click("↩ Back")
	assert.Contains(t, h.panel().Text, "exited", "Back re-reads the container instead of showing the cached state")
	assert.Equal(t, []string{"Start ▶️"}, h.row(1))
}

func TestStaleButtonRefreshesTheList(t *testing.T) {
	h := newHarness(t, fixtures())
	h.send("/containers")

	mark := h.tg.Mark()
	h.press(h.panel().ID, "\fcontStop|9") // a button for a container that is gone
	assert.Equal(t, []string{"The list changed, showing it again"}, h.tg.Toasts(mark))
	assert.Contains(t, h.panel().Text, `web\-1`)
	assert.NotContains(t, filter(h.docker.Calls(), "ContainerStop"), "ContainerStop c1", "nothing is stopped")

	mark = h.tg.Mark()
	h.press(h.panel().ID, "\fcontStop|abc")
	assert.Equal(t, []string{"Invalid button 🤔"}, h.tg.Toasts(mark))

	mark = h.tg.Mark()
	h.press(h.panel().ID, "\fsomethingRetired|1")
	assert.Equal(t, []string{"This button is outdated, send /start"}, h.tg.Toasts(mark))
}

func TestRenameScene(t *testing.T) {
	d := fixtures()
	h := newHarness(t, d)
	h.send("/containers")
	h.click("Rename ✏️")
	assert.Contains(t, h.panel().Text, `Send the new name for *web\-1*`)

	h.send("bad name!")
	assert.Contains(t, h.panel().Text, `Couldn't rename *web\-1*`)
	assert.Contains(t, h.panel().Text, `Invalid container name \(/bad name\!\)`, "daemon errors are escaped")

	h.send("frontend_v2")
	assert.Contains(t, h.panel().Text, `✅ *web\-1* renamed to *frontend\_v2*`)
	assert.Contains(t, h.panel().Text, "`Name:   ` frontend\\_v2")
	assert.Equal(t, "frontend_v2", d.Containers()[0].Name)

	// The scene is over: the next text is not taken as a name.
	h.send("hello")
	assert.Contains(t, h.panel().Text, "Welcome to *Dockeroller*")
	assert.Len(t, filter(d.Calls(), "ContainerRename"), 2)
}

func TestRemoveForm(t *testing.T) {
	d := fixtures()
	h := newHarness(t, d)
	h.send("/containers")

	mark := h.tg.Mark()
	h.click("Remove 🗑")
	assert.Equal(t, []string{"Set the options and press Done"}, h.tg.Toasts(mark))
	assert.Equal(t, [][]string{{"Force: false", "Remove volumes: false"}, {"↩ Back", "Done ✅"}}, h.panel().Labels())

	// Without force a running container can't be removed.
	mark = h.tg.Mark()
	h.click("Done ✅")
	require.Len(t, h.tg.Toasts(mark), 1)
	assert.Contains(t, h.tg.Toasts(mark)[0], "container is running")
	assert.Len(t, d.Containers(), 3)

	h.click("Force: false")
	h.click("Remove volumes: false")
	assert.Equal(t, []string{"Force: true", "Remove volumes: true"}, h.row(0))
	mark = h.tg.Mark()
	h.click("Done ✅")
	assert.Equal(t, []string{"Container removed"}, h.tg.Toasts(mark))
	assert.Contains(t, filter(d.Calls(), "ContainerRemove"), "ContainerRemove c1 force=true volumes=true")
	assert.Len(t, d.Containers(), 2)
	assert.Contains(t, h.panel().Text, `api\_server`, "the refreshed list is shown")

	// The options are reset for the next container.
	h.click("Remove 🗑")
	assert.Equal(t, []string{"Force: false", "Remove volumes: false"}, h.row(0))
}

func TestImagesFlows(t *testing.T) {
	d := fixtures()
	h := newHarness(t, d)
	h.send("/start")
	h.click("Images 🌄")

	p := h.panel()
	assert.True(t, containsAll(p.Text, "`863dd03b5190`", `187\.00 MegaBytes`, `nginx:1\.27, nginx:latest`, "In use"), p.Text)

	h.click("Next ➡")
	assert.True(t, containsAll(h.panel().Text, "`d5318421c392`", `\<none\>`, "Dangling"), h.panel().Text)

	h.click("Tag 🏷")
	h.send("not a tag")
	assert.Contains(t, h.panel().Text, "Couldn't tag")
	h.send("myapp:1.2")
	assert.Contains(t, h.panel().Text, "tagged *myapp:1\\.2*")
	assert.Equal(t, []string{"myapp:1.2"}, d.Images()[1].Tags)

	// An image in use needs force.
	h.click("⬅ Prev")
	h.click("Remove 🗑")
	mark := h.tg.Mark()
	h.click("Done ✅")
	assert.Contains(t, h.tg.Toasts(mark)[0], "must be forced")
	h.click("Force: false")
	h.click("Prune children: false")
	mark = h.tg.Mark()
	h.click("Done ✅")
	assert.Equal(t, []string{"Image removed"}, h.tg.Toasts(mark))
	assert.Contains(t, filter(d.Calls(), "ImageRemove"),
		"ImageRemove sha256:863dd03b51905b838bbe0e32d14295bbe855412b2b23c74b4061551c1d071236 force=true prune=true")
	assert.Len(t, d.Images(), 1)
}

func TestEmptyListsAndDockerDown(t *testing.T) {
	d := testkit.NewDocker(nil, nil)
	h := newHarness(t, d)

	h.send("/containers")
	assert.Equal(t, "There are no containers\\.", h.panel().Text)
	h.send("/images")
	assert.Equal(t, "There are no images\\.", h.panel().Text)

	d.Fail("ContainersList", errors.New("Cannot connect to the Docker daemon at unix:///var/run/docker.sock. Is the docker daemon running?"))
	h.send("/containers")
	assert.Contains(t, h.panel().Text, "Docker didn't answer: Cannot connect to the Docker daemon at unix:///var/run/docker\\.sock\\.")
}

func TestLogsEscaping(t *testing.T) {
	d := fixtures()
	d.Logs = lines(
		"GET /index.html 200 <1ms>",
		"weird `backticks` and \\ backslash",
		"a_b (x) [y] #1 *not bold*",
	)
	h := newHarness(t, d)
	h.send("/containers")
	panelID := h.click("Logs 🪵")

	testkit.Eventually(t, 2*time.Second, func() bool {
		return strings.Contains(h.tg.Message(panelID).Text, "a_b (x) [y] #1 *not bold*")
	}, "logs never appeared: %v", h.tg.Message(panelID))

	text := h.tg.Message(panelID).Text
	assert.True(t, strings.HasPrefix(text, "🪵 *web\\-1*\n```\n"), text)
	assert.Contains(t, text, "weird \\`backticks\\` and \\\\ backslash", "inside pre only ` and \\ are escaped")
	assert.Contains(t, text, "GET /index.html 200 <1ms>", "nothing else is escaped inside pre")
}

func TestBackStopsTheStream(t *testing.T) {
	d := fixtures()
	d.Logs = endless("tick")
	h := newHarness(t, d)
	h.send("/containers")
	panelID := h.click("Logs 🪵")
	testkit.Eventually(t, 2*time.Second, func() bool {
		return strings.Contains(h.tg.Message(panelID).Text, "tick 3")
	}, "logs never appeared")

	h.click("↩ Back")
	assert.Contains(t, h.tg.Message(panelID).Text, "`Name:   ` web\\-1")
	mark := h.tg.Mark()
	time.Sleep(150 * time.Millisecond) // several stream intervals of a busy log
	for _, c := range h.tg.CallsSince(mark) {
		assert.NotEqual(t, "editMessageText", c.Method, "the stream kept editing after Back: %q", c.Text())
	}
	assert.Contains(t, h.tg.Message(panelID).Text, "`Name:   ` web\\-1", "the container view stays")

	// Opening another view also stops a stream.
	panelID = h.click("Logs 🪵")
	testkit.Eventually(t, 2*time.Second, func() bool {
		return strings.Contains(h.tg.Message(panelID).Text, "tick")
	}, "logs never appeared")
	h.send("/images")
	mark = h.tg.Mark()
	time.Sleep(150 * time.Millisecond)
	for _, c := range h.tg.CallsSince(mark) {
		assert.NotEqual(t, "editMessageText", c.Method, "the stream kept editing after /images")
	}
}

func TestLogsKeepTheLastLines(t *testing.T) {
	d := fixtures()
	var many []string
	for i := 1; i <= 50; i++ {
		many = append(many, "line "+strings.Repeat("0", 3-len(itoa(i)))+itoa(i))
	}
	d.Logs = lines(many...)
	h := newHarness(t, d)
	h.send("/containers")
	panelID := h.click("Logs 🪵")

	testkit.Eventually(t, 2*time.Second, func() bool {
		return strings.Contains(h.tg.Message(panelID).Text, "line 050")
	}, "last line never appeared")
	text := h.tg.Message(panelID).Text
	assert.Contains(t, text, "line 031")
	assert.NotContains(t, text, "line 030", "only the last 20 lines are shown")
}

func TestStatsStream(t *testing.T) {
	d := fixtures()
	d.Stats = func(ctx context.Context, id string) io.ReadCloser {
		return io.NopCloser(strings.NewReader(`{"cpu_stats":{"cpu_usage":{"total_usage":400000000},"system_cpu_usage":20000000000,"online_cpus":2},` +
			`"precpu_stats":{"cpu_usage":{"total_usage":300000000},"system_cpu_usage":18000000000,"online_cpus":2},` +
			`"memory_stats":{"usage":52428800,"limit":1048576000,"stats":{"inactive_file":10485760}}}` + "\n"))
	}
	h := newHarness(t, d)
	h.send("/containers")
	panelID := h.click("Stats 📊")

	testkit.Eventually(t, 2*time.Second, func() bool {
		return strings.Contains(h.tg.Message(panelID).Text, "CPU usage")
	}, "stats never appeared: %v", h.tg.Message(panelID))
	text := h.tg.Message(panelID).Text
	assert.Contains(t, text, "📊 *web\\-1*")
	assert.Contains(t, text, "`CPU usage:   ` 10\\.00 %")
	assert.Contains(t, text, "`Memory:      ` 41\\.94 MegaBytes \\(4\\.00 %\\)")
	assert.Equal(t, [][]string{{"↩ Back"}}, h.tg.Message(panelID).Labels())
}

func TestSessionSurvivesRestart(t *testing.T) {
	d := fixtures()
	h := newHarness(t, d)
	h.send("/containers")
	h.click("Next ➡")
	h.click("Rename ✏️")

	h.start([]int64{owner.ID}) // same store, new process

	h.send("gateway")
	assert.Equal(t, "gateway", d.Containers()[1].Name, "the rename scene and its target survived the restart")
	h.click("Next ➡")
	assert.Contains(t, h.panel().Text, `stopped\.one`, "old buttons keep working")
}

func filter(calls []string, methods ...string) []string {
	var out []string
	for _, c := range calls {
		for _, m := range methods {
			if strings.HasPrefix(c, m+" ") || c == m {
				out = append(out, c)
			}
		}
	}
	return out
}

func itoa(i int) string {
	const digits = "0123456789"
	if i < 10 {
		return digits[i : i+1]
	}
	return itoa(i/10) + digits[i%10:i%10+1]
}
