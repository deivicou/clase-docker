// main_1.go
// Ejemplo 1 — Pinta el comando y lo ejecuta en un proceso hijo.
// Uso: go run main_1.go run echo "hola mundo"

package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	default:
		fmt.Println("Comando no reconocido")
	}
}

// run lanza el comando recibido como proceso hijo del proceso actual.
func run() {
	fmt.Printf("Ejecutando (proceso padre PID=%d): %v\n", os.Getpid(), os.Args[2:])

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
