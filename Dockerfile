FROM golang:1.22 AS build
WORKDIR /src
COPY go.mod ./
COPY . .
RUN CGO_ENABLED=0 go build -trimpath -ldflags="-s -w" -o /out/validator-api ./cmd/api

FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=build /out/validator-api /validator-api
EXPOSE 8080
USER nonroot
ENTRYPOINT ["/validator-api"]
