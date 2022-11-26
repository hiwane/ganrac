package ganrac

import (
	"bufio"
	"net"
	"time"
	"github.com/hiwane/ganrac"
)

func testConnectOx(g *ganrac.Ganrac) (net.Conn, net.Conn) {
	cport := "localhost:1234"
	dport := "localhost:4321"
	connc, err := net.Dial("tcp", cport)
	if err != nil {
		return nil, nil
	}

	time.Sleep(time.Second / 20)

	connd, err := net.Dial("tcp", dport)
	if err != nil {
		connc.Close()
		return nil, nil
	}

	dw := bufio.NewWriter(connd)
	dr := bufio.NewReader(connd)
	cw := bufio.NewWriter(connc)
	cr := bufio.NewReader(connc)

	ox, err := NewOpenXM(cw, dw, cr, dr, g.Logger())
	if err != nil {
		connc.Close()
		connd.Close()
		return nil, nil
	}
	g.SetCAS(ox)
	return connc, connd
}
