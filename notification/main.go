package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/smtp"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
)

var (
	rabbitMQHost  string
	exchangeName  string
	routingKey    string
	queueName     string
	smtpServer    string
	smtpPort      string
	emailUser     string
	emailPassword string
)

// Message represents the structure of RabbitMQ message
type Message struct {
	Email   string   `json:"email"`
	Links   []string `json:"links"`
	Summary string   `json:"summary"`
}

func loadEnvVars() {
	err := godotenv.Load()
	if err != nil {
		log.Printf("No .env file found. Continuing with system environment variables.")
	}

	rabbitMQHost = getEnv("AMQP_URL")
	exchangeName = getEnv("EXCHANGE_NAME")
	routingKey = getEnv("ROUTING_KEY")
	queueName = getEnv("QUEUE_NAME")
	smtpServer = getEnv("SMTP_HOST")
	smtpPort = getEnv("SMTP_PORT")
	emailUser = getEnv("FROM_EMAIL")
	emailPassword = getEnv("EMAIL_PASSWORD")
}

func getEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("Environment variable %s is not set", key)
	}
	return value
}

func sendEmail(toEmail string, summary string, links []string) {
	log.Print("Recieved email for sending...")
	body := fmt.Sprintf(
		"Hello,\n\nHere is your summary:\n\n%s\n\nLinks:\n%s\n\nBest regards.",
		summary,
		formatLinks(links),
	)

	log.Print("Doing plain auth for smtp...")
	auth := smtp.PlainAuth("", emailUser, emailPassword, smtpServer)
	msg := []byte(fmt.Sprintf(
		"From: %s\nTo: %s\nSubject: Summary and Links\nContent-Type: text/plain; charset=\"utf-8\"\n\n%s",
		emailUser, toEmail, body,
	))

	log.Print("Sending mail to cleint...")
	err := smtp.SendMail(fmt.Sprintf("%s:%s", smtpServer, smtpPort), auth, emailUser, []string{toEmail}, msg)
	if err != nil {
		log.Printf("Failed to send email to %s: %v", toEmail, err)
	} else {
		log.Printf("Email sent to %s", toEmail)
	}
}

func formatLinks(links []string) string {
	formatted := ""
	for _, link := range links {
		formatted += link + "\n"
	}
	return formatted
}

func handleMessage(d amqp.Delivery) {
	log.Printf("Received a message: %s", string(d.Body))

	var msg Message
	err := json.Unmarshal(d.Body, &msg)
	if err != nil {
		log.Printf("Failed to decode message: %v", err)
		d.Nack(false, false) // Reject the message without requeueing
		return
	}

	log.Printf("Processing email: %s", msg.Email)
	if msg.Email == "" || msg.Summary == "" || len(msg.Links) == 0 {
		log.Printf("Invalid message: Missing email, summary, or links")
		d.Nack(false, false)
		return
	}

	sendEmail(msg.Email, msg.Summary, msg.Links)

	err = d.Ack(false)
	if err != nil {
		log.Printf("Failed to acknowledge message: %v", err)
	}
}

func main() {
	loadEnvVars()

	conn, err := amqp.Dial(rabbitMQHost)
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}
	defer ch.Close()

	err = ch.ExchangeDeclare(exchangeName, "topic", true, false, false, false, nil)
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

	msgs, err := ch.Consume(queueName, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("Failed to register a consumer: %v", err)
	}

	log.Printf("Waiting for messages in queue '%s'...", queueName)

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		for msg := range msgs {
			go handleMessage(msg)
		}
	}()

	<-stop
	log.Println("Shutting down gracefully...")
}
