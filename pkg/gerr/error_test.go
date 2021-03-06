// Copyright 2018 The OpenPitrix Authors. All rights reserved.
// Use of this source code is governed by a Apache license
// that can be found in the LICENSE file.

package gerr

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/status"
)

func TestNew(t *testing.T) {
	var e error
	e = New(Internal, ErrorCreateResourcesFailed)
	assert.Equal(t, e.Error(), "rpc error: code = Internal desc = create resources failed")
	e = New(InvalidArgument, ErrorMissingParameter, "name")
	assert.Equal(t, e.Error(), "rpc error: code = InvalidArgument desc = missing parameter [name]")
}

func TestNewWithDetail(t *testing.T) {
	var e error
	e = NewWithDetail(InvalidArgument, fmt.Errorf("test with error detail"), ErrorCreateResourcesFailed)
	assert.Equal(t, e.Error(), "rpc error: code = InvalidArgument desc = create resources failed")
	ge := status.Convert(e)
	assert.Equal(t, ge.Code().String(), "InvalidArgument")
	assert.Equal(t, ge.Err().Error(), "rpc error: code = InvalidArgument desc = create resources failed")
	assert.Equal(t, fmt.Sprint(ge.Details()), "[error_name:\"create_resources_failed\" cause:\"test with error detail\" ]")
	//t.Log(ge.Code(), ge.Err(), ge.Details())

	e = NewWithDetail(InvalidArgument, errors.New("test with error detail"), ErrorCreateResourcesFailed)
	ge = status.Convert(e)
	assert.Regexp(t, regexp.MustCompile("TestNewWithDetail"), ge.Details())
}

func TestClearErrorCause(t *testing.T) {
	var e error
	e = NewWithDetail(InvalidArgument, fmt.Errorf("test with error detail"), ErrorCreateResourcesFailed)

	ge := status.Convert(e)
	assert.Equal(t, fmt.Sprint(ge.Details()), "[error_name:\"create_resources_failed\" cause:\"test with error detail\" ]")

	e = ClearErrorCause(e)
	ge = status.Convert(e)
	assert.Equal(t, fmt.Sprint(ge.Details()), "[error_name:\"create_resources_failed\" ]")
	assert.True(t, IsGRPCError(e))
	assert.False(t, IsGRPCError(fmt.Errorf("test")))
	assert.False(t, IsGRPCError(func() GRPCError { return nil }()))
}
