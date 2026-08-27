package cli

import (
	"github.com/spf13/cobra"

	"expense-tracker/storage"
)

var id int

func NewDeleteCommand(path string) *cobra.Command {
	deleteCmd := &cobra.Command{
		Use: "delete",
		RunE: func(cmd *cobra.Command, args []string) error {
			return storage.DeleteExpense(path, id)
		},
	}
	deleteCmd.Flags().IntVar(&id, "id", -1, "ID for deletion")
	_ = deleteCmd.MarkFlagRequired("id")
	return deleteCmd
}
