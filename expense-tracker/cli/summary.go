package cli

import (
	"github.com/spf13/cobra"

	"expense-tracker/storage"
)

var duration int

func NewSummaryCommand(path string) *cobra.Command {
	summaryCmd := &cobra.Command{
		Use: "summary",
		RunE: func(cmd *cobra.Command, args []string) error {
			return storage.SummaryExpense(path, duration)
		},
	}
	summaryCmd.Flags().IntVar(&duration, "month", -1, "Summary for the last months")
	return summaryCmd
}
