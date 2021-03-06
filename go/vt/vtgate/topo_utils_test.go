// Copyright 2012, Google Inc. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package vtgate

import (
	"reflect"
	"testing"

	"github.com/youtube/vitess/go/vt/key"
	"github.com/youtube/vitess/go/vt/topo"
	"golang.org/x/net/context"
)

func TestKeyRangeToShardMap(t *testing.T) {
	ts := new(sandboxTopo)
	var testCases = []struct {
		keyspace string
		keyRange string
		shards   []string
	}{
		{keyspace: KsTestSharded, keyRange: "20-40", shards: []string{"20-40"}},
		// check for partial keyrange, spanning one shard
		{keyspace: KsTestSharded, keyRange: "10-18", shards: []string{"-20"}},
		// check for keyrange intersecting with multiple shards
		{keyspace: KsTestSharded, keyRange: "10-40", shards: []string{"-20", "20-40"}},
		// check for keyrange intersecting with multiple shards
		{keyspace: KsTestSharded, keyRange: "1c-2a", shards: []string{"-20", "20-40"}},
		// check for keyrange where kr.End is Max Key ""
		{keyspace: KsTestSharded, keyRange: "80-", shards: []string{"80-a0", "a0-c0", "c0-e0", "e0-"}},
		// test for sharded, non-partial keyrange spanning the entire space.
		{keyspace: KsTestSharded, keyRange: "", shards: []string{"-20", "20-40", "40-60", "60-80", "80-a0", "a0-c0", "c0-e0", "e0-"}},
		// test for unsharded, non-partial keyrange spanning the entire space.
		{keyspace: KsTestUnsharded, keyRange: "", shards: []string{"0"}},
	}

	for _, testCase := range testCases {
		var keyRange key.KeyRange
		var err error
		if testCase.keyRange == "" {
			keyRange = key.KeyRange{Start: "", End: ""}
		} else {
			krArray, err := key.ParseShardingSpec(testCase.keyRange)
			if err != nil {
				t.Errorf("Got error while parsing sharding spec %v", err)
			}
			keyRange = krArray[0]
		}
		_, allShards, err := getKeyspaceShards(context.Background(), ts, "", testCase.keyspace, topo.TYPE_MASTER)
		gotShards, err := resolveKeyRangeToShards(allShards, keyRange)
		if err != nil {
			t.Errorf("want nil, got %v", err)
		}
		if !reflect.DeepEqual(testCase.shards, gotShards) {
			t.Errorf("want \n%#v, got \n%#v", testCase.shards, gotShards)
		}
	}
}

func TestMapExactShards(t *testing.T) {
	ts := new(sandboxTopo)
	var testCases = []struct {
		keyspace string
		keyRange string
		shards   []string
		err      string
	}{
		{keyspace: KsTestSharded, keyRange: "20-40", shards: []string{"20-40"}},
		// check for partial keyrange, spanning one shard
		{keyspace: KsTestSharded, keyRange: "10-18", shards: nil, err: "keyrange {Start: 10, End: 18} does not exactly match shards"},
		// check for keyrange intersecting with multiple shards
		{keyspace: KsTestSharded, keyRange: "10-40", shards: nil, err: "keyrange {Start: 10, End: 40} does not exactly match shards"},
		// check for keyrange intersecting with multiple shards
		{keyspace: KsTestSharded, keyRange: "1c-2a", shards: nil, err: "keyrange {Start: 1c, End: 2a} does not exactly match shards"},
		// check for keyrange where kr.End is Max Key ""
		{keyspace: KsTestSharded, keyRange: "80-", shards: []string{"80-a0", "a0-c0", "c0-e0", "e0-"}},
		// test for sharded, non-partial keyrange spanning the entire space.
		{keyspace: KsTestSharded, keyRange: "", shards: []string{"-20", "20-40", "40-60", "60-80", "80-a0", "a0-c0", "c0-e0", "e0-"}},
	}

	for _, testCase := range testCases {
		var keyRange key.KeyRange
		var err error
		if testCase.keyRange == "" {
			keyRange = key.KeyRange{Start: "", End: ""}
		} else {
			krArray, err := key.ParseShardingSpec(testCase.keyRange)
			if err != nil {
				t.Errorf("Got error while parsing sharding spec %v", err)
			}
			keyRange = krArray[0]
		}
		_, gotShards, err := mapExactShards(context.Background(), ts, "", testCase.keyspace, topo.TYPE_MASTER, keyRange)
		if err != nil && err.Error() != testCase.err {
			t.Errorf("gotShards: %v, want %s", err, testCase.err)
		}
		if !reflect.DeepEqual(testCase.shards, gotShards) {
			t.Errorf("want \n%#v, got \n%#v", testCase.shards, gotShards)
		}
	}
}
