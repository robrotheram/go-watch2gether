package main

import (
	"flag"
	"log"
	"watch2gether/pkg"
)

func main() {
	var addr = flag.String("addr", ":8080", "The addr of the  application.")
	flag.Parse() // parse the flags

	log.Println("Starting web server on", *addr)

	go pkg.CleanUP()
	if err := pkg.StartServer(*addr); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
