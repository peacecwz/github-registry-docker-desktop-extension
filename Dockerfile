FROM golang:1.17-alpine AS builder
ENV CGO_ENABLED=0
WORKDIR /backend
COPY vm/go.* .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go mod download
COPY vm/. .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    go build -trimpath -ldflags="-s -w" -o bin/service

FROM --platform=$BUILDPLATFORM node:17.7-alpine3.14 AS client-builder
WORKDIR /ui
# cache packages in layer
COPY ui/package.json /ui/package.json
COPY ui/package-lock.json /ui/package-lock.json
RUN --mount=type=cache,target=/usr/src/app/.npm \
    npm set cache /usr/src/app/.npm && \
    npm ci
# install
COPY ui /ui
RUN npm run build

FROM alpine
LABEL org.opencontainers.image.title="Github Registry" \
    org.opencontainers.image.description="Docker Desktop Extension for Github Registry, Manage your organization or account registry" \
    org.opencontainers.image.vendor="Baris Ceviz (@Peacecwz)" \
    com.docker.desktop.extension.api.version=">= 0.2.3" \
    com.docker.desktop.extension.icon="https://raw.githubusercontent.com/peacecwz/github-registry-docker-desktop-extension/main/github.svg" \
    com.docker.extension.detailed-description="Docker Desktop Extension for Github Registry, Manage your organization or account registry" \
    com.docker.extension.publisher-url="https://github.com/peacecwz" \
    com.docker.extension.additional-urls="" \
    com.docker.extension.changelog=""

COPY --from=builder /backend/bin/service /
COPY docker-compose.yaml .
COPY metadata.json .
COPY github.svg .
COPY --from=client-builder /ui/build ui
CMD /service -socket /run/guest-services/extension-github-registry.sock
