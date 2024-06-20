# Usar una imagen base oficial de Go para la compilación
FROM golang:1.19 as builder

# Establecer el directorio de trabajo
WORKDIR /app

# Copiar el archivo go.mod y go.sum (si existe) para aprovechar la caché de capas de Docker
COPY go.mod .
COPY go.sum .

# Descargar las dependencias de Go (esto aprovecha la caché si los archivos go.mod y go.sum no cambian)
RUN go mod download

# Copiar el resto del código fuente de la aplicación
COPY . .

# Compilar la aplicación (ajusta el nombre del binario como sea necesario)
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o exporter .

# Usar una imagen base ligera de Alpine para el contenedor final
FROM alpine:latest  

# Instalar ca-certificates para llamadas HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar el ejecutable desde la etapa de construcción
COPY --from=builder /app/exporter .

# Exponer el puerto que tu aplicación utiliza (ajusta si es diferente)
EXPOSE 8080

# Definir ENTRYPOINT para pasar parámetros al ejecutable
ENTRYPOINT ["./exporter"]

# Definir CMD con los argumentos por defecto que pueden ser sobrescritos en tiempo de ejecución
CMD ["--tenant-id", "", "--client-id", "", "--client-secret", "", "--subscription-id", "", "--resource-group", "", "--vnet-name", "", "--interval", "300"]