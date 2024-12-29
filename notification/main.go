package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

// RabbitMQ configuration
var (
	rabbitMQHost = os.Getenv("AMQP_URL")
	exchangeName = os.Getenv("EXCHANGE_NAME")
	routingKey   = os.Getenv("ROUTING_KEY")
	queueName    = os.Getenv("QUEUE_NAME")
)

// Email configuration
var (
	smtpServer    = os.Getenv("SMTP_HOST")
	smtpPort      = os.Getenv("SMTP_PORT")
	emailUser     = os.Getenv("FROM_EMAIL")
	emailPassword = os.Getenv("EMAIL_PASSWORD")
)

// Message represents the structure of RabbitMQ message
type Message struct {
	Email   string   `json:"email"`
	Links   []string `json:"links"`
	Summary string   `json:"summary"`
}

// sendEmail sends an email with the given summary and links to the recipient.
func sendEmail(toEmail string, summary string, links []string) {
	// Create the email body
	body := fmt.Sprintf(
		"Hello,\n\nHere is your summary:\n\n%s\n\nLinks:\n%s\n\nBest regards.",
		summary,
		formatLinks(links),
	)

	// Set up authentication
	auth := smtp.PlainAuth("", emailUser, emailPassword, smtpServer)

	// Create the email headers and body
	msg := []byte(fmt.Sprintf(
		"Subject: Summary and Links\n\n%s",
		body,
	))

	// Send the email
	err := smtp.SendMail(fmt.Sprintf("%s:%s", smtpServer, smtpPort), auth, emailUser, []string{toEmail}, msg)
	if err != nil {
		log.Printf("Failed to send email to %s: %v", toEmail, err)
	} else {
		log.Printf("Email sent to %s", toEmail)
	}
}

// formatLinks formats the list of links as a string with each link on a new line
func formatLinks(links []string) string {
	formatted := ""
	for _, link := range links {
		formatted += link + "\n"
	}
	return formatted
}

// handleMessage processes messages from RabbitMQ
func handleMessage(d amqp.Delivery) {
	var msg Message
	err := json.Unmarshal(d.Body, &msg)
	if err != nil {
		log.Printf("Failed to decode message: %v", err)
		return
	}

	if msg.Email == "" || msg.Summary == "" || len(msg.Links) == 0 {
		log.Printf("Invalid message format: Missing email, summary, or links")
		return
	}

	sendEmail(msg.Email, msg.Summary, msg.Links)

	// Acknowledge the message
	d.Ack(false)
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Failed to load env vars: %v", err)
	}

	// Connect to RabbitMQ
	conn, err := amqp.Dial(rabbitMQHost)
	if err != nil {
		log.Print(rabbitMQHost)
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	// Declare exchange and queue
	err = ch.ExchangeDeclare(exchangeName, "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare exchange: %v", err)
	}

	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	err = ch.QueueBind(queueName, routingKey, exchangeName, false, nil)
	if err != nil {
		log.Fatalf("Failed to bind queue: %v", err)
	}

	// Set up consumer
	msgs, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	log.Printf("Waiting for messages in queue '%s'...", queueName)

	// Consume messages
	for msg := range msgs {
		handleMessage(msg)
	}
}
