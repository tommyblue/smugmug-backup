FROM ubuntu:24.04

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

ADD https://dl.google.com/go/go1.26.0.linux-amd64.tar.gz ./go.tar.gz

RUN echo "aac1b08a0fb0c4e0a7c1555beb7b59180b05dfc5a3d62e40e9de90cd42f88235 go.tar.gz" | sha256sum -c - && \
    tar -C /usr/local -xzf go.tar.gz && \
    rm ./go.tar.gz

ADD . $APP
WORKDIR $APP
RUN go build -mod=vendor -v -o $GOBIN/smugmug-backup ./cmd/smugmug-backup/
