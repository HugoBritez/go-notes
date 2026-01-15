# Nombre del binario final
BINARY_NAME=note

# Directorio de instalación (estándar en macOS/Linux)
INSTALL_DIR=/usr/local/bin

.PHONY: all build install clean run

all: build

# 1. Compilar el proyecto
build:
	@echo "🔨 Compilando go-notes como '$(BINARY_NAME)'..."
	@go build -o $(BINARY_NAME) main.go
	@echo "✅ Compilación exitosa."

# 2. Instalar globalmente (requiere sudo a veces)
install: build
	@echo "📦 Instalando '$(BINARY_NAME)' en $(INSTALL_DIR)..."
	@mv $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "🎉 Instalado! Ahora puedes usar el comando '$(BINARY_NAME)' desde cualquier lugar."

# 3. Limpiar archivos generados
clean:
	@echo "🧹 Limpiando..."
	@rm -f $(BINARY_NAME)
	@rm -f go-notes # por si quedó el viejo
	@echo "✨ Limpio."

# 4. Ejecutar localmente (para dev)
run:
	@go run main.go
