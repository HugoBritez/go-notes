# go-notes 🚀

**go-notes** es un editor de notas para la terminal diseñado para desarrolladores que buscan la velocidad de una CLI con la elegancia visual de herramientas modernas como Obsidian, pero viviendo 100% en tu terminal.

Escrito totalmente en **Go**, utiliza el ecosistema **Charm** (Bubble Tea, Lip Gloss, Glamour) para ofrecer una experiencia TUI (Text User Interface) "Premium".

## ✨ Características Principales

- **Dual Mode UI:**
  - ✏️ **Editor:** Interfaz minimalista para escritura rápida sin distracciones.
  - 👁️ **Reader:** Renderizado Markdown profesional en tiempo real (Tablas, Código Coloreado, Listas, etc.) usando `glamour`.
- **Smart Snippets:** Inserción inteligente de bloques de código (`Ctrl+O`) con autocompletado de lenguaje y posicionamiento de cursor.
- **UI Reactiva:**
  - Bordes dinámicos que cambian de color según el modo (Esmeralda/Azul/Naranja).
  - Footer estilo "LazyVim" con conteo de palabras, caracteres y posición de línea.
  - Header con breadcrumbs estilizados (`carpeta › archivo`).
- **Path Discovery**: Crea estructuras de carpetas dinámicamente al vuelo (`note facu/algebra/clase1`).
- **Zero Lock-in**: Tus notas son archivos `.md` planos estándar.

## 🚀 Instalación y Uso

1. **Clona el repositorio:**
   ```bash
   git clone https://github.com/tu-usuario/go-notes.git
   cd go-notes
   ```

2. **Ejecuta directamente:**
   ```bash
   go run main.go [ruta/de/tu/nota]
   ```
   *Ejemplo:* `go run main.go ideas/app_revolutionaria`

## ⌨️ Atajos de Teclado

| Atajo | Acción | Descripción |
| :--- | :--- | :--- |
| **Ctrl + S** | `Guardar` | Guarda el archivo y muestra confirmación visual. |
| **Ctrl + P** | `Vista` | Alterna entre **Modo Editor** y **Modo Lectura** (Preview). |
| **Ctrl + O** | `Smart Code` | Abre el menú para insertar bloques de código (`go`, `js`, `py`...). |
| **Esc** | `Salir/Volver` | Sale del input/modo lectura o cierra la app. |
| **Ctrl + C** | `Forzar Salida` | Cierra la aplicación inmediatamente. |

## 🛠️ Tech Stack

- **Core:** [Go](https://go.dev/)
- **TUI Framework:** [Bubble Tea](https://github.com/charmbracelet/bubbletea)
- **Estilos:** [Lip Gloss](https://github.com/charmbracelet/lipgloss)
- **Renderizado Markdown:** [Glamour](https://github.com/charmbracelet/glamour)
- **Componentes:** [Bubbles](https://github.com/charmbracelet/bubbles)

---
*Hecho con ❤️ y mucho café.*