# 04 · Docker Compose

Docker Compose permite definir y arrancar **escenarios multi-contenedor** con un único fichero YAML.

---

## 🚀 Comandos principales

```bash
# Arrancar el escenario (busca docker-compose.yml en el directorio actual)
docker compose up -d

# Especificar un fichero YAML concreto
docker compose -f fichero.yml up -d

# Forzar reconstrucción de imágenes antes de arrancar
docker compose up --build

# Construir imágenes sin arrancar
docker compose build

# Listar contenedores del escenario
docker compose ps

# Parar contenedores (sin borrar)
docker compose stop

# Borrar contenedores parados
docker compose rm

# Eliminar contenedores, redes y volúmenes del escenario
docker compose down -v
```

---

## 📋 Estructura del fichero `docker-compose.yml`

```yaml
services:
  # ── Nombre del servicio (contenedor) ─────────────────────────────────────
  cliente:
    image: nombre-imagen:tag          # imagen base
    container_name: mi-cliente        # nombre visible en docker ps
    restart: always                   # política de reinicio:
                                      #   always        → siempre (incluso tras docker stop)
                                      #   no            → nunca (por defecto)
                                      #   on-failure    → solo en caso de error
                                      #   unless-stopped → siempre salvo parada manual
    depends_on:
      - base-de-datos                 # no arranca hasta que base-de-datos esté listo
    ports:
      - "8080:80"                     # puerto_local:puerto_contenedor
    volumes:
      - ./datos:/app/datos            # bind mount local
      - datos-app:/app/cache          # volumen nombrado
    networks:
      - red-interna
    environment:
      APP_ENV: production
      APP_PORT: 80

# ── Declaración de redes ──────────────────────────────────────────────────
networks:
  red-interna:

# ── Declaración de volúmenes ──────────────────────────────────────────────
volumes:
  datos-app:
```

---

## 📌 Políticas de reinicio (`restart`)

| Valor | Comportamiento |
|-------|---------------|
| `no` | No reinicia nunca (por defecto) |
| `always` | Reinicia siempre, incluso después de `docker stop` |
| `on-failure` | Solo reinicia si el contenedor termina con error |
| `unless-stopped` | Siempre salvo que se pare manualmente |

---

## 🌐 Portainer — interfaz gráfica para Docker

Portainer es una UI web para gestionar contenedores.  
Se ejecuta él mismo como un contenedor:

```bash
docker run -d \
  --name portainer \
  --restart=always \
  -p 9000:9000 \
  -v /var/run/docker.sock:/var/run/docker.sock \
  -v portainer_data:/data \
  portainer/portainer
```

Acceder en: **http://localhost:9000**

---

## 🗂️ Ejemplos incluidos

| Fichero | Descripción |
|---------|-------------|
| [`ejemplos/wordpress.yml`](./ejemplos/wordpress.yml) | WordPress + MySQL con volúmenes y red interna |
| [`ejemplos/portainer.yml`](./ejemplos/portainer.yml) | Portainer listo para usar |
