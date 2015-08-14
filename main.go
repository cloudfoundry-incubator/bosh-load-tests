package main

import (
	"os"
	"os/signal"
	"syscall"

	boshlog "github.com/cloudfoundry/bosh-utils/logger"
	boshsys "github.com/cloudfoundry/bosh-utils/system"

	bltaction "github.com/cloudfoundry-incubator/bosh-load-tests/action"
	bltclirunner "github.com/cloudfoundry-incubator/bosh-load-tests/action/clirunner"
	bltconfig "github.com/cloudfoundry-incubator/bosh-load-tests/config"
	bltenv "github.com/cloudfoundry-incubator/bosh-load-tests/environment"
	bltflow "github.com/cloudfoundry-incubator/bosh-load-tests/flow"
)

func main() {
	if len(os.Args) != 2 {
		println("Usage: blt path/to/config.json")
		os.Exit(1)
	}

	logger := boshlog.NewLogger(boshlog.LevelDebug)
	fs := boshsys.NewOsFileSystem(logger)
	cmdRunner := boshsys.NewExecCmdRunner(logger)

	config := bltconfig.NewConfig(fs)
	err := config.Load(os.Args[1])
	if err != nil {
		panic(err)
	}

	logger.Debug("main", "Setting up environment")
	environmentProvider := bltenv.NewProvider(config, fs, cmdRunner)
	environment := environmentProvider.Get()
	err = environment.Setup()
	if err != nil {
		panic(err)
	}
	defer environment.Shutdown()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)
	go func() {
		<-c
		environment.Shutdown()
		os.Exit(1)
	}()

	logger.Debug("main", "Starting deploy")

	cliRunnerFactory := bltclirunner.NewFactory(config.CliCmd, cmdRunner, fs)

	directorInfo, err := bltaction.NewDirectorInfo(environment.DirectorURL(), cliRunnerFactory)
	if err != nil {
		panic(err)
	}

	actionFactory := bltaction.NewFactory(directorInfo, fs)

	actionsFlow := bltflow.NewFlow(1, []string{"prepare"}, actionFactory, cliRunnerFactory)
	err = actionsFlow.Run()
	if err != nil {
		panic(err)
	}

	doneCh := make(chan error)

	for i := 0; i < len(config.Flows); i++ {
		go func(i int) {
			actionNames := config.Flows[i]
			logger.Debug("main", "Creating flow with %#v", actionNames)
			flow := bltflow.NewFlow(i, actionNames, actionFactory, cliRunnerFactory)
			doneCh <- flow.Run()
		}(i)
	}

	for i := 0; i < len(config.Flows); i++ {
		err = <-doneCh
		if err != nil {
			panic(err)
		}
	}

	println("Done!")
}
