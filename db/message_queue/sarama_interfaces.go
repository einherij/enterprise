// For generate mocks on sarama interfaces
package message_queue

import "github.com/IBM/sarama"

type (
	SaramaAsyncProducer interface {
		sarama.AsyncProducer
	}
)
