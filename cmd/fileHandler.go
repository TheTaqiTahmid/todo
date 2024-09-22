// This file is a module to handle the Todo data storage
// In our case, we are using simple CSV file as data storage
package cmd

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

// Main struct for representing the Todo data
type Record struct {
	index int    `csv:"Index"`
	todo  string `csv:"Title"`
	done  string `csv:"Done"`
}

// getDataFilePath get the file path for the Todo application and create the necessary
// directory tree
func getDataFilePath() string {
	basedir := filepath.Join(os.Getenv("HOME"), ".local", "share")
	applicationDir := filepath.Join(basedir, "Todo")
	err := os.MkdirAll(applicationDir, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	return filepath.Join(applicationDir, "todo.csv")
}

// listTodo list the saved todos in a tabular format
func listTodo(filePath string) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	r := csv.NewReader(f)

	records, err := r.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	printTable(records)
}

// printTable is a helper function to print the CSV file into stdout with
// proper formatting
func printTable(records [][]string) {
	colWidth := make([]int, len(records[0]))

	for _, record := range records {
		for i, col := range record {
			if len(col) > colWidth[i] {
				colWidth[i] = len(col)
			}
		}
	}
	fmt.Println()
	for _, record := range records {
		for i, col := range record {
			fmt.Printf("%-*s ", colWidth[i], col)
		}
		fmt.Println()
	}
	fmt.Println()
}

// countLines count number of lines in the csv file
func countLines(filePath string) int {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	reader := bufio.NewReader(f)
	linecount := 0

	for {
		_, _, err := reader.ReadLine()
		if err != nil {
			break
		}
		linecount++
	}
	return linecount

}

// addTodo adds a new todo to the database
func addTodo(filePath string, content []string) []string {
	f, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	w := csv.NewWriter(f)
	defer w.Flush()

	title := []string{
		"Index", "Title", "Done",
	}
	linecount := countLines(filePath)
	if linecount <= 0 {
		if err := w.Write(title); err != nil {
			log.Fatal(err)
		}
	}

	var todo Record

	if linecount <= 1 {
		todo.index = 1
	} else {
		todo.index = linecount
	}

	todo.todo = strings.Join(content, " ")
	todo.done = "false"

	todoString := []string{strconv.Itoa(todo.index), todo.todo, todo.done}

	w.Write(todoString)

	fmt.Println("Todo Added")
	return todoString
}

// deleteTodo deletes a todo from the database by index
func deleteTodo(filePath string, index int) Record {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	reader := csv.NewReader(f)
	r, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var recordArray []Record
	for _, record := range r[1:] {
		todoIndex, err := strconv.Atoi(record[0])
		if err != nil {
			log.Fatal("Invalid index value: ", err)
		}
		recordArray = append(recordArray, Record{index: todoIndex, todo: record[1], done: record[2]})
	}
	removedRecord := recordArray[index-1]
	newRecordArray := append(recordArray[:index-1], recordArray[index:]...)

	f, err = os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	writer := csv.NewWriter(f)
	defer writer.Flush()

	title := []string{
		"Index", "Title", "Done",
	}

	writer.Write(title)

	for i, record := range newRecordArray {
		newIndex := strconv.Itoa(i + 1) // reindexing the Todos
		writeString := []string{newIndex, record.todo, record.done}
		writer.Write(writeString)
	}

	fmt.Println("Todo Removed")
	return removedRecord
}

// markTodoAsDone marks a todo as done by index by index
func markTodoAsDone(filePath string, index int) {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	reader := csv.NewReader(f)
	r, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	var recordArray []Record
	for _, record := range r[1:] {
		index, err := strconv.Atoi(record[0])
		if err != nil {
			log.Fatal(err)
		}
		recordArray = append(recordArray, Record{index: index, todo: record[1], done: record[2]})
	}
	recordToEdit := &recordArray[index-1]
	recordToEdit.done = "true"

	f, err = os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	writer := csv.NewWriter(f)
	defer writer.Flush()
	writer.Write(r[0])

	for _, record := range recordArray {
		index := strconv.Itoa(record.index)
		writeString := []string{index, record.todo, record.done}
		writer.Write(writeString)
	}
}

// updateTodo
func updateTodo(filePath string, index int, todo []string) Record {
	f, err := os.Open(filePath)
	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	reader := csv.NewReader(f)
	r, err := reader.ReadAll()
	if err != nil {
		log.Fatal(err)
	}
	var recordArray []Record
	for _, record := range r[1:] {
		index, err := strconv.Atoi(record[0])
		if err != nil {
			log.Fatal(err)
		}
		recordArray = append(recordArray, Record{index: index, todo: record[1], done: record[2]})
	}

	f, err = os.OpenFile(filePath, os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	writer := csv.NewWriter(f)
	defer writer.Flush()

	recordToEdit := &recordArray[index-1]
	var todoString string
	for _, str := range todo {
		todoString += str + " "
	}
	recordToEdit.todo = todoString

	writer.Write(r[0])

	for _, record := range recordArray {
		index := strconv.Itoa(record.index)
		writeString := []string{index, record.todo, record.done}
		writer.Write(writeString)
	}
	fmt.Println("Todo Updated")
	return *recordToEdit
}
