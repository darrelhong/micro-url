FROM golang:1.22.1-alpine AS build-stage

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -tags netgo -ldflags '-s -w' -o /server

FROM gcr.io/distroless/base-debian11 AS build-release-stage

WORKDIR /

COPY --from=build-stage /server /server
COPY --from=build-stage /app/templates /templates
COPY --from=build-stage /app/static /static

EXPOSE 8080

USER nonroot:nonroot

ENTRYPOINT ["/server"]