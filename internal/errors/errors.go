package errors

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	CacheKeyNotFoundError           = status.Error(codes.NotFound, "key not found")
	InvalidArgumentError            = status.Error(codes.InvalidArgument, "Invalid argument passed")
	UpstreamExchangeRateServerError = status.Error(codes.Internal, "failed to get the exchange rates from upstream")
	InternalCacheError              = status.Error(codes.Internal, "failed to complete a transaction with cache")
	UnImplementedError              = status.Error(codes.Unimplemented, "method not implemented")
)

// IsNotFound returns true if the error is NotFound error.
func IsNotFound(err error) bool {
	if codes.NotFound == status.Code(err) {
		return true
	}

	return false
}

// IsInvalidArgument returns true if the error is NotFound error.
func IsInvalidArgument(err error) bool {
	if codes.InvalidArgument == status.Code(err) {
		return true
	}

	return false
}

// IsInternalServer returns true if the error is NotFound error.
func IsInternalServer(err error) bool {
	if codes.Internal == status.Code(err) {
		return true
	}

	return false
}

// IsUpstreamServerError returns true if the error passed is UpstreamExchangeRateServer error.
func IsUpstreamServerError(err error) bool {
	if UpstreamExchangeRateServerError == err {
		return true
	}

	return false
}

// IsUnImplementedError returns true if the error passed is UnImplemented error.
func IsUnImplementedError(err error) bool {
	if codes.Unimplemented == status.Code(err) {
		return true
	}

	return false
}
