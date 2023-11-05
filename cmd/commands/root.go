package commands

import (
	"fmt"
	"os"

	"github.com/inngest/inngest/pkg/api/tel"
	"github.com/inngest/inngest/pkg/cli"
	isatty "github.com/mattn/go-isatty"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	letters = `
    ____                            __
   /  _/___  ____  ____ ____  _____/ /_
   / // __ \/ __ \/ __ '/ _ \/ ___/ __/
 _/ // / / / / / / /_/ /  __(__  ) /_
/___/_/ /_/_/ /_/\__, /\___/____/\__/
                /____/
`
)

var (
	longDescription = fmt.Sprintf(
		"%s\n%s\n%s%s\n",
		cli.TextStyle.Render(letters),
		cli.TextStyle.Render("Build event-driven queues with zero infra. "),
		cli.TextStyle.Render("Request features, get help, and chat with us: "),
		cli.BoldStyle.Render("https://www.inngest.com/discord"),
	)
)

func Execute() {
	rootCmd := &cobra.Command{
		Use:   "inngest",
		Short: "A serverless event-driven infrastructure platform",
		Long:  longDescription,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			m := tel.NewMetadata(cmd.Context())
			m.SetCobraCmd(cmd)
			tel.SendMetadata(cmd.Context(), m)
		},
		PersistentPostRun: func(cmd *cobra.Command, args []string) {
			// Wait for any events to have been sent.
			tel.Wait()
		},
	}

	var isTTY = isatty.IsTerminal(os.Stdout.Fd())

	rootCmd.PersistentFlags().Bool("prod", false, "Use the production environment for the current command.")
	rootCmd.PersistentFlags().Bool("json", !isTTY, "Output logs as JSON. Defaults to true if stdout is not a TTY. Override with --json=false to force pretty printing.")
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose logging.")

	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		panic(err)
	}

	if err := rootCmd.PersistentFlags().Parse(os.Args[1:]); err == nil {
		if viper.GetBool("json") {
			viper.Set("log.type", "json")
		} else {
			viper.Set("log.type", "tty")
		}
	}

	// Register Top Level Commands
	rootCmd.AddCommand(NewCmdLogin())
	rootCmd.AddCommand(NewCmdDev())
	rootCmd.AddCommand(NewCmdVersion())
	rootCmd.AddCommand(NewCmdServe())

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
