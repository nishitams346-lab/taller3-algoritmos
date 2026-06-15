package pilas

import (
	"fmt"
	"math/rand"
	"os"
	"testing"
)

// ── Stack ────────────────────────────────────────────────────────────────────

func TestPushPop(t *testing.T) {
	var s Stack[int]
	s.Push(10)
	s.Push(20)
	v, ok := s.Pop()
	if !ok || v != 20 {
		t.Error("Pop incorrecto")
	}
}

func TestPilaVacia(t *testing.T) {
	var s Stack[int]
	_, ok := s.Pop()
	if ok {
		t.Error("Deberia estar vacia")
	}
}

// caso limite: Peek en pila vacia
func TestPeekVacia(t *testing.T) {
	var s Stack[int]
	_, ok := s.Peek()
	if ok {
		t.Error("Peek en pila vacia debe retornar false")
	}
}

// caso limite: un solo elemento
func TestUnElemento(t *testing.T) {
	var s Stack[int]
	s.Push(99)
	v, ok := s.Peek()
	if !ok || v != 99 {
		t.Errorf("Peek: esperaba 99, got %d", v)
	}
	v, ok = s.Pop()
	if !ok || v != 99 {
		t.Errorf("Pop: esperaba 99, got %d", v)
	}
	if !s.IsEmpty() {
		t.Error("Pila debe quedar vacia")
	}
}

// ── StockSpan ─────────────────────────────────────────────────────────────────

func TestStockSpan(t *testing.T) {
	precios := []float64{100, 80, 60, 70, 60, 75, 85}
	esperado := []int{1, 1, 1, 2, 1, 4, 6}
	resultado := CalcularStockSpan(precios)
	for i := range esperado {
		if esperado[i] != resultado[i] {
			t.Errorf("Esperado %d obtenido %d", esperado[i], resultado[i])
		}
	}
}

// caso limite: un solo precio -> span = 1
func TestSpanUnElemento(t *testing.T) {
	spans := CalcularStockSpan([]float64{50.0})
	if len(spans) != 1 || spans[0] != 1 {
		t.Errorf("Esperaba [1], got %v", spans)
	}
}

// caso limite: precios siempre crecientes -> span[i] = i+1
func TestSpanCrecientes(t *testing.T) {
	precios := []float64{10, 20, 30, 40, 50}
	spans := CalcularStockSpan(precios)
	for i, s := range spans {
		if s != i+1 {
			t.Errorf("dia %d: esperaba %d, got %d", i, i+1, s)
		}
	}
}

// caso limite: precios siempre decrecientes -> span siempre 1
func TestSpanDecrecientes(t *testing.T) {
	precios := []float64{50, 40, 30, 20, 10}
	for i, s := range CalcularStockSpan(precios) {
		if s != 1 {
			t.Errorf("dia %d: esperaba 1, got %d", i, s)
		}
	}
}

// caso de error: span nunca puede ser < 1 ni > i+1
func TestSpanInvariante(t *testing.T) {
	n := 10_000
	precios := make([]float64, n)
	rng := rand.New(rand.NewSource(42))
	for i := range precios {
		precios[i] = rng.Float64() * 1000
	}
	for i, s := range CalcularStockSpan(precios) {
		if s < 1 || s > i+1 {
			t.Errorf("span invalido en dia %d: %d", i, s)
		}
	}
}

// ── LeerPrecios ───────────────────────────────────────────────────────────────

// caso normal: archivo valido
func TestLeerPrecios_Valido(t *testing.T) {
	contenido := "Date,Open,High,Low,Close,Volume,OpenInt\n" +
		"2000-01-03,88.0,89.0,87.0,88.2,5000,0\n" +
		"2000-01-04,88.2,90.0,87.0,89.5,6000,0\n"

	f, _ := os.CreateTemp("", "precios_*.csv")
	defer os.Remove(f.Name())
	f.WriteString(contenido)
	f.Close()

	registros, err := LeerPrecios(f.Name())
	if err != nil {
		t.Fatalf("error inesperado: %v", err)
	}
	if len(registros) != 2 {
		t.Fatalf("esperaba 2 registros, got %d", len(registros))
	}
	if registros[0].Close != 88.2 {
		t.Errorf("Close[0]: esperaba 88.2, got %v", registros[0].Close)
	}
}

// Caso de error: archivo inexistente
func TestLeerPrecios_Inexistente(t *testing.T) {
	_, err := LeerPrecios("/no/existe.csv")
	if err == nil {
		t.Error("Deberia retornar error con archivo inexistente")
	}
}

// ── Benchmarks ────────────────────────────────────────────────────────────────

func BenchmarkCalcularStockSpan(b *testing.B) {
	for _, n := range []int{1_000, 10_000, 100_000, 1_000_000} {
		precios := make([]float64, n)
		rng := rand.New(rand.NewSource(7))
		for i := range precios {
			precios[i] = rng.Float64() * 500
		}
		b.Run(fmt.Sprintf("n=%d", n), func(b *testing.B) {
			b.ReportAllocs()
			for i := 0; i < b.N; i++ {
				CalcularStockSpan(precios)
			}
		})
	}
}

func BenchmarkStackPushPop(b *testing.B) {
	b.ReportAllocs()
	var s Stack[int]
	for i := 0; i < b.N; i++ {
		s.Push(i)
		s.Pop()
	}
}
