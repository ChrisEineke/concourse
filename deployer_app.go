package concourse

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	"github.com/chigopher/pathlib"
	"github.com/moby/moby/client"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type DeployerApp struct {
	config struct {
		verbose          bool          `mapstructure:"verbose"`
		DocumentFilepath *pathlib.Path `mapstructure:"document"`
	}
	rootCtx         context.Context
	rootCtxCancelFn func()
	vi              *viper.Viper
	rootSlogger     *slog.Logger
	dockerApiClient *client.Client
	deployer        *Deployer
}

func (app *DeployerApp) Configure() {
	var err error

	app.rootCtx, app.rootCtxCancelFn = context.WithCancel(context.Background())

	pflag.Parse()

	app.vi = ViperSetup()
	app.vi.SetDefault("verbose", false)
	app.vi.Set("document", pflag.Arg(0))
	if err = app.vi.Unmarshal(&app.config); err != nil {
		panic(err)
	}

	if app.config.DocumentFilepath.String() == "" {
		fmt.Println("expected 1 argument, got 0")
		os.Exit(1)
	}

	if app.config.verbose {
		SlogSetup(slog.LevelDebug, true)
	} else {
		SlogSetup(slog.LevelInfo, false)
	}

	app.dockerApiClient, err = client.New(
		client.FromEnv,
		client.WithUserAgent("concourse/v1.0"),
	)
	if err != nil {
		panic(err)
	}

	app.deployer = NewDeployer(app.dockerApiClient)
}

func (app *DeployerApp) Run() {
	if exists, err := app.config.DocumentFilepath.Exists(); !exists || err != nil {
		if err != nil {
			panic(err)
		}
		fmt.Println("document filepath does not exist:", app.config.DocumentFilepath.String())
		os.Exit(1)
	}
	slog.Debug("document", "filepath", app.config.DocumentFilepath.String())

	if err := app.deployer.DispatchDocument(app.rootCtx, app.config.DocumentFilepath); err != nil {
		panic(err)
	}
}

func (app *DeployerApp) Stop() {
	defer app.rootCtxCancelFn()
	defer app.dockerApiClient.Close()
}
