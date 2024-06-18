package utils

import (
	"encoding/json"
	"github.com/Andrewalifb/alpha-pos-system-sales-service/dto"
	"github.com/streadway/amqp"
)

func SendDigitalReceipt(receipt dto.DigitalReceipt, ch *amqp.Channel, queueName string) error {
	body, err := json.Marshal(receipt)
	if err != nil {
		return err
	}

	err = ch.Publish(
		"",        // exchange
		queueName, // routing key
		false,     // mandatory
		false,     // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	return err
}
