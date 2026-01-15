package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
)

type EditorModel struct {
	textarea textarea.Model
	filePath string
	err      error
}

func InitialModel(path string, content string) EditorModel {
	ti := textarea.New()
	ti.Placeholder = "Empieza a escribir tu nota..."
	ti.SetValue(content)
	ti.Focus()

	ti.SetWidth(80)
	ti.SetHeight(20)

	return EditorModel{
		ti,
		path,
		nil,
	}
}

func (m EditorModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m EditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "esc": // salir sin guardar
			return m, tea.Quit
		case "ctrl+s":
			err := os.WriteFile(m.filePath, []byte(m.textarea.Value()), 0o644)
			if err != nil {
				m.err = err
			}
			return m, tea.Quit
		}
	}

	m.textarea, cmd = m.textarea.Update(msg)
	return m, cmd
}

func (m EditorModel) View() string {
	return fmt.Sprintf("Editando: %s\n\n%s\n\n%s", m.filePath, m.textarea.View(), "(ctrl+s: guardar • esc: salir)")
}
