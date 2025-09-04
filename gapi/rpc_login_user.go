package gapi

import (
	"context"
	"net"

	db "github.com/hykura1501/simple_bank/db/sqlc"
	"github.com/hykura1501/simple_bank/pb"
	"github.com/hykura1501/simple_bank/util"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// getClientIP extracts the client IP from gRPC context
func getClientIP(ctx context.Context) string {
	if p, ok := peer.FromContext(ctx); ok {
		addr := p.Addr.String()
		// addr thường có dạng "ip:port", nên tách ra chỉ lấy IP
		if host, _, err := net.SplitHostPort(addr); err == nil {
			return host
		}
		return addr
	}
	return ""
}

// getUserAgent extracts the user-agent from gRPC metadata
func getUserAgent(ctx context.Context) string {
	if md, ok := metadata.FromIncomingContext(ctx); ok {
		ua := md.Get("user-agent")
		if len(ua) > 0 {
			return ua[0]
		}
	}
	return ""
}

func (server *Server) LoginUser(ctx context.Context, req *pb.LoginUserRequest) (*pb.LoginUserResponse, error) {
	user, err := server.store.GetUser(ctx, req.GetUsername())
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, status.Errorf(codes.NotFound, "user not found: %s", err)

		}
		return nil, status.Errorf(codes.Internal, "cannot get user: %s", err)
	}

	if err := util.CheckPassword(req.Password, user.HashedPassword); err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid password: %s", err)
	}

	accessToken, accessTokenPayload, err := server.tokenMaker.CreateToken(req.Username, server.config.AccessTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create access token: %s", err)
	}

	refreshToken, refreshTokenPayload, err := server.tokenMaker.CreateToken(req.Username, server.config.RefreshTokenDuration)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create refresh token: %s", err)
	}

	arg := db.CreateSessionParams{
		ID:           refreshTokenPayload.ID,
		Username:     req.GetUsername(),
		RefreshToken: refreshToken,
		UserAgent:    getUserAgent(ctx),
		ClientIp:     getClientIP(ctx),
		IsBlocked:    false,
		ExpiredAt: pgtype.Timestamptz{
			Time:  refreshTokenPayload.ExpiredAt,
			Valid: true,
		},
	}

	session, err := server.store.CreateSession(ctx, arg)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "cannot create session: %s", err)
	}
	loginUserResponse := &pb.LoginUserResponse{
		SessionId:             session.ID.String(),
		AccessToken:           accessToken,
		AccessTokenExpiresAt:  timestamppb.New(accessTokenPayload.ExpiredAt),
		RefreshToken:          refreshToken,
		RefreshTokenExpiresAt: timestamppb.New(refreshTokenPayload.ExpiredAt),
		User:                  convertUser(user),
	}
	return loginUserResponse, nil
}
