# ---- build stage ----
    FROM golang:1.22 AS builder

    WORKDIR /app
    
    # cache de dependências
    COPY go.mod go.sum ./
    RUN go mod download
    
    # código-fonte (inclui seu main em Go com AutoMigrate)
    COPY . .
    
    # binário estático
    RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/app .
    
    # ---- run stage ----
    FROM alpine:3.20
    
    # utilitários mínimos para esperar o DB
    RUN apk add --no-cache bash netcat-openbsd
    
    WORKDIR /app
    
    # binário e entrypoint
    COPY --from=builder /app/app /app/app
    COPY entrypoint.sh /app/entrypoint.sh
    RUN chmod +x /app/entrypoint.sh
    
    ENTRYPOINT ["/app/entrypoint.sh"]
    