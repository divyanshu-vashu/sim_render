FROM golang:1.22 AS builder
    WORKDIR /app
    COPY . .
    RUN go mod download
    RUN go build -o ./sim_render .

FROM alpine:latest AS runner
WORKDIR /app
COPY --from=builder /app/sim_render .
COPY .env .
COPY static/ ./static/
ENV PORT=8080
EXPOSE ${PORT}
ENTRYPOINT ["./sim_render"]