package scp

import (
	"testing"
	"github.com/stretchr/testify/assert"
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
	return false
}

func (o *testDriver) EmitEnvelope(envelope Envelope) {
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

func TestSlot_IsNodeInQuorum(t *testing.T) {
	driver := newTestDriver()
	driver.storeQuorumSet(&quorumSet5T4)
	s1 := newSlot(NewSCP(driver, k1, true, quorumSet5T4), 1)
	s1.statementsHistory = append(
		s1.statementsHistory,
		slotHistoricalStatement{
			statement: Statement{
				NodeId: k1,
				StatementType: StatementTypePrepare,
				Prepare: &StatementPrepare{
					QuorumSetHash: quorumSet5T4.Hash(),
				},
			},
			validated: true,
		},
	)

	assert.Equal(t, TriBoolTrue, s1.IsNodeInQuorum(k1))
	assert.Equal(t, TriBoolTrue, s1.IsNodeInQuorum(k2))
	assert.Equal(t, TriBoolTrue, s1.IsNodeInQuorum(k3))
	assert.Equal(t, TriBoolTrue, s1.IsNodeInQuorum(k4))
	assert.Equal(t, TriBoolTrue, s1.IsNodeInQuorum(k5))
	// 히스토리가 충분하지 않기때문에 maybe
	assert.Equal(t, TriBoolMaybe, s1.IsNodeInQuorum(k6))
}

// node 1-3 이 slot history 쌓였고 각 q = v{1,2,3}T1 일 때
func TestSlot_IsNodeInQuorum_v4_q3t1(t *testing.T) {
	driver := newTestDriver()
	driver.storeQuorumSet(&quorumSet3T1)
	s1 := newSlot(NewSCP(driver, k1, true, quorumSet3T1), 1)

	for _, nodeId := range []PublicKey{k1, k2, k3} {
		s1.statementsHistory = append(
			s1.statementsHistory,
			slotHistoricalStatement{
				statement: Statement{
					NodeId: nodeId,
					StatementType: StatementTypePrepare,
					Prepare: &StatementPrepare{
						QuorumSetHash: quorumSet3T1.Hash(),
					},
				},
				validated: true,
			})
	}

	assert.Equal(t, TriBoolTrue, s1.IsNodeInQuorum(k1))
	assert.Equal(t, TriBoolFalse, s1.IsNodeInQuorum(k5))
}