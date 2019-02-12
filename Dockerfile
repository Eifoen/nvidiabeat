FROM golang:latest as builder
RUN apt-get update
RUN apt-get install -y python2.7 python-virtualenv
WORKDIR /go/src/github.com/Eifoen/nvidiabeat
COPY . .
RUN make update
RUN make

FROM debian:stable-slim
WORKDIR /nvidiabeat
COPY --from=builder /go/src/github.com/Eifoen/nvidiabeat/nvidiabeat .
COPY --from=builder /go/src/github.com/Eifoen/nvidiabeat/fields.yml .
COPY --from=builder /go/src/github.com/Eifoen/nvidiabeat/nvidiabeat.yml .
ENV PATH="/nvidiabeat:${PATH}"
CMD ["./nvidiabeat"]

