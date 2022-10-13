FROM golang:1.19

WORKDIR /usr/local/src

COPY . .

RUN --mount=type=cache,target=/root/.cache/go-build \
  CGO_ENABLED=0 GOARCH=$TARGETARCH go build -o /usr/local/bin ./cmd/...



FROM scratch

COPY --from=0 /usr/local/bin/* /usr/local/bin/

ENTRYPOINT ["/usr/local/bin/powerdns-zone-provisioner"]
