package adifparser

import (
	"log"
	"os"
)

var adiflog *log.Logger

func init() {
	adiflog = log.New(os.Stderr, "adifparser: ", log.Lshortfile|log.LstdFlags)
}
