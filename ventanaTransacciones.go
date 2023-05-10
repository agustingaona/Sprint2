package main

import (
	"database/sql"
	"fmt"
	"strconv"

	"fyne.io/fyne"
	"fyne.io/fyne/container"
	"fyne.io/fyne/layout"
	"fyne.io/fyne/widget"
)

type Transaccion struct {
	id       int
	usuario  string
	producto string
	monto    int
}

var transacciones [][]string

func actualizarTransacciones(query *sql.Rows) {
	transacciones = [][]string{}
	for query.Next() {
		t := new(Transaccion)
		query.Scan(&t.id, &t.usuario, &t.producto, &t.monto)
		aux := []string{strconv.Itoa(t.id), t.usuario, t.producto, strconv.Itoa(t.monto)}
		transacciones = append(transacciones, aux)
	}
}

func (v *Ventana) ventanaTransacciones(base *sql.DB) {
	query, _ := base.Query("SELECT id_transaccion, ci, producto, monto FROM transacciones t INNER JOIN usuarios u ON t.id_usuario = u.id_usuario INNER JOIN productos p ON t.id_producto = p.id_producto;")
	actualizarTransacciones(query)
	fmt.Println(transacciones)
	v.window = fyne.CurrentApp().NewWindow("Transacciones")
	v.window.Resize(fyne.NewSize(900, 600))
	v.window.CenterOnScreen()

	transaccionesTable := widget.NewTable(func() (int, int) { return len(transacciones), len(transacciones[0]) },
		func() fyne.CanvasObject { aux := widget.NewLabel(""); return aux },
		func(tci widget.TableCellID, co fyne.CanvasObject) {
			co.(*widget.Label).SetText(transacciones[tci.Row][tci.Col])
		})
	headers := []string{"Id transaccion", "Usuario CI", "Producto", "Monto"}
	headerTabla := widget.NewTable(func() (int, int) { return 1, 4 }, func() fyne.CanvasObject { return widget.NewLabel("") }, func(tci widget.TableCellID, co fyne.CanvasObject) { co.(*widget.Label).SetText(headers[tci.Col]) })
	transaccionesTable.SetColumnWidth(0, 150)
	headerTabla.SetColumnWidth(0, 150)
	transaccionesTable.SetColumnWidth(1, 80)
	headerTabla.SetColumnWidth(1, 80)
	transaccionesTable.SetColumnWidth(2, 250)
	headerTabla.SetColumnWidth(2, 250)
	transaccionesTable.SetColumnWidth(3, 150)
	headerTabla.SetColumnWidth(3, 150)

	vbox := container.NewVBox(headerTabla, transaccionesTable)

	content := container.NewCenter(vbox)
	content.Layout = layout.NewPaddedLayout()
	v.window.SetContent(content)
	v.window.Show()
}
