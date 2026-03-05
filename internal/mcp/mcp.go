package mcp

import (
	"fmt"
	"os/exec"
	"strings"
)

type MCPServer struct {
	Name        string
	Description string
	Scope       string // "user" or "project"
	Command     string
	Args        []string
	EnvVars     map[string]string
	NeedsBinary string // binary that must exist in PATH
}

var AvailableMCPs = []MCPServer{
	{
		Name:        "dart-flutter",
		Description: "Dart/Flutter: analyzer, pub.dev, hot reload, widget tree",
		Scope:       "user",
		Command:     "dart",
		Args:        []string{"mcp-server", "--force-roots-fallback"},
		NeedsBinary: "dart",
	},
	{
		Name:        "context7",
		Description: "Documentacao up-to-date de qualquer biblioteca",
		Scope:       "user",
		Command:     "npx",
		Args:        []string{"-y", "@upstash/context7-mcp@latest"},
		NeedsBinary: "npx",
	},
	{
		Name:        "github",
		Description: "GitHub: issues, PRs, branches, reviews",
		Scope:       "user",
		Command:     "npx",
		Args:        []string{"-y", "@modelcontextprotocol/server-github"},
		EnvVars:     map[string]string{"GITHUB_TOKEN": ""},
		NeedsBinary: "npx",
	},
	{
		Name:        "playwright",
		Description: "Browser automation e testes de UI",
		Scope:       "user",
		Command:     "npx",
		Args:        []string{"-y", "@playwright/mcp@latest"},
		NeedsBinary: "npx",
	},
	{
		Name:        "postgres",
		Description: "PostgreSQL: queries diretas, schema explorer",
		Scope:       "project",
		Command:     "npx",
		Args:        []string{"-y", "@modelcontextprotocol/server-postgres"},
		NeedsBinary: "npx",
	},
	{
		Name:        "primevue",
		Description: "PrimeVue: componentes, docs, theming",
		Scope:       "user",
		Command:     "npx",
		Args:        []string{"-y", "@anthropic/primevue-mcp@latest"},
		NeedsBinary: "npx",
	},
}

func ListInstalled() ([]string, error) {
	out, err := exec.Command("claude", "mcp", "list").CombinedOutput()
	if err != nil {
		return nil, fmt.Errorf("falha ao listar MCPs: %w", err)
	}
	var installed []string
	for _, line := range strings.Split(string(out), "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) > 0 {
			installed = append(installed, strings.TrimSpace(parts[0]))
		}
	}
	return installed, nil
}

func IsInstalled(name string, installed []string) bool {
	for _, n := range installed {
		if n == name {
			return true
		}
	}
	return false
}

func Install(m MCPServer) error {
	if m.NeedsBinary != "" {
		if _, err := exec.LookPath(m.NeedsBinary); err != nil {
			return fmt.Errorf("%s nao encontrado no PATH", m.NeedsBinary)
		}
	}

	// Resolve full path for command
	cmdPath, err := exec.LookPath(m.Command)
	if err != nil {
		return fmt.Errorf("comando %s nao encontrado: %w", m.Command, err)
	}

	args := []string{"mcp", "add", "--transport", "stdio", "--scope", m.Scope}

	for k, v := range m.EnvVars {
		args = append(args, "-e", fmt.Sprintf("%s=%s", k, v))
	}

	args = append(args, m.Name, "--")
	args = append(args, cmdPath)
	args = append(args, m.Args...)

	cmd := exec.Command("claude", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("falha ao instalar %s: %s — %w", m.Name, string(out), err)
	}
	return nil
}

func Remove(name string) error {
	cmd := exec.Command("claude", "mcp", "remove", name)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("falha ao remover %s: %s — %w", name, string(out), err)
	}
	return nil
}
