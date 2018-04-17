package sxdr

import "github.com/zzim2x/brave/scp"

type Curve25519Public [32]uint8

type Hello struct {
	LedgerVersion     uint32
	OverlayVersion    uint32
	OverlayMinVersion uint32
	NetworkID         scp.Hash
	VersionStr        string
	ListeningPort     int32
	NodeId            scp.PublicKey
	Cert              AuthCert
	Nonce             scp.Uint256
}

type AuthCert struct {
	Pubkey     Curve25519Public
	Expiration uint64
	Signature  scp.Signature
}
