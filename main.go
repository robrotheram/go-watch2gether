package main

import (
	"flag"
	"watch2gether/pkg"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})
	var addr = flag.String("addr", ":8080", "The addr of the  application.")
	flag.Parse() // parse the flags

	log.Println("Starting web server on", *addr)

	go pkg.CleanUP()
	if err := pkg.StartServer(*addr); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
