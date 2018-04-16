package scp

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestForEachNodes(t *testing.T) {
	quorumSet1 := QuorumSet{
		Threshold: 1,
		Validators: []PublicKey{k1, k2, k3},
		InnerSets: []QuorumSet{
				{Threshold: 1, Validators: []PublicKey{k1, k2, k3}},
				{Threshold: 1, Validators: []PublicKey{k1, k2, k3}},
		},
	}

	count := 0
	ForEachNodes(quorumSet1, func(key PublicKey) {
		count += 1
	})
	assert.Equal(t, 3, count)
}
