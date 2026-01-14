package buildtree

import (
	"encoding/json"
	"os"
	"path/filepath"
	"slices"
	"strings"
)

type FileNode struct {
	Path     string
	Type     string
	Language string
	Children []FileNode
}

func BuildTree(root string) (FileNode, error) {
	node := FileNode{
		Path: root,
		Type: "dir",
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		return node, err
	}

	for _, entry := range entries {
		name := entry.Name()
		if shouldIgnore(name) {
			continue
		}

		fullpath := filepath.Join(root, name)

		if entry.IsDir() {
			child, err := BuildTree(fullpath)

			if err == nil {
				node.Children = append(node.Children, child)
			}
		} else {
			node.Children = append(node.Children, FileNode{
				Path:     name,
				Type:     "file",
				Language: detectLanguage(name),
			})
		}
	}

	return node, nil
}

func shouldIgnore(name string) bool {
	ignored := []string{
		".git", "node_modules", "vendor",
		".debugo", "dist", "build",
	}

	return slices.Contains(ignored, name)
}

var extToLang = map[string]string{
	".go":    "go",
	".ts":    "typescript",
	".tsx":   "typescript",
	".js":    "javascript",
	".jsx":   "javascript",
	".py":    "python",
	".java":  "java",
	".kt":    "kotlin",
	".rs":    "rust",
	".cpp":   "cpp",
	".c":     "c",
	".h":     "c",
	".cs":    "csharp",
	".php":   "php",
	".rb":    "ruby",
	".swift": "swift",
	".scala": "scala",
	".lua":   "lua",
	".dart":  "dart",
	".sh":    "shell",
	".ps1":   "powershell",
	".sql":   "sql",
	".yaml":  "yaml",
	".yml":   "yaml",
	".json":  "json",
	".toml":  "toml",
}

func detectLanguage(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))

	if lang, ok := extToLang[ext]; ok {
		return lang
	}
	// if hasShebang(filename) {
	//     return detectFromShebang(filename)
	// }
	return "unknown"
}

func SaveTree(tree FileNode, debugoDir string) error {
	data, err := json.MarshalIndent(tree, "", "  ")
	if err != nil {
		return err
	}

	treePath := filepath.Join(debugoDir, "tree.json")
	return os.WriteFile(treePath, data, 0644)
}
