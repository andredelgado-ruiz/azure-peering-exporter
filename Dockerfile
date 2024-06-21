# Usar una imagen base oficial de Go para la compilación, especificando la plataforma amd64
FROM --platform=linux/amd64 golang:1.22.4 as builder

# Establecer el directorio de trabajo
WORKDIR /app

# Copiar archivos de Go y descargar dependencias (reduce las capas de la imagen)
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código fuente de la aplicación
COPY . .

# Compilar la aplicación asegurando que el binario sea compatible con amd64
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -installsuffix cgo -o azure_peering_exporter .

# Usar una imagen base ligera de Alpine para el contenedor final, especificando la plataforma amd64
FROM --platform=linux/amd64 alpine:latest  

# Instalar ca-certificates para llamadas HTTPS
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Copiar el ejecutable desde la etapa de construcción
COPY --from=builder /app/azure_peering_exporter .

# Exponer el puerto que tu aplicación utiliza
EXPOSE 8080

# Definir ENTRYPOINT para pasar parámetros al ejecutable
ENTRYPOINT ["./azure_peering_exporter"]

# Definir CMD con los argumentos por defecto que pueden ser sobrescritos en tiempo de ejecución
CMD ["--tenant-id", "", "--client-id", "", "--client-secret", "", "--subscription-id", "", "--resource-group", "", "--vnet-name", "", "--interval", "300"]