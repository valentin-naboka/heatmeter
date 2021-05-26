# Heatmeter

The heatmeter is an automated tool to get a measurement from [Sharky 774](https://www.diehl.com/metering/en/portfolio/thermal-energy-metering-/thermal-energy-metering-product/sharky-774-compact/238298/) and submit data to [utilities of heating networks](https://www.hts.kharkov.ua/). It relies on M-Bus protocol to get a measurement from the heatmeter via MBus to USB converter.

Basically, it is designed to work in tandem with Raspberry Pi, but also can be used with an arbitrary operating system and chips that has the support of [Golang](https://golang.org/).
The tool consist of two parts:
1. M-Bus reader.
2. Data submitter.

M-Bus reader is a Go wrapper around [M-bus library](https://github.com/rscada/libmbus) which is written in C. 
Data submitter represents automated user actions such as login to a personal account and submits measurement to heating networks. 

To build for RPi just run build_for_rpi.sh and you'll get an archive heatmeter.tar.gz. During the build process, you'll be asked for some settings like heating networks account credential or a host to submit the data. Copy the archive and install.sh script to RPi and run install.sh to install the heatmeter.

To build on a local machine for a debug purpose:

1. Run ./generate_config.sh in order to generate a config.

2. Run build.sh.
