package main

import (
	"go/upload/app"
)

func main() {
	app.ListenAndServe("127.0.0.1:9091")
}
