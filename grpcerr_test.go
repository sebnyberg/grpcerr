package grpcerr_test

import (
	"errors"
	"fmt"
	"testing"

	"github.com/sebnyberg/grpcerr"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestError(t *testing.T) {
	t.Run("GRPC code is omitted from err.Error()", func(t *testing.T) {
		errStr := "place id is invalid"
		var err error = grpcerr.New(codes.InvalidArgument, errors.New(errStr))
		require.Equal(t, errStr, err.Error())
	})

	t.Run("GRPC code can be parsed via status.FromError", func(t *testing.T) {
		errStr := "place id is invalid"
		code := codes.InvalidArgument
		var err error = grpcerr.New(code, errors.New(errStr))
		parsed, ok := status.FromError(err)
		require.True(t, ok)
		require.Equal(t, code, parsed.Code())
	})

	t.Run("wrapping mimics fmt.Errorf and retains code", func(t *testing.T) {
		errStr := "place id is invalid"
		code := codes.InvalidArgument
		var err error = grpcerr.New(code, errors.New(errStr))
		wrappedErr := grpcerr.Errorf("failed to parse place name, %w", err)
		regularWrap := fmt.Errorf("failed to parse place name, %w", err)
		parsed, ok := status.FromError(wrappedErr)
		require.True(t, ok)
		require.Equal(t, code, parsed.Code())
		require.Equal(t, regularWrap.Error(), wrappedErr.Error())
		_, regularOK := status.FromError(regularWrap)
		require.False(t, regularOK)
	})

	t.Run("errors.Is functions as expected", func(t *testing.T) {
		myErr := grpcerr.New(codes.InvalidArgument, errors.New("place id is invalid"))
		wrapped := fmt.Errorf("failed to parse place name, %w", myErr)
		require.ErrorIs(t, wrapped, myErr)
	})
}
