package cmd

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestParseIDs(t *testing.T) {
	ids, err := parseIDs(" 22, 33,,44 ")
	require.NoError(t, err)
	assert.Equal(t, []int64{22, 33, 44}, ids)

	ids, err = parseIDs("")
	require.NoError(t, err)
	assert.Empty(t, ids)

	_, err = parseIDs("22,abc")
	assert.EqualError(t, err, `"abc" is not a chat ID`)
}

// clearEnv unsets every variable start reads, for the test's duration.
func clearEnv(t *testing.T) {
	for _, k := range []string{"TOKEN", "TELEGRAM__TOKEN", "WHITELISTED_IDS", "DB__DSN", "TELEGRAM__URL"} {
		t.Setenv(k, "") // registers the restore
		os.Unsetenv(k)
	}
}

func TestFromEnvDefaultsAndPrecedence(t *testing.T) {
	clearEnv(t)
	t.Chdir(t.TempDir()) // no .env

	opts := startOptions{}
	assert.EqualError(t, opts.fromEnv(), "telegram can't start because no token is provided (--token or TOKEN)")

	t.Setenv("TOKEN", "legacy")
	t.Setenv("WHITELISTED_IDS", "1,2")
	opts = startOptions{}
	require.NoError(t, opts.fromEnv())
	assert.Equal(t, "legacy", opts.token, "the original TOKEN variable still works")
	assert.Equal(t, []int64{1, 2}, opts.whitelist)
	assert.Equal(t, DefaultDSN, opts.dsn)

	t.Setenv("TELEGRAM__TOKEN", "new")
	t.Setenv("DB__DSN", "postgres://u:p@h/db")
	t.Setenv("TELEGRAM__URL", "http://localhost:8081")
	opts = startOptions{}
	require.NoError(t, opts.fromEnv())
	assert.Equal(t, "new", opts.token, "TELEGRAM__TOKEN wins over TOKEN")
	assert.Equal(t, "postgres://u:p@h/db", opts.dsn)
	assert.Equal(t, "http://localhost:8081", opts.apiURL)

	// Flags win over the environment.
	opts = startOptions{token: "flag", whitelist: []int64{9}, dsn: "file:x.db"}
	require.NoError(t, opts.fromEnv())
	assert.Equal(t, "flag", opts.token)
	assert.Equal(t, []int64{9}, opts.whitelist)
	assert.Equal(t, "file:x.db", opts.dsn)
}

func TestFromEnvReadsDotEnv(t *testing.T) {
	clearEnv(t)
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(dir+"/.env", []byte("TOKEN=from-file\nWHITELISTED_IDS=5\n"), 0o600))
	t.Chdir(dir)
	t.Cleanup(func() { os.Unsetenv("TOKEN"); os.Unsetenv("WHITELISTED_IDS") })

	opts := startOptions{}
	require.NoError(t, opts.fromEnv())
	assert.Equal(t, "from-file", opts.token)
	assert.Equal(t, []int64{5}, opts.whitelist)
}

func TestFromEnvRejectsBadWhitelist(t *testing.T) {
	clearEnv(t)
	t.Chdir(t.TempDir())
	t.Setenv("TOKEN", "x")
	t.Setenv("WHITELISTED_IDS", "1,me")
	opts := startOptions{}
	assert.EqualError(t, opts.fromEnv(), `WHITELISTED_IDS: "me" is not a chat ID`)
}
