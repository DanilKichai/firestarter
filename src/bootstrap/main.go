package main

import (
	"bootstrap/config"
	"flag"
	"log"
)

func main() {
	efivars := flag.String("efivars", "/sys/firmware/efi/efivars", "efivarfs mount path")
	flag.Parse()

	_, err := config.Load(*efivars)
	if err != nil {
		log.Fatal("load config: %w", err)
	}
}
