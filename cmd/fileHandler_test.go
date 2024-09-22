package cmd

import (
	"encoding/csv"
	"fmt"
	"os"
	"testing"
)

func createTempFile(t *testing.T, content [][]string) *os.File {
	t.Helper()

	// Create a temporary file
	tempFile, err := os.CreateTemp("", "todo_test_*.csv")
	if err != nil {
		t.Fatalf("Failed to create temp file %v", err)
	}
	writer := csv.NewWriter(tempFile)
	defer writer.Flush()

	for _, record := range content {
		if err := writer.Write(record); err != nil {
			t.Fatalf("Could not write to temp file %v", err)
		}
	}
	return tempFile
}

func TestAddTodo(t *testing.T) {
	tempFile := createTempFile(t, [][]string{
		{"Index", "Title", "Done"},
	})

	defer os.Remove(tempFile.Name())

	content := []string{"Add", "dummy", "todo"}
	addedTodo := addTodo(tempFile.Name(), content)

	f, err := os.Open(tempFile.Name())
	if err != nil {
		t.Fatalf("Failed to open created temp file %v", err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	records, err := reader.ReadAll()
	fmt.Println(records)
	if err != nil {
		t.Fatalf("Failed to read the todos from temp file %v", err)
	}

	if len(records) != 2 {
		t.Errorf("Expected 2 records, got %d", len(records))
	}

	expectedToDo := []string{"1", "Add dummy todo", "false"}
	if !equalSlices(expectedToDo, addedTodo) {
		t.Errorf("Expected todo %v, got todo %v", expectedToDo, addedTodo)
	}
}

func TestListTodo(t *testing.T) {
	tempFile := createTempFile(t, [][]string{
		{"Index", "Title", "Done"},
		{"1", "Learn Go", "true"},
		{"2", "Write tests", "false"},
	})

	defer os.Remove(tempFile.Name())
	listTodo(tempFile.Name())
}

func TestDeleteTodo(t *testing.T) {
	tempFile := createTempFile(t, [][]string{
		{"Index", "Title", "Done"},
		{"1", "Learn Go", "true"},
		{"2", "Write tests", "false"},
	})

	defer os.Remove(tempFile.Name())
	todoIndex := 1
	got := deleteTodo(tempFile.Name(), todoIndex)
	expeted := Record{
		index: 1,
		todo: "Learn Go",
		done: "true",
	}
	if got != expeted {
		t.Errorf("Expected %v, got %v", expeted, got)
	}
}

// Helper function to compare slice
func equalSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i := range len(a) {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
