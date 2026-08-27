package cli

import (
	"github.com/spf13/cobra"

	"expense-tracker/storage"
)

func NewListCommand(path string) *cobra.Command {
	return &cobra.Command{
		Use: "list",
		RunE: func(cmd *cobra.Command, args []string) error {
			return storage.ListExpense(path)
		},
	}
}
