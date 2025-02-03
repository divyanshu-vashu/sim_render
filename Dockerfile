FROM golang:1.22 AS builder
WORKDIR /sim_render
COPY . .
RUN go mod download
RUN go build -o ./sim_render .

FROM debian:slim AS runner
WORKDIR /sim_render
COPY --from=builder /sim_render/sim_render .
COPY .env .
COPY static/ ./static/
ENV PORT=8080
EXPOSE ${PORT}
ENTRYPOINT ["./sim_render"]