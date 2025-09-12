package gapi

import (
	"context"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GrpcLogger(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	startTime := time.Now()
	result, err := handler(ctx, req)
	duration := time.Since(startTime)

	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}

	logger := log.Info()
	if err != nil {
		logger = log.Error().Err(err)
	}

	logger.Str("protocol", "grpc").
		Str("method", info.FullMethod).
		Int("status_code", int(statusCode)).
		Str("status_text", statusCode.String()).
		Dur("duration", duration).
		Msg("received a gRPC request")
	return result, err
}

type ResponseRecorder struct {
	http.ResponseWriter
	StatusCode int
	Body       []byte
}

func (recorder *ResponseRecorder) WriteHeader(statusCode int) {
	recorder.StatusCode = statusCode
	recorder.ResponseWriter.WriteHeader(statusCode)
}

func (recorder *ResponseRecorder) Write(body []byte) (int, error) {
	recorder.Body = body
	return recorder.ResponseWriter.Write(body)
}

func HttpLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		recorder := &ResponseRecorder{
			ResponseWriter: res,
			StatusCode:     http.StatusOK,
		}
		duration := time.Since(startTime)

		handler.ServeHTTP(recorder, req)

		logger := log.Info()
		if recorder.StatusCode != http.StatusOK {
			logger = log.Error().Bytes("body", recorder.Body)
		}

		logger.Str("protocol", "http").
			Str("path", req.RequestURI).
			Str("method", req.Method).
			Int("status_code", recorder.StatusCode).
			Str("status_text", http.StatusText(recorder.StatusCode)).
			Dur("duration", duration).
			Msg("received a HTTP request")
	})
}
