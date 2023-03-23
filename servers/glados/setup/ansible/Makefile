run:
	ansible-playbook run.yml -i hosts --ask-vault-pass

encrypt:
	ansible-vault encrypt compose/volume/influxdb/init/create-telegraf.iql compose/volume/influxdb/influxdb.conf compose/volume/influxdb/telegraf.conf compose/volume/mosquitto/config/mosquitto.conf other/dynudns/dynu.sh other/homebridge.json

decrypt:
	ansible-vault decrypt compose/volume/influxdb/init/create-telegraf.iql compose/volume/influxdb/influxdb.conf compose/volume/influxdb/telegraf.conf compose/volume/mosquitto/config/mosquitto.conf other/dynudns/dynu.sh other/homebridge.json
