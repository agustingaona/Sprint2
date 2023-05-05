package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
)

type Ventana struct {
	window fyne.Window
}

func (v *Ventana) ventanaUsuario() {
	v.window = fyne.CurrentApp().NewWindow("Usuario")

	label := widget.NewLabel("Creado wow.")
	content := container.NewCenter(label)

	v.window.SetContent(content)
	v.window.Resize(fyne.NewSize(700, 800))
	v.window.CenterOnScreen()
	v.window.Show()

}
