package main

import (
	"database/sql"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

var productosAdmin [][]string

func actualizarProductosAdmin(query *sql.Rows) {
	productosAdmin = [][]string{}
	for query.Next() {
		p := new(Producto)
		query.Scan(&p.id_producto, &p.num_maquina, &p.producto, &p.marca, &p.tipo, &p.precio, &p.sinTAC)
		var sinTAC_str string
		if p.sinTAC == 1 {
			sinTAC_str = "Si"
		} else {
			sinTAC_str = "No"
		}
		aux := []string{strconv.Itoa(p.id_producto), strconv.Itoa(p.num_maquina), p.producto, p.marca, p.tipo, strconv.Itoa(p.precio), sinTAC_str}
		productosAdmin = append(productosAdmin, aux)
	}
}

func (v *Ventana) ventanaAdmin(base *sql.DB, usuario *Usuario) {
	prod, _ := base.Query("SELECT * FROM productos;")
	actualizarProductosAdmin(prod)

	v.window = fyne.CurrentApp().NewWindow("Administrador")
	v.window.Resize(fyne.NewSize(900, 600))
	v.window.CenterOnScreen()

	label1 := widget.NewLabel("Administración de Vending Machine")
	botonReset := widget.NewButton("Restablecer saldos", func() { base.Query("UPDATE usuarios SET saldo = 2000;") })
	campoNumero := widget.NewEntry()
	campoNumero.SetPlaceHolder("Numero en maquina")
	campoProducto := widget.NewEntry()
	campoProducto.SetPlaceHolder("Nombre producto")
	campoMarca := widget.NewEntry()
	campoMarca.SetPlaceHolder("Marca producto")
	campoTipo := widget.NewEntry()
	campoTipo.SetPlaceHolder("Tipo de producto")
	campoPrecio := widget.NewEntry()
	campoPrecio.SetPlaceHolder("Precio")
	campoSinTac := widget.NewEntry()
	campoSinTac.SetPlaceHolder("Sin o con T.A.C. (1 o 0)")
	agreagador := container.NewAdaptiveGrid(6, campoNumero, campoProducto, campoMarca, campoTipo, campoPrecio, campoSinTac)

	productosTable := widget.NewTable(func() (int, int) { return len(productosAdmin), len(productosAdmin[0]) },
		func() fyne.CanvasObject { aux := widget.NewLabel(""); return aux },
		func(tci widget.TableCellID, co fyne.CanvasObject) {
			co.(*widget.Label).SetText(productosAdmin[tci.Row][tci.Col])
		})
	headers := []string{"Id", "Numero", "Producto", "Marca", "Tipo", "Precio", "Sin T.A.C."}
	headerTabla := widget.NewTable(func() (int, int) { return 1, 7 }, func() fyne.CanvasObject { return widget.NewLabel("") }, func(tci widget.TableCellID, co fyne.CanvasObject) { co.(*widget.Label).SetText(headers[tci.Col]) })
	productosTable.SetColumnWidth(0, 20)
	headerTabla.SetColumnWidth(0, 20)
	productosTable.SetColumnWidth(1, 70)
	headerTabla.SetColumnWidth(1, 70)
	productosTable.SetColumnWidth(2, 250)
	headerTabla.SetColumnWidth(2, 250)
	productosTable.SetColumnWidth(3, 150)
	headerTabla.SetColumnWidth(3, 150)
	productosTable.SetColumnWidth(4, 150)
	headerTabla.SetColumnWidth(4, 150)
	productosTable.SetColumnWidth(5, 50)
	headerTabla.SetColumnWidth(5, 50)
	productosTable.SetColumnWidth(6, 150)
	headerTabla.SetColumnWidth(6, 150)

	botonAgregar := widget.NewButton("Agregar producto", func() {
		base.Query("INSERT INTO productos(numero, producto, marca, tipo, precio, sin_tac) VALUES('" + campoNumero.Text + "', '" + campoProducto.Text + "', '" + campoMarca.Text + "', '" + campoTipo.Text + "', " + campoPrecio.Text + ", " + campoSinTac.Text + ");")
		query, _ := base.Query("SELECT * FROM productos;")
		actualizarProductosAdmin(query)
		productosTable.Refresh()
		campoNumero.SetText("")
		campoProducto.SetText("")
		campoMarca.SetText("")
		campoTipo.SetText("")
		campoPrecio.SetText("")
		campoSinTac.SetText("")
	})
	numeroEliminar := widget.NewEntry()
	numeroEliminar.SetPlaceHolder("Id de producto a eliminar")
	botonEliminar := widget.NewButton("Eliminar producto", func() {
		base.Query("DELETE FROM productos WHERE id_producto = " + numeroEliminar.Text + ";")
		query, _ := base.Query("SELECT * FROM productos;")
		actualizarProductosAdmin(query)
		productosTable.Refresh()
		numeroEliminar.SetText("")
	})
	botonTransacciones := widget.NewButton("Ver transacciones", func() {
		ventana := Ventana{}
		ventana.ventanaTransacciones(base)
	})

	box1 := container.NewBorder(nil, nil, label1, botonReset)

	vbox1 := container.NewVBox(box1, headerTabla, productosTable, agreagador, botonAgregar, numeroEliminar, botonEliminar, botonTransacciones)

	content := container.NewCenter(vbox1)
	content.Layout = layout.NewPaddedLayout()
	v.window.SetContent(content)
	v.window.Show()
}
