package main

import (
	"archshell/internal/app/bootstrap/config"
	"archshell/internal/pkg/generator"
	"flag"
	"log"
)

func main() {
	efivars := flag.String("efivars", "/sys/firmware/efi/efivars", "efivarfs mount path")
	genfile := flag.String("generator", "/etc/bootstrap.gen", "generator file")
	dryrun := flag.Bool("dry-run", false, "show generated YAML only")

	flag.Parse()

	cfg, err := config.Load(*efivars)
	if err != nil {
		log.Fatalf("load config from efivars: %v", err)
	}

	batch, err := generator.Load(*genfile, cfg, *dryrun)
	if err != nil {
		log.Fatalf("load generator file: %v", err)
	}

	if *dryrun {
		return
	}

	err = batch.Write()
	if err != nil {
		log.Fatalf("extract files from generated batch: %v", err)
	}
}
