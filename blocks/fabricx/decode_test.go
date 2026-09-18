/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package fabricx

import (
	"testing"
)

func TestDecodeMetadata(t *testing.T) {
	tests := []struct {
		name        string
		metadata    [][]byte
		wantEvents  []byte
		wantPayload []byte
		wantArgs    [][]byte
	}{
		{
			name:     "nil metadata",
			metadata: nil,
		},
		{
			name:     "all absent, zero args",
			metadata: [][]byte{nil, nil, {0}},
		},
		{
			name:        "event and payload only, zero args",
			metadata:    [][]byte{[]byte("evt"), []byte("payload"), {0}},
			wantEvents:  []byte("evt"),
			wantPayload: []byte("payload"),
		},
		{
			name:     "one arg",
			metadata: [][]byte{nil, nil, {1}, []byte("a")},
			wantArgs: [][]byte{[]byte("a")},
		},
		{
			name:     "many args",
			metadata: [][]byte{nil, nil, {3}, []byte("a"), []byte("b"), []byte("c")},
			wantArgs: [][]byte{[]byte("a"), []byte("b"), []byte("c")},
		},
		{
			name:        "mixed: event, payload, and args all present",
			metadata:    [][]byte{[]byte("evt"), []byte("payload"), {2}, []byte("x"), []byte("y")},
			wantEvents:  []byte("evt"),
			wantPayload: []byte("payload"),
			wantArgs:    [][]byte{[]byte("x"), []byte("y")},
		},
		{
			name:     "count present but args truncated",
			metadata: [][]byte{nil, nil, {2}, []byte("only-one")},
			wantArgs: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			events, payload, args := DecodeMetadata(tt.metadata)
			if string(events) != string(tt.wantEvents) {
				t.Errorf("events: got %q, want %q", events, tt.wantEvents)
			}
			if string(payload) != string(tt.wantPayload) {
				t.Errorf("payload: got %q, want %q", payload, tt.wantPayload)
			}
			if len(args) != len(tt.wantArgs) {
				t.Fatalf("args len: got %d, want %d", len(args), len(tt.wantArgs))
			}
			for i, want := range tt.wantArgs {
				if string(args[i]) != string(want) {
					t.Errorf("args[%d]: got %q, want %q", i, args[i], want)
				}
			}
		})
	}
}
