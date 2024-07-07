package options

import (
	"net"
	"reflect"
	"testing"

	"github.com/frannotsleep/lilog/types"
)

func TestServerNameDecode(t *testing.T) {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:4852")
	if err != nil {
		t.Fatal(err)
	}

	server, err := net.ListenTCP("tcp", addr)

	if err != nil {
		t.Fatal(err)
	}

	defer server.Close()

	go func() {
		opts := LigLogOptions{}
		conn, err := server.AcceptTCP()

		if err != nil {
			t.Error(err)
		}

		defer conn.Close()

		typ, err := Decode(&opts, conn)

		if err != nil {
			t.Error(err)
		}

		if typ != ServerNameType {
			t.Errorf("unexpected type. got=%d\n", typ)
		}

		if !reflect.DeepEqual("Web Api", opts.ServerName.String()) {
			t.Errorf("unexpected ServerName. got=%s\n", opts.ServerName)
		}
	}()

	conn, err := net.DialTCP("tcp", nil, addr)

	if err != nil {
		t.Fatal(err)
	}

	sn := ServerName("Web Api")
	var payload types.Payload = &sn

	_, err = payload.WriteTo(conn)

	if err != nil {
		t.Fatal(err)
	}
}
