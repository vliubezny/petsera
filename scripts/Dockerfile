ARG ALPINE_VERSION=3.10
ARG GOLANG_VERSION=1.18.1
ARG NODE_VERSION=16.15.0

FROM node:${NODE_VERSION}-alpine as uiBuilder

ARG NPM_VERSION=8.10.0
RUN npm install -g npm@${NPM_VERSION}

WORKDIR /src

COPY ./ui/package.json ./ui/package-lock.json /src/
RUN npm ci

COPY ./ui /src/
RUN npm run build


FROM golang:${GOLANG_VERSION} as goBuilder
WORKDIR /go/src/github.com/vliubezny/petsera
COPY . .
COPY --from=uiBuilder /src/dist ui/dist/
RUN make linux

FROM alpine:${ALPINE_VERSION}
RUN apk update && apk add --no-cache ca-certificates
ENV POSTGRES_MIGRATIONS="/migrations/postgres"
COPY scripts/migrations /migrations
COPY --from=goBuilder /go/src/github.com/vliubezny/petsera/build/petsera-linux-amd64 /petsera
ENTRYPOINT [ "/petsera" ]