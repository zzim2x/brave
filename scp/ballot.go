package scp

type ballotProtocol struct {
}

func newBallotProtocol() *ballotProtocol {
	return &ballotProtocol{
	}
}

func (o *ballotProtocol) processEnvelope(envelope Envelope, self bool) EnvelopeState {
	return EnvelopeStateValid
}
