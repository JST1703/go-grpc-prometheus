// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

// Forked originally form https://github.com/grpc-ecosystem/go-grpc-prometheus/
// the very same thing with https://github.com/grpc-ecosystem/go-grpc-prometheus/pull/88 integrated
// for the additional functionality to monitore bytes received and send from clients or servers
// eveything that is in between a "---- PR-88 ---- {"  and   "---- PR-88 ---- }" comment is the new addition from the PR88.

package grpc_prometheus

type ServerByteReporter struct {
	metrics     *ServerByteMetrics
	rpcType     grpcType
	serviceName string
	methodName  string
}

func NewServerByteReporter(m *ServerByteMetrics, fullMethod string) *ServerByteReporter {
	r := &ServerByteReporter{
		metrics: m,
		rpcType: Unary,
	}
	r.serviceName, r.methodName = splitMethodName(fullMethod)
	return r
}

// ReceivedMessageSize counts the size of received messages on server-side
func (r *ServerByteReporter) ReceivedMessageSize(rpcStats grpcStats, size float64) {
	r.metrics.serverMsgSizeBytesReceived.WithLabelValues(r.serviceName, r.methodName, rpcStats.String()).Add(size)
	if r.metrics.serverMsgSizeReceivedHistogramEnabled {
		r.metrics.serverMsgSizeReceivedHistogram.WithLabelValues(r.serviceName, r.methodName, rpcStats.String()).Observe(size)
	}
}

// SentMessageSize counts the size of sent messages on server-side
func (r *ServerByteReporter) SentMessageSize(rpcStats grpcStats, size float64) {
	r.metrics.serverMsgSizeBytesSent.WithLabelValues(r.serviceName, r.methodName, rpcStats.String()).Add(size)
	if r.metrics.serverMsgSizeSentHistogramEnabled {
		r.metrics.serverMsgSizeSentHistogram.WithLabelValues(r.serviceName, r.methodName, rpcStats.String()).Observe(size)
	}
}
