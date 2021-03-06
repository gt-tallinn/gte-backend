// Copyright (C) MongoDB, Inc. 2017-present.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at http://www.apache.org/licenses/LICENSE-2.0

package integration

import (
	"context"
	"testing"

	"github.com/mongodb/mongo-go-driver/bson"
	"github.com/mongodb/mongo-go-driver/core/addr"
	"github.com/mongodb/mongo-go-driver/core/command"
	"github.com/mongodb/mongo-go-driver/core/topology"
	"github.com/stretchr/testify/require"
)

func TestCommand(t *testing.T) {
	noerr := func(t *testing.T, err error) {
		// t.Helper()
		if err != nil {
			t.Errorf("Unepexted error: %v", err)
			t.FailNow()
		}
	}
	t.Parallel()

	server, err := topology.ConnectServer(context.Background(), addr.Addr(*host), serveropts(t)...)
	noerr(t, err)

	ctx := context.Background()
	var result *bson.Document

	cmd := &command.Command{
		DB:      "admin",
		Command: bson.NewDocument(bson.EC.Int32("getnonce", 1)),
	}
	rw, err := server.Connection(ctx)
	noerr(t, err)

	rdr, err := cmd.RoundTrip(ctx, server.SelectedDescription(), rw)
	noerr(t, err)

	result, err = bson.ReadDocument(rdr)
	noerr(t, err)

	elem, err := result.Lookup("ok")
	noerr(t, err)
	if got, want := elem.Value().Type(), bson.TypeDouble; got != want {
		t.Errorf("Did not get correct type for 'ok'. got %s; want %s", got, want)
	}
	if got, want := elem.Value().Double(), float64(1); got != want {
		t.Errorf("Did not get correct value for 'ok'. got %f; want %f", got, want)
	}

	elem, err = result.Lookup("nonce")
	require.NoError(t, err)
	require.Equal(t, elem.Value().Type(), bson.TypeString)
	require.NotEqual(t, "", elem.Value().StringValue(), "MongoDB returned empty nonce")

	result.Reset()
	cmd.Command = bson.NewDocument(bson.EC.Int32("ping", 1))

	rdr, err = cmd.RoundTrip(ctx, server.SelectedDescription(), rw)
	noerr(t, err)

	result, err = bson.ReadDocument(rdr)
	require.NoError(t, err)

	elem, err = result.Lookup("ok")
	require.NoError(t, err)
	require.Equal(t, elem.Value().Type(), bson.TypeDouble)
	require.Equal(t, float64(1), elem.Value().Double(), "Unable to ping MongoDB")
}
