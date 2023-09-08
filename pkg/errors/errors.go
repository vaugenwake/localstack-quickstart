package errors

import (
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/table"
)

type Error struct {
	Level   string
	Message string
}

type ErrorsBag struct {
	errors []Error
}

func (eb *ErrorsBag) Add(level string, message string) {
	eb.errors = append(eb.errors, Error{
		Level:   level,
		Message: message,
	})
}

func (eb *ErrorsBag) Any() bool {
	return len(eb.errors) > 0
}

func (eb *ErrorsBag) All() []Error {
	return eb.errors
}

func (eb *ErrorsBag) Get(index int) Error {
	return eb.errors[index]
}

func (eb *ErrorsBag) PrintErrors() {
	if eb.Any() {
		t := table.NewWriter()
		t.SetTitle("Execution Errors")

		t.AppendHeader(table.Row{"#", "Level", "Error"})

		for idx, err := range eb.All() {
			t.AppendRow(table.Row{idx, err.Level, err.Message})
		}

		fmt.Println(t.Render())
		os.Exit(1)
	}
}
