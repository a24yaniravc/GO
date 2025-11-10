package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"os/exec"
	"strconv"
)


// StartAll proporcionada
func StartAll(cmdList []*exec.Cmd) ([]*exec.Cmd, error) {
	for range 20 {
		i, j := rand.Intn(10), rand.Intn(10)
		cmdList[i], cmdList[j] = cmdList[j], cmdList[i]
	}
	for _, cmd := range cmdList {
		err := cmd.Start()
		if err != nil {
			return nil, err
		}
	}
	return cmdList, nil
}

func main() {
	// Si hay argumento, es hijo
	if len(os.Args) > 1 {
		num := os.Args[1]

		// Hijo escribe en un archivo temporal pasado por argumento
		if len(os.Args) > 2 {
			fname := os.Args[2]
			f, err := os.OpenFile(fname, os.O_CREATE|os.O_WRONLY, 0644) // Abrir para escribir
			if err != nil {
				return
			}
			defer f.Close()
			f.WriteString(num + "\n")
		}
		return
	}

	// Padre
	// Eliminar output.txt viejo
	if _, err := os.Stat("output.txt"); err == nil {
		os.Remove("output.txt")
	}

	// Crear lista de hijos y pipes
	cmdList := make([]*exec.Cmd, 10)
	tempFiles := make([]string, 10)

	for i := 0; i < 10; i++ {
		tmp := "tmp_" + strconv.Itoa(i) + ".txt"
		tempFiles[i] = tmp
		cmd := exec.Command(os.Args[0], strconv.Itoa(i), tmp)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmdList[i] = cmd
	}

	// Iniciar todos los procesos con StartAll (código proporcionado)
	cmdList, err := StartAll(cmdList)
	if err != nil {
		log.Fatal("Something went wrong:", err)
	}

	// Esperar que todos terminen
	for i := 0; i < 10; i++ {
		cmdList[i].Wait()
	}

	// Combinar los archivos temporales en orden
	out, err := os.OpenFile("output.txt", os.O_CREATE|os.O_WRONLY, 0644) // Abrir para escribir
	if err != nil {
		log.Fatal(err)
	}
	defer out.Close()

	for i := 0; i < 10; i++ {
		data, err := os.ReadFile(tempFiles[i])
		if err != nil {
			log.Fatal(err)
		}
		out.Write(data)
		os.Remove(tempFiles[i])
	}

	// Leer y mostrar el contenido de output.txt
	// Abrir el archivo
    file, err := os.Open("output.txt")
    if err != nil {
        fmt.Println("Error al abrir el archivo:", err)
        return
    }
    // Asegurarse de cerrarlo al final
    defer file.Close()

    // Leer todo el contenido
    contenido, err := io.ReadAll(file)
    if err != nil {
        fmt.Println("Error al leer el archivo:", err)
        return
    }

    // Imprimir el contenido
    fmt.Println(string(contenido))
}
