package storage

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"text/tabwriter"
	"time"
)

type Expense struct {
	Id          int    `json:"id"`
	Description string `json:"description"`
	Amount      int    `json:"amount"`
	CreatedAt   string `json:"created_at"`
}

func InitStorage() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", fmt.Errorf("get working dir: %w", err)
	}
	path := filepath.Join(dir, "expenses.json")
	f, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0o666)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return path, nil
		}
		return "", fmt.Errorf("create storage file: %w", err)
	}
	if _, err := f.WriteString("[]"); err != nil {
		_ = f.Close()
		return "", fmt.Errorf("initialize empty json array: %w", err)
	}
	if err := f.Close(); err != nil {
		return "", fmt.Errorf("close storage file: %w", err)
	}
	return path, nil
}

func saveAll(path string, expenses []Expense) error {
	data, err := json.MarshalIndent(expenses, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, data, 0o644)
}

func getAll(path string) ([]Expense, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return []Expense{}, err
	}
	var expenses []Expense
	err = json.Unmarshal(file, &expenses)
	return expenses, err
}

func AddExpense(path, description string, amount int) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("path is empty")
	}
	expenses, err := getAll(path)
	if err != nil {
		return err
	}
	expense := Expense{
		Id:          len(expenses) + 1,
		Description: description,
		Amount:      amount,
		CreatedAt:   time.Now().Format(time.DateOnly),
	}
	expenses = append(expenses, expense)
	err = saveAll(path, expenses)
	if err != nil {
		return err
	}
	fmt.Printf("Expense added successfully (ID: %d)\n", expense.Id)
	return nil
}

func ListExpense(path string) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("path is empty")
	}
	expenses, err := getAll(path)
	if err != nil {
		return err
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 0, 3, ' ', 0)
	fmt.Fprintf(w, "ID\tDate\tDescription\tAmount\n")
	for _, e := range expenses {
		fmt.Fprintf(w, "%d\t%s\t%s\t$%d\n", e.Id, e.CreatedAt, e.Description, e.Amount)
	}
	_ = w.Flush()
	return nil
}

func SummaryExpense(path string, dur int) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("path is empty")
	}
	expenses, err := getAll(path)
	if err != nil {
		return err
	}
	cutoffDate := time.Now().AddDate(0, -dur, 0)
	total := 0
	for _, e := range expenses {
		if dur > 0 {
			expenseDate, err := time.Parse("2006-01-02", e.CreatedAt)
			if err != nil {
				return err
			}
			if expenseDate.Before(cutoffDate) {
				continue
			}
		}
		total += e.Amount
	}
	fmt.Printf("Total expenses: $%d", total)
	return nil
}

func DeleteExpense(path string, id int) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("path is empty")
	}
	expenses, err := getAll(path)
	if err != nil {
		return err
	}
	for i, e := range expenses {
		if e.Id == id {
			expenses = append(expenses[:i], expenses[i+1:]...)
			err = saveAll(path, expenses)
			if err != nil {
				return err
			}
			fmt.Println("Expense deleted successfully")
			return nil
		}
	}
	return nil
}
