package scp

import (
	"testing"
	"github.com/stretchr/testify/assert"
)

func TestForEachNodes(t *testing.T) {
	quorumSet1 := QuorumSet{
		Threshold:  1,
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

func Test_IsSan(t *testing.T) {
	st := Statement{
		Nomination: &Nomination{
			Votes:    []Value{},
			Accepted: []Value{},
		},
	}
	// votes & accepted 중 하나는 채워져야 함.
	assert.Equal(t, false, isSane(st))
	st.Nomination.Votes = append(st.Nomination.Votes, []uint8{0})
	assert.Equal(t, true, isSane(st))
	st.Nomination.Accepted = append(st.Nomination.Accepted, []uint8{0})
	assert.Equal(t, true, isSane(st))
	st.Nomination.Votes = append(st.Nomination.Votes, []uint8{0})
	assert.Equal(t, true, isSane(st))
	st.Nomination.Votes = append(st.Nomination.Votes, []uint8{5})
	assert.Equal(t, true, isSane(st))

	// votes, accepted value 정렬되어있어야 함.
	st.Nomination.Votes = append(st.Nomination.Votes, []uint8{3})
	assert.Equal(t, false, isSane(st))
}
