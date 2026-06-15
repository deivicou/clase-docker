# 🐳 Cheatsheet Docker — Referencia rápida

## Contenedores

| Acción | Comando |
|--------|---------|
| Ejecutar (interactivo) | `docker run -it <imagen> bash` |
| Ejecutar (segundo plano) | `docker run -d <imagen>` |
| Ejecutar y borrar al salir | `docker run --rm <imagen>` |
| Listar activos | `docker ps` |
| Listar todos | `docker ps -a` |
| Iniciar | `docker start <nombre\|id>` |
| Parar | `docker stop <nombre\|id>` |
| Borrar instancia | `docker rm <nombre\|id>` |
| Borrar todas las instancias | `docker rm $(docker ps -aq)` |
| Entrar (proceso nuevo) | `docker exec -it <nombre\|id> bash` |
| Entrar (proceso principal) | `docker attach <nombre\|id>` |
| Ver logs | `docker logs -f <nombre\|id>` |
| Copiar al contenedor | `docker cp <local> <cont>:<ruta>` |
| Copiar del contenedor | `docker cp <cont>:<ruta> <local>` |
| Inspeccionar (JSON) | `docker inspect <nombre\|id>` |

## Imágenes

| Acción | Comando |
|--------|---------|
| Listar | `docker images` |
| Descargar | `docker pull <imagen>` |
| Construir | `docker build . -t <usuario>/<imagen>:<tag>` |
| Borrar | `docker rmi <nombre\|id>` |
| Limpiar sin usar | `docker image prune` |
| Subir a Hub | `docker push <usuario>/<imagen>:<tag>` |
| Historial de capas | `docker history <nombre\|id>` |
| Guardar en tar | `docker save <imagen> \| gzip > img.tar.gz` |
| Cargar desde tar | `docker load < img.tar.gz` |
| Añadir etiqueta | `docker tag <imagen\|id> <nueva-etiqueta>` |

## Volúmenes

| Acción | Comando |
|--------|---------|
| Crear | `docker volume create <nombre>` |
| Listar | `docker volume ls` |
| Borrar | `docker volume rm <nombre>` |
| Limpiar sin usar | `docker volume prune` |
| Inspeccionar | `docker inspect <nombre>` |

## Redes

| Acción | Comando |
|--------|---------|
| Crear | `docker network create <nombre>` |
| Listar | `docker network ls` |
| Borrar | `docker network rm <nombre>` |
| Conectar contenedor | `docker network connect <red> <cont>` |
| Desconectar | `docker network disconnect <red> <cont>` |
| Inspeccionar | `docker network inspect <nombre>` |

## Docker Compose

| Acción | Comando |
|--------|---------|
| Arrancar escenario | `docker compose up -d` |
| Arrancar (reconstruir) | `docker compose up --build` |
| Listar contenedores | `docker compose ps` |
| Parar | `docker compose stop` |
| Borrar contenedores | `docker compose rm` |
| Eliminar todo + volúmenes | `docker compose down -v` |

## Docker Swarm

| Acción | Comando |
|--------|---------|
| Inicializar Manager | `docker swarm init --advertise-addr <IP>` |
| Unirse como Worker | `docker swarm join --token <TOKEN> <IP>:2377` |
| Ver token Manager/Worker | `docker swarm join-token {manager\|worker}` |
| Ver nodos | `docker node ls` |
| Desplegar stack | `docker stack deploy -c compose.yml <stack>` |
| Listar stacks | `docker stack ls` |
| Ver servicios del stack | `docker stack services <stack>` |
| Ver réplicas del servicio | `docker service ps <servicio>` |
| Escalar réplicas | `docker service scale <servicio>=N` |
