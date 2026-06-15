// main_0.go
// Ejemplo 0 — Recibe un comando como argumento y lo pinta por pantalla.
// Uso: go run main_0.go run echo "hola mundo"

package main

import (
	"fmt"
	"os"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	default:
		fmt.Println("Comando no reconocido")
	}
}

func run() {
	fmt.Printf("Ejecutando: %v\n", os.Args[2:])
}
