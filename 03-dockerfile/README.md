# 03 · Dockerfile

Un **Dockerfile** es un fichero de texto con las instrucciones para construir una imagen Docker paso a paso.  
Cada instrucción genera una **capa** en la imagen final.

---

## 🔨 Construir una imagen

```bash
# Construir desde el directorio actual (busca ./Dockerfile)
docker build . -t <usuario_dockerhub>/<nombre_imagen>:<tag>

# Ejemplos
docker build . -t miusuario/mi-app:latest
docker build . -t miusuario/mi-app:v1.0
```

---

## 📋 Instrucciones del Dockerfile

```dockerfile
# ── Imagen base ──────────────────────────────────────────────────────────────
FROM php:7.4-apache
# Toda imagen parte de una imagen base. Usar versiones específicas (no latest)
# para garantizar reproducibilidad.

# ── Metadatos ────────────────────────────────────────────────────────────────
LABEL company="iesalixar"
# MAINTAINER está obsoleto; usar LABEL maintainer="email@ejemplo.com"

# ── Variables de entorno ──────────────────────────────────────────────────────
ENV APP_HOME=/app
# Disponibles en tiempo de build Y en tiempo de ejecución.

# ── Directorio de trabajo ─────────────────────────────────────────────────────
WORKDIR /app
# Equivale a `cd /app`. Los comandos RUN, COPY, etc. siguientes parten de aquí.
# Crea el directorio si no existe.

# ── Ejecutar comandos en tiempo de build ──────────────────────────────────────
RUN apt-get update && apt-get install -y curl
# Combinar comandos en un solo RUN reduce el número de capas.

# ── Copiar ficheros locales ───────────────────────────────────────────────────
COPY index.html /var/www/html/
# Copia ficheros del contexto de build (máquina local) al sistema de ficheros
# de la imagen.

# ── Agregar ficheros (con soporte tar y URLs) ─────────────────────────────────
ADD myapp.tar.gz /app/
# ADD descomprime .tar.gz automáticamente.
# Para ficheros locales, preferir COPY (más explícito).

# ── Declarar puerto de escucha ────────────────────────────────────────────────
EXPOSE 80
# Informativo: indica el puerto que usará el contenedor.
# No publica el puerto; eso se hace con -p al hacer docker run.

# ── Volumen ───────────────────────────────────────────────────────────────────
VOLUME /data
# Define un punto de montaje. Docker creará un volumen anónimo si no se
# especifica uno en docker run.

# ── Comando por defecto (sustituible) ────────────────────────────────────────
CMD ["nginx", "-g", "daemon off;"]
# Se ejecuta si no se pasa ningún comando en docker run.
# Puede sobreescribirse: docker run <imagen> /bin/bash

# ── Punto de entrada fijo ────────────────────────────────────────────────────
ENTRYPOINT ["python", "app.py"]
# Siempre se ejecuta. Los argumentos de docker run se añaden como parámetros.
# Combinar ENTRYPOINT + CMD: ENTRYPOINT es el ejecutable, CMD son los args por defecto.
```

---

## ⚠️ Trampa frecuente: los comandos son por capa

Cada línea `RUN` se ejecuta en un shell **nuevo**. El `cd` de una línea no afecta a la siguiente:

```dockerfile
# ❌ MAL — install.sh no se encuentra porque cd solo existe en esa capa
RUN cd /scripts/
RUN ./install.sh

# ✅ BIEN — todo en la misma capa
RUN cd /scripts/ && ./install.sh
```

---

## 📄 Ejemplo completo — App PHP con Apache

```dockerfile
FROM php:7.4-apache

LABEL maintainer="admin@ejemplo.com" \
      version="1.0"

ENV APACHE_DOCUMENT_ROOT=/var/www/html/app

WORKDIR /var/www/html

RUN apt-get update && apt-get install -y \
    curl \
    unzip \
    && rm -rf /var/lib/apt/lists/*

COPY . /var/www/html/

EXPOSE 80

CMD ["apache2-foreground"]
```

```bash
# Construir y ejecutar
docker build . -t miusuario/php-app:1.0
docker run -d -p 8080:80 miusuario/php-app:1.0
```

---

## 💡 Buenas prácticas

| Práctica | Por qué |
|----------|---------|
| Combinar `RUN` con `&&` | Menos capas → imagen más ligera |
| Limpiar caché del gestor de paquetes | `rm -rf /var/lib/apt/lists/*` ahorra espacio |
| Usar versiones concretas en `FROM` | Garantiza reproducibilidad |
| Poner `COPY`/`ADD` al final | Aprovecha la caché de capas en rebuilds |
| Un proceso por contenedor | Facilita el escalado y los logs |
