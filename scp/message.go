package scp

import (
	"github.com/davecgh/go-xdr/xdr2"
	"bytes"
	"crypto/sha256"
)

type Hash [32]uint8
type Uint256 [32]byte
type Value []uint8
type Signature []byte // variable payload max size : 64 : (size + 7) & ~3
type SignatureHint [4]byte
type StatementType int32

const (
	StatementTypePrepare     StatementType = 1
	StatementTypeConfirm     StatementType = 2
	StatementTypeExternalize StatementType = 3
	StatementTypeNomination  StatementType = 4
)

type QuorumSet struct {
	Threshold  uint32
	Validators []PublicKey
	InnerSets  []QuorumSet
}

type PublicKey struct {
	Type    int32
	Ed25519 Uint256
}

type Envelope struct {
	Statement Statement
	Signature Signature
}

type Statement struct {
	NodeId        PublicKey
	SlotIndex     uint64
	StatementType StatementType
	Prepare       *StatementPrepare
	Confirm       *StatementConfirm
	Externalize   *StatementExternalize
	Nomination    *Nomination
}

type StatementPrepare struct {
	QuorumSetHash Hash
	Ballot        Ballot
	Prepared      *Ballot
	PreparedPrime *Ballot
	NC            uint32
	NH            uint32
}

type StatementConfirm struct {
	Ballot        Ballot
	NPrepared     uint32
	NCommit       uint32
	NH            uint32
	QuorumSetHash Hash
}

type StatementExternalize struct {
	Commit              Ballot
	NH                  uint32
	CommitQuorumSetHash Hash
}

type Nomination struct {
	QuorumSetHash Hash
	Votes         []Value
	Accepted      []Value
}

type Ballot struct {
	Counter uint32
	Value   Value
}

func (o *QuorumSet) Hash() (h Hash) {
	var b bytes.Buffer
	xdr.Marshal(&b, o)

	sha := sha256.New()
	sha.Write(b.Bytes())
	copy(h[:], sha.Sum(nil))
	return
}

func ValueComparator(a, b interface{}) int {
	var v1, v2 []uint8

	if _, ok := a.(Value); ok {
		v1 = a.(Value)
	} else if _, ok := a.([]uint8); ok {
		v1 = a.([]uint8)
	}

	if _, ok := b.(Value); ok {
		v2 = b.(Value)
	} else if _, ok := b.([]uint8); ok {
		v2 = b.([]uint8)
	}

	s1, s2 := len(v1), len(v2)

	var minIndex int
	if s1 < s2 {
		minIndex = s1
	} else {
		minIndex = s2
	}

	for i := 0; i < minIndex; i++ {
		if v1[i] < v2[i] {
			return -1
		} else if v1[i] > v2[i] {
			return 1
		}
	}

	if s1 < s2 {
		return -1
	} else if s1 == s2 {
		return 0
	} else {
		return 1
	}
}