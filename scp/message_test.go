package scp

import (
	"testing"
	xdr2 "github.com/davecgh/go-xdr/xdr2"
	"bytes"
	"encoding/hex"
	"github.com/emirpasic/gods/sets/treeset"
	"github.com/stretchr/testify/assert"
)

// fixed & variant byte array marshal 결과 참조하고자
func TestBytes_MarshalXDRInto(t *testing.T) {
	v1 := Uint256{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}
	v2 := Value{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24, 25, 26, 27, 28, 29, 30, 31}

	var b1 bytes.Buffer
	var b2 bytes.Buffer
	xdr2.Marshal(&b1, v1)
	xdr2.Marshal(&b2, v2)

	assert.Equal(t, "000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f", hex.EncodeToString(b1.Bytes()))
	assert.Equal(t, "00000020000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f", hex.EncodeToString(b2.Bytes()))
}

func TestSortedSet(t *testing.T) {
	a := treeset.NewWith(ValueComparator)
	a.Add(Value{0, 1, 2, 9})
	a.Add(Value{0, 1, 2, 3})
	a.Add(Value{0, 1, 2, 4})
	a.Add(Value{0, 1, 2, 4, 5})

	values := a.Values()
	assert.Equal(t, Value{0, 1, 2, 3}, values[0])
	assert.Equal(t, Value{0, 1, 2, 4}, values[1])
	assert.Equal(t, Value{0, 1, 2, 4, 5}, values[2])
	assert.Equal(t, Value{0, 1, 2, 9}, values[3])
}
