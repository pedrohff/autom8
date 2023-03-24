run:
	ansible-playbook run.yml -i hosts --ask-vault-pass
	
runpass:
	ssh-copy-id -i /home/totoro/.ssh/glados glados@192.168.15.18

encrypt:
	ansible-vault encrypt compose/volume/influxdb/init/create-telegraf.iql compose/volume/influxdb/influxdb.conf compose/volume/influxdb/telegraf.conf compose/volume/mosquitto/config/mosquitto.conf other/dynudns/dynu.sh other/homebridge/config.json other/homebridge/auth.json other/installomv.sh

decrypt:
	ansible-vault decrypt compose/volume/influxdb/init/create-telegraf.iql compose/volume/influxdb/influxdb.conf compose/volume/influxdb/telegraf.conf compose/volume/mosquitto/config/mosquitto.conf other/dynudns/dynu.sh other/homebridge/config.json other/homebridge/auth.json other/installomv.sh
