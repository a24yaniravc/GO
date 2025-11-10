package main

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
	"time"
)

// Estructuras
type Producto struct {
	ID        string
	Nombre    string
	Categoria string
	Precio    float64
	Stock     int
}

type Transaccion struct {
	Tipo       string
	IDProducto string
	Cantidad   int
	Fecha      string
}

// Métodos String para logging/reportes
func (p Producto) String() string {
	// Formatear la salida del producto
	return fmt.Sprintf("ID: %s | %s | Stock actual: %d unidades", p.ID, p.Nombre, p.Stock)
}

func (t Transaccion) String() string {
	// Formatea
	return fmt.Sprintf("[%s] %s: Producto %s, Cantidad %d", t.Fecha, t.Tipo, t.IDProducto, t.Cantidad)
}

// Función para leer inventario desde archivo
func leerInventario(nombre string) ([]Producto, error) {
	// Leer archivo
	data, err := os.ReadFile("archivos/" + nombre)
	if err != nil {
		return nil, err
	}

	// Procesar líneas
	lines := strings.Split(string(data), "\n")
	var productos []Producto

	// Saltar encabezado y procesar cada línea
	for i, line := range lines {
		if i == 0 {
			continue
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		campos := strings.Split(line, ",")
		if len(campos) != 5 {
			continue
		}

		// Convertir precio y stock
		precio, err := strconv.ParseFloat(campos[3], 64)
		if err != nil {
			continue
		}

		stock, err := strconv.Atoi(campos[4])
		if err != nil {
			continue
		}

		// Agregar producto al slice
		productos = append(productos, Producto{
			ID:        campos[0],
			Nombre:    campos[1],
			Categoria: campos[2],
			Precio:    precio,
			Stock:     stock,
		})
	}

	return productos, nil
}

// Función para leer transacciones desde archivo
func leerTransacciones(nombre string) ([]Transaccion, error) {
	// Leer archivo
	data, err := os.ReadFile("archivos/" + nombre)
	if err != nil {
		return nil, err
	}

	// Procesar líneas
	lines := strings.Split(string(data), "\n")
	var transacciones []Transaccion

	// Saltar encabezado y procesar cada línea
	for i, line := range lines {
		if i == 0 {
			continue
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		campos := strings.Split(line, ",")
		if len(campos) != 4 {
			continue
		}

		cantidad, err := strconv.Atoi(campos[2])
		if err != nil {
			continue
		}

		// Agregar transacción al slice
		transacciones = append(transacciones, Transaccion{
			Tipo:       strings.ToUpper(campos[0]),
			IDProducto: campos[1],
			Cantidad:   cantidad,
			Fecha:      campos[3],
		})
	}

	return transacciones, nil
}

// Función para registrar errores
func logError(mensaje string) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	log.Printf("[%s] ERROR: %s\n", timestamp, mensaje)
}

// Procesar transacciones y actualizar inventario
func procesarTransacciones(productos []Producto, transacciones []Transaccion) {
	// Crear mapa para acceso rápido a productos por ID
	m := make(map[string]*Producto)
	for i := range productos {
		m[productos[i].ID] = &productos[i]
	}

	// Procesar cada transacción
	for _, t := range transacciones {
		p, existe := m[t.IDProducto]
		if !existe {
			logError(fmt.Sprintf("Producto %s no encontrado en transacción de tipo %s", t.IDProducto, t.Tipo))
			continue
		}

		switch t.Tipo { // Manejar tipos de transacción
		case "VENTA":
			if p.Stock < t.Cantidad {
				logError(fmt.Sprintf("Stock insuficiente para venta. Producto: %s, Stock actual: %d, Cantidad solicitada: %d", p.ID, p.Stock, t.Cantidad))
			} else {
				p.Stock -= t.Cantidad
			}
		case "COMPRA", "DEVOLUCION":
			p.Stock += t.Cantidad
		default:
			logError(fmt.Sprintf("Tipo de transacción desconocido: %s para producto %s", t.Tipo, t.IDProducto))
		}
	}

	// Actualizar slice original con valores modificados
	for i := range productos {
		productos[i] = *m[productos[i].ID]
	}
}

// Escribir inventario actualizado a archivo
func escribirInventario(productos []Producto, nombre string) error {
	// Crear encabezado
	data := []byte("ID,Nombre,Categoría,Precio,Stock\n")
	
	// Agregar cada producto
	for _, p := range productos {
		line := fmt.Sprintf("%s,%s,%s,%.2f,%d\n", p.ID, p.Nombre, p.Categoria, p.Precio, p.Stock)
		data = append(data, []byte(line)...)
	}
	return os.WriteFile("archivos/"+nombre, data, 0644)
}

// Generar reporte de productos con bajo stock
func generarReporteBajoStock(productos []Producto, limite int) error {
	// Crear encabezado del reporte
	data := []byte("ALERTA: PRODUCTOS CON BAJO STOCK\n")
	data = append(data, []byte("================================\n")...)

	// Agregar productos con stock bajo el límite
	bajoStock := 0
	for _, p := range productos {
		if p.Stock < limite {
			data = append(data, []byte(p.String()+"\n")...)
			bajoStock++
		}
	}

	// Agregar total al final del reporte
	data = append(data, []byte(fmt.Sprintf("\nTotal de productos con bajo stock: %d\n", bajoStock))...)
	return os.WriteFile("archivos/productos_bajo_stock.txt", data, 0644)
}

func main() {
	// Configurar logging en archivo
	logFile, err := os.OpenFile("archivos/errores.log", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Println("Error creando log:", err)
		return
	}
	defer logFile.Close()
	log.SetOutput(logFile) // Redirigir logs a archivo
	log.SetFlags(0) // Para que no duplique timestamp, ya que lo ponemos manualmente

	// Leer inventario
	productos, err := leerInventario("inventario.txt")
	if err != nil {
		fmt.Println("Error leyendo inventario:", err)
		return
	}

	// Leer transacciones
	transacciones, err := leerTransacciones("transacciones.txt")
	if err != nil {
		fmt.Println("Error leyendo transacciones:", err)
		return
	}

	// Procesar transacciones
	procesarTransacciones(productos, transacciones)

	// Escribir inventario actualizado
	if err := escribirInventario(productos, "inventario_actualizado.txt"); err != nil {
		fmt.Println("Error escribiendo inventario actualizado:", err)
		return
	}

	// Generar reporte de bajo stock (<10)
	if err := generarReporteBajoStock(productos, 10); err != nil {
		fmt.Println("Error generando reporte bajo stock:", err)
		return
	}

	fmt.Println("Proceso finalizado correctamente.")
}
