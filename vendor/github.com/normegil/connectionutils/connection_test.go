package connectionutils_test

import (
	"net"
	"testing"

	"github.com/normegil/connectionutils"
	"github.com/normegil/interval"
)

func TestAvailable(t *testing.T) {
	ip := "127.0.0.1"
	port := 50900
	available := connectionutils.TCPPortAvalaible(&net.TCPAddr{
		IP:   net.ParseIP(ip),
		Port: port,
	})
	if !available {
		t.Errorf("Port not available {%s:%d}", ip, port)
	}
}

func TestNotAvailable(t *testing.T) {
	ip := "127.0.0.1"
	port := 50900
	tcpAddr := &net.TCPAddr{
		IP:   net.ParseIP(ip),
		Port: port,
	}
	defer Block(t, tcpAddr)()

	available := connectionutils.TCPPortAvalaible(tcpAddr)
	if available {
		t.Errorf("Port available {%s:%d} but should be used", ip, port)
	}
}

func TestSelectPort_SameIP(t *testing.T) {
	ports := interval.Test_ParseIntervalInteger(t, "[50900;50900]")

	tests := []string{
		"127.0.0.1",
		"192.168.0.1",
		"8.8.8.8",
		"::1",
		"2001:0db8:0000:85a3:0000:0000:ac1f:8001",
	}

	for _, test := range tests {
		ip := net.ParseIP(test)
		addr := connectionutils.SelectPort(ip, *ports)
		if addr.IP.String() != ip.String() {
			t.Errorf("Given IP (%+v) doesn't correspond to returned IP (%+v)", ip, addr.IP)
		}
	}
}

func TestSelectPort(t *testing.T) {
	ip := net.ParseIP("127.0.0.1")
	tests := []struct {
		name          string
		portsInterval string
		portsExcluded []int
		portsBlocked  []int
		expectedPort  int
	}{
		{"No port blocked", "[50900;50900]", []int{}, []int{}, 50900},
		{"First port blocked", "[50900;50901]", []int{}, []int{50900}, 50901},
		{"Last port blocked", "[50900;50901]", []int{}, []int{50901}, 50900},
		{"All but last port blocked", "[50900;50904]", []int{}, []int{50900, 50901, 50902, 50903}, 50904},
		{"Excluding first port", "]50900;50901]", []int{}, []int{}, 50901},
		{"Excluding last port", "]50900;50902[", []int{}, []int{}, 50901},
		{"No port avalaible", "[50900;50900]", []int{}, []int{50900}, 0},
		{"Exluding port in interval", "[50900;50902]", []int{50900}, []int{}, 50901},
		{"Exluding port list in interval", "[50900;50902]", []int{50900, 50901}, []int{}, 50902},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			for _, toBlock := range test.portsBlocked {
				defer Block(t, &net.TCPAddr{
					IP:   ip,
					Port: toBlock,
				})()
			}
			ports := interval.Test_ParseIntervalInteger(t, test.portsInterval)
			addr := connectionutils.SelectPortExcluding(ip, *ports, test.portsExcluded)
			if test.expectedPort != addr.Port {
				t.Errorf("%s: Selected port (%d) doesn't correspond to expected port (%d)", test.portsInterval, addr.Port, test.expectedPort)
			}
		})
	}
}

func Block(t *testing.T, tcpAddr *net.TCPAddr) func() error {
	listener, err := net.ListenTCP("tcp", tcpAddr)
	if nil != err {
		t.Fatal(err)
	}
	return listener.Close
}
