package main

import (
	"log"

	"github.com/MagicalTux/n-so/nogl"
)

func main() {
	w, err := nogl.New()
	if err != nil {
		log.Printf("error: %s", err)
		return
	}

	w.Wait()
}
