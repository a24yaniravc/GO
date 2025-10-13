package main

import (
	"fmt"
	"unicode/utf8"
)

func main() {
	str1 := "¿Hola mundo? Aquí todo está bien." // 33 caracteres

	// Formas incorrecta
	fmt.Println("La longitud de la cadena (en bytes) es:", len(str1))

	for i := 0; i < len(str1); i++ {
		fmt.Printf("str1[%d]= %c\n", i, str1[i])
	}

	fmt.Println("")

	for i, r := range str1 {
		fmt.Printf("str1[%d]= %c\n", i, r)
	}

	fmt.Println("")

	// Forma correcta
	// Runas (carácteres)
	fmt.Println("La longitud en runas es:", utf8.RuneCountInString(str1))

	// Convertir a slice de runas:
	rs := []rune(str1)

	for i := 0; i < len(rs); i++ {
		fmt.Printf("str1[%d]= %c\n", i, rs[i])
	}
}
