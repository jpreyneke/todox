package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/cobra"
	"go.uber.org/fx"

	"todox/internal"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	rootCmd := &cobra.Command{
		Use:   "todox",
		Short: "Todo API Server",
	}

	apiCmd := &cobra.Command{
		Use:   "api",
		Short: "Run the API server",
		Run: func(cmd *cobra.Command, args []string) {
			app := fx.New(
				internal.Module,
				fx.NopLogger,
			)
			app.Run()
		},
	}

	rootCmd.AddCommand(apiCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
