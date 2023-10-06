package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/liweiyi88/gti/config"
	"github.com/liweiyi88/gti/database"
	"github.com/liweiyi88/gti/search"
	"github.com/spf13/cobra"
	"golang.org/x/exp/slog"
)

func init() {
	rootCmd.AddCommand(searchCmd)
}

var searchCmd = &cobra.Command{
	Use:   "search [sync|delete]",
	Short: "sync or delete repositories in full text search",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		action := args[0]
		config.Init()
		ctx, stop := context.WithCancel(context.Background())
		db := database.GetInstance(ctx)
		handler := search.NewSearchHandler(db, search.NewSearch())

		defer func() {
			err := db.Close()

			if err != nil {
				slog.Error("failed to close db", slog.Any("error", err))
			}

			stop()
		}()

		appSignal := make(chan os.Signal, 3)
		signal.Notify(appSignal, os.Interrupt, syscall.SIGTERM)

		go func() {
			<-appSignal
			stop()
		}()

		return handler.Handle(ctx, action)
	},
}
