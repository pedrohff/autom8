# Autom8

Conjunto de projetos usados para estudos de automação residencial.

## Arduino

Protótipos ainda em testes

## Golang

Projetos configurados com auxílio do *go workspace*. Ainda não sei se é uma boa ferramenta, e uso o projeto para testá-lo.

### gladoscd

Tentativa de atualizar automaticamente um container docker quando uma tag nova for gerada, usada nos demais projetos daqui.

### homebridgepooler

Tenho conectado ao HomeBridge um sensor de temperatura e umidade da XiaoMI. Essa aplicação usa da API do HomeBridge para entregar seus dados em um broker MQTT, seguindo um intervalo configurado.

### mqtttovolume

Código de testes para ler de um tópico MQTT, e seu valor refletir no volume da máquina que o executa.

### serialtovolume

Projeto dependente de um Arduino que escreve na porta serial o valor de um potenciômetro. O código lê os valores da porta serial e altera o volume de acordo com o valor lido.


## TinyGo

Os projetos usando *tinygo* estão separados dos projetos *go* por problemas de configurações em IDEs, uma vez que o tinygo requer um `GOROOT` manipulado, e isso acaba impacta diretamente interpretação das ferramentas.


## Servers

### glados

Automação das ferramentas usadas em meu Raspberry Pi.
Aqui instalo ferramentas do sistema operacional, o HomeBridge, o Docker e algumas imagens que utilizo para minha stack IoT.