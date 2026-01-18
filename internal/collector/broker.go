package collector

import (
	"context"
	"fmt"
	"log"
	"mqtt-collector/internal/models"
	"mqtt-collector/internal/payload"
	"mqtt-collector/pkg/dbclient"
	"sync"
	"time"

	mqtt "github.com/eclipse/paho.mqtt.golang"
)

type BrokerCollector struct {
	brokerID      string
	brokerURL     string
	mqttClient    mqtt.Client
	dbClient      *dbclient.Client
	sampledTopics map[string]bool
	mu            sync.Mutex
	ctx           context.Context
	wg            *sync.WaitGroup
}

func NewBrokerCollector(brokerID, brokerURL string, dbClient *dbclient.Client, ctx context.Context, wg *sync.WaitGroup) *BrokerCollector {
	clientID := fmt.Sprintf("mqtt-collector-%s-%d", brokerID, time.Now().Unix())

	opts := mqtt.NewClientOptions().
		AddBroker(brokerURL).
		SetClientID(clientID).
		SetCleanSession(true).
		SetAutoReconnect(false).
		SetConnectTimeout(10 * time.Second).
		SetConnectionLostHandler(func(client mqtt.Client, err error) {
			log.Printf("[%s] Connection lost: %v", brokerID, err)
		})

	client := mqtt.NewClient(opts)

	bc := &BrokerCollector{
		brokerID:      brokerID,
		brokerURL:     brokerURL,
		mqttClient:    client,
		dbClient:      dbClient,
		sampledTopics: make(map[string]bool),
		ctx:           ctx,
		wg:            wg,
	}

	opts.SetDefaultPublishHandler(bc.messageHandler)

	return bc
}

func (bc *BrokerCollector) messageHandler(client mqtt.Client, msg mqtt.Message) {
	topic := msg.Topic()

	bc.mu.Lock()
	if bc.sampledTopics[topic] {
		bc.mu.Unlock()
		return
	}
	bc.sampledTopics[topic] = true
	bc.mu.Unlock()

	payloadData := msg.Payload()
	payloadType := payload.DetectType(payloadData)

	sample := models.Sample{
		BrokerID:    bc.brokerID,
		Topic:       topic,
		PayloadType: payloadType,
		Payload:     payloadData,
		Timestamp:   time.Now(),
	}

	if err := bc.dbClient.SendSample(bc.ctx, sample); err != nil {
		log.Printf("[%s] Error sending sample for topic %s: %v", bc.brokerID, topic, err)
	} else {
		log.Printf("[%s] Sampled topic: %s (type: %s, size: %d bytes)", bc.brokerID, topic, payloadType, len(payloadData))
	}
}

func (bc *BrokerCollector) Run(duration time.Duration) error {
	defer bc.wg.Done()

	log.Printf("[%s] Connecting to MQTT broker at %s...", bc.brokerID, bc.brokerURL)
	if token := bc.mqttClient.Connect(); token.Wait() && token.Error() != nil {
		return fmt.Errorf("[%s] connect error: %w", bc.brokerID, token.Error())
	}
	defer bc.mqttClient.Disconnect(250)

	log.Printf("[%s] Connected. Subscribing to all topics (#)...", bc.brokerID)
	if token := bc.mqttClient.Subscribe("#", 0, nil); token.Wait() && token.Error() != nil {
		return fmt.Errorf("[%s] subscribe error: %w", bc.brokerID, token.Error())
	}

	log.Printf("[%s] Collecting samples for %v...", bc.brokerID, duration)

	timer := time.NewTimer(duration)
	defer timer.Stop()

	select {
	case <-timer.C:
		log.Printf("[%s] Collection period completed", bc.brokerID)
	case <-bc.ctx.Done():
		log.Printf("[%s] Context canceled", bc.brokerID)
	}

	bc.mu.Lock()
	count := len(bc.sampledTopics)
	bc.mu.Unlock()

	log.Printf("[%s] Collection finished. Sampled %d unique topics", bc.brokerID, count)

	return nil
}
