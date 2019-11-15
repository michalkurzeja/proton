package main

import (
	"github.com/spf13/cobra"
)

func init() {
	rootCmd.PersistentFlags().StringSliceVarP(&importPaths, flagImportPaths, "I", nil, "proto import paths")
	rootCmd.PersistentFlags().StringVarP(&protoFile, flagProtoFile, "f", "", "proto file path (required)")

	rootCmd.MarkFlagRequired(flagProtoFile)

	rootCmd.AddCommand(convertCmd)
}

// Flag names.
const (
	flagImportPaths = "import-paths"
	flagProtoFile   = "proto-file"
)

// Config values.
var (
	importPaths []string
	protoFile   string
)

var rootCmd = &cobra.Command{
	Use: "proton",
}

// Execute executes the root command.
func Execute() error {
	return rootCmd.Execute()
}
