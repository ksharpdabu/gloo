package main

import (
	fdssetup "github.com/solo-io/gloo/projects/discovery/pkg/fds/setup"
	uds "github.com/solo-io/gloo/projects/discovery/pkg/uds/setup"
	"github.com/solo-io/solo-kit/pkg/utils/log"
	"github.com/solo-io/solo-kit/pkg/utils/stats"
	"os"
)

const (
	START_STATS_SERVER = "START_STATS_SERVER"
)

func main() {
	if os.Getenv(START_STATS_SERVER) != "" {
		stats.StartStatsServer()
	}
	if err := run(); err != nil {
		log.Fatalf("err in main: %v", err.Error())
	}
}

func run() error {
	errs := make(chan error)
	go func() {
		errs <- uds.Main()
	}()
	go func() {
		errs <- fdssetup.Main()
	}()
	return <-errs
}
