FROM golang:1.13.8-buster as builder

COPY . /tmp/scr

RUN cd /tmp/scr/ && make build

FROM scratch

COPY --from=builder /tmp/scr/scr /usr/local/bin/scr
