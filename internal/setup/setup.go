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

var DefaultPermissions = []string{
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
	"Bash(npm run *)",
	"Bash(npx *)",
	"Bash(npm test *)",
	"Bash(next lint *)",
	"Bash(nuxi typecheck *)",
	"Bash(npx supabase *)",
	"Bash(npx tsc *)",
}

type InstallResult struct {
	Existed    bool
	BackupPath string
	DestPath   string
}

func GlobalCLAUDEPath() string {
	return filepath.Join(config.ClaudeDir(), "CLAUDE.md")
}

func SettingsPath() string {
	return filepath.Join(config.ClaudeDir(), "settings.json")
}

func GlobalCLAUDEExists() bool {
	return fileExists(GlobalCLAUDEPath())
}

func SettingsExists() bool {
	return fileExists(SettingsPath())
}

func ReadGlobalCLAUDE() (string, error) {
	data, err := os.ReadFile(GlobalCLAUDEPath())
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func GlobalCLAUDETemplate() (string, error) {
	data, err := templates.FS.ReadFile("files/global.md")
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func InstallGlobalCLAUDE() (*InstallResult, error) {
	claudeDir := config.ClaudeDir()
	if err := os.MkdirAll(claudeDir, 0o755); err != nil {
		return nil, fmt.Errorf("erro ao criar %s: %w", claudeDir, err)
	}

	dest := GlobalCLAUDEPath()
	result := &InstallResult{DestPath: dest}

	if fileExists(dest) {
		result.Existed = true
		backup := dest + ".bak"
		data, err := os.ReadFile(dest)
		if err != nil {
			return nil, fmt.Errorf("erro ao ler CLAUDE.md existente: %w", err)
		}
		if err := os.WriteFile(backup, data, 0o644); err != nil {
			return nil, fmt.Errorf("erro ao criar backup: %w", err)
		}
		result.BackupPath = backup
	}

	content, err := templates.FS.ReadFile("files/global.md")
	if err != nil {
		return nil, fmt.Errorf("erro ao ler template global: %w", err)
	}

	if err := os.WriteFile(dest, content, 0o644); err != nil {
		return nil, fmt.Errorf("erro ao escrever CLAUDE.md: %w", err)
	}

	return result, nil
}

func WriteGlobalCLAUDE(content string) error {
	dest := GlobalCLAUDEPath()
	claudeDir := config.ClaudeDir()
	if err := os.MkdirAll(claudeDir, 0o755); err != nil {
		return fmt.Errorf("erro ao criar %s: %w", claudeDir, err)
	}
	return os.WriteFile(dest, []byte(content), 0o644)
}

func InstallSettings() (*InstallResult, error) {
	settingsPath := SettingsPath()
	result := &InstallResult{DestPath: settingsPath}

	if fileExists(settingsPath) {
		result.Existed = true
		return result, nil
	}

	settings := Settings{
		Permissions: Permissions{Allow: DefaultPermissions},
	}

	data, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return nil, err
	}

	if err := os.WriteFile(settingsPath, data, 0o644); err != nil {
		return nil, fmt.Errorf("erro ao escrever settings.json: %w", err)
	}

	return result, nil
}

func ReadSettings() (*Settings, error) {
	data, err := os.ReadFile(SettingsPath())
	if err != nil {
		return nil, err
	}
	var s Settings
	if err := json.Unmarshal(data, &s); err != nil {
		return nil, err
	}
	return &s, nil
}

func MergeSettings(existing *Settings) *Settings {
	permSet := make(map[string]bool)
	for _, p := range existing.Permissions.Allow {
		permSet[p] = true
	}

	added := false
	for _, p := range DefaultPermissions {
		if !permSet[p] {
			existing.Permissions.Allow = append(existing.Permissions.Allow, p)
			permSet[p] = true
			added = added || true
		}
	}

	return existing
}

func WriteSettings(s *Settings) error {
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(SettingsPath(), data, 0o644)
}

func InstallProjectCLAUDE(dir string, stack Stack) (*InstallResult, error) {
	tmplName := stack.TemplateName()
	if tmplName == "" {
		return nil, fmt.Errorf("stack desconhecida, nao ha template disponivel")
	}

	content, err := templates.FS.ReadFile("files/" + tmplName)
	if err != nil {
		return nil, fmt.Errorf("erro ao ler template %s: %w", tmplName, err)
	}

	dest := filepath.Join(dir, "CLAUDE.md")
	result := &InstallResult{DestPath: dest}

	if fileExists(dest) {
		result.Existed = true
		return result, nil
	}

	if err := os.WriteFile(dest, content, 0o644); err != nil {
		return nil, fmt.Errorf("erro ao escrever CLAUDE.md: %w", err)
	}

	return result, nil
}
