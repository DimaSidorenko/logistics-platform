package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"

	"github.com/IBM/sarama"

	"route256/notifier/internal/infra/config"
	"route256/notifier/internal/infra/kafka/consumer_group"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	cfg, err := config.ReadConfig()
	if err != nil {
		log.Fatalf("read config: %v", err)
		return
	}

	handler := consumer_group.NewConsumerGroupHandler()
	cg, err := consumer_group.NewConsumerGroup(
		[]string{cfg.Kafka.Brokers},
		cfg.Kafka.ConsumerGroupId,
		[]string{cfg.Kafka.OrderTopic},
		handler,
		consumer_group.WithOffsetsInitial(sarama.OffsetOldest),
	)
	if err != nil {
		log.Fatal(err)
	}
	defer cg.Close()

	var wg sync.WaitGroup

	runCGErrorHandler(ctx, cg, &wg)

	cg.Run(ctx, &wg)

	wg.Wait()
	log.Println("graceful shutdown complete")
}

func runCGErrorHandler(ctx context.Context, cg sarama.ConsumerGroup, wg *sync.WaitGroup) {
	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case chErr, ok := <-cg.Errors():
				if !ok {
					fmt.Println("[cg-error] error: chan closed")
					return
				}

				fmt.Printf("[cg-error] error: %s\n", chErr)
			case <-ctx.Done():
				fmt.Printf("[cg-error] ctx closed: %s\n", ctx.Err().Error())
				return
			}
		}
	}()
}
