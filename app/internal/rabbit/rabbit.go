package rabbit

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"
	"user-auth/internal/models"

	amqp "github.com/rabbitmq/amqp091-go"
)

var (
	exchange    = "user_deletion"
	exchangeDLX = "user_deletion_timeout"
	queueDLX    = "user_deletion_timeout_queue"
	tasksKey    = "tasks_key"
	usersKey    = "users_key"
	expiration  = "259200000" //72h
	queue       = "user_deletion_queue"
)

type Message struct {
	Id int
}

type DeletionProducer struct {
	amqpChan *amqp.Channel
}

func NewProducer(conn amqp.Connection) (*DeletionProducer, error) {
	pr := &DeletionProducer{}
	ch, err := conn.Channel()
	pr.amqpChan = ch
	if err != nil {
		return nil, fmt.Errorf("failed to open a channel: %v", err)
	}
	err = setChannel(ch)
	if err != nil {
		return nil, err
	}
	return pr, nil
}

func setChannel(ch *amqp.Channel) error {
	err := ch.ExchangeDeclare(exchange, "direct", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare Exchange: %v", err)
	}
	err = ch.ExchangeDeclare(exchangeDLX, "fanout", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("failed to declare Exchange: %v", err)
	}

	_, err = ch.QueueDeclare(queueDLX, false, false, false, false,
		amqp.Table{"x-dead-letter-exchange": exchange},
	)
	if err != nil {
		return fmt.Errorf("failed to declare Queue: %v", err)
	}

	err = ch.QueueBind(queueDLX, tasksKey, exchangeDLX, false, nil)
	if err != nil {
		return fmt.Errorf("failed to bind Queue: %v", err)
	}
	return nil
}

func (p *DeletionProducer) Delete(id int) error {

	message := Message{id}
	err := p.publish(message, usersKey)
	if err != nil {
		return fmt.Errorf("failed to publish: %v", err)
	}
	err = p.publish(message, tasksKey)
	if err != nil {
		return fmt.Errorf("failed to publish: %v", err)
	}
	return nil
}

func (p *DeletionProducer) publish(message Message, key string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	body, _ := json.Marshal(message)
	err := p.amqpChan.PublishWithContext(ctx, exchangeDLX, usersKey, false, false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(body),
			Expiration:  expiration,
		})
	if err != nil {
		return err
	}
	return nil
}

func ConsumeDeletion(model models.UserModel, conn *amqp.Connection) {

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(exchange, "direct", true, false, false, false, nil)
	failOnError(err, "Failed to declare Exchange")

	err = ch.ExchangeDeclare(exchangeDLX, "fanout", true, false, false, false, nil)
	failOnError(err, "Failed to declare Exchange")

	_, err = ch.QueueDeclare(queue, false, false, false, false, nil)
	failOnError(err, "Failed to declare a queue")

	err = ch.QueueBind(queue, usersKey, exchange, false, nil)
	failOnError(err, "Failed to bind Queue")

	msgs, err := ch.Consume(queue, "", true, false, false, false, nil)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)
			message := Message{}
			err := json.Unmarshal(d.Body, &message)
			if err != nil {
				log.Printf("Error:unmarshal message:%v", err)
			}
			model.Delete(message.Id)
			if err != nil {
				log.Printf("Error: deleting user %d:%v", message.Id, err)
			}
		}
	}()

	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}
