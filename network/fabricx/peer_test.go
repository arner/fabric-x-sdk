/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package fabricx

import (
	"testing"

	"github.com/hyperledger/fabric-x-common/api/committerpb"
	"github.com/hyperledger/fabric-x-sdk/blocks"
)

func TestConvertTxEventBatch_FieldMapping(t *testing.T) {
	protoBatch := &committerpb.TxEventBatch{
		BlockNumber: 42,
		Events: []*committerpb.TxEvent{
			{
				Ref:    &committerpb.TxRef{TxId: "txABC", BlockNum: 42, TxNum: 3},
				Status: committerpb.Status_COMMITTED,
			},
			{
				Ref:    &committerpb.TxRef{TxId: "txDEF", BlockNum: 42, TxNum: 7},
				Status: committerpb.Status_ABORTED_MVCC_CONFLICT,
			},
		},
	}

	got := convertTxEventBatch(protoBatch)

	if got.BlockNumber != 42 {
		t.Errorf("BlockNumber: want 42, got %d", got.BlockNumber)
	}
	if len(got.Events) != 2 {
		t.Fatalf("want 2 events, got %d", len(got.Events))
	}

	first := got.Events[0]
	if first.ID != "txABC" || first.BlockNum != 42 || first.Number != 3 {
		t.Errorf("first event ref mismatch: %+v", first)
	}
	if first.Status != blocks.StatusCommitted || first.Reason != "COMMITTED" ||
		first.RawCode != int32(committerpb.Status_COMMITTED) {
		t.Errorf("first event status: want COMMITTED, got %v (%q, code %d)", first.Status, first.Reason, first.RawCode)
	}

	second := got.Events[1]
	if second.ID != "txDEF" || second.Number != 7 || second.Status != blocks.StatusMVCCConflict ||
		second.Reason != "ABORTED_MVCC_CONFLICT" || second.RawCode != int32(committerpb.Status_ABORTED_MVCC_CONFLICT) {
		t.Errorf("second event mismatch: %+v", second)
	}
}

func TestConvertTxEventBatch_EmptyEvents(t *testing.T) {
	got := convertTxEventBatch(&committerpb.TxEventBatch{BlockNumber: 7})
	if got.BlockNumber != 7 {
		t.Errorf("BlockNumber: want 7, got %d", got.BlockNumber)
	}
	if len(got.Events) != 0 {
		t.Errorf("expected empty events slice, got %d", len(got.Events))
	}
}

func TestToProtoFilterStatus(t *testing.T) {
	tests := []struct {
		name string
		in   []blocks.Status
		want []committerpb.Status
	}{
		{
			name: "nil in means nil out",
			in:   nil,
			want: nil,
		},
		{
			name: "committed",
			in:   []blocks.Status{blocks.StatusCommitted},
			want: []committerpb.Status{committerpb.Status_COMMITTED},
		},
		{
			name: "endorsement policy failure has no distinct Fabric-X code: expands to the signature-invalid code",
			in:   []blocks.Status{blocks.StatusEndorsementPolicyFailure},
			want: []committerpb.Status{committerpb.Status_ABORTED_SIGNATURE_INVALID},
		},
		{
			name: "invalid signature and endorsement policy failure collapse onto the same code",
			in:   []blocks.Status{blocks.StatusInvalidSignature, blocks.StatusEndorsementPolicyFailure},
			want: []committerpb.Status{committerpb.Status_ABORTED_SIGNATURE_INVALID, committerpb.Status_ABORTED_SIGNATURE_INVALID},
		},
		{
			name: "unrecognized has no fixed code and is dropped",
			in:   []blocks.Status{blocks.StatusUnrecognized},
			want: nil,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := toProtoFilterStatus(tc.in)
			if len(got) != len(tc.want) {
				t.Fatalf("toProtoFilterStatus(%v) = %v, want %v", tc.in, got, tc.want)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Errorf("toProtoFilterStatus(%v)[%d] = %v, want %v", tc.in, i, got[i], tc.want[i])
				}
			}
		})
	}
}
