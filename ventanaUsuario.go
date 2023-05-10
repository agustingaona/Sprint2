package main

import (
	"database/sql"
	"fmt"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
	_ "github.com/go-sql-driver/mysql"
)

func estaEn(arr []string, pal string) bool {
	for _, v := range arr {
		fmt.Println("sisis")
		if v == pal {
			return true
		}
	}
	return false
}

var tipos []string
var productos [][]string

func actualizarProductos(query *sql.Rows) {
	productos = [][]string{}
	for query.Next() {
		p := new(Producto)
		query.Scan(&p.id_producto, &p.num_maquina, &p.producto, &p.marca, &p.tipo, &p.precio, &p.sinTAC)
		var sinTAC_str string
		if p.sinTAC == 1 {
			sinTAC_str = "Si"
		} else {
			sinTAC_str = "No"
		}
		aux := []string{strconv.Itoa(p.num_maquina), p.producto, p.marca, p.tipo, strconv.Itoa(p.precio), sinTAC_str}
		productos = append(productos, aux)
	}
}

func (v *Ventana) ventanaUsuario(base *sql.DB, usuario *Usuario) {
	tipos = append(tipos, "Todos")

	prod, err := base.Query("SELECT * FROM productos ORDER BY precio;")
	actualizarProductos(prod)
	if err != nil {
		fmt.Println("error al pedir productos")
	}
	aux, _ := base.Query("SELECT * FROM productos ORDER BY precio;")
	for aux.Next() {
		p := new(Producto)
		aux.Scan(&p.id_producto, &p.num_maquina, &p.producto, &p.marca, &p.tipo, &p.precio, &p.sinTAC)
		if !estaEn(tipos, p.tipo) {
			tipos = append(tipos, p.tipo)
		}
	}

	v.window = fyne.CurrentApp().NewWindow("Usuario")
	v.window.Resize(fyne.NewSize(900, 600))
	v.window.CenterOnScreen()

	label1 := widget.NewLabel("Bienvenido " + usuario.nombre)
	label2 := widget.NewLabel("Saldo disponible: " + strconv.Itoa(usuario.saldo))
	campoNumProd := widget.NewEntry()
	campoNumProd.SetPlaceHolder("Ingrese numero de producto")
	botonCompra := widget.NewButton("Comprar", func() {
		numMaquina := campoNumProd.Text
		var id int
		var precio int
		prod := base.QueryRow("SELECT id_producto, precio FROM productos WHERE numero = " + numMaquina + ";")
		prod.Scan(&id, &precio)
		fmt.Println(id, precio, numMaquina)
		base.Query("INSERT INTO transacciones(id_usuario, id_producto, monto) VALUES (" + strconv.Itoa(usuario.id_usuario) + ", " + strconv.Itoa(id) + ", " + strconv.Itoa(precio) + ");")
		base.Query("UPDATE usuarios SET saldo = saldo - " + strconv.Itoa(precio) + " WHERE id_usuario = " + strconv.Itoa(usuario.id_usuario) + ";")
		fmt.Println("comprado")
		usuario.saldo = usuario.saldo - precio
		label2.Refresh()
		campoNumProd.SetText("")
	})

	tipo_label := widget.NewLabel("Tipo: ")
	tipo := widget.NewSelect(tipos, func(s string) {})
	tipo.SetSelected("Todos")
	tac := widget.NewCheck("Sin T.A.C.", func(b bool) {})
	productosTable := widget.NewTable(func() (int, int) { return len(productos), len(productos[0]) },
		func() fyne.CanvasObject { aux := widget.NewLabel(""); return aux },
		func(tci widget.TableCellID, co fyne.CanvasObject) {
			co.(*widget.Label).SetText(productos[tci.Row][tci.Col])
		})
	headers := []string{"Numero", "Producto", "Marca", "Tipo", "Precio", "Sin T.A.C."}
	headerTabla := widget.NewTable(func() (int, int) { return 1, 6 }, func() fyne.CanvasObject { return widget.NewLabel("") }, func(tci widget.TableCellID, co fyne.CanvasObject) { co.(*widget.Label).SetText(headers[tci.Col]) })
	productosTable.SetColumnWidth(0, 70)
	headerTabla.SetColumnWidth(0, 70)
	productosTable.SetColumnWidth(1, 250)
	headerTabla.SetColumnWidth(1, 250)
	productosTable.SetColumnWidth(2, 150)
	headerTabla.SetColumnWidth(2, 150)
	productosTable.SetColumnWidth(3, 150)
	headerTabla.SetColumnWidth(3, 150)
	productosTable.SetColumnWidth(4, 150)
	headerTabla.SetColumnWidth(4, 150)
	productosTable.SetColumnWidth(5, 150)
	headerTabla.SetColumnWidth(5, 150)

	actualiza := widget.NewButton("Actualizar", func() {
		if tipo.Selected == "Todos" && !tac.Checked {
			query, _ := base.Query("SELECT * FROM productos ORDER BY precio;")
			actualizarProductos(query)
		} else if tipo.Selected != "Todos" && !tac.Checked {
			query, _ := base.Query("SELECT * FROM productos WHERE tipo = '" + tipo.Selected + "' ORDER BY precio;")
			actualizarProductos(query)
		} else if tipo.Selected == "Todos" && tac.Checked {
			query, _ := base.Query("SELECT * FROM productos WHERE sin_tac = 1 ORDER BY precio;")
			actualizarProductos(query)
		} else {
			query, _ := base.Query("SELECT * FROM productos WHERE tipo = '" + tipo.Selected + "' AND sin_tac = 1 ORDER BY precio")
			actualizarProductos(query)
		}
		productosTable.Refresh()
	})
	filtros := container.NewHBox(tipo_label, tipo, tac, actualiza)
	hbox1 := container.NewGridWithColumns(4, label1, label2, campoNumProd, botonCompra)
	hbox2 := container.NewHBox(filtros)
	hboxTable := container.NewHBox(productosTable)
	hboxTable.Layout = layout.NewCenterLayout()
	hboxTable.Layout = layout.NewMaxLayout()
	vbox := container.NewVBox(hbox1, hbox2, headerTabla, hboxTable)

	content := container.NewCenter(vbox)
	content.Layout = layout.NewPaddedLayout()
	v.window.SetContent(content)
	v.window.Show()

}
