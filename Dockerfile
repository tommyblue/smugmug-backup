FROM ubuntu:20.04

ENV DEBIAN_FRONTEND=noninteractive \
    LANG=en_US.UTF-8 \
    LANGUAGE=en_US.UTF-8 \
    LC_ALL=en_US.UTF-8 \
    TERM=xterm \
    TZ=:/etc/localtime \
    PATH=$PATH:/usr/local/go/bin \
    GOBIN=/go/bin \
    APP=/go/src/smugmug-backup

RUN sed -e "/deb-src/d" -i /etc/apt/sources.list \
    && apt-get update \
    && apt-get install --no-install-recommends --yes \
        ca-certificates \
    && apt-get clean \
    && rm -rf /.root/cache \
    && rm -rf /var/lib/apt/lists/*

ADD https://dl.google.com/go/go1.14.4.linux-amd64.tar.gz ./go.tar.gz

RUN echo "aed845e4185a0b2a3c3d5e1d0a35491702c55889192bb9c30e67a3de6849c067 go.tar.gz" | sha256sum -c - && \
    tar -C /usr/local -xzf go.tar.gz && \
    rm ./go.tar.gz

ADD . $APP
WORKDIR $APP
RUN go build -mod=vendor -i -v -o $GOBIN/smugmug-backup .
