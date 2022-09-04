package examples

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/muhammad-fakhri/go-libs/messaging"
	"github.com/muhammad-fakhri/go-libs/messaging/pubsub"
)

func main() {
	p := pubsub.NewPublisher(pubsub.RabbitMQ, &pubsub.PubConfig{
		ServerURL: fmt.Sprintf("amqp://guest:guest@localhost:5672/"),
		VHost:     "vhost",
	})

	ticker := time.NewTicker(time.Second)
	go func() {
		for t := range ticker.C {
			err := p.Publish("test.fanout", "alalalalala"+t.String())
			logs.Error(err)
		}
	}()

	s := pubsub.NewSubscriber(pubsub.RabbitMQ, &pubsub.SubConfig{
		ServerURL: fmt.Sprintf("amqp://guest:guest@localhost:5672/"),
		List: []*pubsub.Sub{
			{
				Name:    "test",
				Topic:   "test.fanout",
				Handler: test,
			},
		},
		VHost: "vhost",
	})

	s.Start()

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c
	s.Stop()
	log.Println("shutting down")
	os.Exit(0)
}

func test(msg *messaging.Message) error {
	log.Println("woww", string(msg.Body))
	return nil
}
