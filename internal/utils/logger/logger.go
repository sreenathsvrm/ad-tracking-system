package logger

import (
	"log"
	"os"
)

var Log = log.New(os.Stdout, "ad-tracking-system: ", log.LstdFlags|log.Lshortfile)
