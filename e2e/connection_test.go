package e2e

import (
	"testing"
	"net"
	"fmt"
	"github.com/zzim2x/brave/sxdr"
	"github.com/davecgh/go-xdr/xdr2"
	"bytes"
	"github.com/zzim2x/brave/scp"
	"encoding/binary"
)

func Test_main(t *testing.T) {
	conn, err := net.Dial("tcp", "127.0.0.1:11625")
	if err != nil {
		panic(err)
	}

	hello := sxdr.Message{
		MessageType: sxdr.MessageTypeHello,
		Hello: &sxdr.Hello{
			LedgerVersion:     uint32(0),
			OverlayVersion:    uint32(0),
			OverlayMinVersion: uint32(0),
			NetworkID:         [32]byte{},
			VersionStr:        "",
			ListeningPort:     11625,
			NodeId: scp.PublicKey{
				Type:    0,
				Ed25519: [32]byte{},
			},
			Cert: sxdr.AuthCert{
				Pubkey:     [32]uint8{},
				Expiration: uint64(0),
				// 야는 가변길이이면서 뒤에 padding 붙여야 함.
				// (length + 7) & ~3
				Signature: []uint8{
					0, 0, 0, 0,
				},
			},
			Nonce: [32]uint8{},
		},
	}

	var payload bytes.Buffer

	xdr.Marshal(&payload, sxdr.AuthenticatedMessage{
		Type:     0,
		Sequence: uint64(1),
		Message:  hello,
		Mac:      [32]uint8{},
	})

	binary.Write(conn, binary.BigEndian, uint32(payload.Len()))
	n, e := conn.Write(payload.Bytes())

	fmt.Println(n, e)
}
