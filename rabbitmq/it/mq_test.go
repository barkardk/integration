// +build rabbitintegration

package it

import (
	"fmt"
	"github.com/streadway/amqp"
	"os"
	"testing"
)

type mqResources struct {
	queueName string
	messageBody string
	expectedMessageCount int
}

func defaultConnStr() string{
	defaultConnStr, exists := os.LookupEnv("RABBITMQ_AMQP_CONN_STR")
	if exists {
		return defaultConnStr
	}
	return "amqp://guest:guest@127.0.0.1:5672/"
}


func RabbitMQAMQPConnection() (*amqp.Connection, error)  {
	conn, err := amqp.Dial(defaultConnStr())
	if err != nil {
		return nil, fmt.Errorf("Could not connect to rabbitMQ %s", err)
	}
	return conn, nil
}

func TestRabbitMQWithQueuesAndMessage(t *testing.T) {
	qn := &mqResources{
		queueName: "MQTestQueue",
		messageBody: "This is a test message",
		expectedMessageCount: 1,
	}

	conn, err := RabbitMQAMQPConnection()
	if err != nil {
		t.Fatalf("Could not create RabbitMQ connection %v", err)
	}
	t.Run("create queue", func(t *testing.T){
		if err := RabbitMQCreateQueue(conn, qn); err != nil {
			t.Errorf( "Error connect to RabbitMQ and create a queue %v", err)
		}
	})

	t.Run("publish message", func(t *testing.T){
		if err := RabbitMQPublishMessage(conn, qn); err != nil {
			t.Errorf("Error publish message to RabbitMQ %v", err)
		}
	} )
	t.Run("consume message", func(t *testing.T){
		if err := RabbitMQConsumeMessage(conn, qn); err != nil {
			t.Errorf("Error consuming message from RabbitMQ %v", err)
		}
	})
}

func RabbitMQCreateQueue(conn *amqp.Connection, m *mqResources) error {
	ch, err := conn.Channel()
	if err != nil {
		fmt.Errorf("could not open channel %v", err)
	}
	q, err := ch.QueueDeclare(m.queueName, false, false, false, false, nil)
	if err != nil {
		fmt.Errorf("failed to declare RabbitMQ queue %v", err)
	}
	if m.queueName != q.Name {
		return fmt.Errorf("expected queue name (%v) got (%v)", m.queueName, q.Name)
	}
	return nil
}

func  RabbitMQPublishMessage(conn *amqp.Connection, m *mqResources) error {
	ch, err := conn.Channel()
	if err != nil {
		fmt.Errorf("could not open channel %s", err)
	}
	if err = ch.Publish(
		"",
		m.queueName,
		false,
		false,
		amqp.Publishing{
			ContentType: "text/plain",
			Body: []byte(m.messageBody),
		}); err != nil {
		fmt.Errorf("failed to publish message %s", err)
	}
	if ch.Confirm(false); err != nil {
		fmt.Errorf("message published could not be confirmed %s", err)
	}
	result, _:= ch.QueueInspect(m.queueName)
	if result.Messages != m.expectedMessageCount {
		return fmt.Errorf("expected messagecount %d got messageCount %d", m.expectedMessageCount, result.Messages)
	}
	return nil
}

func RabbitMQConsumeMessage(conn *amqp.Connection, m *mqResources) error{
	ch, err := conn.Channel()
	if err != nil {
		fmt.Errorf("fould not open channel %s", err)
	}
	msgs, err := ch.Consume(
		m.queueName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)

	if err != nil {
		fmt.Errorf("failed to register a consumer %s", err)
	}
	for d := range msgs {
		if string(d.Body) !=  m.messageBody {
			fmt.Errorf("Expected message body %s got %s", d.Body, m.messageBody)
		}
		d.Ack(false)
	}
	return nil
}

