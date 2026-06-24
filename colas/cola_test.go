package colas

import (
	"testing"
)

// tests de la cola

// test fifo
func TestColaFIFO(t *testing.T) {
	c := &Cola{}
	c.Enqueue(100)
	c.Enqueue(200)
	c.Enqueue(300)

	for _, esperado := range []int64{100, 200, 300} {
		v, ok := c.Dequeue()
		if !ok || v != esperado {
			t.Errorf("esperaba %d, obtuve %d (ok=%v)", esperado, v, ok)
		}
	}
}

// comprueba comportamiento en cola vacia
func TestColaVacia(t *testing.T) {
	c := &Cola{}
	_, ok := c.Dequeue()
	if ok {
		t.Error("dequeue en cola vacia deberia devolver false")
	}
	_, ok = c.Front()
	if ok {
		t.Error("front en cola vacia deberia devolver false")
	}
}

// caso limite cola con un unico elemento
func TestColaUnElemento(t *testing.T) {
	c := &Cola{}
	c.Enqueue(42)
	if c.Len() != 1 {
		t.Errorf("len esperado 1, obtenido %d", c.Len())
	}
	v, ok := c.Dequeue()
	if !ok || v != 42 {
		t.Errorf("esperaba 42, obtuve %d", v)
	}
	if c.Len() != 0 {
		t.Errorf("len esperado 0 tras dequeue, obtenido %d", c.Len())
	}
}

// verifica que front no elimine elementos
func TestFront(t *testing.T) {
	c := &Cola{}
	c.Enqueue(10)
	c.Enqueue(20)

	f, ok := c.Front()
	if !ok || f != 10 {
		t.Errorf("front esperaba 10, obtuvo %d", f)
	}
	if c.Len() != 2 {
		t.Error("front no deberia modificar el tamano de la cola")
	}
}

// tests del rate limiter

// bloquear peticiones que superen el limite maximo (m)
func TestPermitirPeticionLimite(t *testing.T) {
	colasIP := make(map[string]*Cola)
	ip := "1.2.3.4"
	M, T := 3, int64(60)
	ts := int64(1000)

	for i := 0; i < M; i++ {
		if !PermitirPeticion(colasIP, ip, ts, M, T) {
			t.Errorf("peticion %d deberia aceptarse", i+1)
		}
	}
	if PermitirPeticion(colasIP, ip, ts, M, T) {
		t.Error("la peticion m+1 deberia rechazarse")
	}
}

// expirar registros viejos al pasar el tiempo de la ventana (t)
func TestVentanaDeslizante(t *testing.T) {
	colasIP := make(map[string]*Cola)
	ip := "5.6.7.8"
	M, T := 2, int64(10)

	// llena la ventana en t=0
	PermitirPeticion(colasIP, ip, 0, M, T)
	PermitirPeticion(colasIP, ip, 0, M, T)

	// en t=5 sigue llena
	if PermitirPeticion(colasIP, ip, 5, M, T) {
		t.Error("deberia rechazarse en t=5 (ventana llena)")
	}

	// en t=11 expiro t=0, vuelve a aceptar
	if !PermitirPeticion(colasIP, ip, 11, M, T) {
		t.Error("deberia aceptarse en t=11 (timestamps expirados)")
	}
}

// validar que los limites apliquen por separado para cada ip
func TestIPsIndependientes(t *testing.T) {
	colasIP := make(map[string]*Cola)
	M, T := 1, int64(60)
	ts := int64(1000)

	if !PermitirPeticion(colasIP, "a.a.a.a", ts, M, T) {
		t.Error("primera peticion de ip a deberia aceptarse")
	}
	if !PermitirPeticion(colasIP, "b.b.b.b", ts, M, T) {
		t.Error("primera peticion de ip b deberia aceptarse")
	}
	if PermitirPeticion(colasIP, "a.a.a.a", ts, M, T) {
		t.Error("segunda peticion de ip a deberia rechazarse")
	}
}

// -------------------------------------------------------
// tests del parser
// -------------------------------------------------------

// procesar log valido en formato apache
func TestParsearLineaOK(t *testing.T) {
	linea := `83.149.9.216 - - [17/May/2015:10:05:03 +0000] "GET /presentations HTTP/1.1" 200 5678`
	ip, ts, err := ParsearLinea(linea)
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if ip != "83.149.9.216" {
		t.Errorf("ip esperada 83.149.9.216, obtenida %s", ip)
	}
	if ts <= 0 {
		t.Errorf("timestamp debe ser positivo, obtenido %d", ts)
	}
}

// controlar error en logs malformados
func TestParsearLineaInvalida(t *testing.T) {
	_, _, err := ParsearLinea("esto no es un log")
	if err == nil {
		t.Error("se esperaba error con linea invalida")
	}
}

// benchmarks

// rendimiento de insercion y extraccion
func BenchmarkEnqueueDequeue(b *testing.B) {
	c := &Cola{}
	for i := 0; i < b.N; i++ {
		c.Enqueue(int64(i))
		c.Dequeue()
	}
}

// rendimiento del limitador con carga continua
func BenchmarkPermitirPeticion(b *testing.B) {
	colasIP := make(map[string]*Cola)
	M, T := 100, int64(60)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		PermitirPeticion(colasIP, "1.1.1.1", int64(i), M, T)
	}
}
