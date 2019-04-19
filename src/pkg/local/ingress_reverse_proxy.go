package local

import (
	"context"
	"log"
	"time"

	diodes "code.cloudfoundry.org/go-diodes"
	"github.com/cloudfoundry/metric-store-release/src/pkg/persistence/transform"
	rpc "github.com/cloudfoundry/metric-store-release/src/pkg/rpc/metricstore_v1"
	"google.golang.org/grpc"
)

// IngressReverseProxy is a reverse proxy for Ingress requests.
type IngressReverseProxy struct {
	localClient rpc.IngressClient
	log         *log.Logger
	buffer      *diodes.OneToOne
}

// NewIngressReverseProxy returns a new IngressReverseProxy.
func NewIngressReverseProxy(
	localClient rpc.IngressClient,
	log *log.Logger,
) *IngressReverseProxy {

	p := &IngressReverseProxy{
		localClient: localClient,
		log:         log,
	}

	p.buffer = diodes.NewOneToOne(10000, diodes.AlertFunc(func(missed int) {
		p.log.Printf("Ingress buffer dropped %d points", missed)
	}))

	go p.WritePoints()

	return p
}

func (p *IngressReverseProxy) Send(ctx context.Context, r *rpc.SendRequest) (*rpc.SendResponse, error) {
	p.buffer.Set(diodes.GenericDataType(r.Batch))
	return &rpc.SendResponse{}, nil
}

func (p *IngressReverseProxy) WritePoints() {
	var points []*rpc.Point
	poller := diodes.NewPoller(p.buffer)

	for {
		data, found := poller.TryNext()

		if !found {
			time.Sleep(time.Millisecond)
			continue
		}

		batch := (*rpc.Points)(data)

		for _, point := range batch.Points {
			point.Name = transform.SanitizeMetricName(point.GetName())

			sanitizedLabels := make(map[string]string)
			for label, value := range point.GetLabels() {
				sanitizedLabels[transform.SanitizeLabelName(label)] = value
			}
			if len(sanitizedLabels) > 0 {
				point.Labels = sanitizedLabels
			}

			points = append(points, point)
		}

		p.localClient.Send(context.Background(), &rpc.SendRequest{
			Batch: &rpc.Points{
				Points: points,
			},
		})

		// reset our slice to save on allocations
		points = points[:0]
	}
}

// IngressClientFunc transforms a function into an IngressClient.
type IngressClientFunc func(ctx context.Context, r *rpc.SendRequest, opts ...grpc.CallOption) (*rpc.SendResponse, error)

// Send implements an IngressClient.
func (f IngressClientFunc) Send(ctx context.Context, r *rpc.SendRequest, opts ...grpc.CallOption) (*rpc.SendResponse, error) {
	return f(ctx, r, opts...)
}
