package tui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
)

var (
	Title = lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FF6600")).
		MarginBottom(1)

	Success = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00CC00"))

	Error = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#FF0000"))

	Info = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#00AAFF"))

	Dim = lipgloss.NewStyle().
		Foreground(lipgloss.Color("#666666"))

	Separator = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#444444"))
)

func timestamp() string {
	return Dim.Render(time.Now().Format("15:04:05"))
}

func LogInfo(msg string) {
	fmt.Printf("%s %s %s\n", timestamp(), Info.Render("[INFO]"), msg)
}

func LogSuccess(msg string) {
	fmt.Printf("%s %s %s\n", timestamp(), Success.Render("[ OK ]"), msg)
}

func LogError(msg string) {
	fmt.Printf("%s %s %s\n", timestamp(), Error.Render("[ERRO]"), msg)
}

func LogWarn(msg string) {
	fmt.Printf("%s %s %s\n", timestamp(), lipgloss.NewStyle().Foreground(lipgloss.Color("#FFAA00")).Render("[WARN]"), msg)
}

func LogStep(msg string) {
	fmt.Printf("%s %s %s\n", timestamp(), Dim.Render("[----]"), msg)
}

func PrintHeader(title string) {
	fmt.Println()
	fmt.Println(Title.Render(title))
	fmt.Println(Separator.Render(strings.Repeat("─", 50)))
}

func PrintSeparator() {
	fmt.Println(Separator.Render(strings.Repeat("─", 50)))
}

func WaitForEnter() {
	fmt.Println()
	fmt.Print(Dim.Render("  Pressione Enter para continuar..."))
	fmt.Scanln()
}
