
#!/bin/bash

sudo useradd -r -s /bin/false --no-create-home heatmeter
mkdir -p out & tar -zxf heatmeter.tar.gz -C out
cd out 
sudo cp heatmeter /usr/local/bin
sudo chown heatmeter:heatmeter /usr/local/bin/heatmeter
sudo chmod 500 /usr/local/bin/heatmeter

sudo cp heatmeter.service /etc/systemd/system/
sudo chmod 644 /etc/systemd/system/heatmeter.service

sudo mkdir -p /var/log/heatmeter
sudo chown heatmeter:heatmeter /var/log/heatmeter
sudo chmod 644 /etc/systemd/system/heatmeter.service

systemctl daemon-reload
systemctl enable heatmeter.service

