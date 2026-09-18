/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

package fabricx

import (
	"github.com/hyperledger/fabric-x-common/api/applicationpb"
	"github.com/hyperledger/fabric-x-sdk/blocks"
)

// DecodeMetadata extracts Events, Payload, and InputArgs from the transaction metadata.
// The layout is purely positional: metadata[0] = event, metadata[1] = payload,
// metadata[2] = arg count (1 byte), metadata[3:3+count] = args.
func DecodeMetadata(metadata [][]byte) (events []byte, payload []byte, inputArgs [][]byte) {
	if len(metadata) > 0 {
		events = metadata[0]
	}
	if len(metadata) > 1 {
		payload = metadata[1]
	}
	if len(metadata) > 2 && len(metadata[2]) > 0 {
		count := int(metadata[2][0])
		if end := 3 + count; end <= len(metadata) {
			inputArgs = metadata[3:end]
		}
	}
	return events, payload, inputArgs
}

// DecodeNamespaces converts Fabric-X TxNamespace protos into the SDK's
// protocol-neutral per-namespace read/write sets.
func DecodeNamespaces(namespaces []*applicationpb.TxNamespace) []blocks.NsReadWriteSet {
	nsRWS := make([]blocks.NsReadWriteSet, len(namespaces))
	for i, ns := range namespaces {
		nsrws := blocks.NsReadWriteSet{
			Namespace: ns.NsId,
			RWS: blocks.ReadWriteSet{
				Reads:  []blocks.KVRead{},
				Writes: []blocks.KVWrite{},
			},
		}

		for _, r := range ns.ReadsOnly {
			read := blocks.KVRead{Key: string(r.Key)}
			// Version nil means "no constraint" (new key / blind-write semantics).
			// Version 0 is a valid MVCC constraint: the key was first written at block 0.
			if r.Version != nil {
				read.Version = &blocks.Version{
					BlockNum: *r.Version,
				}
			}
			nsrws.RWS.Reads = append(nsrws.RWS.Reads, read)
		}
		for _, bw := range ns.BlindWrites {
			// All blind writes are now normal world state writes
			// (events and inputs are in metadata)
			nsrws.RWS.Writes = append(nsrws.RWS.Writes, blocks.KVWrite{
				Key:   string(bw.Key),
				Value: bw.Value,
			})
		}
		for _, rw := range ns.ReadWrites {
			read := blocks.KVRead{Key: string(rw.Key)}
			// Version nil means "no constraint" (new key / blind-write semantics).
			// Version 0 is a valid MVCC constraint: the key was first written at block 0.
			if rw.Version != nil {
				read.Version = &blocks.Version{
					BlockNum: *rw.Version,
				}
			}
			nsrws.RWS.Reads = append(nsrws.RWS.Reads, read)
			nsrws.RWS.Writes = append(nsrws.RWS.Writes, blocks.KVWrite{
				Key:   string(rw.Key),
				Value: rw.Value,
			})
		}

		nsRWS[i] = nsrws
	}
	return nsRWS
}
