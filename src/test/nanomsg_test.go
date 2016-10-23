// nanomsg_test
package test

import (
	"fmt"
	"os"

	"testing"
	"time"

	"github.com/go-mangos/mangos"
	"github.com/go-mangos/mangos/protocol/sub"
	"github.com/go-mangos/mangos/transport/ipc"
	"github.com/go-mangos/mangos/transport/tcp"
)

func TestSub(t *testing.T) {
	fmt.Println("Runing")
	client("tcp://127.0.0.1:40899", "d", t)
}
func die(format string, v ...interface{}) {
	fmt.Fprintln(os.Stderr, fmt.Sprintf(format, v...))
	os.Exit(1)
}
func client(url string, name string, t *testing.T) {
	var sock mangos.Socket
	var err error
	var msg []byte

	fmt.Println("Connecting")
	if sock, err = sub.NewSocket(); err != nil {
		die("can't get new sub socket: %s", err.Error())
	}
	sock.AddTransport(ipc.NewTransport())
	sock.AddTransport(tcp.NewTransport())

	if err = sock.Dial(url); err != nil {
		die("can't dial on sub socket: %s", err.Error())
	}
	// Empty byte array effectively subscribes to everything
	err = sock.SetOption(mangos.OptionReconnectTime, time.Millisecond*10)
	if err != nil {
		die("cannot set reconn: %s", err.Error())
	}
	err = sock.SetOption(mangos.OptionSubscribe, []byte(""))
	if err != nil {
		die("cannot subscribe: %s", err.Error())
	}

	fmt.Println("Waiting")
	for {
		if msg, err = sock.Recv(); err != nil {
			die("Cannot recv: %s", err.Error())
		}
		fmt.Printf("CLIENT(%s): RECEIVED %s\n", name, string(msg))
		t.Log("CLIENT(%s): RECEIVED %s\n", name, string(msg))
	}
}
