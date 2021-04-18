# Heatmeter

The heatmeter is an automated tool to get measurement from [Sharky 774](https://www.diehl.com/metering/en/portfolio/thermal-energy-metering-/thermal-energy-metering-product/sharky-774-compact/238298/) and sumbit data to [utilities of heating networks](https://www.hts.kharkov.ua/). It relies on M-Bus protocol to get measurement from the heatmeter via MBus to USB converter.
	Basically, it is designed to work in tandem with Raspberry Pi, but also can be used with arbitrary operating system and chips that has support of [Golang](https://golang.org/).
The tool consist of two parts:
1. M-Bus reader.
2. Data submitter.

M-Bus reader is a Go wraper around [M-bus library](https://github.com/rscada/libmbus) written in C. 
Data submitter represents automated user actions such as login to personal account and submit measurement to heating networks. 
  
//TODO:

Build steps:

  

1. In order to generate a config run ./generate_config.sh

2. run docker . from the source directory