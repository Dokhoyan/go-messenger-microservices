# syntax=docker/dockerfile:experimental

# MULTI STAGE EXAMPLE

#######################################
# STAGE 1. BUILD STAGE
#######################################

# Для сборки нашего приложения на go требуется образ ОС, в котором установлен golang нужной нам версии.
# Alpine выбран из-за его небольшого размера по сравнению с Ubuntu.
FROM golang:1.24.1-alpine AS build

# Устанавливаем место назначения для COPY
WORKDIR /app

# Копируем файлы go.mod и go.sum в WORKDIR
COPY go.mod go.sum ./
# Скачиваем необходимые Go модули (зависимости нашего проетка)
RUN go mod download

# Копируем все исходные go файлы нашего проекта в образ
# https://docs.docker.com/reference/dockerfile/#copy
COPY *.go ./
# Собираем бинарный файл нашего приложения
RUN CGO_ENABLED=0 GOOS=linux go build -o /bin/main

#######################################
# STAGE 2. FINAL STAGE
#######################################

FROM scratch AS final

WORKDIR /

COPY --from=build /bin/main /main

# Указываем какой порт необходимо слушать
# https://docs.docker.com/reference/dockerfile/#expose
EXPOSE 8080

# Точка входа
ENTRYPOINT ["/main"]