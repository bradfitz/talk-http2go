// +build ignore

package main

import (
	"bytes"
	"log"

	"camlistore.org/third_party/github.com/bradfitz/http2"
	"camlistore.org/third_party/github.com/bradfitz/http2/hpack"
)

func main() {
	var buf bytes.Buffer
	f := http2.NewFramer(&buf, nil)

	var hpbuf bytes.Buffer
	hpe := hpack.NewEncoder(&hpbuf)
	hpe.WriteField(hpack.HeaderField{Name: ":method", Value: "GET"})
	hpe.WriteField(hpack.HeaderField{Name: ":path", Value: "/"})
	hpe.WriteField(hpack.HeaderField{Name: "host", Value: "example.com"})

	f.WriteHeaders(http2.HeadersFrameParam{
		StreamID:      1,
		EndStream:     true,
		EndHeaders:    true,
		BlockFragment: hpbuf.Bytes(),
	})

	log.Printf("Got: %q", buf.Bytes())
	log.Printf("Got: %x", buf.Bytes())
	log.Printf("Got: len %x type %x flags %x stream %x payload %x",
		buf.Bytes()[:3],
		buf.Bytes()[3:4],
		buf.Bytes()[4:5],
		buf.Bytes()[5:9],
		buf.Bytes()[9:])

	//  "\x00\x00\f\x01\x05\x00\x00\x00\x01\x82\x84f\x88/\x91\xd3]\x05\\\x87\xa7"
	//  "\x00\x00\f" \x01\x05\x00\x00\x00\x01\x82\x84f\x88/\x91\xd3]\x05\\\x87\xa7"

}
