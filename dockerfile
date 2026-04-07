FROM golang:1.25.5-alpine AS builder

WORKDIR /arch

COPY ./go.mod ./go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-s -w" -o /arch/white_archive .

FROM alpine:3.20

ENV USER=runner

WORKDIR /arch

RUN adduser -D -h /home/$USER $USER

RUN chown -R $USER:$USER /arch

COPY --from=builder /arch/white_archive .

COPY entrypoint.sh .

ENTRYPOINT [ "./entrypoint.sh" ]