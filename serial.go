package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var (
	sFilename = flag.String("filename", "32u2.eep", "Name fo the eeprom image we insert the serial number in")
	bCreate   = flag.Bool("create", true, "Create the file if it does not exist")
)

func main() {
	flag.Parse()

	var image []byte
	var err error
	if image, err = ioutil.ReadFile(*sFilename); err != nil {
		if err != os.ErrNotExist && !*bCreate {
			fmt.Printf("FATAL %v\n", err)
			return
		}
		image = bytes.Repeat([]byte{0xFF}, 1024)
	}
	if len(image) != 1024 {
		fmt.Printf("Expected 1024 byte eprom image, not %v. Exiting\n", len(image))
		return
	}

	if flag.NArg() > 1 {
		fmt.Printf("Usage: %s [--create] [--filename filename] [Serialnumber]\n")
		flag.PrintDefaults()
	}
	if flag.NArg() == 1 {
		if len(flag.Arg(0)) > 19 {
			fmt.Printf("FATAL: Serial numbers need to be smaller than 20 characters\n")
			return
		}
		copy(image[256:256+32], bytes.Repeat([]byte{0xFF}, 20))
		for i, ch := range flag.Arg(0) {
			image[256+i] = byte(ch) + 0x80
		}

		if err := ioutil.WriteFile(*sFilename, image, os.FileMode(0666)); err != nil {
			fmt.Printf("Error writing file: %v\n", err)
			return
		}
	}

	var buf bytes.Buffer
	for p := 0x100; image[p] != 0xFF && p < 0x11A; p++ {
		buf.WriteByte(image[p] & 0x7F)
	}
	fmt.Printf("Serial number: %v\n", buf.String())
}
