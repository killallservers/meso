package main

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var agents = []string{
	"claude-code",
	"opencode",
	"cursor",
	"vscode",
}

var agentDirs = map[string]string{
	"claude-code": ".claude",
	"opencode":   ".opencode",
	"cursor":     ".cursor",
	"vscode":     ".vscode",
}

// PromptModel is the Bubble Tea model for the scaffold prompts
type PromptModel struct {
	step      int // 0=agent, 1=name, 2=description, 3=confirm
	agent     string
	agentIdx  int
	name      string
	desc      string
	inputs    [2]string
	focused   int
	err       string
}

func NewPromptModel() PromptModel {
	return PromptModel{
		step:    0,
		focused: 0,
		inputs:  [2]string{"", ""},
	}
}

func (m PromptModel) Init() tea.Cmd {
	return nil
}

func (m PromptModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "up", "k":
			if m.step == 0 {
				if m.agentIdx > 0 {
					m.agentIdx--
				}
			} else if m.step < 3 && m.step > 0 {
				m.focused = (m.focused + 1) % 2
			}

		case "down", "j":
			if m.step == 0 {
				if m.agentIdx < len(agents)-1 {
					m.agentIdx++
				}
			} else if m.step < 3 && m.step > 0 {
				m.focused = (m.focused + 1) % 2
			}

		case "tab", "shift+tab":
			if m.step > 0 && m.step < 3 {
				m.focused = (m.focused + 1) % 2
			}

		case "enter":
			if m.step == 0 {
				// Select agent
				m.agent = agents[m.agentIdx]
				m.step = 1
				m.focused = 0
				m.err = ""
			} else if m.step == 1 {
				// Validate project name
				name := strings.TrimSpace(m.inputs[0])
				if name == "" {
					m.err = "Project name cannot be empty"
					return m, nil
				}
				m.name = name
				m.step = 2
				m.focused = 0
				m.err = ""
			} else if m.step == 2 {
				// Description is optional
				m.desc = strings.TrimSpace(m.inputs[1])
				m.step = 3
			} else if m.step == 3 {
				// Confirmed
				return m, tea.Quit
			}

		default:
			// Type character
			if m.step > 0 && m.step < 3 {
				m.inputs[m.step-1] += msg.String()
			} else if m.step == 3 && (msg.String() == "y" || msg.String() == "Y") {
				return m, tea.Quit
			} else if m.step == 3 && (msg.String() == "n" || msg.String() == "N") {
				m.step = 0
				m.agentIdx = 0
				m.focused = 0
				m.inputs = [2]string{"", ""}
			}
		}
	}

	return m, nil
}

func (m PromptModel) View() string {
	var output strings.Builder

	// Header
	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("11")).
		Render("✨ Scaffold Your Project")

	output.WriteString(title)
	output.WriteString("\n\n")

	if m.step == 0 {
		output.WriteString(m.renderAgentStep())
	} else if m.step == 1 {
		output.WriteString(m.renderNameStep())
	} else if m.step == 2 {
		output.WriteString(m.renderDescriptionStep())
	} else {
		output.WriteString(m.renderConfirmStep())
	}

	// Error message
	if m.err != "" {
		errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
		output.WriteString("\n")
		output.WriteString(errorStyle.Render("❌ " + m.err))
	}

	return output.String()
}

func (m PromptModel) renderAgentStep() string {
	var output strings.Builder

	output.WriteString("Which agent/IDE are you using?\n\n")

	selectedStyle := lipgloss.NewStyle().
		Background(lipgloss.Color("11")).
		Foreground(lipgloss.Color("0")).
		Padding(0, 1)

	normalStyle := lipgloss.NewStyle().
		Padding(0, 1)

	for i, agent := range agents {
		var style lipgloss.Style
		prefix := "  "

		if i == m.agentIdx {
			style = selectedStyle
			prefix = "> "
		} else {
			style = normalStyle
		}

		output.WriteString(prefix)
		output.WriteString(style.Render(agent))
		output.WriteString("\n")
	}

	output.WriteString("\n")
	output.WriteString(lipgloss.NewStyle().Faint(true).Render("(Use ↑/↓ or j/k to navigate, Enter to continue)"))

	return output.String()
}

func (m PromptModel) renderNameStep() string {
	var output strings.Builder

	output.WriteString("What's your project name?\n\n")

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("4")).
		Padding(0, 1)

	focusedStyle := inputStyle.
		BorderForeground(lipgloss.Color("11"))

	style := inputStyle
	if m.focused == 0 {
		style = focusedStyle
	}

	input := m.inputs[0]
	if input == "" {
		input = "my-project"
	}

	output.WriteString(style.Render(input))
	output.WriteString("\n\n")
	output.WriteString(lipgloss.NewStyle().Faint(true).Render("(Enter to continue)"))

	return output.String()
}

func (m PromptModel) renderDescriptionStep() string {
	var output strings.Builder

	output.WriteString("Describe your project (optional):\n\n")

	inputStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("4")).
		Padding(0, 1)

	focusedStyle := inputStyle.
		BorderForeground(lipgloss.Color("11"))

	style := inputStyle
	if m.focused == 0 {
		style = focusedStyle
	}

	input := m.inputs[1]
	if input == "" {
		input = "A new Claude Code project..."
	}

	output.WriteString(style.Render(input))
	output.WriteString("\n\n")
	output.WriteString(lipgloss.NewStyle().Faint(true).Render("(Enter to continue)"))

	return output.String()
}

func (m PromptModel) renderConfirmStep() string {
	var output strings.Builder

	// Summary card
	summaryStyle := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("2")).
		Padding(1, 2)

	desc := m.desc
	if desc == "" {
		desc = "(none)"
	}

	summary := fmt.Sprintf("Agent:       %s\nProject:     %s\nDescription: %s",
		m.agent,
		m.name,
		desc)

	output.WriteString(summaryStyle.Render(summary))
	output.WriteString("\n\n")

	output.WriteString("Proceed with scaffolding? ")
	output.WriteString(lipgloss.NewStyle().Bold(true).Render("[y/n]"))
	output.WriteString("\n")

	return output.String()
}
