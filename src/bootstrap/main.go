package main

import (
	"bootstrap/config"
	"bootstrap/generator"
	"flag"
	"log"
)

func main() {
	efivars := flag.String("efivars", "/sys/firmware/efi/efivars", "efivarfs mount path")
	genfile := flag.String("generator", "/etc/bootstrap.gen", "bootstrap generator file")

	flag.Parse()

	cfg, err := config.Load(*efivars)
	if err != nil {
		log.Fatal("load config from efivars: %w", err)
	}

	var tc = struct {
		Config *config.Config
	}{
		Config: cfg,
	}

	batch, err := generator.Load(*genfile, &tc)
	if err != nil {
		log.Fatal("load generator file: %w", err)
	}

	err = batch.Write()
	if err != nil {
		log.Fatal("extract files from generated batch: %w", err)
	}
}
