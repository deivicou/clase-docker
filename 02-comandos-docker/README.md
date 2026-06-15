# 02 · Comandos Docker

Referencia rápida de los comandos más usados, organizados por categoría.

---

## 📦 Instalación

```bash
sudo apt-get update
sudo curl -sSL https://get.docker.com | sh        # descarga e instala Docker

# Añade el usuario actual al grupo docker para no necesitar sudo
sudo usermod -aG docker ${USER}
```

---

## ▶️ Ejecución de contenedores

```bash
# Ejecuta una imagen (la descarga si no existe); cada llamada crea una instancia nueva
docker run <imagen|id>

# Ejecuta un comando puntual y la instancia muere (pero NO se borra)
docker run <imagen|id> ls

# Opciones más habituales
docker run \
  --name mi-contenedor \        # nombre de la instancia
  -p 8080:80 \                  # mapeo puerto_local:puerto_contenedor
  -h mi-hostname \              # hostname del contenedor
  -u 0 \                        # ejecutar como root
  -d \                          # modo detached (segundo plano)
  --network mi-red \            # conectar a una red concreta
  <imagen|id>

# Ejecutar en modo interactivo (shell dentro del contenedor)
docker run -it <imagen|id> bash
# → Ctrl+P  Ctrl+Q  para salir sin matar el contenedor (detach)

# Ejecutar y borrar la instancia al terminar
docker run --rm <imagen|id>

# Ejecutar un comando complejo
docker run -d <imagen|id> /bin/bash -c "sleep infinity"
```

---

## 🔌 Entrar en un contenedor en ejecución

```bash
# attach: entra como proceso principal (bloquea si hay algo corriendo)
docker attach <nombre|id>

# exec: abre una shell adicional sin interrumpir el proceso principal
docker exec -it <nombre|id> bash
# → salir: Ctrl+C  o  exit
```

---

## 🗂️ Gestión de contenedores e imágenes

```bash
# --- Contenedores ---
docker ps                          # contenedores activos
docker ps -a                       # todos (incluidos parados)
docker ps | grep <nombre|id>       # filtrar activos
docker ps -a | grep <nombre|id>    # filtrar todos

docker start <nombre|id>           # iniciar
docker stop  <nombre|id>           # parar (graceful)
docker rm    <nombre|id>           # borrar instancia
docker rm $(docker ps -aq)         # borrar TODAS las instancias

# --- Imágenes ---
docker images                      # imágenes descargadas
docker rmi <nombre|id>             # borrar imagen (las instancias deben borrarse antes)
docker image prune                 # borrar imágenes sin usar

# --- Información ---
docker inspect <nombre|id>         # detalles en JSON
docker history <nombre|id>         # historial de capas
docker logs -f <nombre|id>         # logs en tiempo real
docker system df -v                # uso de recursos del host

# --- Copiar ficheros ---
docker cp <ruta_local>  <contenedor>:<ruta_destino>   # host → contenedor
docker cp <contenedor>:<ruta>  <ruta_local>            # contenedor → host (funciona en parado)
```

---

## 🐳 Docker Hub

```bash
docker login -u <usuario>                            # autenticarse

docker search <imagen>                               # buscar en Docker Hub
docker pull  <imagen>                                # descargar imagen

docker push <usuario>/<imagen>:<tag>                 # subir imagen al registro

# Importar / exportar
docker save <imagen>:<tag> | gzip > imagen.tar.gz   # guardar imagen en fichero
docker load < imagen.tar.gz                          # cargar imagen desde fichero

docker export <contenedor> > contenedor.tar          # exportar contenedor (sin metadatos)
docker import contenedor.tar                         # importar contenedor

docker commit <contenedor>                           # guardar estado actual como imagen
docker tag <imagen|id> ejemplo:v1.0                  # añadir etiqueta a una imagen
```

---

## 💾 Volúmenes

Los volúmenes permiten persistir datos independientemente del ciclo de vida del contenedor.

```bash
# Crear un volumen con nombre
docker volume create mi-volumen

# Montar un volumen al arrancar
docker run -v mi-volumen:/ruta/contenedor <imagen>         # volumen nombrado
docker run -v /ruta/local:/ruta/contenedor <imagen>        # bind mount (directorio del host)
docker run -v /ruta/contenedor <imagen>                    # volumen anónimo (sin host explícito)
docker run -v mi-volumen:/datos:ro <imagen>                # solo lectura (read-only)

# Gestión
docker volume ls                          # listar volúmenes
docker volume rm <nombre|id>              # borrar volumen
docker volume prune                       # borrar volúmenes sin usar
docker inspect <nombre|id>               # ver ruta real en el host (JSON)
```

---

## 🌐 Redes

```bash
# Tipos de red en Docker:
#   bridge   → red por defecto; todos los contenedores conectados
#   host     → usa directamente la red del host
#   none     → contenedor totalmente aislado
#   overlay  → para Docker en diferentes hosts (Swarm)
#   macvlan  → cada contenedor obtiene su propia MAC

docker network create <nombre>                              # red bridge
docker network create --subnet 192.168.100.0/24 <nombre>   # con IP fija

docker network connect    <red> <contenedor>               # conectar contenedor a red
docker network disconnect <red> <contenedor>               # desconectar

docker network ls                                          # listar redes
docker network rm <nombre>                                 # eliminar red
docker network inspect <nombre>                            # detalles (JSON)

# Listar contenedores con su IP
docker inspect $(docker ps -q) \
  --format='{{ printf "%-50s" .Name}} {{range .NetworkSettings.Networks}}{{.IPAddress}} {{end}}'
```
