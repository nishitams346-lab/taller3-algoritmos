package colas

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"
)

// 1. nodo

type nodo struct {
	ts   int64
	next *nodo
}

// 2. cola fifo

type Cola struct {
	head *nodo
	tail *nodo
	len  int
}

// insertar al final de la cola
func (c *Cola) Enqueue(ts int64) {
	n := &nodo{ts: ts}
	if c.tail != nil {
		c.tail.next = n
	}
	c.tail = n
	if c.head == nil {
		c.head = n
	}
	c.len++
}

// extraer del frente de la cola
func (c *Cola) Dequeue() (int64, bool) {
	if c.head == nil {
		return 0, false
	}
	ts := c.head.ts
	c.head = c.head.next
	if c.head == nil {
		c.tail = nil
	}
	c.len--
	return ts, true
}

// ver el elemento del frente sin sacarlo
func (c *Cola) Front() (int64, bool) {
	if c.head == nil {
		return 0, false
	}
	return c.head.ts, true
}

// obtener cantidad de elementos
func (c *Cola) Len() int {
	return c.len
}

// 3. rate limiter

// controlar limite de peticiones con ventana deslizante
func PermitirPeticion(colas map[string]*Cola, ip string, ts int64, M int, T int64) bool {
	c, existe := colas[ip]
	if !existe {
		c = &Cola{}
		colas[ip] = c
	}

	// eliminar registros fuera de la ventana de tiempo
	limite := ts - T
	for {
		frente, ok := c.Front()
		if !ok || frente > limite {
			break
		}
		c.Dequeue()
	}

	// validar cupo disponible
	if c.Len() >= M {
		return false
	}

	// registrar nueva peticion
	c.Enqueue(ts)
	return true
}

// 4. parseo de linea

type Registro struct {
	IP string
	TS int64
}

// extraer ip y timestamp de una linea de log apache
func ParsearLinea(linea string) (ip string, ts int64, err error) {
	partes := strings.Fields(linea)
	if len(partes) < 4 {
		return "", 0, fmt.Errorf("linea invalida: %q", linea)
	}

	ip = partes[0]

	// limpiar y procesar fecha
	fechaRaw := strings.TrimPrefix(partes[3], "[")
	t, errT := time.Parse("02/Jan/2006:15:04:05", fechaRaw)
	if errT != nil {
		return "", 0, fmt.Errorf("timestamp invalido %q: %w", fechaRaw, errT)
	}

	return ip, t.Unix(), nil
}

// 5. procesamiento del log completo

type Resultado struct {
	TotalPeticiones int
	TotalRechazos   int
	RechazosPorIP   map[string]int
}

// procesar archivo, aplicar limites y generar resumen
func ProcesarLog(ruta string, M int, T int64, muestra int) (Resultado, error) {
	f, err := os.Open(ruta)
	if err != nil {
		return Resultado{}, fmt.Errorf("no se pudo abrir el log: %w", err)
	}
	defer f.Close()

	colasIP := make(map[string]*Cola)
	rechazosPorIP := make(map[string]int)
	totalPeticiones, totalRechazos, mostradas := 0, 0, 0

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		linea := scanner.Text()
		ip, ts, err := ParsearLinea(linea)
		if err != nil {
			continue
		}

		totalPeticiones++
		aceptada := PermitirPeticion(colasIP, ip, ts, M, T)

		if !aceptada {
			totalRechazos++
			rechazosPorIP[ip]++
		}

		// mostrar las primeras decisiones en consola
		if mostradas < muestra {
			estado := "ACEPTADA "
			if !aceptada {
				estado = "RECHAZADA"
			}
			fmt.Printf("IP: %-15s  ts: %d  → %s\n", ip, ts, estado)
			mostradas++
		}
	}

	if err := scanner.Err(); err != nil {
		return Resultado{}, fmt.Errorf("error leyendo log: %w", err)
	}

	return Resultado{
		TotalPeticiones: totalPeticiones,
		TotalRechazos:   totalRechazos,
		RechazosPorIP:   rechazosPorIP,
	}, nil
}
