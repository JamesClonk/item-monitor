FROM ubuntu:22.04

ARG package_args='--allow-downgrades --allow-remove-essential --allow-change-held-packages --no-install-recommends'
RUN echo "debconf debconf/frontend select noninteractive" | debconf-set-selections && \
  export DEBIAN_FRONTEND=noninteractive && \
  apt-get -y $package_args update && \
  apt-get -y $package_args dist-upgrade && \
  apt-get -y $package_args install curl ca-certificates gnupg tzdata git
RUN curl --location --output go.tar.gz "https://go.dev/dl/go1.21.3.linux-amd64.tar.gz" && \
  echo "1241381b2843fae5a9707eec1f8fb2ef94d827990582c7c7c32f5bdfbfd420c8  go.tar.gz" | sha256sum -c  && \
  tar -C /usr/local -xzf go.tar.gz && \
  rm go.tar.gz

ENV PATH=$PATH:/usr/local/go/bin

WORKDIR /go/src/github.com/JamesClonk/item-monitor
COPY . .
RUN go build -o item-monitor

FROM ubuntu:22.04

LABEL maintainer="JamesClonk <jamesclonk@jamesclonk.ch>"

ARG package_args='--allow-downgrades --allow-remove-essential --allow-change-held-packages --no-install-recommends'
RUN echo "debconf debconf/frontend select noninteractive" | debconf-set-selections && \
  export DEBIAN_FRONTEND=noninteractive && \
  apt-get -y $package_args update && \
  apt-get -y $package_args dist-upgrade && \
  apt-get -y $package_args install curl ca-certificates gnupg tzdata bash && \
  apt-get clean && \
  find /usr/share/doc/*/* ! -name copyright | xargs rm -rf && \
  rm -rf \
  /usr/share/man/* /usr/share/info/* \
  /var/lib/apt/lists/* /tmp/*

RUN useradd -u 1000 -mU -s /bin/bash -d /item-monitor item-monitor && \
  mkdir /item-monitor/app && \
  chown item-monitor:item-monitor /item-monitor/app

ENV PATH=$PATH:/item-monitor/app
WORKDIR /item-monitor/app
COPY public ./public/
COPY static ./static/
COPY --from=0 /go/srcgithub.com/JamesClonk/item-monitor/item-monitor ./item-monitor

RUN chmod +x /item-monitor/app/item-monitor && \
  chown -R item-monitor:item-monitor /item-monitor/app
USER item-monitor

EXPOSE 8080

CMD ["./item-monitor"]
