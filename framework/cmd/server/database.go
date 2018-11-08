package server

import "github.com/spf13/cobra"

var databaseCmd = &cobra.Command{
	Use:   "database",
	Short: "Database operations for MLModelScope evaluations",
}

func init() {
	databaseCmd.PersistentFlags().StringVar(&databaseAddress, "database_address", "", "address of the database")
	databaseCmd.PersistentFlags().StringVar(&databaseName, "database_name", "", "name of the database to use")

	databaseCmd.AddCommand(divergenceCmds...)
}
