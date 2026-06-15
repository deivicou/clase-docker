// main_7.go
// Ejemplo 7 — Añade namespace Mount (CLONE_NEWNS) para que el montaje de
//             /proc sea privado al contenedor y no afecte al host padre.
//
// Uso: go run main_7.go run /bin/bash
//      ps        ← solo procesos del contenedor
//      exit      ← el host sigue funcionando con normalidad

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
		// CLONE_NEWNS → namespace de montajes propio (los montajes no se propagan al host)
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
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

	// Ahora el montaje de /proc es privado gracias a CLONE_NEWNS.
	syscall.Mount("proc", "proc", "proc", 0, "")

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	syscall.Unmount("proc", 0)
}
