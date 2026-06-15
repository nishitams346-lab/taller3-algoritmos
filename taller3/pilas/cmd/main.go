package main

import (
	"fmt"
	"log"
	"taller3/pilas"
)

func main() {
	ruta := "data/xog.us.txt"

	registros, err := pilas.LeerPrecios(ruta)
	if err != nil {
		log.Fatal(err)
	}

	var precios []float64
	for _, r := range registros {
		precios = append(precios, r.Close)
	}

	span := pilas.CalcularStockSpan(precios)

	fmt.Printf("Total de registros: %d\n\n", len(registros))

	// primeros 10
	fmt.Println("── Primeros 10 días ──")
	fmt.Println("Fecha\t\tPrecio\t\tSpan")
	for i := 0; i < len(registros) && i < 10; i++ {
		fmt.Printf("%s\t%.2f\t%d\n", registros[i].Fecha, registros[i].Close, span[i])
	}

	// ultimos 10
	fmt.Println("\n── Últimos 10 días ──")
	fmt.Println("Fecha\t\tPrecio\t\tSpan")
	inicio := len(registros) - 10
	if inicio < 0 {
		inicio = 0
	}
	for i := inicio; i < len(registros); i++ {
		fmt.Printf("%s\t%.2f\t%d\n", registros[i].Fecha, registros[i].Close, span[i])
	}

	// dia con span maximo
	maxSpan, indice := 0, 0
	for i, s := range span {
		if s > maxSpan {
			maxSpan = s
			indice = i
		}
	}
	fmt.Println("\n── Día con span máximo ──")
	fmt.Printf("Fecha: %s\nPrecio: %.2f\nSpan: %d dias\n",
		registros[indice].Fecha, registros[indice].Close, maxSpan)
}
