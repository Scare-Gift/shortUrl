package main

import "shorturl/go/src/router"

func main() {
	r := router.Router()
	err := r.Run(":9900")
	if err != nil {
		return
	}
}
