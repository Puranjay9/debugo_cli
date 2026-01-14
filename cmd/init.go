/*
Copyright © 2026 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"bufio"
	"debugo_cli/api"
	"debugo_cli/buildtree"
	"debugo_cli/metadata"
	"fmt"
	"os"
	"path/filepath"
	"strings"

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

		reader := bufio.NewReader(os.Stdin)
		fmt.Println("\nEnter project name: ")
		projectName, err := reader.ReadString('\n')
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
			os.Exit(1)
		}

		projectName = strings.TrimSpace(projectName)

		if projectName == "" {
			fmt.Fprintf(os.Stderr, "Project name cannot be empty\n")
			os.Exit(1)
		}

		//API call
		fmt.Printf("\nRegistering project '%s'...\n", projectName)
		projectID, err := api.CreateProjectOnServer(projectName)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error creating project on server: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("✓ Project registered successfully (ID: %s)\n", projectID)

		if err := metadata.SaveMetadata(debugoDir, projectID, projectName); err != nil {
			fmt.Fprintf(os.Stderr, "Error saving metadata: %v\n", err)
			os.Exit(1)
		}
		fmt.Println("✓ Metadata saved to .debugo/metadata.json")

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
	initCmd.Flags().StringP("output", "o", ".debugo", "Output directory for debugo files")
}
