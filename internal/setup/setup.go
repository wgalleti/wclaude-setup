package setup

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/wgalleti/wclaude-setup/internal/config"
	"github.com/wgalleti/wclaude-setup/internal/templates"
)

type Permissions struct {
	Allow []string `json:"allow"`
}

type Settings struct {
	Permissions Permissions `json:"permissions"`
}

var defaultPermissions = []string{
	"Bash(git diff:*)",
	"Bash(git log:*)",
	"Bash(git status:*)",
	"Bash(git branch:*)",
	"Bash(python manage.py shell_plus)",
	"Bash(python manage.py showmigrations)",
	"Bash(pytest *)",
	"Bash(uv run pytest *)",
	"Bash(dart analyze)",
	"Bash(dart format *)",
	"Bash(flutter analyze)",
	"Bash(flutter test *)",
	"Bash(go test *)",
	"Bash(go vet *)",
	"Bash(go build *)",
}

func InstallGlobalCLAUDE() error {
	claudeDir := config.ClaudeDir()
	if err := os.MkdirAll(claudeDir, 0o755); err != nil {
		return fmt.Errorf("erro ao criar %s: %w", claudeDir, err)
	}

	dest := filepath.Join(claudeDir, "CLAUDE.md")

	if fileExists(dest) {
		backup := dest + ".bak"
		data, err := os.ReadFile(dest)
		if err != nil {
			return fmt.Errorf("erro ao ler CLAUDE.md existente: %w", err)
		}
		if err := os.WriteFile(backup, data, 0o644); err != nil {
			return fmt.Errorf("erro ao criar backup: %w", err)
		}
		fmt.Printf("  Backup criado em %s\n", backup)
	}

	content, err := templates.FS.ReadFile("files/global.md")
	if err != nil {
		return fmt.Errorf("erro ao ler template global: %w", err)
	}

	if err := os.WriteFile(dest, content, 0o644); err != nil {
		return fmt.Errorf("erro ao escrever CLAUDE.md: %w", err)
	}

	fmt.Printf("  CLAUDE.md instalado em %s\n", dest)
	return nil
}

func InstallSettings() error {
	settingsPath := filepath.Join(config.ClaudeDir(), "settings.json")

	if fileExists(settingsPath) {
		fmt.Printf("  settings.json ja existe — nao alterado\n")
		return nil
	}

	settings := Settings{
		Permissions: Permissions{Allow: defaultPermissions},
	}

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(settingsPath, data, 0o644); err != nil {
		return fmt.Errorf("erro ao escrever settings.json: %w", err)
	}

	fmt.Printf("  settings.json criado em %s\n", settingsPath)
	return nil
}

func InstallProjectCLAUDE(dir string, stack Stack) error {
	tmplName := stack.TemplateName()
	if tmplName == "" {
		return fmt.Errorf("stack desconhecida, nao ha template disponivel")
	}

	content, err := templates.FS.ReadFile("files/" + tmplName)
	if err != nil {
		return fmt.Errorf("erro ao ler template %s: %w", tmplName, err)
	}

	dest := filepath.Join(dir, "CLAUDE.md")
	if fileExists(dest) {
		return fmt.Errorf("CLAUDE.md ja existe em %s — use 'claude-setup merge' para atualizar", dir)
	}

	if err := os.WriteFile(dest, content, 0o644); err != nil {
		return fmt.Errorf("erro ao escrever CLAUDE.md: %w", err)
	}

	fmt.Printf("  CLAUDE.md criado em %s (template: %s)\n", dest, stack.Label())
	return nil
}
