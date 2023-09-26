FROM golang:1.21.1 AS BUILD
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 make build

FROM gcr.io/distroless/static-debian12
ENV GIN_MODE=release
USER nonroot:nonroot
WORKDIR /app
COPY --from=BUILD /app/bin/backend .
CMD ["./backend"]
