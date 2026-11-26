FROM golang:1.25-alpine AS build

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG GO_BUILD_TAGS=""
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -tags "${GO_BUILD_TAGS}" -o /out/docs-svc ./cmd/docs-svc

FROM alpine:3.20
RUN apk add --no-cache ca-certificates git openssh-client \
    && adduser -D -u 10001 appuser
WORKDIR /app
COPY --from=build /out/docs-svc /app/docs-svc
USER appuser
EXPOSE 8084
ENTRYPOINT ["/app/docs-svc"]
