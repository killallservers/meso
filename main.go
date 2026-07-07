package main

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const repoURL = "https://github.com/killallservers/templates/archive/refs/heads/main.zip"

func main() {
	// Parse command-line arguments
	templateName := ""
	for i, arg := range os.Args[1:] {
		if arg == "--template" && i+1 < len(os.Args)-1 {
			templateName = os.Args[i+2]
			break
		}
	}

	// Create temp directory for downloading
	tmpDir, err := os.MkdirTemp("", "meso-*")
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to create temp directory: %v\n", err)
		os.Exit(1)
	}
	defer os.RemoveAll(tmpDir)

	// Fetch templates
	fmt.Fprintf(os.Stderr, "📥 Fetching templates from GitHub...\n")
	templatesDir, err := fetchTemplates(tmpDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to fetch templates: %v\n", err)
		os.Exit(1)
	}

	// List available templates
	templates, err := listTemplates(templatesDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to list templates: %v\n", err)
		os.Exit(1)
	}

	if len(templates) == 0 {
		fmt.Fprintf(os.Stderr, "❌ No templates found\n")
		os.Exit(1)
	}

	// If template specified via flag, use it directly
	if templateName != "" {
		found := false
		for _, t := range templates {
			if t == templateName {
				found = true
				break
			}
		}
		if !found {
			fmt.Fprintf(os.Stderr, "❌ Template \"%s\" not found\n", templateName)
			fmt.Fprintf(os.Stderr, "Available templates: %s\n", strings.Join(templates, ", "))
			os.Exit(1)
		}

		// Use minimal info when using flag
		minimalInfo := PromptModel{
			inputs: [2]string{templateName, ""},
		}

		if err := copyTemplate(templatesDir, templateName, minimalInfo); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to copy template: %v\n", err)
			os.Exit(1)
		}
		return
	}

	// Show interactive selector
	selectorModel := NewSelector(templates)
	p := tea.NewProgram(selectorModel)
	selectorResult, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}

	selectedTemplates := selectorResult.(Selector).Selected()
	if len(selectedTemplates) == 0 {
		fmt.Fprintf(os.Stderr, "❌ No templates selected\n")
		os.Exit(1)
	}

	// Show fancy prompts for scaffold info
	promptModel := NewPromptModel()
	p = tea.NewProgram(promptModel)
	promptResult, err := p.Run()
	if err != nil {
		fmt.Fprintf(os.Stderr, "❌ Error: %v\n", err)
		os.Exit(1)
	}

	promptInfo := promptResult.(PromptModel)

	// Create agent directory first
	agentDir, ok := agentDirs[promptInfo.agent]
	if !ok {
		agentDir = ".claude"
	}
	if err := os.MkdirAll(agentDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "❌ Failed to create agent directory: %v\n", err)
		os.Exit(1)
	}

	// Scaffold each selected template
	renamed := false
	for i, templateName := range selectedTemplates {
		if err := copyTemplate(templatesDir, templateName, promptInfo); err != nil {
			fmt.Fprintf(os.Stderr, "❌ Failed to copy template %s: %v\n", templateName, err)
			os.Exit(1)
		}

		// Rename .claude to agent dir after first template
		if i == 0 && !renamed && agentDir != ".claude" {
			claudeDir := ".claude"
			if _, err := os.Stat(claudeDir); err == nil {
				if err := os.Rename(claudeDir, agentDir); err == nil {
					fmt.Fprintf(os.Stderr, "📁 Renamed .claude → %s\n", agentDir)
					renamed = true
				}
			}
		}
	}

	// Create local config file
	if err := createProjectConfig(promptInfo); err != nil {
		fmt.Fprintf(os.Stderr, "⚠️  Warning: Could not create config: %v\n", err)
	}

	// Print summary
	fmt.Fprintf(os.Stderr, "\n✨ Template scaffolded successfully!\n")
	fmt.Fprintf(os.Stderr, "\nAgent:       %s\n", promptInfo.agent)
	fmt.Fprintf(os.Stderr, "Project:     %s\n", promptInfo.inputs[0])
	if promptInfo.inputs[1] != "" {
		fmt.Fprintf(os.Stderr, "Description: %s\n", promptInfo.inputs[1])
	}

	configDir := agentDir
	if configDir == "" {
		configDir = ".claude"
	}
	fmt.Fprintf(os.Stderr, "Config dir:  %s/\n", configDir)

	fmt.Fprintf(os.Stderr, "\nNext steps:\n")
	fmt.Fprintf(os.Stderr, "  1. cd %s\n", promptInfo.inputs[0])
	fmt.Fprintf(os.Stderr, "  2. go mod tidy (or bun install if using Bun)\n")
	if promptInfo.agent == "claude-code" {
		fmt.Fprintf(os.Stderr, "  3. claude code\n")
	} else if promptInfo.agent == "opencode" {
		fmt.Fprintf(os.Stderr, "  3. opencode\n")
	} else if promptInfo.agent == "vscode" {
		fmt.Fprintf(os.Stderr, "  3. code .\n")
	} else {
		fmt.Fprintf(os.Stderr, "  3. Open in your IDE\n")
	}
}

// Selector is the interactive template selector model
type Selector struct {
	templates []string
	checked   map[int]bool
	cursor    int
}

func NewSelector(templates []string) Selector {
	return Selector{
		templates: templates,
		checked:   make(map[int]bool),
		cursor:    0,
	}
}

func (s Selector) Init() tea.Cmd {
	return nil
}

func (s Selector) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "up", "k":
			if s.cursor > 0 {
				s.cursor--
			}
		case "down", "j":
			if s.cursor < len(s.templates)-1 {
				s.cursor++
			}
		case " ":
			// Toggle current item
			s.checked[s.cursor] = !s.checked[s.cursor]
		case "enter":
			// Require at least one selection
			if len(s.checked) == 0 {
				return s, nil
			}
			return s, tea.Quit
		case "ctrl+c":
			return s, tea.Quit
		}
	}
	return s, nil
}

func (s Selector) View() string {
	var output strings.Builder
	output.WriteString("Select templates to scaffold:\n\n")

	for i, template := range s.templates {
		cursor := "  "
		checkbox := "☐"

		if i == s.cursor {
			cursor = "> "
		}

		if s.checked[i] {
			checkbox = "☑"
		}

		line := cursor + checkbox + " " + template

		if i == s.cursor {
			line = lipgloss.NewStyle().
				Background(lipgloss.Color("4")).
				Foreground(lipgloss.Color("15")).
				Padding(0, 1).
				Render(line)
		}

		output.WriteString(line)
		output.WriteString("\n")
	}

	output.WriteString("\n(Use ↑/↓ or j/k to navigate, Space to toggle, Enter to confirm)\n")

	// Show selected count
	selectedCount := len(s.checked)
	if selectedCount > 0 {
		countStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color("2")).
			Bold(true)
		output.WriteString("\n" + countStyle.Render(fmt.Sprintf("Selected: %d template(s)", selectedCount)))
	}

	return output.String()
}

func (s Selector) Selected() []string {
	var selected []string
	for i := 0; i < len(s.templates); i++ {
		if s.checked[i] {
			selected = append(selected, s.templates[i])
		}
	}
	return selected
}

func fetchTemplates(tmpDir string) (string, error) {
	// Download the zip file
	zipPath := filepath.Join(tmpDir, "templates.zip")
	resp, err := http.Get(repoURL)
	if err != nil {
		return "", fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	out, err := os.Create(zipPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	if _, err := io.Copy(out, resp.Body); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	// Extract the zip file
	fmt.Fprintf(os.Stderr, "📦 Extracting templates...\n")
	reader, err := zip.OpenReader(zipPath)
	if err != nil {
		return "", fmt.Errorf("failed to open zip: %w", err)
	}
	defer reader.Close()

	// Find the templates-main directory and extract it
	var extractDir string
	for _, file := range reader.File {
		if strings.HasPrefix(file.Name, "templates-main/") && !strings.HasSuffix(file.Name, "/") {
			// Extract file
			path := filepath.Join(tmpDir, file.Name)
			if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
				return "", fmt.Errorf("failed to create dir: %w", err)
			}

			rc, err := file.Open()
			if err != nil {
				return "", fmt.Errorf("failed to open file in zip: %w", err)
			}
			defer rc.Close()

			out, err := os.Create(path)
			if err != nil {
				return "", fmt.Errorf("failed to create extracted file: %w", err)
			}
			defer out.Close()

			if _, err := io.Copy(out, rc); err != nil {
				return "", fmt.Errorf("failed to write extracted file: %w", err)
			}
		}

		if extractDir == "" && strings.HasPrefix(file.Name, "templates-main/") {
			extractDir = filepath.Join(tmpDir, "templates-main")
		}
	}

	if extractDir == "" {
		return "", fmt.Errorf("could not find templates directory in zip")
	}

	return extractDir, nil
}

func listTemplates(templatesDir string) ([]string, error) {
	entries, err := os.ReadDir(templatesDir)
	if err != nil {
		return nil, err
	}

	var templates []string
	for _, entry := range entries {
		if entry.IsDir() && !strings.HasPrefix(entry.Name(), ".") {
			templates = append(templates, entry.Name())
		}
	}

	sort.Strings(templates)
	return templates, nil
}

func copyTemplate(sourceDir, templateName string, info PromptModel) error {
	source := filepath.Join(sourceDir, templateName)
	target := "."

	fmt.Fprintf(os.Stderr, "📋 Copying %s to current directory...\n", templateName)

	// Copy the template directory
	if err := copyDir(source, target); err != nil {
		return err
	}

	// Remove .git directory if it exists
	gitDir := filepath.Join(target, ".git")
	if _, err := os.Stat(gitDir); err == nil {
		if err := os.RemoveAll(gitDir); err == nil {
			fmt.Fprintf(os.Stderr, "🗑️  Removed .git directory\n")
		}
	}

	return nil
}

func createProjectMetadata(target string, info PromptModel) error {
	agentDir, ok := agentDirs[info.agent]
	if !ok {
		agentDir = ".claude"
	}

	metadata := fmt.Sprintf(`# Project Metadata

**Name:** %s
**Agent:** %s
**Created:** %s
**Description:** %s
`,
		info.inputs[0],
		info.agent,
		"Scaffolded with Meso",
		info.inputs[1],
	)

	metaDir := filepath.Join(target, agentDir)
	metaFile := filepath.Join(metaDir, "PROJECT.md")

	return os.WriteFile(metaFile, []byte(metadata), 0o644)
}

func copyDir(src, dst string) error {
	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Calculate relative path
		rel, err := filepath.Rel(src, path)
		if err != nil {
			return err
		}

		target := filepath.Join(dst, rel)

		if info.IsDir() {
			return os.MkdirAll(target, 0o755)
		}

		// Copy file
		src, err := os.Open(path)
		if err != nil {
			return err
		}
		defer src.Close()

		dst, err := os.Create(target)
		if err != nil {
			return err
		}
		defer dst.Close()

		_, err = io.Copy(dst, src)
		return err
	})
}
