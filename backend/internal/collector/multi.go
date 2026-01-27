package collector

import (
	"context"
	"log"
	"mqtt-catalog/internal/config"
	"mqtt-catalog/pkg/dbclient"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

type MultiCollector struct {
	brokers  []config.BrokerConfig
	dbClient *dbclient.Client
	ctx      context.Context
	cancel   context.CancelFunc
}

func NewMultiCollector(brokers []config.BrokerConfig, dbServiceURL string) *MultiCollector {
	ctx, cancel := context.WithCancel(context.Background())

	return &MultiCollector{
		brokers:  brokers,
		dbClient: dbclient.New(dbServiceURL),
		ctx:      ctx,
		cancel:   cancel,
	}
}

func (mc *MultiCollector) Run(duration time.Duration) error {
	var wg sync.WaitGroup

	log.Printf("Starting collection from %d brokers...", len(mc.brokers))

	for _, broker := range mc.brokers {
		wg.Add(1)
		bc := NewBrokerCollector(
			broker.ID,
			broker.URL,
			broker.ClientID,
			broker.Username,
			broker.Password,
			mc.dbClient,
			mc.ctx,
			&wg,
		)
		go func(collector *BrokerCollector) {
			if err := collector.Run(duration); err != nil {
				log.Printf("Error in collector: %v", err)
			}
		}(bc)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-timer.C:
		log.Printf("Collection timer expired")
	case <-sigChan:
		log.Printf("Received interrupt signal")
		mc.cancel()
	}

	log.Printf("Waiting for all broker collectors to finish...")
	wg.Wait()

	log.Printf("All collectors finished")
	return nil
}
