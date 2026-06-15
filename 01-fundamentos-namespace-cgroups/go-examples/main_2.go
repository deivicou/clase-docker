// main_2.go
// Ejemplo 2 — Ejecuta el comando hijo en un namespace UTS aislado.
//
// El namespace UTS permite asignar un hostname propio al proceso hijo
// sin afectar al hostname del sistema anfitrión.
//
// Uso: go run main_2.go run /bin/bash
//      hostname david      ← solo cambia dentro del namespace
//      exit

package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
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
	fmt.Printf("Proceso padre PID=%d ejecutando: %v\n", os.Getpid(), os.Args[2:])

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	// Cloneflags crea el proceso hijo dentro de un nuevo namespace UTS.
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS,
	}

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
