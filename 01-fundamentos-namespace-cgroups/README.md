# 01 · Fundamentos — Namespace, cgroups y chroot

> **¿Cómo funciona un contenedor por dentro?**  
> En este módulo construimos un contenedor *a mano*, sin Docker, usando solo Go y primitivas del kernel Linux.

---

## 🧠 Conceptos clave

| Mecanismo | Función |
|-----------|---------|
| **Namespace** | Aísla los procesos (PID, UTS, Mount…) del resto del sistema |
| **cgroups** | Limita los recursos (CPU, memoria, procesos) de cada proceso |
| **chroot** | Cambia la raíz del sistema de ficheros del proceso |

> Los **contenedores funcionan por capas**: si las reutilizamos ocupamos menos espacio en disco.

---

## 🔬 Progresión de ejemplos en Go

Cada fichero añade una capa de aislamiento sobre el anterior:

| Fichero | Qué hace |
|---------|----------|
| `main_0.go` | Recibe un comando y lo pinta por pantalla |
| `main_1.go` | Ejecuta el comando en un proceso hijo |
| `main_2.go` | Aísla el proceso hijo en un namespace UTS (hostname propio) |
| `main_3.go` | El proceso padre pinta, el hijo ejecuta aislado y cambia el hostname |
| `main_4.go` | Renumera el árbol de procesos del hijo (de forma incompleta) |
| `main_5.go` | Asigna un nuevo root con `chroot` |
| `main_6.go` | Renumera correctamente los procesos aislados |
| `main_7.go` | Igual que el anterior pero sin necesidad de salir para corregirlo |
| `main_final.go` | Limita a 20 procesos con **cgroup** |
| `main_completo.go` | Ejemplo completo: namespace + chroot + cgroups |

---

## 🚀 Cómo ejecutar los ejemplos

### 1. Conectarse a la VM Debian

```bash
ssh -p 2222 user@localhost   # usuario: user  |  contraseña: user
sudo -i                      # escalar a root  (contraseña: user)
```

### 2. Preparar el entorno

```bash
cd script_go
./nuevo_namespace.sh         # crea el esquema de directorios Debian con debootstrap

cd ~/docker/namespace/contenedor
touch CONTENEDOR             # marca el nuevo root del namespace
```

### 3. Ejecutar los ejemplos

```bash
cd script_go

go run main_0.go run echo "hola mundo"
go run main_1.go run echo "hola mundo"
go run main_2.go run /bin/bash
go run main_3.go run /bin/bash
go run main_4.go run /bin/bash
go run main_5.go run /bin/bash
go run main_6.go run /bin/bash
go run main_7.go run /bin/bash
go run main_final.go run /bin/bash
go run main_completo.go run /bin/bash
```

---

## 🧪 Comandos útiles durante las pruebas

```bash
# Ver árbol de procesos del sistema
ps af

# Dentro del contenedor: cambiar el hostname
hostname
hostname david
hostname          # comprobamos que solo cambia dentro del namespace

# Prueba de fork bomb (limitada a 20 procesos por cgroup)
:(){ :|:& };:

# Ver los límites del cgroup creado
cat /sys/fs/cgroup/xavi    # muestra límites de CPU, procesos, etc.
                           # cgroups.procs → PID de los procesos del host padre
```
