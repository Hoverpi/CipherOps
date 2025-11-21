package utils

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"bytes"
)

func runCmd(ctx context.Context, sudo bool, name string, args ...string) (string, error) {
	if DryRun {
		cmdLine := name + " " + strings.Join(args, " ")
		if sudo {
			cmdLine = "sudo " + cmdLine
		}
		fmt.Println("[DRYRUN]", cmdLine)
		return cmdLine, nil
	}

	if sudo && os.Geteuid() != 0 {
		// intentamos con sudo
		args = append([]string{name}, args...)
		name = "sudo"
	}
	cmd := exec.CommandContext(ctx, name, args...)
	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr
	err := cmd.Run()
	outStr := strings.TrimSpace(out.String())
	if err != nil {
		return outStr, fmt.Errorf("%w: %s", err, strings.TrimSpace(stderr.String()))
	}
	return outStr, nil
}