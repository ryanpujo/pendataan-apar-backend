# Stage 1: Build binary aplikasi Go
FROM golang:1.27-alpine AS builder

WORKDIR /app

# Copy go.mod dan go.sum terlebih dahulu agar layer dependency ter-cache
COPY go.mod go.sum ./
RUN go mod download

# Copy seluruh source code dan compile menjadi binary statis
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o apar-backend ./cmd/api/main.go
# Catatan: Sesuaikan "./cmd/api/main.go" dengan letak file main.go Anda (misal "." jika di root)

# Stage 2: Image tipis untuk menjalankan binary
FROM alpine:3.20

# Tambahkan tzdata agar zona waktu Asia/Jakarta berfungsi dan ca-certificates untuk HTTPS
RUN apk --no-cache add ca-certificates tzdata

WORKDIR /app

# Ambil hanya file binary hasil compile dari stage builder
COPY --from=builder /app/apar-backend .

EXPOSE 8080

CMD ["./apar-backend"]