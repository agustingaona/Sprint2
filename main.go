package main

// import "database/sql"
import (
	"database/sql"
	"fmt"

	"fyne.io/fyne"
	"fyne.io/fyne/app"
	"fyne.io/fyne/container"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	_ "github.com/go-sql-driver/mysql"
)

type Ventana struct {
	window fyne.Window
}

type Usuario struct {
	id_usuario int
	admin      int
	nombre     string
	apellido   string
	ci         string
	contraseña string
	telefono   string
	saldo      int
}

type Producto struct {
	id_producto int
	num_maquina int
	producto    string
	marca       string
	tipo        string
	precio      int
	sinTAC      int
}

func main() {
	db, err := sql.Open("mysql", "root:123456@tcp(localhost:3306)/prueba")
	if err != nil {
		fmt.Println("error al conectar.")
		fmt.Println(err.Error())
	}
	defer db.Close()

	a := app.New()
	window := a.NewWindow("Inicio de sesión")

	window.CenterOnScreen()
	window.Resize(fyne.NewSize(350, 250))

	text := widget.NewLabel("Bienvenido al sistema de Vending Machine")

	campoUsuario := widget.NewEntry()
	campoUsuario.SetPlaceHolder("Usuario: ")

	campoPass := widget.NewPasswordEntry()
	campoPass.SetPlaceHolder("Contraseña: ")

	info := widget.NewLabel("")
	info.Hide()

	boton := widget.NewButton("Ingresar", func() {
		info.Hide()
		user := campoUsuario.Text
		pass := campoPass.Text
		row := db.QueryRow("SELECT * FROM usuarios WHERE ci='" + user + "';")
		var okUser bool
		var okAdmin bool

		u := new(Usuario)
		row.Scan(&u.id_usuario, &u.admin, &u.nombre, &u.apellido, &u.ci, &u.contraseña, &u.telefono, &u.saldo)

		if u.ci != "" {
			if pass == u.contraseña {
				if u.admin == 1 {
					okAdmin = true
				} else {
					okUser = true
				}
			} else {
				info.SetText("Contraseña incorrecta")
				info.Show()
			}

		} else {
			info.SetText("Usuario no existe")
			info.Show()
		}

		//var admin bool
		if okUser {
			ventana := Ventana{}
			ventana.ventanaUsuario(db, u)
		} else if okAdmin {
			ventana := Ventana{}
			ventana.ventanaAdmin(db, u)
		}
		campoUsuario.SetText("")
		campoPass.SetText("")
	})

	vbox := container.NewVBox(text, campoUsuario, campoPass, boton, info)
	content := container.NewCenter(vbox)
	content.Layout = layout.NewPaddedLayout()

	window.SetContent(content)
	window.ShowAndRun()

}
