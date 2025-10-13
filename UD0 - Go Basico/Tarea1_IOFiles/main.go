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

// Función para leer inventario desde archivo
func leerInventario(nombre string) ([]Producto, error) {
	// Leer el archivo de texto
	data, err := os.ReadFile("archivos/" + nombre)

	if err != nil {
		log.Fatal(err)
	}

	datosComoString := string(data)
	fmt.Println(datosComoString)
	print("\n")

	// Parsear CSV
	lines := strings.Split(datosComoString, "\n")

	// Crear slice de productos
	var productos []Producto

	// Iterar sobre las líneas (saltando la primera que es el encabezado)
	for i, line := range lines {
		// Saltar encabezado
		if i == 0 {
			continue
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		campos := strings.Split(line, ",")
		if len(campos) != 5 {
			continue // línea inválida, saltar
		}

		precio, err := strconv.ParseFloat(campos[3], 64)
		if err != nil {
			continue // si falla parseo, saltar línea
		}

		stock, err := strconv.Atoi(campos[4])
		if err != nil {
			continue
		}

		// Crear producto y agregar al slice
		p := Producto{
			ID:        campos[0],
			Nombre:    campos[1],
			Categoria: campos[2],
			Precio:    precio,
			Stock:     stock,
		}
		productos = append(productos, p)
	}

	return productos, nil
}

// Función para leer transacciones desde archivo
func leerTransacciones(nombre string) ([]Transaccion, error) {
	// Leer el archivo de texto
	data, err := os.ReadFile("archivos/" + nombre)
	if err != nil {
		return nil, err // Mejor retornar el error para que se maneje afuera
	}

	// Mostrar contenido del archivo
	datosComoString := string(data)
	fmt.Println(datosComoString)
	fmt.Println()

	// Parsear CSV
	lines := strings.Split(datosComoString, "\n")

	// Crear slice de transacciones
	var transacciones []Transaccion
	// Iterar sobre las líneas (saltando la primera que es el encabezado)
	for i, lineas := range lines {
		// Saltar encabezado
		if i == 0 {
			continue
		}

		lineas = strings.TrimSpace(lineas)
		if lineas == "" {
			continue
		}

		campos := strings.Split(lineas, ",")
		if len(campos) != 4 {
			continue // línea inválida, saltar
		}

		cantidad, err := strconv.Atoi(campos[2])
		if err != nil {
			continue
		}

		// Crear transacción y agregar al slice
		t := Transaccion{
			Tipo:       strings.ToUpper(campos[0]),
			IDProducto: campos[1],
			Cantidad:   cantidad,
			Fecha:      campos[3],
		}
		transacciones = append(transacciones, t)
	}

	return transacciones, nil
}

// Función para registrar errores con timestamp
func logError(fecha string, mensaje string) string {
	// Agregar timestamp actual para registro
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	return fmt.Sprintf("[%s] ERROR: %s en fecha %s", timestamp, mensaje, fecha)
}

// Procesar transacciones y actualizar inventario
func procesarTransacciones(productos []Producto, transacciones []Transaccion) []string {
	var errores []string

	// Mapa para acceso rápido a productos por ID
	m := make(map[string]*Producto)
	for i := range productos {
		m[productos[i].ID] = &productos[i]
	}

	for _, t := range transacciones {
		p, existe := m[t.IDProducto]
		if !existe {
			errores = append(errores, logError(t.Fecha, fmt.Sprintf("Producto %s no encontrado en transacción de tipo %s", t.IDProducto, t.Tipo)))
			continue
		}

		switch t.Tipo {
		case "VENTA":
			if p.Stock < t.Cantidad {
				errores = append(errores, logError(t.Fecha, fmt.Sprintf("Stock insuficiente para venta. Producto: %s, Stock actual: %d, Cantidad solicitada: %d", p.ID, p.Stock, t.Cantidad)))
			} else {
				p.Stock -= t.Cantidad
			}
		case "COMPRA":
			p.Stock += t.Cantidad
		case "DEVOLUCION":
			p.Stock += t.Cantidad
		default:
			errores = append(errores, logError(t.Fecha, fmt.Sprintf("Tipo de transacción desconocido: %s para producto %s", t.Tipo, t.IDProducto)))
		}
	}

	// Actualizar productos con valores modificados
	for i := range productos {
		productos[i] = *m[productos[i].ID]
	}

	return errores
}

// Escribir inventario actualizado a archivo
func escribirInventario(productos []Producto, nombre string) error {
	// Crear o truncar el archivo
	data := []byte("ID,Nombre,Categoría,Precio,Stock\n")
	for _, p := range productos {
		line := fmt.Sprintf("%s,%s,%s,%.2f,%d\n",
			p.ID,
			p.Nombre,
			p.Categoria,
			p.Precio,
			p.Stock,
		)
		data = append(data, []byte(line)...)
	}
	return os.WriteFile("archivos/"+nombre, data, 0644)
}

func generarReporteBajoStock(productos []Producto, limite int) error {
	data := []byte("ALERTA: PRODUCTOS CON BAJO STOCK\n")
	data = append(data, []byte("================================\n")...)

	for _, p := range productos {
		if p.Stock < limite {
			line := fmt.Sprintf("ID: %s | %s | Stock actual: %d unidades\n", p.ID, p.Nombre, p.Stock)
			data = append(data, []byte(line)...)
		}
	}

	return os.WriteFile("archivos/productos_bajo_stock.txt", data, 0644)
}

func escribirLog(errores []string, nombreArchivo string) error {
	err := os.WriteFile("archivos/"+ nombreArchivo, []byte(strings.Join(errores, "\n")), 0644)
	return err
}

func main() {
	// Leer inventario inicial
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
	errores := procesarTransacciones(productos, transacciones)

	// Escribir inventario actualizado
	err = escribirInventario(productos, "inventario_actualizado.txt")
	if err != nil {
		fmt.Println("Error escribiendo inventario actualizado:", err)
		return
	}

	// Generar reporte de productos con bajo stock (menos de 10)
	err = generarReporteBajoStock(productos, 10)
	if err != nil {
		fmt.Println("Error generando reporte bajo stock:", err)
		return
	}

	// Escribir log de errores
	err = escribirLog(errores, "errores.log")
	if err != nil {
		fmt.Println("Error escribiendo log de errores:", err)
		return
	}

	fmt.Println("Proceso finalizado correctamente.")
}
