package scp

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestSlot_IsNodeInQuorum(t *testing.T) {
	driver := newTestDriver()
	driver.storeQuorumSet(&quorumSet5T4)
	s1 := newSlot(NewSCP(driver, k1.PublicKey, true, quorumSet5T4), 1)
	s1.statementsHistory = append(
		s1.statementsHistory,
		slotHistoricalStatement{
			statement: Statement{
				NodeId:        k1.PublicKey,
				StatementType: StatementTypePrepare,
				Prepare: &StatementPrepare{
					QuorumSetHash: quorumSet5T4.Hash(),
				},
			},
			validated: true,
		},
	)

	assert.Equal(t, TriBoolTrue, s1.IsNodeInQuorum(k1.PublicKey))
	assert.Equal(t, TriBoolTrue, s1.IsNodeInQuorum(k2.PublicKey))
	assert.Equal(t, TriBoolTrue, s1.IsNodeInQuorum(k3.PublicKey))
	assert.Equal(t, TriBoolTrue, s1.IsNodeInQuorum(k4.PublicKey))
	assert.Equal(t, TriBoolTrue, s1.IsNodeInQuorum(k5.PublicKey))
	// 히스토리가 충분하지 않기때문에 maybe
	assert.Equal(t, TriBoolMaybe, s1.IsNodeInQuorum(k6.PublicKey))
}

// node 1-3 이 slot history 쌓였고 각 q = v{1,2,3}T1 일 때
func TestSlot_IsNodeInQuorum_v4_q3t1(t *testing.T) {
	driver := newTestDriver()
	driver.storeQuorumSet(&quorumSet3T1)
	s1 := newSlot(NewSCP(driver, k1.PublicKey, true, quorumSet3T1), 1)

	for _, nodeId := range []PublicKey{k1.PublicKey, k2.PublicKey, k3.PublicKey} {
		s1.statementsHistory = append(
			s1.statementsHistory,
			slotHistoricalStatement{
				statement: Statement{
					NodeId:        nodeId,
					StatementType: StatementTypePrepare,
					Prepare: &StatementPrepare{
						QuorumSetHash: quorumSet3T1.Hash(),
					},
				},
				validated: true,
			})
	}

	assert.Equal(t, TriBoolTrue, s1.IsNodeInQuorum(k1.PublicKey))
	assert.Equal(t, TriBoolFalse, s1.IsNodeInQuorum(k5.PublicKey))
}
