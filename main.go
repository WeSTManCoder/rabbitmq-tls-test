package main

import (
	"crypto/tls"
	"fmt"
	"log"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	// TLS-конфигурация (отключаем проверку сертификата)
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	// AMQPS URI
	amqpURL := "amqps://USERNAME:PASSWORD@IP:PORT/VHOST"

	// Подключение по TLS
	conn, err := amqp.DialTLS(amqpURL, tlsConfig)
	if err != nil {
		log.Fatalf("Ошибка подключения: %v", err)
	}
	defer conn.Close()
	fmt.Println("✅ Успешно подключено к RabbitMQ (TLS)")

	// Открытие канала
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Ошибка открытия канала: %v", err)
	}
	defer ch.Close()

	// Подписка на очередь
	queue := "test_queue"
	msgs, err := ch.Consume(
		queue, // имя очереди
		"",    // consumer tag
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // args
	)
	if err != nil {
		log.Fatalf("Ошибка подписки: %v", err)
	}

	// Чтение сообщений
	fmt.Println("⏳ Ожидание сообщений...")
	for msg := range msgs {
		fmt.Printf("📨 %s\n", msg.Body)
	}
}
