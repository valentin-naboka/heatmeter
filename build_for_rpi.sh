
#!/bin/bash

#Generating config.go file
./generate_config.sh

#Build crosscompile build inside docker
env DOCKER_BUILDKIT=1 docker build --output out .

#Create daemon service file.
read -p "Enter start month: " HM_START_MONTH
read -p "Enter end month: " HM_END_MONTH
read -p "Enter day of the report: " HM_DAY_OF_REPORT
read -p "Enter end month day of the report: " HM_END_MONTH_DAY_OF_REPORT
read -p "Enter device: " HM_DEVICE
read -p "Enter mbus ID : " HM_MBUS_ID
read -p "Enter badurate: " HM_BAUDRATE

pushd out
echo "[Unit]
Description=Heatmeter report submitter
After=network-online.target
After=time-sync.target

[Install]
WantedBy=multi-user.target

[Service]
Environment=\"HM_START_MONTH=$HM_START_MONTH\"
Environment=\"HM_END_MONTH=$HM_END_MONTH\"
Environment=\"HM_DAY_OF_REPORT=$HM_DAY_OF_REPORT\"
Environment=\"HM_END_MONTH_DAY_OF_REPORT=$HM_END_MONTH_DAY_OF_REPORT\"
Environment=\"HM_DEVICE=\"$HM_DEVICE\"\"
Environment=\"HM_MBUS_ID=$HM_MBUS_ID\"
Environment=\"HM_BAUDRATE=$HM_BAUDRATE\"
Environment=\"HM_HTML_DUMP_PATH=\"/var/log/heatmeter/\"\"

ExecStartPre=shutdown -P +10
ExecStart=/usr/local/bin/heatmeter" > heatmeter.service

#Copy build artifacts to host and archive it
tar czf ../heatmeter.tar.gz .
popd
