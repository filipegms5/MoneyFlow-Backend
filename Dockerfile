ARG GO_VERSION=1.25.0
FROM golang:${GO_VERSION}-alpine AS builder


# Instala as dependencias necessarias para compilar a aplicação.
RUN apk add --no-cache gcc musl-dev make

WORKDIR /app

# Copia o go.mod e faz o download das dependencias.
COPY go.mod go.sum ./
RUN go mod download

# Copia o código da aplicação e compila o binario.
COPY . .

ENV PORT :8000

RUN apk add --no-cache build-base sqlite-dev
RUN CGO_ENABLED=1 GOOS=linux go build -o /server
################################################

### Step 2: Copiar o binario do stage anterior para a imagem final.
FROM scratch

# Copy required libraries
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /lib/ld-musl-x86_64.so.1 /lib
COPY --from=builder /app/server /


# Define o ponto de entrada para o container como /server.
# O binario será executado quando o container for iniciado.
ENTRYPOINT ["/server"]