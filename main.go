package main

import (
	"fmt"
	"os"

	"go-notes/internal/config"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Printf("Uso: go-notes init <ruta>", "go-notes <nota>")
		return
	}

	command := os.Args[1]

	if command == "init" {
		if len(os.Args) < 3 {
			fmt.Println("Error: Debes especificar una ruta. Ej: go-notes init ~/Notas")
			return
		}

		path := os.Args[2]
		cfg := &config.Config{
			NotesRoot: path,
			Editor:    "nvim", // por defecto
		}

		err := cfg.Save()
		if err != nil {
			fmt.Printf("Error al guardar la configuracion %v\n", err)
			return
		}

		fmt.Printf("Proceso inicializado. Notas se guardaran en: %s \n", path)
	}
}
