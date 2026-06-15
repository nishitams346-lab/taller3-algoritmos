package pilas

func CalcularStockSpan(precios []float64) []int {

	n := len(precios)

	span := make([]int, n)

	var pila Stack[int]

	for i := 0; i < n; i++ {

		for !pila.IsEmpty() {

			top, _ := pila.Peek()

			if precios[top] <= precios[i] {
				pila.Pop()
			} else {
				break
			}
		}

		if pila.IsEmpty() {
			span[i] = i + 1
		} else {

			top, _ := pila.Peek()

			span[i] = i - top
		}

		pila.Push(i)
	}

	return span
}
