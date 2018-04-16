package scp

type slot struct {
	scp                *SCP
	slotId             uint64
	ballotProtocol     *ballotProtocol
	nominationProtocol *nominationProtocol
}

func newSlot(scp *SCP, slotId uint64) *slot {
	return &slot{
		scp:                scp,
		slotId:             slotId,
		ballotProtocol:     newBallotProtocol(),
		nominationProtocol: newNominationProtocol(),
	}
}

func (o *slot) nominate(value *Value, previousValue *Value, timeout bool) bool {
	return o.nominationProtocol.nominate(value, previousValue, timeout)
}

func (o *slot) stopNomination() {
	o.nominationProtocol.stopNomination()
}

func (o *slot) processEnvelope(envelope Envelope, self bool) EnvelopeState {
	if envelope.Statement.StatementType == StatementTypeNomination {
		return o.nominationProtocol.processEnvelope(envelope)
	} else {
		return o.ballotProtocol.processEnvelope(envelope, self)
	}
}