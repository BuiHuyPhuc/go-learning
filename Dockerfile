# FROM golang:alpine

# WORKDIR /build

# COPY  . .

# RUN go build -o main ./cmd/server/main.go

# WORKDIR /dist

# RUN cp /build/main .

# EXPOSE 8888

# CMD [ "/dist/main" ]

### Image size 567MB


# -------- Build stage --------
FROM golang:alpine AS builder

WORKDIR /build

RUN apk add --no-cache git


COPY go.mod go.sum ./
RUN go mod download

COPY  . .

# Build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o main ./cmd/server/main.go


# -------- Run stage --------
FROM scratch

WORKDIR /

COPY --from=builder /build/main /main
COPY ./config /config

EXPOSE 8888

ENTRYPOINT [ "/main" ]
CMD [ "config/local.yaml" ]

### Image size 10MB
