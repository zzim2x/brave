package scp

type nominationProtocol struct {
	started bool

	roundNumber int32

	// std::set 의 경우 기본이 sorted set 이다. 그에맞게 바꿔야한다. golang 에는 set container 도 없다.
	votes       []Value
	accepted    []Value
}

func newNominationProtocol() *nominationProtocol {
	return &nominationProtocol{
		started: false,
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
