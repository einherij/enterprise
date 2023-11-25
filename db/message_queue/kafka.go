package message_queue

import (
	"fmt"
	"math"
	"strings"
	"time"

	"github.com/IBM/sarama"
)

type KafkaConfig struct {
	Brokers []string `mapstructure:"brokers"`
}

const (
	maxRetries      = 10
	returnSuccess   = true
	returnErrors    = true
	producerTimeout = 10 * time.Second
)

var (
	version     = sarama.V2_1_0_0
	backoffFunc = expRetry(time.Millisecond * 100)
)

// ExpRetry returns function which calculates time between retries.
// The period changes exponentially, starting with initialTime
func expRetry(initialTime time.Duration) func(int, int) time.Duration {
	return func(retries, maxRetries int) time.Duration {
		if retries > maxRetries {
			panic("retries must be not greater than maxRetries")
		}
		t := math.Pow(2, float64(retries)) * float64(initialTime)
		return time.Duration(t)
	}
}

// NewKafkaClient return new sarama.Client interface
func NewKafkaClient(cfg KafkaConfig) (sarama.Client, error) {
	c := sarama.NewConfig()
	c.Producer.Retry.Max = maxRetries
	c.Producer.Return.Successes = returnSuccess
	c.Producer.Return.Errors = returnErrors
	c.Producer.Timeout = producerTimeout
	c.Version = version
	c.Producer.Retry.BackoffFunc = backoffFunc
	// TODO Also async producer use golang channels which default length if 256 when publish message
	// So maybe we can increase c.ChannelBufferSize from config?
	//c.ChannelBufferSize = cfg.KafkaChannelBufferSize
	client, err := sarama.NewClient(cfg.Brokers, c)
	if err != nil {
		return nil, fmt.Errorf("[NewKafkaClient] Brokers [%v]. Error: %v", strings.Join(cfg.Brokers, ", "), err)
	}
	return client, nil
}
