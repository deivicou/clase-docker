// main_5.go
// Ejemplo 5 — Añade chroot para cambiar el sistema de ficheros raíz del hijo.
//
// El proceso hijo verá ~/docker/namespace/contenedor como su "/".
// No podrá acceder al sistema de ficheros del host padre.
//
// Uso: go run main_5.go run /bin/bash
//      ls        ← verás el contenido de ~/docker/namespace/contenedor

package main

import (
	"fmt"
	"os"
	"os/exec"
	"syscall"
)

// rootfs es el directorio que se convertirá en la nueva raíz del contenedor.
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
	fmt.Printf("[HIJO] PID=%d — cambiando root a %s\n", os.Getpid(), rootfs)

	if err := syscall.Sethostname([]byte("contenedor")); err != nil {
		fmt.Fprintf(os.Stderr, "Sethostname: %v\n", err)
	}

	// Cambia el directorio raíz del proceso al rootfs definido.
	if err := syscall.Chroot(rootfs); err != nil {
		fmt.Fprintf(os.Stderr, "Chroot: %v\n", err)
		os.Exit(1)
	}

	// Después de chroot nos movemos a la nueva raíz.
	if err := os.Chdir("/"); err != nil {
		fmt.Fprintf(os.Stderr, "Chdir: %v\n", err)
		os.Exit(1)
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
