package scp

import (
	"github.com/davecgh/go-xdr/xdr2"
	"bytes"
	"crypto/sha256"
)

type Hash [32]byte
type Uint256 [32]byte
type Value []uint8
type Signature [64]byte
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
	innerSets  []QuorumSet
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
