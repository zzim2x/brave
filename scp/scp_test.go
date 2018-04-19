package scp

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"fmt"
)

var k1, _ = randomSecret()
var k2, _ = randomSecret()
var k3, _ = randomSecret()
var k4, _ = randomSecret()
var k5, _ = randomSecret()
var k6, _ = randomSecret()
var k7, _ = randomSecret()

var quorumSet3T1 = QuorumSet{Threshold: 1, Validators: []PublicKey{k1.PublicKey, k2.PublicKey, k3.PublicKey}}
var quorumSet5T4 = QuorumSet{Threshold: 4, Validators: []PublicKey{k1.PublicKey, k2.PublicKey, k3.PublicKey, k4.PublicKey, k5.PublicKey}}
var quorumSet7T5 = QuorumSet{Threshold: 5, Validators: []PublicKey{k1.PublicKey, k2.PublicKey, k3.PublicKey, k4.PublicKey, k5.PublicKey, k6.PublicKey, k7.PublicKey}}

func TestSCP_GetSlot(t *testing.T) {
	scp := NewSCP(nil, k1.PublicKey, true, quorumSet3T1)

	s1 := scp.GetSlot(uint64(1), false)
	assert.Nil(t, s1)

	s2 := scp.GetSlot(uint64(1), true)
	assert.NotNil(t, s2)

	assert.Equal(t, uint64(1), scp.GetLowSlotIndex())
	assert.Equal(t, uint64(1), scp.GetHighSlotIndex())
}

func TestSCP_PurgeSlots(t *testing.T) {
	scp := NewSCP(nil, k1.PublicKey, true, quorumSet3T1)
	for i := 1; i <= 10; i++ {
		scp.GetSlot(uint64(i), true)
	}

	assert.Equal(t, uint64(1), scp.GetLowSlotIndex())
	assert.Equal(t, uint64(10), scp.GetHighSlotIndex())

	scp.PurgeSlots(6)
	assert.Equal(t, uint64(6), scp.GetLowSlotIndex())
	assert.Equal(t, uint64(10), scp.GetHighSlotIndex())
}

// quorum 5 threshold 4
func TestSCP_Simple(t *testing.T) {
	var envs []Envelope
	driver := &testDriver{}
	qs1 := quorumSet5T4

	scp := NewSCP(driver, qs1.Validators[0], true, qs1)

	votes := make([]Value, 0)

	votes = append(votes, Value{})

	scp.Nominate(1, votes[0], Value{})

	scp.ReceiveEnvelope(newNomination(1, *k2, qs1.Hash()))
	scp.ReceiveEnvelope(newNomination(1, *k3, qs1.Hash()))

	fmt.Println(envs, scp)

}