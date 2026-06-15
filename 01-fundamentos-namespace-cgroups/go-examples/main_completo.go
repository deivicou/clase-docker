// main_completo.go
// Ejemplo completo — Contenedor mínimo autocontenido.
//
// Combina todos los mecanismos vistos:
//   • Namespace UTS  → hostname propio
//   • Namespace PID  → árbol de procesos aislado
//   • Namespace Mount → sistema de ficheros propio
//   • chroot          → cambia la raíz al rootfs del contenedor
//   • cgroups         → limita a 20 procesos
//   • Copia automática de los binarios necesarios
//
// Uso: go run main_completo.go run /bin/bash

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
	cgroupDir = "/sys/fs/cgroup/xavi"
)

func main() {
	switch os.Args[1] {
	case "run":
		run()
	case "child":
		child()
	default:
		panic("Se esperaba 'run' o 'child'")
	}
}

// run se ejecuta en el proceso padre y lanza el hijo en namespaces aislados.
func run() {
	fmt.Printf("[PADRE] PID=%d iniciando contenedor…\n", os.Getpid())

	cmd := exec.Command("/proc/self/exe", append([]string{"child"}, os.Args[2:]...)...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS |
			syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS,
	}

	must(cmd.Run())
}

// child configura el entorno aislado y ejecuta el comando solicitado.
func child() {
	fmt.Printf("[HIJO] PID=%d configurando entorno…\n", os.Getpid())

	configurarCgroup()

	must(syscall.Sethostname([]byte("contenedor")))
	must(syscall.Chroot(rootfs))
	must(os.Chdir("/"))
	must(syscall.Mount("proc", "proc", "proc", 0, ""))

	cmd := exec.Command(os.Args[2], os.Args[3:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	must(cmd.Run())

	must(syscall.Unmount("proc", 0))
}

// configurarCgroup aplica límites de recursos al proceso actual.
func configurarCgroup() {
	os.MkdirAll(cgroupDir, 0755)
	escribirCgroup("pids.max", "20")
	escribirCgroup("cgroup.procs", strconv.Itoa(os.Getpid()))
}

func escribirCgroup(fichero, valor string) {
	ruta := filepath.Join(cgroupDir, fichero)
	must(os.WriteFile(ruta, []byte(valor), 0700))
}

// must termina el programa si err != nil.
func must(err error) {
	if err != nil {
		panic(err)
	}
}
