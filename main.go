package main

import (
	"bufio"
	"crypto/tls"
	"flag"
	"fmt"
	"os"
	"strings"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	queueFlag := flag.String("q", "", "Имя очереди (по умолчанию: test_queue)")
	flag.Parse()

	reader := bufio.NewReader(os.Stdin)

	// Проверка, был ли явно передан флаг -q
	queueWasProvided := false
	for _, arg := range os.Args[1:] {
		if strings.HasPrefix(arg, "-q") || strings.HasPrefix(arg, "--q") {
			queueWasProvided = true
			break
		}
	}

	defaultURL := "user:pass@host:port/vhost"
	if len(flag.Args()) > 0 {
		defaultURL = strings.TrimSpace(flag.Args()[0])
	}

	defaultQueue := strings.TrimSpace(*queueFlag)
	if defaultQueue == "" {
		defaultQueue = "test_queue"
	}

	lastURL := defaultURL
	lastQueue := defaultQueue

	for {
		// === ВВОД URL ===
		if len(flag.Args()) == 0 {
			fmt.Println("\nRabbitMQ TLS тест v1.0.0 | Автор: ChatGPT | Идея: WeSTManCoder")
			fmt.Printf("🔐 Введите адрес подключения [по умолчанию: %s]: ", lastURL)
			inputURL, _ := reader.ReadString('\n')
			inputURL = strings.TrimSpace(inputURL)
			if inputURL != "" {
				lastURL = inputURL
			}
			if lastURL == "" {
				fmt.Println("❌ Адрес подключения не указан. Повтор.")
				continue
			}
		}

		// Добавляем amqps:// если нужно
		amqpURL := lastURL
		if !strings.HasPrefix(amqpURL, "amqps://") {
			amqpURL = "amqps://" + amqpURL
		}

		// === ВВОД ОЧЕРЕДИ ===
		if !queueWasProvided {
			fmt.Printf("📦 Введите имя очереди [по умолчанию: %s]: ", lastQueue)
			inputQueue, _ := reader.ReadString('\n')
			inputQueue = strings.TrimSpace(inputQueue)
			if inputQueue != "" {
				lastQueue = inputQueue
			}
			if lastQueue == "" {
				lastQueue = "test_queue"
			}
		}

		// Подключение
		tlsConfig := &tls.Config{InsecureSkipVerify: true}
		conn, err := amqp.DialTLS(amqpURL, tlsConfig)
		if err != nil {
			fmt.Printf("❌ Ошибка подключения к RabbitMQ: %v\n", err)
			if len(flag.Args()) > 0 {
				os.Exit(1)
			}
			continue
		}

		ch, err := conn.Channel()
		if err != nil {
			fmt.Printf("❌ Ошибка открытия канала: %v\n", err)
			conn.Close()
			if len(flag.Args()) > 0 {
				os.Exit(1)
			}
			continue
		}

		fmt.Printf("✅ Подключено с TLS к %s\n", amqpURL)

		msgs, err := ch.Consume(
			lastQueue,
			"",
			false,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			fmt.Printf("❌ Ошибка подписки на очередь '%s': %v\n", lastQueue, err)
			ch.Close()
			conn.Close()
			if len(flag.Args()) > 0 {
				os.Exit(1)
			}
			continue
		}

		for msg := range msgs {
			fmt.Printf("📨 %s\n", msg.Body)
			break
		}

		ch.Close()
		conn.Close()
		fmt.Println("✅ Подключение закрыто.")

		if len(flag.Args()) > 0 {
			break
		}
	}
}
