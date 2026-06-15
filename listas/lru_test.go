package listas

import (
	"os"
	"reflect"
	"testing"
)

// TestNewLRU verifica que la función NewLRU inicializa correctamente la estructura LRU.
func TestNewLRU(t *testing.T) {
	lru := NewLRU(3)

	if lru.cap != 3 {
		t.Errorf("Capacidad esperada 3, obtenida %d", lru.cap)
	}

	if lru.mapa == nil {
		t.Error("Mapa cache no debería ser nil")
	}

	if lru.head != nil {
		t.Error("Head debería ser nil al inicio")
	}

	if lru.tail != nil {
		t.Error("Tail debería ser nil al inicio")
	}
}

// TestCacheHit verifica que Get devuelve el valor correcto y actualiza el orden de los nodos.
func TestCacheHit(t *testing.T) {
	lru := NewLRU(3)
	lru.Put(1, 100)

	valor, ok := lru.Get(1)
	if !ok || valor != 100 {
		t.Errorf("Get(1) esperaba 100, obtuvo %d (ok=%v)", valor, ok)
	}

	if len(lru.mapa) != 1 {
		t.Errorf("Mapa cache debería tener tamaño 1, tiene %d", len(lru.mapa))
	}

	if lru.head == nil || lru.tail == nil {
		t.Error("Head y Tail no deberían ser nil después de agregar un elemento")
	} else if lru.head != lru.tail {
		t.Error("Head y Tail deberían apuntar al mismo nodo cuando hay un solo elemento")
	}
}

// TestUpdateExistingKey verifica que actualizar el valor de una clave existente funciona correctamente.
func TestUpdateExistingKey(t *testing.T) {
	lru := NewLRU(3)
	lru.Put(1, 100)
	lru.Put(1, 200)

	valor, ok := lru.Get(1)
	if !ok || valor != 200 {
		t.Errorf("Get(1) esperaba 200 después de actualización, obtuvo %d (ok=%v)", valor, ok)
	}

	if len(lru.mapa) != 1 {
		t.Errorf("Mapa cache debería tener tamaño 1 después de actualización, tiene %d", len(lru.mapa))
	}
}

// TestEviccionLRU verifica que el elemento menos recientemente usado es evictado cuando se alcanza la capacidad máxima.
func TestEviccionLRU(t *testing.T) {
	lru := NewLRU(2)
	lru.Put(1, 10)
	lru.Put(2, 20)

	if _, ok := lru.Get(1); !ok {
		t.Error("Get(1) debería devolver true")
	}

	lru.Put(3, 30)

	if _, ok := lru.Get(2); ok {
		t.Error("Get(2) debería devolver false porque fue evictado")
	}

	if valor, ok := lru.Get(1); !ok || valor != 10 {
		t.Errorf("Get(1) esperaba 10, obtuvo %d (ok=%v)", valor, ok)
	}

	if valor, ok := lru.Get(3); !ok || valor != 30 {
		t.Errorf("Get(3) esperaba 30, obtuvo %d (ok=%v)", valor, ok)
	}
}

// TestCacheMiss verifica que Get devuelva false y el valor por defecto cuando la clave no existe.
func TestCacheMiss(t *testing.T) {
	lru := NewLRU(2)

	valor, ok := lru.Get(999)
	if ok {
		t.Error("Get(999) debería devolver false porque la clave no existe")
	}
	if valor != 0 {
		t.Errorf("Get(999) debería devolver 0 como valor por defecto, obtuvo %d", valor)
	}
}

// TestPutUpdateKeyNoHead verifica que actualizar una clave que no es la cabeza funciona correctamente.
func TestPutUpdateKeyNoHead(t *testing.T) {
    lru := NewLRU(3)
    lru.Put(1, 10)
    lru.Put(2, 20)

    lru.Put(1, 99)

    valor, ok := lru.Get(1)
    if !ok || valor != 99 {
        t.Errorf("Se esperaba valor 99, obtenido %d", valor)
    }
    if lru.head.clave != 1 {
        t.Error("El nodo 1 debió moverse al frente tras ser actualizado")
    }
}

// TestGetItemIsHead verifica que Get de un item lo mueve a la cabeza de la lista.
func TestGetItemIsHead(t *testing.T) {
    lru := NewLRU(3)
    lru.Put(1, 10)

    valor, ok := lru.Get(1)
    if !ok || valor != 10 {
        t.Errorf("Esperado 10, obtenido %d", valor)
    }
}

// TestCargarSecuencia verifica que la función cargue y ordene correctamente los IDs desde un archivo CSV.
func TestCargarSecuencia(t *testing.T) {
    contenidoCSV := `userId,movieId,rating,timestamp
1,101,5.0,1000000005
1,102,4.0,1000000000
1,103,3.0,1000000010
`
    tempFile, err := os.CreateTemp("", "test_ratings_*.csv") // Crea un archivo temporal para la prueba
    if err != nil {
        t.Fatalf("No se pudo crear el archivo temporal: %v", err)
    }

    defer os.Remove(tempFile.Name()) 

    if _, err := tempFile.Write([]byte(contenidoCSV)); err != nil {
        t.Fatalf("No se pudo escribir en el archivo temporal: %v", err)
    }
    tempFile.Close()

    movieIds, err := CargarSecuencia(tempFile.Name())
    if err != nil {
        t.Fatalf("CargarSecuencia falló: %v", err)
    }

    esperado := []int{102, 101, 103}

    if !reflect.DeepEqual(movieIds, esperado) {
        t.Errorf("Orden de películas incorrecto. Se esperaba %v, pero se obtuvo %v", esperado, movieIds)
    }
}

// Benchmark para medir la velocidad de inserción y evicción
func BenchmarkLRUPut(b *testing.B) {
    lru := NewLRU(500) 
    
    b.ResetTimer() 

    for i := 0; i < b.N; i++ {
        lru.Put(i, i)
    }
}

// Benchmark para medir la velocidad de lectura (Cache Hit)
func BenchmarkLRUGet(b *testing.B) {
    lru := NewLRU(500)
    for i := 0; i < 500; i++ {
        lru.Put(i, i)
    }

    b.ResetTimer()

    for i := 0; i < b.N; i++ {
        lru.Get(i % 500) // El módulo %500 asegura que siempre leamos algo que existe
    }
}