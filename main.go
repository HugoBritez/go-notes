package main

import (
	"fmt"
	"os"

	"go-notes/internal/config"
	"go-notes/internal/storage"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Uso: go-notes init <ruta>", "go-notes <nota>")
		return
	}

	cfg, err := config.Load()

	if os.Args[1] == "init" {
		handleInit()
		return
	}

	if err != nil {
		fmt.Println("Error: Primero debes iniciarlizar la app con ''go-notes init <ruta>")
		return
	}

	noteInput := os.Args[1]

	fullPath, err := storage.CreateNotePath(cfg.NotesRoot, noteInput)
	if err != nil {
		fmt.Printf("Error %v\n", err)
		return
	}

	if _, err := os.Stat(fullPath); os.IsNotExist(err) {
		err = os.WriteFile(fullPath, []byte("# "+noteInput+"\n\n"), 0o644)
		if err != nil {
			fmt.Printf("No se pudo crear el archivo: %v\n", err)
		}
		fmt.Printf("Nueva nota creada en... %s\n", fullPath)
	}
	fmt.Printf("Abriendo nota: %s\n", fullPath)
}

func handleInit() {
	if len(os.Args) < 3 {
		fmt.Println("Uso: go-notes init ~/MisNotas")
		return
	}
	path := os.Args[2]
	cfg := &config.Config{NotesRoot: path, Editor: "nvim"}
	if err := cfg.Save(); err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}
	fmt.Printf("✅ Inicializado en: %s\n", path)
}
