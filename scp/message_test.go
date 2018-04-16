package scp

import (
	"testing"
	xdr2 "github.com/davecgh/go-xdr/xdr2"
	"bytes"
	"fmt"
	"encoding/hex"
)

// fixed & variant byte array marshal 결과 참조하고자 (일단은 대충 println 으로)
func TestBytes_MarshalXDRInto(t *testing.T) {
	v1 := Uint256{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31}
	v2 := Value{0,1,2,3,4,5,6,7,8,9,10,11,12,13,14,15,16,17,18,19,20,21,22,23,24,25,26,27,28,29,30,31}

	var b1 bytes.Buffer
	var b2 bytes.Buffer
	xdr2.Marshal(&b1, v1)
	fmt.Println(hex.EncodeToString(b1.Bytes()))
	xdr2.Marshal(&b2, v2)
	fmt.Println(hex.EncodeToString(b2.Bytes()))
}
