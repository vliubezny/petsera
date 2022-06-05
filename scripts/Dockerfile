ARG ALPINE_VERSION=3.10
ARG GOLANG_VERSION=1.18.1

FROM golang:${GOLANG_VERSION}
WORKDIR /go/src/github.com/vliubezny/petsera
COPY . .
RUN mkdir -p ui/dist && cp index.html ui/dist/
RUN make linux

FROM alpine:${ALPINE_VERSION}
RUN apk update && apk add --no-cache ca-certificates
ENV POSTGRES_MIGRATIONS="/migrations/postgres"
COPY scripts/migrations /migrations
COPY --from=0 /go/src/github.com/vliubezny/petsera/build/petsera-linux-amd64 /petsera
ENTRYPOINT [ "/petsera" ]