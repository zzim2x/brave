package scp

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestForEachNodes(t *testing.T) {
	quorumSet1 := QuorumSet{
		Threshold: 1,
		Validators: []PublicKey{k1.PublicKey, k2.PublicKey, k3.PublicKey},
		InnerSets: []QuorumSet{
				{Threshold: 1, Validators: []PublicKey{k1.PublicKey, k2.PublicKey, k3.PublicKey}},
				{Threshold: 1, Validators: []PublicKey{k1.PublicKey, k2.PublicKey, k3.PublicKey}},
		},
	}

	count := 0
	ForEachNodes(quorumSet1, func(key PublicKey) {
		count += 1
	})
	assert.Equal(t, 3, count)
}
