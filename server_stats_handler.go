// Copyright 2016 Michal Witkowski. All Rights Reserved.
// See LICENSE for licensing terms.

// Forked originally form https://github.com/grpc-ecosystem/go-grpc-prometheus/
// the very same thing with https://github.com/grpc-ecosystem/go-grpc-prometheus/pull/88 integrated
// for the additional functionality to monitore bytes received and send from clients or servers
// everything in this file is only from the PR-88

package grpc_prometheus

import (
	"context"

	"google.golang.org/grpc/stats"
)

type ServerByteStatsHandler struct {
	serverByteMetrics *ServerByteMetrics
}

// TagRPC implements the stats.Hanlder interface.
func (h *ServerByteStatsHandler) TagRPC(ctx context.Context, info *stats.RPCTagInfo) context.Context {
	rpcInfo := newRPCInfo(info.FullMethodName)
	return context.WithValue(ctx, &rpcInfoKey, rpcInfo)
}

// HandleRPC implements the stats.Hanlder interface.
func (h *ServerByteStatsHandler) HandleRPC(ctx context.Context, s stats.RPCStats) {
	v, ok := ctx.Value(&rpcInfoKey).(*rpcInfo)
	if !ok {
		return
	}
	monitor := NewServerByteReporter(h.serverByteMetrics, v.fullMethodName)
	switch s := s.(type) {
	case *stats.InPayload:
		monitor.ReceivedMessageSize(Payload, float64(len(s.Data)))
	case *stats.OutPayload:
		monitor.SentMessageSize(Payload, float64(len(s.Data)))
	}
}

// TagConn implements the stats.Hanlder interface.
func (h *ServerByteStatsHandler) TagConn(ctx context.Context, info *stats.ConnTagInfo) context.Context {
	return ctx
}

// HandleConn implements the stats.Hanlder interface.
func (h *ServerByteStatsHandler) HandleConn(ctx context.Context, s stats.ConnStats) {
}
