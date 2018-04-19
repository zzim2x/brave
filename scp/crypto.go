package scp

import (
	"io"
	"crypto/rand"
	"bytes"
	"github.com/agl/ed25519"
)

const (
	KeyTypeEd25519   KeyType = 0
	KeyTypePreAuthTx KeyType = 1
	KeyTypeHashX     KeyType = 2
)

type SecretKey struct {
	rawSeed   Uint256
	PublicKey PublicKey
}

func randomSecret() (*SecretKey, error) {
	var rawSeed [32]uint8
	if _, err := io.ReadFull(rand.Reader, rawSeed[:]); err != nil {
		return nil, err
	}

	if pubkey, _, err := ed25519.GenerateKey(bytes.NewReader(rawSeed[:])); err != nil {
		return nil, err
	} else {
		return &SecretKey{
			rawSeed: rawSeed,
			PublicKey: PublicKey{
				Type:    KeyTypeEd25519,
				Ed25519: *pubkey,
			},
		}, nil
	}
}

func (o *SecretKey) Sign(bytes []byte) Signature {
	return nil
}

func (o *SecretKey) VerifySignature(publicKey PublicKey, signature Signature, bytes []byte) bool {
	return false
}

// 눈으로 식별이 용이하고 수기할 수 있는
// ABCDEFGHIJKLMNOPQRSTUVWXYZ234567 base32
func (o *SecretKey) GetStrPublicKey() string {
	return ""
}

// 눈으로 식별이 용이하고 수기할 수 있는
// ABCDEFGHIJKLMNOPQRSTUVWXYZ234567 base32
func (o *SecretKey) GetStrKeySeed() string {
	return ""
}
