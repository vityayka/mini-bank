package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

type Metadata struct {
	ClientIP  string
	UserAgent string
}

const (
	userAgentHeader     = "grpcgateway-user-agent"
	userAgentHeaderGRPC = "user-agent"
	clientIPHeader      = "x-forwarded-for"
)

func (server *Server) extractMedadata(ctx context.Context) *Metadata {
	mtdt := &Metadata{}
	if meta, ok := metadata.FromIncomingContext(ctx); ok {
		if clientIPs := meta.Get(clientIPHeader); len(clientIPs) > 0 {
			mtdt.ClientIP = clientIPs[0]
		}
		if userAgents := meta.Get(userAgentHeader); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		} else if gRPCuserAgents := meta.Get(userAgentHeaderGRPC); len(gRPCuserAgents) > 0 {
			mtdt.UserAgent = gRPCuserAgents[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIP = p.Addr.String()
	}

	return mtdt
}
