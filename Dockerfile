FROM golang AS builder
WORKDIR /smartmeter-reader
COPY ./src/ /smartmeter-reader/
ENV GOARM 7
ENV GOOS linux
ENV GOARCH arm
ENV GO111MODULE auto
RUN mkdir bin && cd /smartmeter-reader/run && go build -v -o ../bin/smartmeter_reader

FROM hypriot/rpi-alpine-scratch

# RUN apk update && \
# apk upgrade && \
# apk add bash && \
# rm -rf /var/cache/apk/*

WORKDIR /app/
COPY --from=builder /smartmeter-reader/bin/smartmeter_reader /app/
VOLUME [ "/config" ]
ENTRYPOINT [ "/app/smartmeter_reader " ]