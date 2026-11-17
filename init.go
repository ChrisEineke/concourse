package concourse

import (
	"log/slog"
	"os"
	"strings"

	"github.com/ChrisEineke/concourse/utils"
	"github.com/go-viper/mapstructure/v2"
	spretty "github.com/mickamy/slog-pretty"
	"github.com/spf13/viper"
)

func SlogSetup(lvl slog.Level, pretty bool) *slog.Logger {
	var handler slog.Handler
	if pretty {
		handler = spretty.NewHandler(os.Stdout, &spretty.HandlerOptions{
			Level:     lvl,
			AddSource: true,
		})
	} else {
		handler = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})
	}
	slogger := slog.New(handler)
	slog.SetDefault(slogger)
	return slogger
}

func ViperSetup() *viper.Viper {
	vi := viper.NewWithOptions(
		viper.ExperimentalBindStruct(),
		viper.WithDecodeHook(mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
			utils.StringToPathHookFunc(),
			utils.StringsToMultiaddrsHookFunc())))
	vi.SetEnvPrefix("CONCOURSE")
	vi.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	vi.AutomaticEnv()
	return vi
}
