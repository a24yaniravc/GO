// Tarea 1: Group anagrams

// Our program will have as imput a list of words. We ned to generate a slice of slices of Strings. Each slide should have anagrams fo the same word.

// Examples:
// Input:
// ["hola", "gato", "peso","alho","toga"]

// Expected output:
// [["hola, "alho"], ["gato, toga"],["peso"]]

package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	words := []string{"hola", "gato", "peso", "alho", "toga"}

	// Creamos un mapa para agrupar los anagramas
	anagramMap := make(map[string][]string)
	
	// Rellenamos el mapa con Clave: Las letras de una palabra ordenadas alfabéticamente y Valor: Un slice de todas las palabras que tienen esas mismas letras
	for _, word := range words {
		// Ordenamos la palabra para usarla como clave
		sortedWord := sortString(word)
		// Añadimos la palabra al slice correspondiente en el mapa
		anagramMap[sortedWord] = append(anagramMap[sortedWord], word)
	}

	// Convertimos el mapa a un slice de slices
	var result [][]string
	for _, anagrams := range anagramMap {
		// Añadimos el slice de anagramas al resultado
		result = append(result, anagrams)
	}

	// Imprimimos el resultado
	fmt.Println("Output:")
	for _, group := range result {
		fmt.Print("[")
		for i, word := range group {
			if i > 0 {
				fmt.Print(", ")
			}
			fmt.Printf("\"%s\"", word)
		}
		fmt.Println("]")
	}
}

// Función para ordenar las letras de una palabra
func sortString(s string) string {
	// Convertimos la palabra en un slice de letras, las ordenamos y luego las unimos de nuevo en una cadena
	letters := strings.Split(s, "")
	sort.Strings(letters)
	return strings.Join(letters, "")
}
