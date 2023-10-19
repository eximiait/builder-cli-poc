# Etapa de compilación
FROM golang:1.21.3 as builder

# Establece argumentos con valores predeterminados para OS y EXT
ARG TARGETOS=linux
ARG EXT=""

# Establece el directorio de trabajo
WORKDIR /app

# Copia el go.mod y go.sum para instalar dependencias
COPY go.mod go.sum ./

# Instala las dependencias del proyecto
RUN go mod download

# Copia el código fuente del proyecto al contenedor
COPY . .

# Compila la aplicación
RUN CGO_ENABLED=0 GOOS=${TARGETOS} go build -a -installsuffix cgo -o builder-cli${EXT} .

# Etapa de ejecución
FROM alpine:3.18.4

# Establece el directorio de trabajo
WORKDIR /root/

# Copia el ejecutable compilado de la etapa de compilación a la etapa de ejecución
COPY --from=builder /app/builder-cli* /root/

# Ejecuta el programa al iniciar el contenedor (solo para Linux)
ENTRYPOINT ["./builder-cli"]
