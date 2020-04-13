package main

import (
	"router"
)

func main() {
	e := router.New()
	e.Logger.Fatal(e.Start(":1323"))
}
