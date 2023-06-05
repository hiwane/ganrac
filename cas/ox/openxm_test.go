package ganrac

import (
	"fmt"
	"github.com/hiwane/ganrac"
	"net"
	"time"
)

func testConnectOx(g *ganrac.Ganrac) *OpenXM {
	time.Sleep(time.Second / 2)
	for i := 1; i <= 10; i++ {
		ox := _testConnectOx(g)
		if ox != nil {
			return ox
		}
		fmt.Printf("waiting for openxm server... %d sec\n", i*3)
		time.Sleep(time.Second * 3)
	}
	return nil
}

func _testConnectOx(g *ganrac.Ganrac) *OpenXM {
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
