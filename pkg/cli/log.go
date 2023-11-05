package cli

import (
	"os"

	isatty "github.com/mattn/go-isatty"
	"github.com/spf13/viper"
)

func IsTtyLogger() bool {
	switch f := viper.GetString("log.type"); f {
	case "tty":
		return true
	case "json":
		return false
	default:
		return isatty.IsTerminal(os.Stdout.Fd())
	}
}
