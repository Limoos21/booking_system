FROM golang:1.24.1-alpine

RUN apk update && apk add --no-cache build-base musl-dev ca-certificates tzdata

ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64 \
    GOFLAGS="-mod=vendor"

WORKDIR /build
COPY go.mod .
COPY go.sum .
COPY . .

## Копируем .env файл
#COPY .env /build/.env

RUN go build -buildvcs=false -trimpath -a -mod=mod -o main ./cmd/main.go

# Перемещаемся в /dist каталог для хранения результирующего бинарника
WORKDIR /dist

# Копируем бинарник из каталога сборки
RUN cp /build/main .

# Собираем минимальный образi8
FROM scratch

COPY --from=0 /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=0 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=0 /dist/main /



# Копируем .env файл в конечный образ
COPY --from=0 /build/.env /.

CMD ["/main","-p"]