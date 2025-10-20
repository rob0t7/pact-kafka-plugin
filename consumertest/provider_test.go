package consumertest

import (
	"testing"

	"github.com/pact-foundation/pact-go/v2/message"
	"github.com/pact-foundation/pact-go/v2/models"
	"github.com/pact-foundation/pact-go/v2/provider"
)

func TestProvider(t *testing.T) {
	verifier := provider.NewVerifier()

	functionMappings := message.Handlers{
		"a schema registery enabled Kafka message encoded with AVRO": func([]models.ProviderState) (message.Body, message.Metadata, error) {
			return nil, message.Metadata{"contentType": "application/vnd.kafka.avro.v2"}, nil
		},
	}

	if err := verifier.VerifyProvider(t, provider.VerifyRequest{
		Provider:        "Provider",
		PactFiles:       []string{"pacts/Consumer-Provider.json"},
		MessageHandlers: functionMappings,
	}); err != nil {
		t.Error(err)
	}
}
