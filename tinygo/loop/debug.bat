tinygo build -scheduler tasks -o loop -target arduino ./main.go
tinygo flash -scheduler tasks -target arduino -port COM3 loop
.\arduino-cli.exe monitor -p COM3