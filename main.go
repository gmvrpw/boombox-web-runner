package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"

	"gmvr.pw/boombox-web-runner/config"
	"gmvr.pw/boombox-web-runner/internal/controller/http"
	runtimeModuleRepository "gmvr.pw/boombox-web-runner/internal/repository/module/runtime"
	runtimeRunnerRepository "gmvr.pw/boombox-web-runner/internal/repository/runner/runtime"
	"gmvr.pw/boombox-web-runner/internal/service/runner"
)

func main() {
	var err error

	logger := slog.Default()

	cfg := &config.Config{}

	path := os.Getenv("BOOMBOX_WEB_RUNNER_CONFIG_FILE")
	cfg, err = config.NewConfig(path)
	if err != nil {
		log.Fatalf("cannot get config. %s", err)
	}

	runnerRepository, err := runtimeRunnerRepository.NewRuntimeRunnerRepository()
	if err != nil {
		log.Fatalf("cannot create runner repository. %s", err)
	}

	moduleRepository, err := runtimeModuleRepository.NewRuntimeModuleRepository(cfg.Modules)
	if err != nil {
		log.Fatalf("cannot create module repository. %s", err)
	}

	runnerService, err := runner.NewRunnerService(logger)
	if err != nil {
		log.Fatalf("cannot create runner service. %s", err)
	}

	err = runnerService.Init(runnerRepository, moduleRepository)
	if err != nil {
		log.Fatalf("cannot initialize runner service. %s", err)
	}

	httpController, err := http.NewHttpRunnerController(&cfg.Controllers.HTTP, logger)
	if err != nil {
		log.Fatalf("cannot create http controller. %s", err)
	}

	err = httpController.Init(runnerService)
	if err != nil {
		log.Fatalf("cannot initialize http controller. %s", err)
	}

	go func() {
		err = httpController.Serve()
		fmt.Println("shutdown")
		if err != nil {
			log.Fatalf("failed when serving http. %s", err)
		}
	}()

	stop := make(chan os.Signal, 1)

	signal.Notify(stop, os.Interrupt)
	<-stop
}
