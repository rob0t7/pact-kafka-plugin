package consumertest

import (
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/hamba/avro/v2"
	pactlog "github.com/pact-foundation/pact-go/v2/log"
	message "github.com/pact-foundation/pact-go/v2/message/v4"
)

const (
	AVRO_SCHEMA_CONTENT_TYPE = "application/vnd.kafka.avro.v2"
	PLUGIN_NAME              = "kafkaplugin"
)

var PLUGIN_VERSION string = "0.0.1"

func init() {
	if version, ok := os.LookupEnv("VERSION"); ok {
		PLUGIN_VERSION = version
	}
}

func TestConsumer(t *testing.T) {
	schemaID := int64(16)
	magicByte := 0

	// nolint:errcheck
	pactlog.SetLogLevel("DEBUG")
	provider, err := message.NewAsynchronousPact(message.Config{
		Consumer: "Consumer",
		Provider: "Provider",
	})
	if err != nil {
		t.Fatalf("Error creating Pact Provider: %v", err)
	}

	schema, err := avro.ParseFiles("./schema.avsc")
	if err != nil {
		t.Fatalf("Error parsing Avro schema: %v", err)
	}
	expected := User{
		ID:        "94af717d-1b04-4fad-9879-01dc828e410d",
		FirstName: "Jane",
		LastName:  "Doe",
		Email:     "jane.doe@example.com",
	}
	data, err := avro.Marshal(schema, expected)
	if err != nil {
		t.Fatalf("Error marshalling Avro message: %v", err)
	}

	err = provider.
		AddAsynchronousMessage().
		ExpectsToReceive("a schema registery enabled Kafka message encoded with AVRO").
		UsingPlugin(message.PluginConfig{
			Plugin:  PLUGIN_NAME,
			Version: PLUGIN_VERSION,
		}).
		WithContents(
			fmt.Sprintf(`{"schemaId": %d, "message": "%s"}`, schemaID, base64.StdEncoding.EncodeToString(data)),
			AVRO_SCHEMA_CONTENT_TYPE,
		).
		ExecuteTest(t, func(m message.AsynchronousMessage) error {
			return verifyKafkaMessage(t, m, schema, expected, schemaID, magicByte)
		})
	if err != nil {
		t.Fatalf("Error during message test: %v", err)
	}
}

// verifyKafkaMessage checks the contents of the Kafka message.
func verifyKafkaMessage(t *testing.T, m message.AsynchronousMessage, schema avro.Schema, expected User, schemaID int64, magicByte int) error {
	t.Helper()
	if l := len(m.Contents); l < 5 {
		t.Fatalf("Message contents too short, expected at least 5 bytes, got %d", l)
	}
	if v, err := strconv.Atoi(string(m.Contents[0])); err != nil || v != magicByte {
		if err != nil {
			t.Fatalf("Error parsing magic byte: %v", err)
		} else {
			t.Errorf("Expected first byte to be a magicByte (0x0), got %d", v)
		}
	}
	if v, err := strconv.ParseInt(string(m.Contents[1:5]), 16, 64); err != nil || v != schemaID {
		if err != nil {
			t.Fatalf("Error parsing schema ID: %v", err)
		} else {
			t.Errorf("Expected bytes 1-4 to be schema ID 16, got %d", v)
		}
	}
	var actual User
	if err := avro.Unmarshal(schema, m.Contents[5:], &actual); err != nil {
		t.Fatalf("Error unmarshalling Avro message: %v", err)
	}
	if diff := cmp.Diff(expected, actual); diff != "" {
		t.Errorf("Message content mismatch (-want +got):\n%s", diff)
	}
	return nil
}
