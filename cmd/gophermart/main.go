package main

import "github.com/IgorViskov/33_loyalty/internal/app"

func main() {
	app.Create().
		Configure().
		ApplyMigrations().
		Build().
		Start()
}
