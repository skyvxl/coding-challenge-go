package cli

import (
	"github.com/spf13/cobra"
)

func NewRootCommand(path string) *cobra.Command {
	rootCmd := &cobra.Command{
		Use:           "expense-tracker",
		SilenceErrors: true,
		CompletionOptions: cobra.CompletionOptions{
			DisableDefaultCmd: true,
		},
	}
	rootCmd.AddCommand(NewAddCommand(path))
	rootCmd.AddCommand(NewListCommand(path))
	rootCmd.AddCommand(NewSummaryCommand(path))
	rootCmd.AddCommand(NewDeleteCommand(path))
	return rootCmd
}
