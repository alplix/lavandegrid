package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/alplix/lavandegrid/internal/ui"
)

func main() {
	a := app.NewWithID("dev.alplix.lavandegrid")
	ui.Run(a)
}
