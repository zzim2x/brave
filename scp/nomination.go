package scp

import (
	"github.com/emirpasic/gods/sets/treeset"
)

type nominationProtocol struct {
	started bool

	roundNumber   int32
	previousValue Value
	votes         *treeset.Set
	accepted      *treeset.Set
	// votes, accepted = sorted set
	// marshaling 때, stellar-core 에서는 4byte 크기로 맞추고 실제 크기보다 4~7바이트 정도 크게 만들어진다.
}

func newNominationProtocol() *nominationProtocol {
	return &nominationProtocol{
		started:     false,
		roundNumber: 0,
		votes:       treeset.NewWith(ValueComparator),
		accepted:    treeset.NewWith(ValueComparator),
	}
}

func (o *nominationProtocol) nominate(value *Value, previousValue *Value, timeout bool) bool {
	o.started = true
	return true
}

func (o *nominationProtocol) processEnvelope(envelope Envelope) EnvelopeState {
	return EnvelopeStateValid
}

func (o *nominationProtocol) stopNomination() {
	o.started = false
}

func (o *nominationProtocol) getNewValueFromNomination() Value {
	return nil
}
