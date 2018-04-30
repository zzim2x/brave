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
		slot:              s,
		started:           false,
		roundNumber:       0,
		latestNominations: make(map[PublicKey]Envelope),
		roundLeaders:      treeset.NewWith(ValueComparator),
		votes:             treeset.NewWith(ValueComparator),
		accepted:          treeset.NewWith(ValueComparator),
		candidates:        treeset.NewWith(ValueComparator),
	}
}

func (o *nominationProtocol) isNewerStatement(nodeId PublicKey, nomination Nomination) bool {
	isNewer := false
	if old, exists := o.latestNominations[nodeId]; exists {
		isNewer = isNewerStatement(*old.Statement.Nomination, nomination)
	} else {
		isNewer = true
	}

	return isNewer
}

func isNewerStatement(old Nomination, st Nomination) bool {
	return true
}

func isSane(statement Statement) bool {
	nom := statement.Nomination

	if len(nom.Votes)+len(nom.Accepted) == 0 {
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
		o.roundLeaders.Each(func(index int, value interface{}) {
			fmt.Println("latestNominations", value, o.latestNominations[value.(PublicKey)])
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

	nominationTimeout := o.slot.getDriver().ComputeTimeout(uint32(o.roundNumber))
	o.slot.getDriver().NominatingValue(o.slot.slotIndex, nominatingValue)

	o.slot.getDriver().SetupTimer(o.slot.slotIndex, 0, nominationTimeout, func() {
		o.slot.nominate(value, previousValue, true)
	})

	if updated {
		o.emitNomination()
	}

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
	return o.slot.getDriver().ComputeHashNode(o.slot.slotIndex, o.previousValue, isPriority, o.roundNumber, nodeId)
}

func (o *nominationProtocol) hashValue(value Value) uint64 {
	return o.slot.getDriver().ComputeHashValue(o.slot.slotIndex, o.previousValue, o.roundNumber, value)
}

func (o *nominationProtocol) processEnvelope(envelope Envelope) EnvelopeState {
	st := envelope.Statement
	nom := st.Nomination

	isNewer := o.isNewerStatement(st.NodeId, *nom)
	res := EnvelopeStateInvalid


	if isNewer {
		if isSane(st) {
			o.recordEnvelope(envelope)
			res = EnvelopeStateValid

			if o.started {
				modified := false
				newCandidates := false

				for _, v := range nom.Votes {
					if o.accepted.Contains(v) {
						continue
					}

					if o.slot.federatedAccept(votedPredicate(v), acceptedPredicate(v), o.latestNominations) {
						vl := o.validateValue(v)
						if vl == ValidationLevelFullyValidatedValue {
							o.accepted.Add(v)
							o.votes.Add(v)
							modified = true
						} else {
							toVote := o.extractValidValue(v)
							if len(toVote) != 0 {
								if !o.votes.Contains(toVote) {
									o.votes.Add(toVote)
									modified = true
								}
							}
						}
					}
				}

				for _, v := range nom.Accepted {
					if o.candidates.Contains(v) {
						continue
					}

					if o.slot.federatedRatify(acceptedPredicate(v), o.latestNominations) {
						o.candidates.Add(v)
						newCandidates = true
					}
				}

				if o.candidates.Empty() && o.roundLeaders.Contains(st.NodeId) {
					newVote := o.getNewValueFromNomination(*nom)
					if len(newVote) > 0 {
						o.votes.Add(newVote)
						modified = true
						o.slot.getDriver().NominatingValue(o.slot.slotIndex, newVote)
					}
				}

				if modified {
					o.emitNomination()
				}

				if newCandidates {
					o.latestCompositeCandidate = o.slot.getDriver().CombineCandidates(o.slot.slotIndex, o.candidates)
					o.slot.getDriver().UpdatedCandidateValue(o.slot.slotIndex, o.latestCompositeCandidate)
					o.slot.bumpState(o.latestCompositeCandidate, false)
				}
			}

			return res
		}
	}

	return res
}

func votedPredicate(v Value) func(Statement) bool {
	return func(statement Statement) bool {
		n := statement.Nomination
		for _, nomVote := range n.Votes {
			if ValueComparator(v, nomVote) == 0 {
				return true
			}
		}
		return false
	}
}

func acceptedPredicate(v Value) func(Statement) bool {
	return func(statement Statement) bool {
		n := statement.Nomination
		for _, nomVote := range n.Accepted {
			if ValueComparator(v, nomVote) == 0 {
				return true
			}
		}
		return false
	}
}

func (o *nominationProtocol) validateValue(value Value) ValidationLevel {
	return o.slot.scp.validateValue(o.slot.slotIndex, value, true)
}

func (o *nominationProtocol) extractValidValue(value Value) Value {
	return o.slot.getDriver().ExtractValidValue(o.slot.slotIndex, value)
}

func (o *nominationProtocol) emitNomination() {
	st := Statement{}
	st.NodeId = o.slot.getLocalNode().nodeId
	st.StatementType = StatementTypeNomination

	votes := make([]Value, 0)
	accepted := make([]Value, 0)

	for _, v := range o.votes.Values() {
		votes = append(votes, v.(Value))
	}
	for _, v := range o.accepted.Values() {
		accepted = append(accepted, v.(Value))
	}

	st.Nomination = &Nomination{
		QuorumSetHash: o.slot.getLocalNode().quorumSetHash,
		Votes:         votes,
		Accepted:      accepted,
	}

	env := o.slot.createEnvelope(st)
	if o.slot.processEnvelope(env, true) == EnvelopeStateValid {
		if o.lastEnvelope == nil || isNewerStatement(*o.lastEnvelope.Statement.Nomination, *st.Nomination) {
			o.lastEnvelope = &env
			if o.slot.fullyValidated {
				o.slot.getDriver().EmitEnvelope(env)
			}
		}
	} else {
		fmt.Println("ERROR!!!!!")
	}
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

func (o *nominationProtocol) recordEnvelope(envelope Envelope) {
	o.latestNominations[envelope.Statement.NodeId] = envelope
	o.slot.recordStatement(envelope.Statement)
}

func (o *nominationProtocol) getCurrentState() []Envelope {
	return nil
}

func (o *nominationProtocol) getNewValueFromNomination(nomination Nomination) Value {
	return nil
}
