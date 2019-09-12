package main

import (
	"fay_admin/router"
	"github.com/henrylee2cn/faygo"
)

func main() {
	router.Route(faygo.New("fay_admin"))
	faygo.Run()
}
