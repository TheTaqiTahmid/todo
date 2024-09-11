/*
Copyright © 2024 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"log"
	"strconv"

	"github.com/spf13/cobra"
)

// updateCmd represents the update command
var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a todo",
	Long: `Update a todo by index. The update string has to be wrapped around "" or ''`,
	Run: func(cmd *cobra.Command, args []string) {
    if len(args) > 2 {
      log.Fatal("Too many input arguments")
    }
		index, err := strconv.Atoi(args[0])
    if err != nil {
      log.Fatal(err)
    }
    updateTodo(todoFileDir, index, args[1])
	},
}

func init() {
	rootCmd.AddCommand(updateCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// updateCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// updateCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
