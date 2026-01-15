package ui

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

// --- Paleta de Colores & Estilos ---

const (
	colorCarbon       = "#1a1a1a"
	colorEmerald      = "#04B575"
	colorPastelYellow = "#FFFDF5"
	colorGray         = "#767676"
	colorLightGray    = "#A8A8A8"
	colorDarkGray     = "#2e2e2e"
	colorHighlight    = "#2B2B2B"
)

var (
	// Estilo base de la aplicación
	appStyle = lipgloss.NewStyle().
		Padding(0)

	// 1. Header Inmersivo
	headerStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorEmerald)).
		Background(lipgloss.Color(colorCarbon)).
		Bold(true).
		Padding(0, 1).
		BorderBottom(true).
		BorderStyle(lipgloss.ThickBorder()).
		BorderForeground(lipgloss.Color(colorHighlight))

	// 2. Styles for Manual Editor (Minimalist)
	cursorLineStyle = lipgloss.NewStyle().
		Background(lipgloss.Color(colorHighlight)).
		Foreground(lipgloss.Color("#FFFFFF"))

	placeholderStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorGray))

	cursorStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorPastelYellow))

	// 3. Footer Styles (LazyVim / Airline inspired)
	statusBarStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorLightGray)).
		Background(lipgloss.Color(colorDarkGray))

	modeStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorCarbon)).
		Background(lipgloss.Color(colorEmerald)).
		Bold(true).
		Padding(0, 1).
		MarginRight(1)

	shortcutStyle = lipgloss.NewStyle().
		Foreground(lipgloss.Color(colorEmerald)).
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
	renderer     *glamour.TermRenderer
	filePath     string
	err          error
	notification string // Mensaje temporal para el usuario
	width        int
	height       int
	renderMode   bool
}

func InitialModel(path string, content string) EditorModel {
	// Configurar Editor
	ti := textarea.New()
	ti.Placeholder = "# Empieza tu nueva nota..."
	ti.SetValue(content)
	ti.Focus()

	ti.Prompt = " "
	ti.ShowLineNumbers = true // "IDE Feel" needs line numbers
	ti.FocusedStyle.CursorLine = cursorLineStyle
	ti.FocusedStyle.Placeholder = placeholderStyle
	ti.Cursor.Style = cursorStyle
	ti.FocusedStyle.Base = lipgloss.NewStyle().Foreground(lipgloss.Color("#E0E0E0"))
	ti.FocusedStyle.LineNumber = lipgloss.NewStyle().Foreground(lipgloss.Color(colorGray))

	// Configurar Viewport
	vp := viewport.New(0, 0)
	vp.Style = lipgloss.NewStyle().Padding(0, 1) // Padding para lectura cómoda

	// Inicializar renderer de Glamour (Dark Mode)
	renderer, _ := glamour.NewTermRenderer(
		glamour.WithStandardStyle("dark"),
		glamour.WithWordWrap(80), // Default safe wrap
	)

	return EditorModel{
		textarea: ti,
		viewport: vp,
		renderer: renderer,
		filePath: path,
	}
}

func (m EditorModel) Init() tea.Cmd {
	return textarea.Blink
}

func (m EditorModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Ajuste de layout dinámico
	headerHeight := 2
		footerHeight := 2
		newHeight := msg.Height - headerHeight - footerHeight

		if newHeight < 0 {
			newHeight = 0
		}

		m.textarea.SetWidth(msg.Width)
		m.textarea.SetHeight(newHeight)

		m.viewport.Width = msg.Width
		m.viewport.Height = newHeight

		// Re-iniciar renderer con el ancho correcto para word-wrapping
		m.renderer, _ = glamour.NewTermRenderer(
			glamour.WithStandardStyle("dark"),
			glamour.WithWordWrap(msg.Width-4), // -4 por paddings
		)

		if m.renderMode {
			content, _ := m.renderer.Render(m.textarea.Value())
			m.viewport.SetContent(content)
		}

	case tea.KeyMsg:
		// Limpiar notificación al escribir
		m.notification = ""

		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "esc":
			if m.renderMode {
				m.renderMode = false
				m.textarea.Focus()
				return m, nil
			}
			return m, tea.Quit // Doble esc para salir
		case "ctrl+s":
			err := os.WriteFile(m.filePath, []byte(m.textarea.Value()), 0o644)
			if err != nil {
				m.err = err
				m.notification = "❌ Error al guardar"
			} else {
				m.notification = "✅ Archivo Guardado"
			}
			// No salimos (tea.Quit), solo notificamos
			return m, nil
		case "ctrl+p":
			m.renderMode = !m.renderMode
			if m.renderMode {
				// Renderizado PRO con Glamour
				content, err := m.renderer.Render(m.textarea.Value())
				if err != nil {
					m.viewport.SetContent("Error rendering markdown: " + err.Error())
				} else {
					m.viewport.SetContent(content)
				}
			} else {
				m.textarea.Focus()
			}
			return m, nil
		}

		if !m.renderMode {
			switch msg.String() {
			case "ctrl+t":
				m.textarea.InsertString("# ")
				return m, nil // Stop propagation
			case "ctrl+l":
				m.textarea.InsertString("- ")
				return m, nil // Stop propagation
			case "ctrl+k":
				m.textarea.InsertString("\n```go\n\n```") // Default to go
				m.textarea.CursorUp()
				m.textarea.CursorUp()
				return m, nil // Stop propagation
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
	// 1. Header
	headerText := fmt.Sprintf("📝 GO-NOTES  |  %s", m.filePath)
	header := headerStyle.Width(m.width).Render(headerText)

	// 2. Content
	var content string
	if m.renderMode {
		content = m.viewport.View()
	} else {
		content = m.textarea.View()
	}

	// 3. Smart Footer
	modeText := "EDITOR"
	if m.renderMode {
		modeText = "READER"
	}
	mode := modeStyle.Render(modeText)
	fileInfo := statusBarStyle.Render(m.filePath)

	// Posición del Cursor (Línea)
	cursorRow := m.textarea.Line() + 1
	totalLines := m.textarea.LineCount()
	posText := fmt.Sprintf("Ln %d  %d%%", cursorRow, int(float64(cursorRow)/float64(totalLines)*100))
	position := positionStyle.Render(posText)
	
	// Área Central del Footer (Notificaciones o Atajos)
	var centerContent string
	
	if m.notification != "" {
		// Mostrar Notificación Prioritaria
		centerContent = lipgloss.NewStyle().
			Foreground(lipgloss.Color(colorEmerald)).
			Bold(true).
			Render(m.notification)
	} else {
		// Mostrar Atajos
		baseShortcuts := fmt.Sprintf("%s Guardar %s Salir %s Vista",
			shortcutStyle.Render("Ctrl+S"),
			shortcutStyle.Render("Esc"),
			shortcutStyle.Render("Ctrl+P"),
		)
		
		var extraShortcuts string
		if !m.renderMode {
			extraShortcuts = fmt.Sprintf(" | %s Título %s Lista %s Código",
				shortcutStyle.Render("^T"),
				shortcutStyle.Render("^L"),
				shortcutStyle.Render("^K"),
			)
		}
		centerContent = baseShortcuts + extraShortcuts
	}

	// Layout del Footer: Mode | File ... Center ... Position
	leftSide := lipgloss.JoinHorizontal(lipgloss.Center, mode, fileInfo)
	
	// Espaciador flexible
	// Calculamos el espacio restante para centrar el contenido o simplemente rellenar
	availWidth := m.width - lipgloss.Width(leftSide) - lipgloss.Width(position)
	if availWidth < 0 { availWidth = 0 }
	
	// Renderizamos el contenido central en el medio del espacio disponible
	centerArea := statusBarStyle.Width(availWidth).Align(lipgloss.Center).Render(centerContent)

	footer := lipgloss.JoinHorizontal(lipgloss.Center, 
		statusBarStyle.Render(leftSide), 
		centerArea, 
		statusBarStyle.Render(position),
	)

	if m.err != nil {
		footer = statusBarStyle.Width(m.width).
			Background(lipgloss.Color("#FF0000")).
			Render(fmt.Sprintf("ERROR: %v", m.err))
	}

	return lipgloss.JoinVertical(
		lipgloss.Left,
		header,
		content,
		footer,
	)
}
