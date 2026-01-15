package ui

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

// --- Paleta de Colores & Estilos ---

const (
	colorCarbon       = "#1a1a1a"
	colorEmerald      = "#04B575"
	colorBlue         = "#51AFEF"
	colorPastelYellow = "#FFFDF5"
	colorOrange       = "#FF8F00" // Nuevo color para Input Mode
	colorGray         = "#767676"
	colorLightGray    = "#A8A8A8"
	colorDarkGray     = "#2e2e2e"
	colorHighlight    = "#2B2B2B"
)

var (
	// Estilos de Breadcrumbs
	dirStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorGray))

	fileStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorPastelYellow)).
			Bold(true)

	// Estilos de Modos (Indicadores)
	editModeColor  = lipgloss.Color(colorEmerald)
	readModeColor  = lipgloss.Color(colorBlue)
	inputModeColor = lipgloss.Color(colorOrange)

	// 2. Styles for Manual Editor
	cursorLineStyle = lipgloss.NewStyle().
			Background(lipgloss.Color(colorHighlight)).
			Foreground(lipgloss.Color("#FFFFFF"))

	placeholderStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(colorGray))

	cursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorPastelYellow))

	// 3. Footer Styles
	statusBarStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorLightGray)).
			Background(lipgloss.Color(colorDarkGray))

	baseModeStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorCarbon)).
			Bold(true).
			Padding(0, 1).
			MarginRight(1)

	shortcutStyle = lipgloss.NewStyle().
			Bold(true)

	positionStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorPastelYellow)).
			Background(lipgloss.Color(colorDarkGray)).
			Padding(0, 1).
			MarginLeft(1)
)

// --- Modelo ---

type EditorModel struct {
	textarea     textarea.Model
	viewport     viewport.Model
	langInput    textinput.Model // Input para el lenguaje
	renderer     *glamour.TermRenderer
	filePath     string
	err          error
	notification string
	width        int
	height       int
	renderMode   bool
	askingLang   bool // ¿Estamos pidiendo el lenguaje?
}

func InitialModel(path string, content string) EditorModel {
	ti := textarea.New()
	ti.Placeholder = "# Empieza tu nueva nota..."
	ti.SetValue(content)
	ti.Focus()

	ti.Prompt = " "
	ti.ShowLineNumbers = true
	ti.FocusedStyle.CursorLine = cursorLineStyle
	ti.FocusedStyle.Placeholder = placeholderStyle
	ti.Cursor.Style = cursorStyle
	ti.FocusedStyle.Base = lipgloss.NewStyle().Foreground(lipgloss.Color("#E0E0E0"))
	ti.FocusedStyle.LineNumber = lipgloss.NewStyle().Foreground(lipgloss.Color(colorGray))

	vp := viewport.New(0, 0)
	vp.Style = lipgloss.NewStyle().Padding(0, 1)

	// Configurar Input de Lenguaje
	li := textinput.New()
	li.Placeholder = "go, js, py..."
	li.Prompt = "Lenguaje: "
	li.CharLimit = 20
	li.TextStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorOrange))
	li.PromptStyle = lipgloss.NewStyle().Foreground(lipgloss.Color(colorOrange)).Bold(true)

	renderer, _ := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithWordWrap(80),
	)

	return EditorModel{
		textarea:  ti,
		viewport:  vp,
		langInput: li,
		renderer:  renderer,
		filePath:  path,
	}
}

func (m EditorModel) Init() tea.Cmd {
	return tea.Batch(textarea.Blink, textinput.Blink)
}

func (m EditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		verticalMargins := 5
		newHeight := msg.Height - verticalMargins
		if newHeight < 0 {
			newHeight = 0
		}
		newWidth := msg.Width - 4

		m.textarea.SetWidth(newWidth)
		m.textarea.SetHeight(newHeight)
		m.viewport.Width = newWidth
		m.viewport.Height = newHeight

		m.renderer, _ = glamour.NewTermRenderer(
			glamour.WithStandardStyle("dark"),
			glamour.WithWordWrap(newWidth),
		)

		if m.renderMode {
			content, _ := m.renderer.Render(m.textarea.Value())
			m.viewport.SetContent(content)
		}

	case tea.KeyMsg:
		// Lógica exclusiva del Input Mode
		if m.askingLang {
			switch msg.String() {
			case "enter":
				lang := m.langInput.Value()
				if lang == "" {
					lang = "text"
				}
				
				// Inserción Inteligente con bloque cerrado
				snippet := fmt.Sprintf("\n```%s\n\n```", lang) 

				// Insertamos el snippet completo. El cursor quedará al final.
				m.textarea.InsertString(snippet)
				
				// Movemos el cursor ARRIBA una vez para quedar dentro del bloque
				m.textarea.CursorUp()

				m.askingLang = false
				m.langInput.Reset()
				
				return m, m.textarea.Focus()
				
			case "esc":
				m.askingLang = false
				m.langInput.Reset()
				return m, m.textarea.Focus()
			}
			m.langInput, cmd = m.langInput.Update(msg)
			return m, cmd
		}

		m.notification = ""
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.renderMode {
				m.renderMode = false
				// Devolver foco al salir de modo lectura
				return m, m.textarea.Focus()
			}
			return m, tea.Quit
		case "ctrl+s":
			err := os.WriteFile(m.filePath, []byte(m.textarea.Value()), 0o644)
			if err != nil {
				m.err = err
				m.notification = "❌ Error"
			} else {
				m.notification = "✅ Guardado"
			}
			return m, nil
		case "ctrl+p":
			m.renderMode = !m.renderMode
			if m.renderMode {
				content, err := m.renderer.Render(m.textarea.Value())
				if err != nil {
					m.viewport.SetContent("Error: " + err.Error())
				} else {
					m.viewport.SetContent(content)
				}
				// Al entrar a lectura, quizás queramos limpiar foco, o no importa.
				return m, nil
			} else {
				// Al volver, recuperar foco
				return m, m.textarea.Focus()
			}
		case "ctrl+k", "ctrl+o": // Aceptamos ambos por si acaso
			if !m.renderMode {
				m.askingLang = true
				return m, m.langInput.Focus()
			}
		}
	}

	if m.renderMode {
		m.viewport, cmd = m.viewport.Update(msg)
		cmds = append(cmds, cmd)
	} else {
		m.textarea, cmd = m.textarea.Update(msg)
		cmds = append(cmds, cmd)
	}

	return m, tea.Batch(cmds...)
}

func (m EditorModel) View() string {
	var accentColor lipgloss.Color
	var modeName string
	var icon string

	// Lógica de colores y estados
	if m.askingLang {
		accentColor = inputModeColor
		modeName = "INPUT"
		icon = "⌨️"
	} else if m.renderMode {
		accentColor = readModeColor
		modeName = "READER"
		icon = "👁️"
	} else {
		accentColor = editModeColor
		modeName = "EDITOR"
		icon = "✏️"
	}

	// 1. Breadcrumb Header
	dir, file := filepath.Split(m.filePath)
	if dir == "" {
		dir = "./"
	}

	headerContent := fmt.Sprintf("%s %s %s",
		dirStyle.Render(dir),
		lipgloss.NewStyle().Foreground(lipgloss.Color(colorDarkGray)).Render("›"),
		fileStyle.Render(file))

	header := lipgloss.NewStyle().
		Width(m.width - 4).
		Padding(0, 1).
		BorderBottom(true).
		BorderForeground(lipgloss.Color(colorDarkGray)).
		Render(headerContent)

	// 2. Content
	var content string
	if m.renderMode {
		content = m.viewport.View()
	} else {
		content = m.textarea.View()
	}

	// 3. Footer
	modeBadge := baseModeStyle.
		Background(accentColor).
		Render(fmt.Sprintf("%s %s", icon, modeName))

	cursorRow := m.textarea.Line() + 1
	totalLines := m.textarea.LineCount()
	textValue := m.textarea.Value()
	wordCount := len(strings.Fields(textValue))

	statsText := fmt.Sprintf("Ln %d/%d │ %dw", cursorRow, totalLines, wordCount)
	stats := lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorGray)).
		Render(statsText)

	// --- Centro del Footer (Dinámico) ---
	var centerMsg string
	
	if m.askingLang {
		// Mostrar INPUT FIELD
		centerMsg = m.langInput.View()
	} else if m.notification != "" {
		centerMsg = lipgloss.NewStyle().Foreground(accentColor).Bold(true).Render(m.notification)
	} else {
		centerMsg = lipgloss.NewStyle().Foreground(lipgloss.Color(colorDarkGray)).Render("^S Save  ^O Code  ^P View")
	}

	footerWidth := m.width - 4
	leftWidth := lipgloss.Width(modeBadge)
	rightWidth := lipgloss.Width(stats)
	centerWidth := lipgloss.Width(centerMsg)

	gap1 := (footerWidth - leftWidth - rightWidth - centerWidth) / 2
	if gap1 < 1 {
		gap1 = 1
	}

	gap2 := footerWidth - leftWidth - rightWidth - centerWidth - gap1
	if gap2 < 1 {
		gap2 = 1
	}

	footer := lipgloss.JoinHorizontal(lipgloss.Center,
		modeBadge,
		lipgloss.NewStyle().Width(gap1).Render(""),
		centerMsg,
		lipgloss.NewStyle().Width(gap2).Render(""),
		stats,
	)

	mainView := lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		content,
		lipgloss.NewStyle().BorderTop(true).BorderForeground(lipgloss.Color(colorDarkGray)).Width(footerWidth).Render(footer),
	)

	return lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(accentColor).
		Width(m.width - 2).
		Render(mainView)
}