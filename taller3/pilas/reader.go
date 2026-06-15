package pilas

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Registro struct {
	Fecha string
	Close float64
}

// leerPrecios lee el archivo del dataset (formato Kaggle: Date,Open,High,Low,Close,...)
// y devuelve los registros con fecha y precio de cierre. Complejidad: O(n).
func LeerPrecios(ruta string) ([]Registro, error) {
	file, err := os.Open(ruta)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	datos, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}
	if len(datos) < 2 {
		return nil, fmt.Errorf("archivo sin datos")
	}

	// detectar indice de la columna Close desde el encabezado
	colClose := -1
	for i, h := range datos[0] {
		if strings.ToLower(strings.TrimSpace(h)) == "close" {
			colClose = i
			break
		}
	}
	if colClose == -1 {
		return nil, fmt.Errorf("columna Close no encontrada")
	}

	var registros []Registro
	for i := 1; i < len(datos); i++ {
		if len(datos[i]) <= colClose {
			continue
		}
		closePrice, err := strconv.ParseFloat(strings.TrimSpace(datos[i][colClose]), 64)
		if err != nil {
			continue
		}
		registros = append(registros, Registro{
			Fecha: datos[i][0],
			Close: closePrice,
		})
	}
	return registros, nil
}
