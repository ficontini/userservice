FROM golang:1.22 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/bin/usersvc ./usersvc/cmd
RUN ls -la /app/bin


FROM alpine:latest
WORKDIR /root/
RUN apk add --no-cache libc6-compat
COPY --from=builder /app/bin/usersvc /root/bin/usersvc
RUN ls -la /root/bin
EXPOSE 3004
CMD ["/root/bin/usersvc"]
