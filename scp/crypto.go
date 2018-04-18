package scp

const (
	KeyTypeEd25519   KeyType = 0
	KeyTypePreAuthTx KeyType = 1
	KeyTypeHashX     KeyType = 2
)

type KeyType int32
type Uint512 [64]uint8

type SecretKey struct {
	KeyType   KeyType
	secretKey Uint256
}

func randomSecret() (SecretKey, error) {
	return SecretKey{
		KeyType: KeyTypeEd25519,
	}, nil
}

func (o *SecretKey) Sign(bytes []byte) Signature {
	return nil
}

func (o *SecretKey) VerifySignature(publicKey PublicKey, signature Signature, bytes []byte) bool {
	return false
}

func (o *SecretKey) GetPublicKey() PublicKey {
	return PublicKey{}
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

