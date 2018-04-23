package scp

import (
	"hash"
	"github.com/davecgh/go-xdr/xdr"
	"github.com/emirpasic/gods/sets/treeset"
	"fmt"
)

type testDriver struct {
	quorumSets map[Hash]*QuorumSet
}

var _, _ Driver = &testDriver{}, (*testDriver)(nil)

func newTestDriver() *testDriver {
	return &testDriver{
		quorumSets: make(map[Hash]*QuorumSet),
	}
}

func (o *testDriver) VerifyEnvelope(envelope Envelope) bool {
	return true
}

func (o *testDriver) EmitEnvelope(envelope Envelope) {
	fmt.Println("EMIT", envelope)
}

func (o *testDriver) SignEnvelope(envelope *Envelope) {
}

func (o *testDriver) GetQuorumSet(hash Hash) *QuorumSet {
	return o.quorumSets[hash]
}

func (o *testDriver) storeQuorumSet(quorumSet *QuorumSet) {
	o.quorumSets[quorumSet.Hash()] = quorumSet
}

func (o *testDriver) ValidateValue(slotId uint64, value Value, nomination bool) ValidationLevel {
	return 0
}

func (o *testDriver) NominatingValue(slotIndex uint64, value Value) {
}

func (o *testDriver) CombineCandidates(slotIndex uint64, candidates *treeset.Set) Value {
	return nil
}

func (o *testDriver) UpdatedCandidateValue(slotIndex uint64, value Value) {
}

func (o *testDriver) SetupTimer(slotIndex uint64, timerId int32, timeout uint64, fn func()) {
}

func (o *testDriver) ComputeTimeout(roundNumber uint32) uint64 {
	return 1000
}

func (o *testDriver) ComputeHashNode(slotIndex uint64, prev Value, isPriority bool, roundNumber int32, nodeId PublicKey) uint64 {
	return hashHelper(slotIndex, prev, func(hash hash.Hash) {
		var priority uint32
		if isPriority {
			priority = 2
		} else {
			priority = 1
		}
		if b, err := xdr.Marshal(priority); err == nil {
			hash.Write(b)
		}
		if b, err := xdr.Marshal(roundNumber); err == nil {
			hash.Write(b)
		}
		if b, err := xdr.Marshal(nodeId); err == nil {
			hash.Write(b)
		}
	})
}

func (o *testDriver) ComputeHashValue(slotIndex uint64, prev Value, roundNumber int32, value Value) uint64 {
	return uint64(0)
}

func newNomination(slotIndex uint64, secretKey SecretKey, quorumSetHash Hash, votes []Value, accepted []Value) Envelope {
	return makeEnvelope(secretKey, Statement{
		SlotIndex:     slotIndex,
		NodeId:        secretKey.PublicKey,
		StatementType: StatementTypeNomination,
		Nomination: &Nomination{
			QuorumSetHash: quorumSetHash,
			Votes:         votes,
			Accepted:      accepted,
		},
	})
}

func makeEnvelope(secretKey SecretKey, statement Statement) Envelope {
	return Envelope{
		Statement: statement,
	}
}
