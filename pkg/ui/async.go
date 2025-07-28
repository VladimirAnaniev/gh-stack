package ui

import (
	"fmt"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type AsyncModel struct {
	spinner    spinner.Model
	message    string
	done       bool
	err        error
	successMsg string
	failureMsg string
}

type AsyncCompleteMsg struct {
	err error
}

func (m AsyncModel) Init() tea.Cmd {
	return m.spinner.Tick
}

func (m AsyncModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case AsyncCompleteMsg:
		m.done = true
		m.err = msg.err
		return m, tea.Quit
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
}

func (m AsyncModel) View() string {
	if m.done {
		if m.err != nil {
			errorStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
			return errorStyle.Render(fmt.Sprintf("❌ %s: %v", m.failureMsg, m.err)) + "\n"
		}
		return formatSuccess(m.successMsg)
	}

	return fmt.Sprintf("%s %s", m.spinner.View(), m.message)
}

func NewAsyncModel(message, successMsg, failureMsg string) AsyncModel {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return AsyncModel{
		spinner:    s,
		message:    message,
		done:       false,
		err:        nil,
		successMsg: successMsg,
		failureMsg: failureMsg,
	}
}

func RunAsync(message, successMsg, failureMsg string, operation func() error) error {
	p := tea.NewProgram(NewAsyncModel(message, successMsg, failureMsg))

	// Run operation in background
	go func() {
		err := operation()
		p.Send(AsyncCompleteMsg{err: err})
	}()

	model, err := p.Run()
	if err != nil {
		return err
	}

	// Return the operation error if it failed
	if asyncModel, ok := model.(AsyncModel); ok && asyncModel.err != nil {
		return asyncModel.err
	}

	return nil
}

func formatSuccess(message string) string {
	successStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("46"))
	return successStyle.Render(fmt.Sprintf("✅ %s", message)) + "\n"
}

func ShowSuccess(message string) {
	fmt.Print(formatSuccess(message))
}
