package rabbit

import (
	"encoding/json"

	"Crawler/env"
	"Crawler/modules/errors"
	"github.com/streadway/amqp"
)

type RConn struct {
	Conn         *amqp.Channel
	Queue        amqp.Queue
	MessagesChan <-chan amqp.Delivery
}

type Message struct {
	Rk   string
	Mess interface{}
}

func (r *RConn) InitRabbit() {
	err := r.conn(env.HostRebbit)
	errors.HandlerError(err)
	err = r.addQueue(env.Queue)
	errors.HandlerError(err)
	err = r.addExchange(env.Exchange)
	errors.HandlerError(err)
	err = r.addBind("*.urls.#", env.Exchange)
	errors.HandlerError(err)
}

func (r *RConn) conn(url string) error {
	conn, err := amqp.Dial(url)
	if err != nil {
		return err
	}
	amqpChannel, err := conn.Channel()
	if err != nil {
		return err
	}
	r.Conn = amqpChannel
	return nil
}

func (r *RConn) addQueue(name string) error {
	queue, err := r.Conn.QueueDeclare(name, true, false, false, false, nil)
	if err != nil {
		return err
	}
	r.Queue = queue
	return nil
}

func (r *RConn) addExchange(name string) error {
	err := r.Conn.ExchangeDeclare(name, "topic", true, false, false, false, nil)
	if err != nil {
		return err
	}
	return nil
}

func (r *RConn) addBind(rk string, exchange string) error {
	err := r.Conn.QueueBind(r.Queue.Name, rk, exchange, false, nil)
	if err != nil {
		return err
	}
	return nil
}

func (r *RConn) Publisher(m chan Message) {
	for {
		ms := <-m
		body, err := json.Marshal(ms.Mess)
		if err != nil {
			errors.HandlerError(err)
		}
		err = r.Conn.Publish(env.Exchange, ms.Rk, false, false, amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/plain",
			Body:         body,
		})
		if err != nil {
			errors.HandlerError(err)
		}
	}
}
