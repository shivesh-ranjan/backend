package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"

	"github.com/streadway/amqp"
	"shivesh-ranjan.github.io/backend/notification/utils"
)

// EmailNotification represents the structure of an email
type EmailNotification struct {
	To      string `json:"to"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

func main() {
	config, err := utils.LoadConfig(".")
	if err != nil {
		log.Fatal("Can't load config:", err)
	}

	// RabbitMQ connection settings
	amqpURL := config.AMQPURL
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	queueName := config.QueueName
	q, err := ch.QueueDeclare(
		queueName,
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare a queue: %v", err)
	}

	msgs, err := ch.Consume(
		q.Name,
		"",
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	done := make(chan bool)

	go func() {
		for d := range msgs {
			var email EmailNotification
			if err := json.Unmarshal(d.Body, &email); err != nil {
				log.Printf("Error decoding message: %v", err)
				continue
			}

			log.Printf("Sending email to: %s", email.To)
			if err := sendEmail(email, config); err != nil {
				log.Printf("Failed to send email: %v", err)
			} else {
				log.Printf("Email sent successfully to: %s", email.To)
			}
		}
	}()

	log.Println("Waiting for email messages...")
	<-done
}

func sendEmail(email EmailNotification, config utils.Config) error {
	// SMTP server configuration
	from := config.FromEmail
	password := config.EmailPassword
	smtpHost := config.SMTPHost
	smtpPort := config.SMTPPort

	// Message body
	msg := fmt.Sprintf("From: %s\nTo: %s\nSubject: %s\n\n%s", from, email.To, email.Subject, email.Body)

	// Authentication
	auth := smtp.PlainAuth("", from, password, smtpHost)

	// Sending the email
	return smtp.SendMail(
		smtpHost+":"+smtpPort,
		auth,
		from,
		[]string{email.To},
		[]byte(msg),
	)
}
