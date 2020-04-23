// Copyright (c) 2019, 2020, Maxime Soulé
// All rights reserved.
//
// This source code is licensed under the BSD-style license found in the
// LICENSE file in the root directory of this source tree.

// Package tdutil allows to write unit tests for go-testdeep helpers
// and so provides some helpful functions.
//
// It is not intended to be used in tests outside go-testdeep and its
// helpers perimeter.
package tdutil

import (
	"reflect"
	"testing"
)

// T can be used in tests, to test testing.T behavior as it overrides
// Run() method.
type T struct {
	testing.T
	name string
}

// NewT returns a new *T instance. "name" is the string returned by
// method Name.
func NewT(name string) *T {
	return &T{name: name}
}

// Run is a simplified version of testing.T.Run() method, without edge
// cases.
func (t *T) Run(name string, f func(*testing.T)) bool {
	f(&t.T)
	return !t.Failed()
}

// Name returns the name of the running test (in fact the one set by NewT).
func (t *T) Name() string {
	return t.name
}

// LogBuf is an ugly hack allowing to access internal testing.T log
// buffer. Keep cool, it is only used for internal unit tests.
func (t *T) LogBuf() string {
	return string(reflect.ValueOf(t.T).FieldByName("output").Bytes()) // nolint: govet
}
