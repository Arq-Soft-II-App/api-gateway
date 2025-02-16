#!/bin/bash
set -e  # Detiene el script si hay algún error
# update_nginx.sh
# Este script genera un nuevo nginx.conf basándose en las instancias activas de cada servicio y recarga Nginx.

# Definimos los servicios y puertos usando arrays simples
services_names=("api-gateway" "users-api" "courses-api" "inscription-api" "search-api")
services_ports=("8000" "4001" "4002" "4003" "4004")

# Iniciamos la plantilla del archivo.
NGINX_CONF="events {}

http {
"

# Para cada servicio, consultamos los contenedores activos y armamos el bloque upstream.
for i in "${!services_names[@]}"; do
    service="${services_names[$i]}"
    port="${services_ports[$i]}"
    echo "Procesando servicio $service..."
    
    # Filtramos contenedores cuyo nombre contenga el identificador del servicio
    containers=$(docker ps --filter "name=${service}" --format "{{.Names}}")
    # Si no se encuentra ningún contenedor, usamos el nombre base
    if [ -z "$containers" ]; then
       containers="$service"
    fi

    NGINX_CONF+="    upstream ${service}_backend {\n"
    for c in $containers; do
        echo " container - $c"
       NGINX_CONF+="        server ${c}:${port};\n"
    done
    NGINX_CONF+="    }\n\n"
done

# Armamos el bloque server. Aquí definimos las rutas para balancear cada servicio.
NGINX_CONF+="    server {
        listen 80;

        # Por defecto, redirige al api_gateway
        location / {
            proxy_pass http://api-gateway_backend/;
            proxy_set_header Host \$host;
            proxy_set_header X-Real-IP \$remote_addr;
            proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        }

        # Para usuarios
        location /users/ {
            proxy_pass http://users-api_backend/;
        }

        # Para cursos
        location /courses/ {
            proxy_pass http://courses-api_backend/;
        }

        # Para inscripciones
        location /inscriptions/ {
            proxy_pass http://inscription-api_backend/;
        }

        # Para búsqueda (search)
        location /search/ {
            proxy_pass http://search-api_backend/;
        }
    }
}
"

# Escribimos el contenido generado en el archivo nginx.conf.
echo -e "$NGINX_CONF" > ./nginx.conf

# Recargamos Nginx para que tome la nueva configuración.
docker exec nginx nginx -s reload

echo "nginx.conf actualizado y Nginx recargado."