package main

import (
	"os"

	"github.com/elastic/beats/libbeat/beat"

	"github.com/gamegos/etcdbeat/beater"
)

func main() {
	err := beat.Run("etcdbeat", "", beater.New)
	if err != nil {
		os.Exit(1)
	}
}
