package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"syscall"

	"github.com/docker/libcontainer/system"
)

func main() {
	netnsPath := flag.String("net", "", "Path to the NET namespace to enter")
	flag.Parse()
	if flag.NArg() == 0 {
		fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] PATH [ARGS...]\n", os.Args[0])
		os.Exit(1)
	}
	if *netnsPath != "" {
		netns, err := os.Open(*netnsPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "exec: %v\n", err)
			os.Exit(1)
		}
		if err := system.Setns(netns.Fd(), syscall.CLONE_NEWNET); err != nil {
			fmt.Fprintf(os.Stderr, "setns: %v\n", err)
			os.Exit(1)
		}
	}
	var (
		cmdName = flag.Args()[0]
	)
	if !strings.Contains(cmdName, "/") {
		var err error
		cmdName, err = exec.LookPath(cmdName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%v\n", err)
			os.Exit(1)
		}
	}
	if err := syscall.Exec(cmdName, flag.Args(), os.Environ()); err != nil {
		fmt.Fprintf(os.Stderr, "exec: %v\n", err)
		os.Exit(1)
	}
}
