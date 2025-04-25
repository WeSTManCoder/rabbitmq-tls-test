package main

import (
	"crypto/tls"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func main() {
	// TLS-конфигурация с отключённой проверкой
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true, // ⚠️ Не использовать в проде!
	}

	// URL подключения — логин, пароль, хост, порт 5671 (TLS)
	amqpURL := "amqps://USERNAME:PASSWORD@IP:PORT/VHOST"

	// Подключаемся с TLS
	conn, err := amqp.DialTLS(amqpURL, tlsConfig)
	if err != nil {
		log.Fatalf("Не удалось подключиться к RabbitMQ: %v", err)
	}
	defer conn.Close()
	fmt.Println("✅ Подключение к RabbitMQ установлено (TLS)")

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("Не удалось открыть канал: %v", err)
	}
	defer ch.Close()

	queue := "test_queue"

	msgs, err := ch.Consume(
		queue, // имя очереди
		"",    // consumer tag
		true,  // auto-ack
		false, // exclusive
		false, // no-local
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		log.Fatalf("Ошибка при подписке на очередь: %v", err)
	}

	fmt.Println("⏳ Ожидание сообщений. Для выхода — Ctrl+C")
	for msg := range msgs {
		fmt.Printf("📨 Получено сообщение: %s\n", msg.Body)
	}
}
