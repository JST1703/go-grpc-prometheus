// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

// Forked originally form https://github.com/grpc-ecosystem/go-grpc-prometheus/
// the very same thing with https://github.com/grpc-ecosystem/go-grpc-prometheus/pull/88 integrated
// for the additional functionality to monitore bytes received and send from clients or servers
// eveything that is in between a "---- PR-88 ---- {"  and   "---- PR-88 ---- }" comment is the new addition from the PR88.

package grpc_prometheus

import (
	prom "github.com/prometheus/client_golang/prometheus"

	"google.golang.org/grpc/stats" // PR-88
)

// ServerMetrics represents a collection of metrics to be registered on a
// Prometheus metrics registry for a gRPC server.
type ServerByteMetrics struct {
	serverMsgSizeBytesSent     *prom.CounterVec
	serverMsgSizeBytesReceived *prom.CounterVec

	serverMsgSizeReceivedHistogramEnabled bool
	serverMsgSizeReceivedHistogramOpts    prom.HistogramOpts
	serverMsgSizeReceivedHistogram        *prom.HistogramVec

	serverMsgSizeSentHistogramEnabled bool
	serverMsgSizeSentHistogramOpts    prom.HistogramOpts
	serverMsgSizeSentHistogram        *prom.HistogramVec
}

// NewServerMetrics returns a ServerMetrics object. Use a new instance of
// ServerMetrics when not using the default Prometheus metrics registry, for
// example when wanting to control which metrics are added to a registry as
// opposed to automatically adding metrics via init functions.
func NewServerByteMetrics() *ServerByteMetrics {
	return &ServerByteMetrics{
		serverMsgSizeBytesSent: prom.NewCounterVec(
			prom.CounterOpts{
				Name: "serverMsgSizeBytesSent",
				Help: "Total number of bytes sent by server.",
			}, []string{"grpc_service", "grpc_method", "grpc_stats"}),

		serverMsgSizeBytesReceived: prom.NewCounterVec(
			prom.CounterOpts{
				Name: "serverMsgSizeBytesReceived",
				Help: "Total number of bytes received by server.",
			}, []string{"grpc_service", "grpc_method", "grpc_stats"}),

		serverMsgSizeReceivedHistogramEnabled: false,
		serverMsgSizeReceivedHistogramOpts: prom.HistogramOpts{
			Name:    "grpc_server_msg_size_received_bytes",
			Help:    "Histogram of message sizes received by the server.",
			Buckets: defMsgBytesBuckets,
		},
		serverMsgSizeReceivedHistogram: nil,

		serverMsgSizeSentHistogramEnabled: false,
		serverMsgSizeSentHistogramOpts: prom.HistogramOpts{
			Name:    "grpc_server_msg_size_sent_bytes",
			Help:    "Histogram of message sizes sent by the server.",
			Buckets: defMsgBytesBuckets,
		},
		serverMsgSizeSentHistogram: nil,
	}
}

// EnableMsgSizeReceivedBytesHistogram turns on recording of received message size of RPCs.
// Histogram metrics can be very expensive for Prometheus to retain and query. It takes
// options to configure histogram options such as the defined buckets.
func (m *ServerByteMetrics) EnableMsgSizeReceivedBytesHistogram() {
	if !m.serverMsgSizeReceivedHistogramEnabled {
		m.serverMsgSizeReceivedHistogram = prom.NewHistogramVec(
			m.serverMsgSizeReceivedHistogramOpts,
			[]string{"grpc_service", "grpc_method", "grpc_stats"},
		)
	}
	m.serverMsgSizeReceivedHistogramEnabled = true
}

// EnableMsgSizeSentBytesHistogram turns on recording of sent message size of RPCs.
// Histogram metrics can be very expensive for Prometheus to retain and query. It takes
// options to configure histogram options such as the defined buckets.
func (m *ServerByteMetrics) EnableMsgSizeSentBytesHistogram() {
	if !m.serverMsgSizeSentHistogramEnabled {
		m.serverMsgSizeSentHistogram = prom.NewHistogramVec(
			m.serverMsgSizeSentHistogramOpts,
			[]string{"grpc_service", "grpc_method", "grpc_stats"},
		)
	}
	m.serverMsgSizeSentHistogramEnabled = true
}

// Describe sends the super-set of all possible descriptors of metrics
// collected by this Collector to the provided channel and returns once
// the last descriptor has been sent.
func (m *ServerByteMetrics) Describe(ch chan<- *prom.Desc) {
	m.serverMsgSizeBytesSent.Describe(ch)
	m.serverMsgSizeBytesReceived.Describe(ch)
	if m.serverMsgSizeReceivedHistogramEnabled {
		m.serverMsgSizeReceivedHistogram.Describe(ch)
	}
	if m.serverMsgSizeSentHistogramEnabled {
		m.serverMsgSizeSentHistogram.Describe(ch)
	}
}

// Collect is called by the Prometheus registry when collecting
// metrics. The implementation sends each collected metric via the
// provided channel and returns once the last metric has been sent.
func (m *ServerByteMetrics) Collect(ch chan<- prom.Metric) {
	m.serverMsgSizeBytesSent.Collect(ch)
	m.serverMsgSizeBytesReceived.Collect(ch)
	if m.serverMsgSizeReceivedHistogramEnabled {
		m.serverMsgSizeReceivedHistogram.Collect(ch)
	}
	if m.serverMsgSizeSentHistogramEnabled {
		m.serverMsgSizeSentHistogram.Collect(ch)
	}
}

// NewServerStatsHandler is a gRPC server-side stats.Handler that providers Prometheus monitoring for RPCs.
func (m *ServerByteMetrics) NewServerByteStatsHandler() stats.Handler {
	return &ServerByteStatsHandler{
		serverByteMetrics: m,
	}
}
