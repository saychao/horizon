package main

import (
	"github.com/spf13/cobra"
)

var dbCmd = &cobra.Command{
	Use:   "db [command]",
	Short: "commands to manage horizon's dbs",
}

func init() {
	dbCmd.AddCommand(historyCmd)
}
