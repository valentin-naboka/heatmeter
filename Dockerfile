FROM golang:1.16.2

ARG DEBIAN_FRONTEND=noninteractive

RUN apt-get update
RUN apt-get upgrade -y
RUN apt-get install gcc-arm-linux-gnueabihf -y sudo -y make -y git -y autoconf -y libtool -y

RUN apt-get clean

RUN echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

RUN useradd -rm -d /home/debian -s /bin/bash -g root -G sudo debian
USER debian

ARG HOME_DIR="/home/debian"
WORKDIR ${HOME_DIR}

ADD . heatmeter
RUN sudo chown -R debian:root heatmeter

WORKDIR ${HOME_DIR}/heatmeter

RUN bash build_mbus.sh --arm
