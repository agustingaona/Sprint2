package main

import (
	"fyne.io/fyne"
)

type Ventana struct {
	window fyne.Window
}

func (v *Ventana) ventanaUsuario() {
	v.window = fyne.CurrentApp().NewWindow("Usuario")
	v.window.Resize(fyne.NewSize(700, 800))
	v.window.CenterOnScreen()
	v.window.Show()

}
