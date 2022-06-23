package kafka

import (
	"encoding/json"
	route2 "github.com/codeedu/imersaofsfc2-simulator/application/route"
	"github.com/codeedu/imersaofsfc2-simulator/infra/kafka"
	ckafka "github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"os"
	"time"
)

 
func(k *KafkaConsumer) Consume() {
	configMap:=&ckafka.ConfigMap{
		"bootstrap.servers":os.Getenv(key:"KafkaBootstrapServers"),
		"group.id":os.Getenv(key:"KafkaConsumerGroupId"),

	}

	c, err := ckafka.NewConsumer(configMap)
	if err != nil {
		log.Fatalf("Failed to create consumer: %s" + err.Error())
	}
	//defer c.Close()
	topics := []string{os.Getenv(key:"KafkaConsumerTopic")}
	c.SubscribeTopics(topics, nil)
	fmt.Println("Subscribed to topic: " + topics[0])
	for {
		msg, err := c.ReadMessage(-1)
		if err == nil {
			k.MsgChan <- msg
		}
	}

}

func Produce(msg *ckafka.Message) {
	producer := kafka.NewKafkaProducer()
	route := route2.NewRoute()
	json.Unmarshal(msg.Value, &route)
	route.LoadPositions()
	positions, err := route.ExportJsonPositions()
	if err != nil {
		log.Println(err.Error())
	}
	for _, p := range positions {
		kafka.Publish(p, os.Getenv("KafkaProduceTopic"), producer)
		time.Sleep(time.Millisecond * 500)
	}
}