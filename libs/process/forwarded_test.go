package process

import (
	"bytes"
	"context"
	"os/exec"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestForwarded(t *testing.T) {
	ctx := context.Background()
	buf := bytes.NewBufferString("")
	err := Forwarded(ctx, []string{
		"python3", "-c", "print(input('input: '))",
	}, strings.NewReader("abc\n"), buf)
	assert.NoError(t, err)

	assert.Equal(t, "input: abc\n", buf.String())
}

func TestForwardedFails(t *testing.T) {
	ctx := context.Background()
	buf := bytes.NewBufferString("")
	err := Forwarded(ctx, []string{
		"_non_existent_",
	}, strings.NewReader("abc\n"), buf)
	assert.NotNil(t, err)
}

func TestForwardedFailsOnStdinPipe(t *testing.T) {
	ctx := context.Background()
	buf := bytes.NewBufferString("")
	err := Forwarded(ctx, []string{
		"_non_existent_",
	}, strings.NewReader("abc\n"), buf, func(c *exec.Cmd) error {
		c.Stdin = strings.NewReader("x")
		return nil
	})
	assert.NotNil(t, err)
}
