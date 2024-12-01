package main

import (
	"firestarter/pkg/systemd/notify"
	"flag"
	"log"
	"os"
	"time"

	"k8s.io/utils/inotify"
)

func main() {
	path := flag.String("path", "/etc/default/bootfile", "path to wait for inotify events")
	timeout := flag.Uint("timeout", 30, "timeout seconds")

	flag.Parse()

	watcher, err := inotify.NewWatcher()
	if err != nil {
		log.Fatalf("construct inotify watcher: %v", err)
	}
	defer watcher.Close()

	err = watcher.Watch(*path)
	if err != nil {
		log.Fatalf("assign path to inotify watcher: %v", err)
	}

	err = notify.Send("READY=1")
	if err != nil {
		log.Fatalf("send systemd-notify: %v", err)
	}

	if *timeout != 0 {
		go func() {
			time.Sleep(time.Duration(*timeout) * time.Second)
			os.Exit(0)
		}()

	}

	for {
		select {
		case event := <-watcher.Event:
			log.Printf("watched event: %s", event)

			return
		case err := <-watcher.Error:
			log.Fatalf("watch error occurred: %v", err)
		}
	}
}
