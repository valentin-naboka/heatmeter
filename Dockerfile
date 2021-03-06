FROM golang:1.16.2 AS buld_heatmeter

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

COPY . heatmeter
RUN sudo chown -R debian:root heatmeter

WORKDIR ${HOME_DIR}/heatmeter

RUN bash build.sh --arm

FROM scratch AS export-stage
COPY --from=buld_heatmeter /home/debian/heatmeter/heatmeter .
