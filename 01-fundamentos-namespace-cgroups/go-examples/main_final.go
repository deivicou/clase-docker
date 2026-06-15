// main_final.go
// Ejemplo final — Añade cgroups para limitar el número máximo de procesos
//                del contenedor a 20.
//
// Uso: go run main_final.go run /bin/bash
//      :(){ :|:& };:    ← fork bomb limitada a 20 procesos
//      cat /sys/fs/cgroup/xavi   ← ver límites aplicados

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"syscall"
)

const (
	rootfs    = "/root/docker/namespace/contenedor"
	cgroupDir = "/sys/fs/cgroup/xavi" // directorio del cgroup propio
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
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func child() {
	fmt.Printf("[HIJO] PID=%d — aplicando cgroup con límite de 20 procesos\n", os.Getpid())

	configurarCgroup()

	syscall.Sethostname([]byte("contenedor"))
	syscall.Chroot(rootfs)
	os.Chdir("/")
	syscall.Mount("proc", "proc", "proc", 0, "")

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	syscall.Unmount("proc", 0)
}

// configurarCgroup crea un cgroup y aplica las restricciones de recursos.
func configurarCgroup() {
	// Crear el directorio del cgroup (cgroups v2)
	os.MkdirAll(cgroupDir, 0755)

	// Límite máximo de tareas (procesos + hilos) = 20
	escribirCgroup("pids.max", "20")

	// Añadir el proceso actual al cgroup
	escribirCgroup("cgroup.procs", strconv.Itoa(os.Getpid()))
}

// escribirCgroup escribe un valor en un fichero de control del cgroup.
func escribirCgroup(fichero, valor string) {
	ruta := filepath.Join(cgroupDir, fichero)
	if err := os.WriteFile(ruta, []byte(valor), 0700); err != nil {
		fmt.Fprintf(os.Stderr, "cgroup %s: %v\n", fichero, err)
	}
}
