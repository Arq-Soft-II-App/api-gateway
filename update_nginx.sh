#!/bin/bash
set -e

services_names=("api-gateway" "users-api" "courses-api" "inscription-api" "search-api")
services_ports=("8000" "4001" "4002" "4003" "4004")

NGINX_CONF="events {}

http {
"

for i in "${!services_names[@]}"; do
    service="${services_names[$i]}"
    port="${services_ports[$i]}"
    echo "Procesando servicio $service..."
    
    containers=$(docker ps --filter "name=${service}" --format "{{.Names}}")
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
        location /api_users/ {
            proxy_pass http://users-api_backend/;
        }

        # Para cursos
        location /api_courses/ {
             proxy_pass http://courses-api_backend/api_courses/;
        }

        # Para inscripciones
        location /api_inscriptions/ {
            proxy_pass http://inscription-api_backend/;
        }

        # Para búsqueda (search)
        location /api_search/ {
            proxy_pass http://search-api_backend/search/;
        }
    }
}
"

echo -e "$NGINX_CONF" > ./nginx.conf

docker exec nginx nginx -s reload

echo "nginx.conf actualizado y Nginx recargado."