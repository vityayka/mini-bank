package gapi

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/rs/zerolog/log"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func GRPCLogger(
	ctx context.Context,
	req any,
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp any, err error) {
	startTime := time.Now()

	res, err := handler(ctx, req)

	logger := log.Info()
	if err != nil {
		logger = log.Err(err)
	}

	logger.Str("method", info.FullMethod).
		Dur("duration", time.Since(startTime))

	statusCode := codes.Unknown
	if st, ok := status.FromError(err); ok {
		statusCode = st.Code()
	}

	logger.Int("status_code", int(statusCode)).
		Str("status_text", statusCode.String())

	if message, isOk := req.(fmt.Stringer); isOk {
		logger.Str("message", message.String())
	}

	logger.Msg("incoming gRPC request")

	return res, err
}

type ResponseRecorder struct {
	http.ResponseWriter
	statusCode int
	body       []byte
}

func (rec *ResponseRecorder) WriteHeader(code int) {
	rec.ResponseWriter.WriteHeader(code)
	rec.statusCode = code
}

func (rec *ResponseRecorder) Write(bytes []byte) (int, error) {
	rec.body = bytes
	return rec.ResponseWriter.Write(bytes)
}

func HTTPLogger(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		startTime := time.Now()
		recorder := &ResponseRecorder{ResponseWriter: res, statusCode: http.StatusOK}
		handler.ServeHTTP(recorder, req)

		logger := log.Info()
		if recorder.statusCode != http.StatusOK {
			logger = log.Error().Bytes("body", recorder.body)
		}

		logger.Str("protocol", "http").
			Str("method", req.Method).
			Str("path", req.RequestURI).
			Int("status_code", recorder.statusCode).
			Dur("duration", time.Since(startTime)).
			Msg("incoming http request")

	})
}
