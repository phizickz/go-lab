# argument for Go version
ARG GO_VERSION=1.18
 
# STAGE 1: building the executable
FROM golang:${GO_VERSION}-alpine AS build
 
RUN apk add --no-cache git
WORKDIR /src
COPY src/go.mod ./
# COPY ./go.sum ./
RUN go mod download
COPY src/ ./

# Build the executable
RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /app
 
# STAGE 2: build the container to run
FROM gcr.io/distroless/static AS final
 
USER nonroot:nonroot
 
# copy compiled app
COPY --from=build --chown=nonroot:nonroot /app /app
 
# run binary; use vector form
ENTRYPOINT ["/app"]