FROM golang:1.19.2-alpine3.16 AS build
WORKDIR /build
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . ./
RUN go build -o /app

FROM gcr.io/distroless/base-debian10
WORKDIR /
COPY --from=build /app /app
EXPOSE 8080
USER user:user
ENTRYPOINT ["/app"]