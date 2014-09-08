package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/docker/libcontainer/system"
)

func main() {
	if err := serveCommands(os.Stdin); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

func serveCommands(input io.Reader) error {
	cmds := json.NewDecoder(input)
	for {
		var cmd []string
		if err := cmds.Decode(&cmd); err == io.EOF {
			return nil
		} else if err != nil {
			return fmt.Errorf("json decode: %v", err)
		}
		if len(cmd) == 0 {
			fmt.Fprintf(os.Stderr, "skipping empty command\n", cmd)
			continue
		}
		if err := do(cmd[0], cmd[1:]...); err != nil {
			return err
		}
	}
}

func do(name string, args ...string) error {
	switch name {
	case "netns":
		{
			return doNetns(args...)
		}
	case "mntns":
		{
			return doMntns(args...)
		}
	case "exec":
		{
			return doExec(args...)
		}
	case "setenv":
		{
			return doSetenv(args...)
		}
	}
	return fmt.Errorf("no such command: %s", name)
}

func doSetenv(args ...string) error {
	if len(args) != 2 {
		return fmt.Errorf("usage: setenv KEY VALUE")
	}
	os.Setenv(args[0], args[1])
	return nil
}

func doExec(args ...string) error {
	if len(args) < 1 {
		return fmt.Errorf("usage: exec NAME [ARGS...]")
	}
	var cmdName = args[0]
	if !strings.Contains(cmdName, "/") {
		var err error
		cmdName, err = exec.LookPath(cmdName)
		if err != nil {
			return err
		}
	}
	return syscall.Exec(cmdName, args, os.Environ())
}

func doNetns(args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: netns PATH")
	}
	var nsPath = args[0]
	ns, err := os.Open(nsPath)
	if err != nil {
		return err
	}
	return system.Setns(ns.Fd(), syscall.CLONE_NEWNET)
}

func doMntns(args ...string) error {
	if len(args) != 1 {
		return fmt.Errorf("usage: mntns PATH")
	}
	var nsPath = args[0]
	ns, err := os.Open(nsPath)
	if err != nil {
		return err
	}
	return system.Setns(ns.Fd(), syscall.CLONE_NEWNS)
}
