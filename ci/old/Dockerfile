FROM ubuntu:latest

WORKDIR /app

# Basics
RUN apt-get update
RUN apt-get -y upgrade  
RUN apt-get install -y \
    sudo \
    apt-transport-https \
    ca-certificates \
    curl \
    wget \
    gnupg-agent \
    software-properties-common
# Docker-CE
RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
RUN add-apt-repository \
   "deb [arch=amd64] https://download.docker.com/linux/ubuntu \
   bionic \
   stable"
RUN apt-get update
RUN apt-get install -y docker-ce docker-ce-cli containerd.io
# RUN docker version

# Golang
#RUN apt-get -y install sudo wget
#RUN wget https://golang.org/dl/go1.15.3.linux-amd64.tar.gz
#RUN tar -xzf go1.15.3.linux-amd64.tar.gz
#RUN mv go ~/go
#RUN sudo chmod +x ~/go
#RUN alias go='~/go'
#RUN export PATH=$PATH:/usr/local/go/bin
#RUN cd ~/ && ./go version

# Docker daemon config
# RUN rm -rf /etc/docker/daemon.json
# COPY hpc/container/images/base/daemon.json /etc/docker/daemon.json