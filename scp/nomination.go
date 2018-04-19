package scp

import (
	"github.com/emirpasic/gods/sets/treeset"
	"fmt"
)

type nominationProtocol struct {
	slot                     *slot
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

func newNominationProtocol(s *slot) *nominationProtocol {
	return &nominationProtocol{
		slot:         s,
		started:      false,
		roundNumber:  0,
		roundLeaders: treeset.NewWith(ValueComparator),
		votes:        treeset.NewWith(ValueComparator),
		accepted:     treeset.NewWith(ValueComparator),
	}
}

func (o *nominationProtocol) isNewerStatement(nodeId PublicKey, nomination Nomination) bool {
	return true
}

func isNewerStatement(old Nomination, st Nomination) bool {
	return true
}

func isSane(statement Statement) bool {
	nom := statement.Nomination

	if len(nom.Votes) + len(nom.Accepted) == 0 {
		return false
	}

	for i := 1; i < len(nom.Votes); i++ {
		if ValueComparator(nom.Votes[i-1], nom.Votes[i]) == 1 {
			return false
		}
	}

	for i := 1; i < len(nom.Accepted); i++ {
		if ValueComparator(nom.Accepted[i-1], nom.Accepted[i]) == 1 {
			return false
		}
	}

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

	var updated = false
	var nominatingValue Value

	if o.roundLeaders.Contains(o.slot.getLocalNode().nodeId) {
		if !o.votes.Contains(value) {
			o.votes.Add(value)
			updated = true
		}
		nominatingValue = value
	} else {
		fmt.Println("o.roundLeaders")
		o.roundLeaders.Each(func(index int, value interface{}) {
			if it, exists := o.latestNominations[value.(PublicKey)]; exists {
				nom := it.Statement.Nomination
				nominatingValue = o.getNewValueFromNomination(*nom)
				if len(nominatingValue) == 0 {
					o.votes.Add(nominatingValue)
					updated = true
				}
			}
		})
	}

	// wip

	return updated
}

func (o *nominationProtocol) updateRoundLeaders() {
	o.roundLeaders.Clear()
	topPriority := uint64(0)
	qset := o.slot.getLocalNode().quorumSet

	forEachNodes(qset, func(nodeId PublicKey) {
		w := o.getNodePriority(nodeId, qset)
		if w > topPriority {
			topPriority = w
			o.roundLeaders.Clear()
		}
		if w == topPriority && w > 0 {
			o.roundLeaders.Add(nodeId)
		}
	})

	fmt.Println("updateRoundLeaders:", o.roundLeaders.Size())
	for _, leader := range o.roundLeaders.Values() {
		fmt.Println("    leader", leader.(PublicKey).Address())
	}
}

func (o *nominationProtocol) getNodePriority(nodeId PublicKey, quorumSet QuorumSet) uint64 {
	w := GetNodeWeight(nodeId, quorumSet)

	if o.hashNode(false, nodeId) < w {
		return o.hashNode(true, nodeId)
	} else {
		return uint64(0)
	}
}

func (o *nominationProtocol) hashNode(isPriority bool, nodeId PublicKey) uint64 {
	return o.slot.getDriver().ComputeHashNode(o.slot.slotId, o.previousValue, isPriority, o.roundNumber, nodeId)
}

func (o *nominationProtocol) hashValue(value Value) uint64 {
	return o.slot.getDriver().ComputeHashValue(o.slot.slotId, o.previousValue, o.roundNumber, value)
}

func (o *nominationProtocol) processEnvelope(envelope Envelope) EnvelopeState {
	st := envelope.Statement
	nom := st.Nomination

	isNewer := false
	if old, exists := o.latestNominations[st.NodeId]; exists {
		isNewer = isNewerStatement(*old.Statement.Nomination, *nom)
	} else {
		isNewer = true
	}

	if isNewer {
		if isSane(st) {
		}
	}

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

func (o *nominationProtocol) getNewValueFromNomination(nomination Nomination) Value {
	return nil
}
