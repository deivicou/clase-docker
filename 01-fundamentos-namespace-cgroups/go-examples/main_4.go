// main_4.go
// Ejemplo 4 — Añade namespace PID al hijo para renumerar su árbol de procesos.
//
// NOTA: la renumeración es incompleta en este ejemplo porque /proc aún no
// está montado dentro del namespace. Se corregirá en main_5.go y sucesores.
//
// Uso: go run main_4.go run /bin/bash
//      ps        ← verás PIDs del host padre (comportamiento esperado aquí)

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
	case "child":
		child()
	default:
		fmt.Println("Comando no reconocido")
	}
}

func run() {
	fmt.Printf("[PADRE] PID=%d\n", os.Getpid())

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		// CLONE_NEWUTS → hostname aislado
		// CLONE_NEWPID → árbol de procesos aislado (renumerado desde PID 1)
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	}

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func child() {
	fmt.Printf("[HIJO] PID=%d (debería verse como 1 dentro del namespace)\n", os.Getpid())

	if err := syscall.Sethostname([]byte("contenedor")); err != nil {
		fmt.Fprintf(os.Stderr, "Sethostname: %v\n", err)
	}

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
