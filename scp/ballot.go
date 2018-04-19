package scp

type ballotProtocol struct {
	slot *slot
}

func newBallotProtocol(s *slot) *ballotProtocol {
	return &ballotProtocol{
		slot: s,
	}
}

func (o *ballotProtocol) processEnvelope(envelope Envelope, self bool) EnvelopeState {
	return EnvelopeStateValid
}
