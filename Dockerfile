FROM golang:1.22 AS builder
WORKDIR /app
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 go build -o ./sim_render .

FROM alpine:latest AS runner
RUN apk add --no-cache ca-certificates libc6-compat
WORKDIR /app
COPY --from=builder /app/sim_render .
COPY .env .
COPY static/ ./static/
ENV PORT=8080
EXPOSE ${PORT}
ENTRYPOINT ["./sim_render"]