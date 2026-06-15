// main_6.go
// Ejemplo 6 — Monta /proc dentro del namespace para que `ps` muestre
//             únicamente los procesos del contenedor.
//
// AVISO: al salir del contenedor el /proc del host queda desmontado.
//        Esto se corrige en main_7.go usando un namespace Mount propio.
//
// Uso: go run main_6.go run /bin/bash
//      ps        ← verás solo los procesos del contenedor
//      exit      ← el host puede quedar con /proc desmontado

package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

const rootfs = "/root/docker/namespace/contenedor"

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
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID,
	}

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func child() {
	fmt.Printf("[HIJO] PID=%d\n", os.Getpid())

	syscall.Sethostname([]byte("contenedor"))
	syscall.Chroot(rootfs)
	os.Chdir("/")

	// Monta el pseudo-sistema de ficheros /proc para que `ps` funcione.
	syscall.Mount("proc", "proc", "proc", 0, "")
	defer syscall.Unmount("proc", 0) // limpieza al salir

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
