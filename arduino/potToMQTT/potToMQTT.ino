
#include <Ethernet.h>
#include <MQTT.h>

byte mac[] = {0xDE, 0xAD, 0xBE, 0xEF, 0xFE, 0xED};
byte ip[] = {192, 168, 15, 177};  // <- change to match your network

EthernetClient net;
MQTTClient client;

//set interval for sending messages (milliseconds)
const long interval = 100;
unsigned long previousMillis = 0;

int lastVol = 999;

const char broker[] = "192.168.15.9";
int        port     = 1883;
const char topic[]  = "mytest";

void connect() {
  Serial.print("connecting...");
  while (!client.connect("arduino")) {
    Serial.print(".");
    delay(1000);
  }

  Serial.println("\nconnected!");
}

void setup() {
  Serial.begin(9600);
  Ethernet.begin(mac, ip);

  client.begin(broker, port, net);

  Serial.println("You're connected to the MQTT broker!");
  Serial.println();
  connect();
}

bool volumeThreshold(int last, int now) {
  const int threshold = 1;
  int high = last + threshold;
  int low = last - threshold;
  return now > high || now < low;
}

void loop() {
  client.loop();

  if (!client.connected()) {
    connect();
  }

  unsigned long currentMillis = millis();

  if (currentMillis - previousMillis >= interval) {
    // save the last time a message was sent
    previousMillis = currentMillis;


    int aVal = analogRead(A0);
    int currentVolume = map(aVal, 0, 1023, 0, 100);

    if (volumeThreshold(lastVol, currentVolume)) {
      char message[16];
      sprintf(message, "%d", currentVolume);
      client.publish(topic, message);
      Serial.print("message sent: ");
      Serial.println(message);
      
      lastVol = currentVolume;
    }

  }
}