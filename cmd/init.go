/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"debugo_cli/buildtree"
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize debugo in your project repo",
	Long: `initialize debugo in your project repo by 

- Creating a .debugo directory
- Building a file tree of your project
- Saving the tree structure to .debugo/tree.json`,
	Run: func(cmd *cobra.Command, args []string) {
		rootDir := "."

		if len(args) > 0 {
			rootDir = args[0]
		}

		absPath, err := filepath.Abs(rootDir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error resolving path: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Initializing debugo in: %s\n", absPath)

		debugoDir := filepath.Join(absPath, ".debugo")
		if err := os.MkdirAll(debugoDir, 0755); err != nil {
			fmt.Fprintf(os.Stderr, "Error creating .debugo directory: %v\n", err)
			os.Exit(1)
		}

		fmt.Println("✓ Created .debugo directory")

		fmt.Println("Building project tree...")
		tree, err := buildtree.BuildTree(absPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error building tree: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✓ Project tree built successfully")

		if err := buildtree.SaveTree(tree, debugoDir); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving tree: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✓ Tree saved to .debugo/tree.json")

		fmt.Println("\n🎉 Debugo initialized successfully!")

	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	initCmd.Flags().StringP("output", "o", ".debugo", "Output directory for debugo files")
}
