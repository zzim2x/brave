package scp

import (
	"github.com/emirpasic/gods/sets/treeset"
)

type nominationProtocol struct {
	started                  bool
	roundNumber              int32
	roundLeaders             *treeset.Set
	previousValue            Value
	lastEnvelope             *Envelope
	latestCompositeCandidate Value

	votes             *treeset.Set           // X
	accepted          *treeset.Set           // Y
	candidates        *treeset.Set           // Z
	latestNominations map[PublicKey]Envelope // N
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

func (o *nominationProtocol) isNewerStatement(nodeId PublicKey, nomination Nomination) bool {
	return true
}

func isNewerStatement(old Nomination, st Nomination) bool {
	return true
}

func isSubsetHelper(p []Value, v []Value, notEqual bool) bool {
	return true
}

func (o *nominationProtocol) nominate(value Value, previousValue Value, timeout bool) bool {
	if timeout && !o.started {
		return false
	}

	o.started = true
	o.previousValue = previousValue
	o.roundNumber++

	o.updateRoundLeaders()

	return true
}

func (o *nominationProtocol) updateRoundLeaders() {

}

func (o *nominationProtocol) processEnvelope(envelope Envelope) EnvelopeState {
	return EnvelopeStateValid
}

func getStatementValues(statement Statement) []Value {
	return nil
}

func (o *nominationProtocol) stopNomination() {
	o.started = false
}

func (o *nominationProtocol) getLatestCompositeCandidate() Value {
	return nil
}

func (o *nominationProtocol) getLastMessageSend() *Envelope {
	return o.lastEnvelope
}

func (o *nominationProtocol) setStateFromEnvelope(envelope Envelope) {
}

func (o *nominationProtocol) getCurrentState() []Envelope {
	return nil
}

func (o *nominationProtocol) getNewValueFromNomination() Value {
	return nil
}
