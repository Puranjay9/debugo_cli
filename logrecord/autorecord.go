package logrecord

import (
	"debugo_cli/metadata"
	_ "embed"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

//go:embed scipts/hook.sh
var hookContent string

func CreateHookShell() error {

	cwd, _ := os.Getwd()
	debugoDir := filepath.Join(cwd, ".debugo")
	metaData, err := metadata.LoadMetadata(debugoDir)
	if err != nil {
		return err
	}

	if !metaData.IsAutoRecord {
		return nil
	}

	home, _ := os.UserHomeDir()
	hookDir := filepath.Join(home, ".debugo")
	hookPath := filepath.Join(hookDir, "shell-hook.sh")

	if err := os.MkdirAll(hookDir, 0755); err != nil {
		return err
	}

	if err := os.WriteFile(hookPath, []byte(hookContent), 0755); err != nil {
		return err
	}

	shell := os.Getenv("SHELL")
	var rc string

	if strings.Contains(shell, "zsh") {
		rc = filepath.Join(home, ".zshrc")
	} else if strings.Contains(shell, "bash") {
		rc = filepath.Join(home, ".bashrc")
	} else {
		return fmt.Errorf("unsupported shell")
	}

	sourceLine := "\n# Debugo shell integration\n[ -f \"$HOME/.debugo/shell-hook.sh\" ] && source \"$HOME/.debugo/shell-hook.sh\"\n"
	rcBytes, _ := os.ReadFile(rc)

	if !strings.Contains(string(rcBytes), ".debugo/shell-hook.sh") {
		f, err := os.OpenFile(rc, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			return err
		}
		defer f.Close()
		f.WriteString(sourceLine)
	}

	fmt.Println("Auto record initialized. Restart yout shell or run source ~/.bashrc for bash")

	return nil
}
