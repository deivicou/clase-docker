// main_3.go
// Ejemplo 3 — El proceso padre imprime información; el hijo se ejecuta en
// un namespace UTS aislado y puede cambiar su hostname sin afectar al host.
//
// Uso: go run main_3.go run /bin/bash

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

// run se ejecuta en el proceso padre: relanza el propio binario como "child"
// dentro de un nuevo namespace UTS.
func run() {
	fmt.Printf("[PADRE] PID=%d\n", os.Getpid())

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS,
	}

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

// child se ejecuta dentro del nuevo namespace.
func child() {
	fmt.Printf("[HIJO] PID=%d dentro del namespace UTS\n", os.Getpid())

	// Asigna un hostname propio al namespace
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
