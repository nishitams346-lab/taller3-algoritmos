package main

import (
	"fmt"
	"os"
	"strconv"
	"taller3/listas"
)

func main() {
	if len(os.Args) < 3 {
		fmt.Println("Uso: go run ./listas/cmd/main.go <ruta_csv> <capacidad_1> <capacidad_2> ...")
		fmt.Println("Ejemplo: go run ./listas/cmd/main.go ./data/ratings.csv 50 100 500 1000")

		os.Exit(1)
	}

	path := os.Args[1]

	ids, err := listas.CargarSecuencia(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error al cargar secuencia:", err)
		os.Exit(1)
	}

	totalHits := len(ids)

	fmt.Println("\n=======================================================")
	fmt.Printf("%-15s %-15s %-10s %-12s\n", "Capacidad", "Total Accesos", "Hits", "Hit Ratio")
	fmt.Println("=======================================================")

	for _, arg := range os.Args[2:] {
		capacidad, err := strconv.Atoi(arg)
		if err != nil || capacidad <= 0 {
			fmt.Fprintln(os.Stderr, "Error: La capacidad debe ser un entero positivo")
			os.Exit(1)
		}

		hits := 0

		lru := listas.NewLRU(capacidad)
		for _, id := range ids {
			if _, ok := lru.Get(id); ok {
				hits++
			} else {
				lru.Put(id, id)
			}
		}

		hitRatio := float64(hits) / float64(totalHits) * 100

		fmt.Printf("%-15d %-15d %-10d %.2f%%\n", capacidad, totalHits, hits, hitRatio)
	}
}
