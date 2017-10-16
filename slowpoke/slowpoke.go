package slowpoke

import (
	"net"
	"net/http"
	"time"
)

type Conn struct {
	net.Conn
	nr, nw           int
	start            time.Time
	PacketSize       int
	TargetThroughput float64
}

func (s *Conn) sleepif() {
	n := s.nr + s.nw
	if n < 1 {
		s.start = time.Now()
		return
	}
	since := float64(time.Since(s.start))
	bps := float64(n) / since
	if bps > s.TargetThroughput {
		d := float64(n)/s.TargetThroughput - since
		time.Sleep(time.Duration(d))
	}
}

func (s *Conn) Read(p []byte) (int, error) {
	if len(p) > s.PacketSize {
		p = p[:s.PacketSize]
	}
	s.sleepif()
	n, err := s.Conn.Read(p)
	s.nr += n
	return n, err
}

func (s *Conn) Write(p []byte) (int, error) {
	if len(p) > s.PacketSize {
		p = p[:s.PacketSize]
	}
	s.sleepif()
	n, err := s.Conn.Write(p)
	s.nw += n
	return n, err
}

func (s *Conn) BytesPer(unit time.Duration) float64 {
	return float64(s.nr+s.nw) / float64(time.Since(s.start)) * float64(unit)
}

func (s *Conn) BytesRead() int {
	return s.nr
}

func (s *Conn) BytesWritten() int {
	return s.nw
}

func Dialer(sc *Conn) func(network, address string) (net.Conn, error) {
	return func(network, address string) (net.Conn, error) {
		conn, err := net.Dial(network, address)
		sc.Conn = conn
		return sc, err
	}
}

func NewClient() (*http.Client, *Conn) {
	slowconn := &Conn{}

	var transport = &http.Transport{
		Proxy:              http.ProxyFromEnvironment,
		Dial:               Dialer(slowconn),
		DisableCompression: true,
	}

	var client = &http.Client{
		Transport: transport,
	}

	return client, slowconn
}
