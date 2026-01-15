# go-notes 🚀

**go-notes** es un editor de notas para la terminal diseñado para desarrolladores que buscan la velocidad de una CLI con la elegancia visual de herramientas como Notion u Obsidian.

Escrito totalmente en **Go**, utiliza el framework de TUI **Bubble Tea** para ofrecer una experiencia interactiva, fluida y moderna.

## ✨ Características (Roadmap)

- [x] **Path Discovery**: Crea notas y carpetas dinámicamente (`note facu/algebra/clase1`).
- [x] **Persistencia Local**: Tus notas son archivos `.md` estándar, tú eres el dueño de tus datos.
- [x] **Zero Config**: Inicialización rápida con `go-notes init`.
- [p] **Notion-look (En progreso)**: Renderizado de Markdown en tiempo real con estilos y colores.
- [ ] **Buscador Integrado**: Integración nativa con FZF para encontrar notas al instante.
- [ ] **Exportación**: Convertir notas a PDF o HTML desde la CLI.

## 🚀 Instalación rápida

1. Clona el repositorio:
   ```bash
   git clone [https://github.com/tu-usuario/go-notes.git](https://github.com/tu-usuario/go-notes.git)
   cd go-notes
Instala las dependencias y compila:

Bash

go mod tidy
go build -o go-notes main.go
🛠️ Cómo usarlo
1. Inicializar
Define dónde quieres que vivan tus notas:

Bash

./go-notes init ~/Documents/notas
2. Crear o Editar una nota
Bash

./go-notes facu/matematica/clase1
Si las carpetas no existen, go-notes las creará por ti.

3. Atajos dentro del editor
Ctrl + S: Guardar y salir.

Esc / Ctrl + C: Salir sin guardar.

🛠️ Tech Stack
Lenguaje: Go

TUI Framework: Bubble Tea

Estilos: Lip Gloss

CLI Helpers: Bubbles

🤝 Contribuir
¡Este es un proyecto Open Source! Si tienes ideas para el renderizado tipo Notion, integración con bases de datos o simplemente quieres mejorar la UI, las Pull Requests son bienvenidas.
