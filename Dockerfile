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

ADD https://dl.google.com/go/go1.14.7.linux-amd64.tar.gz ./go.tar.gz

RUN echo "0f215de06019a054a3da46a0722989986c956d719c7a0a8fc38a5f3c216d6f6b go.tar.gz" | sha256sum -c - && \
    tar -C /usr/local -xzf go.tar.gz && \
    rm ./go.tar.gz

ADD . $APP
WORKDIR $APP
RUN go build -mod=vendor -i -v -o $GOBIN/smugmug-backup ./cmd/smugmug-backup/
