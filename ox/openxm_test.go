package ganrac

import (
	"bufio"
	"github.com/hiwane/ganrac"
	"net"
	"time"
)

func testConnectOx(g *ganrac.Ganrac) *OpenXM {
	cport := "localhost:1234"
	dport := "localhost:4321"
	connc, err := net.Dial("tcp", cport)
	if err != nil {
		return nil
	}

	time.Sleep(time.Second / 20)

	connd, err := net.Dial("tcp", dport)
	if err != nil {
		connc.Close()
		return nil
	}

	ox, err := NewOpenXM(connc, connd, g.Logger())
	if err != nil {
		connc.Close()
		connd.Close()
		return nil
	}
	g.SetCAS(ox)
	return ox
}
