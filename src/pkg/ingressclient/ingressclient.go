package ingressclient

import (
	"crypto/tls"
	"io/ioutil"
	"log"
	"time"

	"github.com/cloudfoundry/metric-store-release/src/pkg/leanstreams"
	rpc "github.com/cloudfoundry/metric-store-release/src/pkg/rpc/metricstore_v1"
	"github.com/gogo/protobuf/proto"
)

const (
	MAX_BATCH_SIZE_IN_BYTES           = 32 * 1024
	MAX_INGRESS_PAYLOAD_SIZE_IN_BYTES = 2 * MAX_BATCH_SIZE_IN_BYTES
)

type IngressClient struct {
	connection *leanstreams.TCPClient
	log        *log.Logger
}

func New(ingressAddress string, tlsConfig *tls.Config, opts ...IngressClientOption) (*IngressClient, error) {
	clientConfig := &leanstreams.TCPClientConfig{
		MaxMessageSize: MAX_INGRESS_PAYLOAD_SIZE_IN_BYTES,
		Address:        ingressAddress,
		TLSConfig:      tlsConfig,
	}

	connection, err := leanstreams.DialTCPUntilConnected(clientConfig, 5*time.Minute)
	if err != nil {
		return nil, err
	}

	client := &IngressClient{
		connection: connection,
		log:        log.New(ioutil.Discard, "", 0),
	}

	for _, o := range opts {
		o(client)
	}

	return client, nil
}

type IngressClientOption func(*IngressClient)

func WithIngressClientLogger(log *log.Logger) IngressClientOption {
	return func(client *IngressClient) {
		client.log = log
	}
}

func (c *IngressClient) Write(points []*rpc.Point) error {
	payload, err := proto.Marshal(&rpc.SendRequest{
		Batch: &rpc.Points{
			Points: points,
		},
	})

	if err != nil {
		log.Printf("failed to marshal metric points: %s\n", err)
		return err
	}

	// TODO: consider adding back in a timeout (i.e. 3 seconds)
	bytesWritten, err := c.connection.Write(payload)

	if err == nil {
		c.log.Printf("Wrote %d of %d bytes\n", bytesWritten, len(payload))
	}

	return err
}
