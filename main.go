package main

// import "database/sql"
import (
	"database/sql"
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/widget"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/prueba")
	if err != nil {
		fmt.Println("error al conectar.")
		fmt.Println(err.Error())
	}
	defer db.Close()

	a := app.New()
	window := a.NewWindow("Inicio de sesión")

	text := widget.NewLabel("Inicio de sesión")

	campoUsuario := widget.NewEntry()
	campoUsuario.SetPlaceHolder("Usuario: ")

	campoPass := widget.NewEntry()
	campoPass.SetPlaceHolder("Contraseña: ")

	boton := widget.NewButton("Ingresar", func() {
		var user string
		user = campoUsuario.Text
		pass := campoPass.Text
		fmt.Println(user)
		fmt.Println(pass)

		//var admin bool
		var okUser bool
		var okAdmin bool
		if okUser {
			ventana := Ventana{}
			ventana.ventanaUsuario()
		} else if okAdmin {

		}
	})

	vbox := container.NewVBox(text, campoUsuario, campoPass, boton)
	content := container.NewCenter(vbox)

	window.CenterOnScreen()
	window.Resize(fyne.NewSize(350, 250))
	window.SetContent(content)
	window.ShowAndRun()

}
