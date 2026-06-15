# 05 · Docker Swarm

Docker Swarm convierte un conjunto de máquinas Docker en un **clúster** con alta disponibilidad y balanceo de carga automático.

---

## 🧩 Conceptos clave

| Concepto | Descripción |
|----------|-------------|
| **Manager** | Nodo que gestiona el clúster y delega tareas a los Workers |
| **Worker** | Nodo que ejecuta las tareas asignadas por los Managers |
| **Stack** | Conjunto de servicios y sus dependencias (definido en un Compose) |
| **Servicio** | Componente del stack: imagen, réplicas, puertos y política de despliegue |
| **Task** | Unidad atómica: un contenedor por tarea |

> 💡 Se recomienda tener entre **3 y 5 nodos Manager** para garantizar tolerancia a fallos.

---

## 🔌 Puertos necesarios

| Puerto | Protocolo | Uso |
|--------|-----------|-----|
| **2377** | TCP | Gestión del clúster y comunicación con el Manager |
| **7946** | TCP + UDP | Comunicación entre nodos del enjambre |
| **4789** | UDP | Red de superposición (overlay) |

---

## 🚀 Inicializar el clúster

```bash
# ── En el nodo MANAGER ────────────────────────────────────────────────────
# Inicializar Swarm (genera el token para unir Workers)
docker swarm init --advertise-addr <IP_DEL_MANAGER>

# Obtener el token para unir nuevos Managers o Workers
docker swarm join-token manager
docker swarm join-token worker

# ── En los nodos WORKER ───────────────────────────────────────────────────
docker swarm join --token <TOKEN> <IP_NODO_MANAGER>:2377
```

---

## 🗂️ Gestión de nodos

```bash
docker node ls                        # estado del clúster y sus nodos

docker node promote <nombre-worker>   # ascender Worker → Manager
docker node demote  <nombre-manager>  # degradar Manager → Worker
```

---

## 🌐 Red overlay

Necesaria para que los contenedores de distintos nodos se comuniquen entre sí:

```bash
docker network create -d overlay <nombre-red>
```

---

## 📦 Servicios

```bash
# Crear un servicio con réplicas
docker service create \
  --name mi-servicio \
  --publish published=8080,target=80 \
  --replicas 3 \
  <imagen>

# Escalar el número de réplicas
docker service scale mi-servicio=5

# Ver estado de cada réplica (en qué nodo está, si está corriendo…)
docker service ps mi-servicio
```

---

## 📋 Stacks (despliegue con docker-compose)

Un **Stack** en Swarm es equivalente a un `docker compose up`, pero distribuido entre nodos:

```bash
# Desplegar un stack desde un fichero Compose
docker stack deploy -c docker-compose.yml <nombre-stack>

# Listar stacks desplegados
docker stack ls

# Ver los servicios de un stack
docker stack services <nombre-stack>
```

> ⚠️ En Swarm, el fichero Compose debe usar la clave `deploy:` para definir réplicas y políticas.

---

## 📄 Ejemplo de Compose con Swarm

```yaml
services:
  web:
    image: nginx:alpine
    ports:
      - "80:80"
    networks:
      - red-publica
    deploy:
      replicas: 3                    # 3 instancias distribuidas entre los nodos
      update_config:
        parallelism: 1               # actualizar de 1 en 1
        delay: 10s
      restart_policy:
        condition: on-failure

networks:
  red-publica:
    driver: overlay
```

```bash
# Desplegar
docker stack deploy -c swarm-ejemplo.yml mi-app

# Ver réplicas y su ubicación
docker service ps mi-app_web
```

---

## ☁️ Recursos adicionales

- 🌐 [Play with Docker — Introducción a Swarm](https://training.play-with-docker.com/ops-s1-swarm-intro/)
- 📺 [Docker Swarm desde cero](https://www.youtube.com/watch?v=9cpE-vNDK7A)

---

## 💡 Notas sobre AWS (Launch Templates)

Al desplegar Swarm en AWS:

1. Definir un **Security Group** con los puertos 2377, 7946 y 4789 abiertos entre nodos, y 22 para SSH.
2. Crear un **Launch Template** para lanzar varias instancias EC2 del mismo tipo.
3. Conectarse a cada instancia:
   ```bash
   ssh -i certificado.pem usuario@servidor
   ```
4. Instalar Docker en cada nodo e inicializar el Swarm desde el Manager.
