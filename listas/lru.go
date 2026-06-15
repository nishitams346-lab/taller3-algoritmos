package listas

import (
	"cmp"
	"encoding/csv"
	"io"
	"os"
	"slices"
	"strconv"
)

type Nodo struct { 
	clave	int
	valor	int
	prev	*Nodo
	next	*Nodo
}

type LRU struct {
	cap 	int;
	mapa 	map[int]*Nodo;
	head 	*Nodo
	tail 	*Nodo
}

func (l *LRU) Get(clave int) (int, bool) {
	node, ok := l.mapa[clave]

	if !ok {
		return 0, false
	}

	if node != l.head {
		l.eliminar(node)
		l.agregarAlFrente(node)
	}

	return node.valor, true
}

func (l *LRU) Put(clave, valor int) {
	node, ok := l.mapa[clave]

	if ok {
		node.valor = valor
		if node != l.head {
			l.eliminar(node)
			l.agregarAlFrente(node)
		}
	} else {
		if len(l.mapa) >= l.cap {
			delete(l.mapa, l.tail.clave) // Eliminar del mapa
			l.eliminar(l.tail) // Eliminar el nodo menos recientemente usado
		}

		newNode := &Nodo{
			clave: clave,
			valor: valor,
		}
		l.mapa[clave] = newNode
		l.agregarAlFrente(newNode)
	}
}

func NewLRU(capacidad int) *LRU {
	return &LRU{
		cap: capacidad,
		mapa: make(map[int]*Nodo),
	}
}

// Funciones auxiliares para manejar la lista doblemente enlazada
func (l *LRU) eliminar(node *Nodo) {
	if node.prev != nil {
		node.prev.next = node.next
	} else {
		l.head = node.next
	}

	if node.next != nil {
		node.next.prev = node.prev
	} else {
		l.tail = node.prev
	}
}

func (l *LRU) agregarAlFrente(node *Nodo) {
	node.prev = nil // Limpiar el enlace previo del nodo
	if l.head != nil {
		l.head.prev = node
		node.next = l.head
	}
	l.head = node
	if l.tail == nil {
		l.tail = node
	}
}

// Estructura auxiliar para almacenar los registros de rating con su timestamp.
type RegistroRating struct {
	movieId 	int
	timestamp 	int64
}

// Leer el archivo CSV y cargar la secuencia de IDs ordenados por timestamp.
func CargarSecuencia(ruta string) ([]int, error) {
	file, err := os.Open(ruta)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	csvReader := csv.NewReader(file)

	if _, err := csvReader.Read(); err != nil && err != io.EOF {
		return nil, err
	}

	var ratings []RegistroRating
	linea := 1

	for {
		linea++
		record, err := csvReader.Read()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, err
		}

		// Estructura: userId, movieId, rating, timestamp
		// Indices:	   0       1        2       3

		movieId, err := strconv.Atoi(record[1])
		if err != nil {
			return nil, err
		}

		timestamp, err := strconv.ParseInt(record[3], 10, 64) // Tiempo UNIX
		if err != nil {
			return nil, err
		}

		ratings = append(ratings, RegistroRating{
			movieId:   movieId,
			timestamp: timestamp,
		})
	}

	slices.SortFunc(ratings, func(a, b RegistroRating) int {
		return cmp.Compare(a.timestamp, b.timestamp)
	})

	movieIds := make([]int, 0, len(ratings))
	for _, rating := range ratings {
		movieIds = append(movieIds, rating.movieId)
	}

	return movieIds, nil
}
