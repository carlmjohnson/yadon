package main

import (
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/carlmjohnson/yadon/prettyprint"
	"github.com/carlmjohnson/yadon/slowpoke"
)

func main() {
	throughput := flag.Float64("throughput", 1024, "target throughput in bytes per second")
	flag.Parse()
	for _, url := range flag.Args() {
		if err := MeasureBPS(url, *throughput); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v", err)
			os.Exit(1)
		}
	}
}

func MeasureBPS(url string, throughput float64) error {
	client, slowconn := slowpoke.NewClient()
	slowconn.PacketSize = 1024
	slowconn.TargetThroughput = throughput / float64(time.Second)

	req, _ := http.NewRequest("GET", url, nil)

	res, err := client.Do(req)
	if err == nil {
		_, err = io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	}

	bps := slowconn.BytesPer(time.Second)
	nr, nw := slowconn.BytesRead(), slowconn.BytesWritten()
	fmt.Printf("GET %s\n", url)
	fmt.Printf("Read: %d, Wrote: %d, Total: %d\n",
		nr, nw, (nr + nw))
	fmt.Printf("Bytes per second: %v/s\n", prettyprint.Size(bps))
	return err
}
