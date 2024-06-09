package queue

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strconv"

	amqp "github.com/rabbitmq/amqp091-go"
)

const (
	exchangeName = "exchangeEvent"
	queueName    = "eventsNotify"
)

type Queue struct {
	conn *amqp.Connection
	ch   *amqp.Channel
}

func New(host string, port uint16, userName string, password string) (*Queue, error) {
	var (
		queue = &Queue{}
		err   error
	)

	queue.conn, err = amqp.Dial(fmt.Sprintf(
		"amqp://%s:%s@%s/",
		userName, password,
		net.JoinHostPort(host, strconv.Itoa(int(port))),
	))
	if err != nil {
		return nil, err
	}

	queue.ch, err = queue.conn.Channel()
	if err != nil {
		return nil, errors.Join(err, queue.conn.Close())
	}

	if err := queue.ch.ExchangeDeclare(
		exchangeName, // exchangName
		"direct",     // kind
		true,         // durable
		false,        // autoDelete
		false,        // internal
		true,         // noWait
		nil,          // args
	); err != nil {
		return nil, err
	}

	if _, err := queue.ch.QueueDeclare(
		queueName, // queueName
		false,     // durable
		false,     // autoDelete
		false,     // exclusive
		true,      // noWait
		nil,       // args
	); err != nil {
		return nil, err
	}

	if err := queue.ch.QueueBind(
		queueName,    // queueName
		queueName,    // key
		exchangeName, //exchangeName
		true,         // noWait
		nil,          // args
	); err != nil {
		return nil, err
	}

	return queue, nil
}

func (q *Queue) Send(ctx context.Context, body []byte) error {
	return q.ch.PublishWithContext(
		ctx,
		exchangeName,
		queueName,
		false,
		true,
		amqp.Publishing{
			Body: body,
		},
	)
}

func (q *Queue) Consume(ctx context.Context) (<-chan []byte, error) {
	delivery, err := q.ch.ConsumeWithContext(
		ctx,
		queueName, // queueName
		"",        // consumer
		true,      // autoAck
		false,     // exclusive
		false,     //noLocal
		true,      // noWait
		nil,       // args
	)
	if err != nil {
		return nil, err
	}
	out := make(chan []byte)
	go func() {
		defer close(out)
		for {
			select {
			case <-ctx.Done():
				return
			case message := <-delivery:
				out <- message.Body
			}
		}
	}()
	return out, nil
}

func (q *Queue) Close() error {
	return errors.Join(q.ch.Close(), q.conn.Close())
}
