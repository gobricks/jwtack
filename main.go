package main

import (
	jwtack "github.com/gobricks/jwtack/src"
	app "github.com/gobricks/jwtack/src/app"
)

func main() {
	jwtack.RunServer(app.NewApp())
}