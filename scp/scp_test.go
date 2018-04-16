package scp

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

var k1 = PublicKey{Type: 0, Ed25519: [32]byte{1}}
var k2 = PublicKey{Type: 0, Ed25519: [32]byte{2}}
var k3 = PublicKey{Type: 0, Ed25519: [32]byte{3}}
var quorumSet1 = QuorumSet{Threshold: 1, Validators: []PublicKey{k1, k2, k3}}

func TestSCP_GetSlot(t *testing.T) {
	scp := NewSCP(k1, true, quorumSet1)

	s1 := scp.GetSlot(uint64(1), false)
	assert.Nil(t, s1)

	s2 := scp.GetSlot(uint64(1), true)
	assert.NotNil(t, s2)

	assert.Equal(t, uint64(1), scp.GetLowSlotIndex())
	assert.Equal(t, uint64(1), scp.GetHighSlotIndex())
}

func TestSCP_PurgeSlots(t *testing.T) {
	scp := NewSCP(k1, true, quorumSet1)
	for i := 1; i <= 10; i++ {
		scp.GetSlot(uint64(i), true)
	}

	assert.Equal(t, uint64(1), scp.GetLowSlotIndex())
	assert.Equal(t, uint64(10), scp.GetHighSlotIndex())

	scp.PurgeSlots(6)
	assert.Equal(t, uint64(6), scp.GetLowSlotIndex())
	assert.Equal(t, uint64(10), scp.GetHighSlotIndex())
}
