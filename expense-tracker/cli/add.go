package cli

import (
	"github.com/spf13/cobra"

	"expense-tracker/storage"
)

var (
	description string
	amount      int
)

func NewAddCommand(path string) *cobra.Command {
	addCmd := &cobra.Command{
		Use: "add",
		RunE: func(cmd *cobra.Command, args []string) error {
			return storage.AddExpense(path, description, amount)
		},
	}
	addCmd.Flags().StringVar(&description, "description", "", "Description of expenses")
	addCmd.Flags().IntVar(&amount, "amount", 0, "Amount of expenses")

	_ = addCmd.MarkFlagRequired("description")
	_ = addCmd.MarkFlagRequired("amount")
	return addCmd
}
