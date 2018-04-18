package scp

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

var k1 = PublicKey{Type: 0, Ed25519: [32]byte{1}}
var k2 = PublicKey{Type: 0, Ed25519: [32]byte{2}}
var k3 = PublicKey{Type: 0, Ed25519: [32]byte{3}}
var k4 = PublicKey{Type: 0, Ed25519: [32]byte{4}}
var k5 = PublicKey{Type: 0, Ed25519: [32]byte{5}}
var k6 = PublicKey{Type: 0, Ed25519: [32]byte{6}}
var k7 = PublicKey{Type: 0, Ed25519: [32]byte{7}}
var quorumSet3T1 = QuorumSet{Threshold: 1, Validators: []PublicKey{k1, k2, k3}}
var quorumSet5T4 = QuorumSet{Threshold: 4, Validators: []PublicKey{k1, k2, k3, k4, k5}}
var quorumSet7T5 = QuorumSet{Threshold: 5, Validators: []PublicKey{k1, k2, k3, k4, k5, k6, k7}}

func TestSCP_GetSlot(t *testing.T) {
	scp := NewSCP(nil, k1, true, quorumSet3T1)

	s1 := scp.GetSlot(uint64(1), false)
	assert.Nil(t, s1)

	s2 := scp.GetSlot(uint64(1), true)
	assert.NotNil(t, s2)

	assert.Equal(t, uint64(1), scp.GetLowSlotIndex())
	assert.Equal(t, uint64(1), scp.GetHighSlotIndex())
}

func TestSCP_PurgeSlots(t *testing.T) {
	scp := NewSCP(nil, k1, true, quorumSet3T1)
	for i := 1; i <= 10; i++ {
		scp.GetSlot(uint64(i), true)
	}

	assert.Equal(t, uint64(1), scp.GetLowSlotIndex())
	assert.Equal(t, uint64(10), scp.GetHighSlotIndex())

	scp.PurgeSlots(6)
	assert.Equal(t, uint64(6), scp.GetLowSlotIndex())
	assert.Equal(t, uint64(10), scp.GetHighSlotIndex())
}

// quorum 5 threshold 4 & nomination test
func TestSCP_Nominate(t *testing.T) {
	var envs []Envelope

	driver := &testDriver{}
	scp := NewSCP(driver, k1, true, quorumSet5T4)
	quorumSetHash := quorumSet5T4.Hash()

	votes, accepted := make([]Value, 0), make([]Value, 0)
	votes = append(votes, Value{})

	scp.ReceiveEnvelope()

}