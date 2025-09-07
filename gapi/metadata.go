package gapi

import (
	"context"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
)

const (
	grpcgatewayUserAgentHeader = "grpcgateway-user-agent"
	userAgent                  = "user-agent"
	xForwardedFor              = "x-forwarded-for"
)

type Metadata struct {
	UserAgent string
	ClientIP  string
}

func extractMetadataFromContext(ctx context.Context) *Metadata {
	mtdt := &Metadata{}
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		if userAgents := md.Get(grpcgatewayUserAgentHeader); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		if userAgents := md.Get(userAgent); len(userAgents) > 0 {
			mtdt.UserAgent = userAgents[0]
		}

		if clientIPs := md.Get(xForwardedFor); len(clientIPs) > 0 {
			mtdt.ClientIP = clientIPs[0]
		}
	}

	if p, ok := peer.FromContext(ctx); ok {
		mtdt.ClientIP = p.Addr.String()
	}

	return mtdt
}
