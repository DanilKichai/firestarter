package main

import (
	"firestarter/internal/app/bootstrap/config"
	"firestarter/internal/pkg/generator"
	"flag"
	"log"
)

func main() {
	efivars := flag.String("efivars", "/sys/firmware/efi/efivars", "efivarfs mount path")
	genfile := flag.String("generator", "/etc/bootstrap.gen", "generator file")

	flag.Parse()

	cfg, err := config.Load(*efivars)
	if err != nil {
		log.Fatalf("load config from efivars: %v", err)
	}

	batch, err := generator.Load(*genfile, cfg)
	if err != nil {
		log.Fatalf("load generator file: %v", err)
	}

	err = batch.Write()
	if err != nil {
		log.Fatalf("extract files from generated batch: %v", err)
	}
}
