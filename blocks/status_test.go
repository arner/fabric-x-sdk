/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package blocks_test

import (
	"testing"

	"github.com/hyperledger/fabric-x-sdk/blocks"
)

func TestStatus_String(t *testing.T) {
	cases := []struct {
		status blocks.Status
		want   string
	}{
		{blocks.StatusUnknown, "UNKNOWN"},
		{blocks.StatusCommitted, "COMMITTED"},
		{blocks.StatusInvalidSignature, "INVALID_SIGNATURE"},
		{blocks.StatusMVCCConflict, "MVCC_CONFLICT"},
		{blocks.StatusDuplicateTxID, "DUPLICATE_TX_ID"},
		{blocks.StatusMalformed, "MALFORMED"},
		{blocks.StatusUnrecognized, "UNRECOGNIZED"},
		{blocks.StatusEndorsementPolicyFailure, "ENDORSEMENT_POLICY_FAILURE"},
		{blocks.Status(99), "UNKNOWN"}, // arbitrary raw value with no constant falls back to the default label
	}
	for _, c := range cases {
		if got := c.status.String(); got != c.want {
			t.Errorf("Status(%d).String() = %q, want %q", c.status, got, c.want)
		}
	}
}

func TestStatus_IsFinal(t *testing.T) {
	cases := []struct {
		status blocks.Status
		want   bool
	}{
		{blocks.StatusUnknown, false},
		{blocks.StatusCommitted, true},
		{blocks.StatusInvalidSignature, true},
		{blocks.StatusMVCCConflict, true},
		{blocks.StatusDuplicateTxID, true},
		{blocks.StatusMalformed, true},
		{blocks.StatusUnrecognized, true},
		{blocks.StatusEndorsementPolicyFailure, true},
	}
	for _, c := range cases {
		if got := c.status.IsFinal(); got != c.want {
			t.Errorf("Status(%d).IsFinal() = %v, want %v", c.status, got, c.want)
		}
	}
}

func TestStatus_Valid(t *testing.T) {
	cases := []struct {
		status blocks.Status
		want   bool
	}{
		{blocks.StatusUnknown, false},
		{blocks.StatusCommitted, true},
		{blocks.StatusInvalidSignature, false},
		{blocks.StatusMVCCConflict, false},
		{blocks.StatusDuplicateTxID, false},
		{blocks.StatusMalformed, false},
		{blocks.StatusUnrecognized, false},
		{blocks.StatusEndorsementPolicyFailure, false},
	}
	for _, c := range cases {
		if got := c.status.Valid(); got != c.want {
			t.Errorf("Status(%d).Valid() = %v, want %v", c.status, got, c.want)
		}
	}
}
