package concourse

import (
	"log/slog"
	"os"
	"strings"

	"github.com/go-viper/mapstructure/v2"
	"github.com/spf13/viper"
)

func SlogSetup(lvl slog.Level) {
	slog.SetDefault(
		slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			Level: slog.LevelDebug,
		})),
	)
}

func ViperSetup() *viper.Viper {
	vi := viper.NewWithOptions(
		viper.WithDecodeHook(mapstructure.ComposeDecodeHookFunc(
			mapstructure.StringToTimeDurationHookFunc(),
			mapstructure.StringToSliceHookFunc(","),
			StringToPathHookFunc(),
			StringsToMultiaddrsHookFunc())))
	vi.SetEnvPrefix("concourse")
	vi.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	vi.AutomaticEnv()
	return vi
}
