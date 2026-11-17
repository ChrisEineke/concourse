package concourse

import (
	"context"
	"fmt"
	"log/slog"
	"os"

	docs "github.com/ChrisEineke/concourse/documents"
	"github.com/chigopher/pathlib"
	"github.com/moby/moby/client"
	"github.com/rogpeppe/rjson"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

type DeployerApp struct {
	config struct {
		verbose          bool          `mapstructure:"verbose"`
		DocumentFilepath *pathlib.Path `mapstructure:"document"`
	}
	vi          *viper.Viper
	rootSlogger *slog.Logger
	deployer    Deployer
}

func (app *DeployerApp) Configure() {
	pflag.Parse()

	app.vi = ViperSetup()
	app.vi.SetDefault("verbose", false)
	app.vi.SetDefault("document", "")
	app.vi.Set("document", pflag.Arg(0))
	if err := app.vi.Unmarshal(&app.config); err != nil {
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

	rawDocument, err := app.config.DocumentFilepath.ReadFile()
	if err != nil {
		panic(err)
	}
	containerSetDoc := docs.ContainerSetDocumentV_1_0{}
	if err = rjson.Unmarshal(rawDocument, &containerSetDoc); err != nil {
		panic(err)
	}
	fmt.Println(containerSetDoc)
	/*
	 *res := SchemaContainerSetV_1_0.Validate(containerSetDoc)
	 *if !res.Valid {
	 *    slog.Debug("containerSet document failed to validate", "errors", res.LocalizedDetailedErrors)
	 *    os.Exit(1)
	 *}
	 */

	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	apiClient, err := client.New(
		client.FromEnv,
		client.WithUserAgent("concourse/v1.0"),
	)
	if err != nil {
		panic(err)
	}
	defer apiClient.Close()

	if err := app.deployer.DeployContainerSet_v1_0(ctx, &containerSetDoc, apiClient); err != nil {
		panic(err)
	}
}
