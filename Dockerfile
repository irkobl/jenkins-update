FROM golang:1.23
RUN apt-get update -qq && apt-get -y install apt-transport-https ca-certificates lsb-release \
    sudo make gcc zip unzip git curl wget vim build-essential dh-make default-jdk default-jre && \
    rm -rf /var/lib/apt/lists/* rm -rf /var/cache/apt/*
RUN useradd -ms /bin/bash admin && usermod -aG sudo admin && echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers && echo 'root:root' | chpasswd
USER admin

WORKDIR /home/admin
COPY . .

# RUN GO_ENABLED=0 go install -ldflags "-s -w -extldflags '-static'" github.com/go-delve/delve/cmd/dlv@latest



# CMD [ "/go/bin/dlv", "--wd=./cmd", "--listen=:2345", "--headless=true", "--log=true", "--accept-multiclient", "--api-version=2", "debug", ".", "-- backupFile" ]
