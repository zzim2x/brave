package scp

import "time"

type slot struct {
	scp                *SCP
	slotIndex          uint64
	ballotProtocol     *ballotProtocol
	nominationProtocol *nominationProtocol
	fullyValidated     bool

	statementsHistory []slotHistoricalStatement
}

type slotHistoricalStatement struct {
	when      uint64
	statement Statement
	validated bool
}

func newSlot(scp *SCP, slotId uint64) *slot {
	s := &slot{
		scp:               scp,
		slotIndex:         slotId,
		statementsHistory: make([]slotHistoricalStatement, 0),
	}
	s.ballotProtocol = newBallotProtocol(s)
	s.nominationProtocol = newNominationProtocol(s)
	return s
}

func (o *slot) nominate(value Value, previousValue Value, timeout bool) bool {
	return o.nominationProtocol.nominate(value, previousValue, timeout)
}

func (o *slot) stopNomination() {
	o.nominationProtocol.stopNomination()
}

func (o *slot) federatedAccept(votedPredicate func(Statement) bool, acceptedPredicate func(Statement) bool, envelopes map[PublicKey]Envelope) bool {
	return true
}

func (o *slot) federatedRatify(votedPredicate func(Statement) bool, envelopes map[PublicKey]Envelope) bool {
	return false
}

func (o *slot) getLatestCompositeCandidate() Value {
	return o.nominationProtocol.getLatestCompositeCandidate()
}

func (o *slot) recordStatement(statement Statement) {
	o.statementsHistory = append(o.statementsHistory, slotHistoricalStatement{
		when:      uint64(time.Now().Unix()),
		statement: statement,
		validated: o.fullyValidated,
	})
}

func (o *slot) processEnvelope(envelope Envelope, self bool) EnvelopeState {
	if envelope.Statement.StatementType == StatementTypeNomination {
		return o.nominationProtocol.processEnvelope(envelope)
	} else {
		return o.ballotProtocol.processEnvelope(envelope, self)
	}
}

func (o *slot) createEnvelope(statement Statement) Envelope {
	env := Envelope{
		Statement: statement,
	}

	env.Statement.NodeId = o.getLocalNode().nodeId
	env.Statement.SlotIndex = o.slotIndex
	o.getDriver().SignEnvelope(&env)

	return env
}

func (o *slot) IsNodeInQuorum(nodeId PublicKey) TriBool {
	history := make(map[PublicKey][]Statement)
	for _, h := range o.statementsHistory {
		if _, ok := history[h.statement.NodeId]; !ok {
			history[h.statement.NodeId] = make([]Statement, 0)
		}
		history[h.statement.NodeId] = append(history[h.statement.NodeId], h.statement)
	}

	return o.scp.localNode.IsNodeInQuorum(nodeId, o.getQuorumSet, history)
}

func (o *slot) getQuorumSet(statement Statement) *QuorumSet {
	return o.scp.driver.GetQuorumSet(o.getCompanionQuorumSetHashFromStatement(statement))
}

func (o *slot) getCompanionQuorumSetHashFromStatement(statement Statement) Hash {
	switch statement.StatementType {
	case StatementTypePrepare:
		return statement.Prepare.QuorumSetHash
	case StatementTypeConfirm:
		return statement.Confirm.QuorumSetHash
	case StatementTypeExternalize:
		return statement.Externalize.CommitQuorumSetHash
	case StatementTypeNomination:
		return statement.Nomination.QuorumSetHash
	}
	return Hash{}
}

func (o *slot) getLocalNode() *LocalNode {
	return o.scp.localNode
}

func (o *slot) getDriver() Driver {
	return o.scp.driver
}
