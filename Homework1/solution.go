package main

func Repeater(s, sep string) func(int) string {
	return func(count int) string {
		result := ""
		for i := 0; i < count; i++ {
			result += s
			if i != count-1 {
				result += sep
			}
		}
		return result
	}
}

func Generator(gen func(int) int, initial int) func() int {
	return func() int {
		defer func() {
			initial = gen(initial)
		}()
		return initial
	}
}

func MapReducer(mapper func(int) int, reducer func(int, int) int, initial int) func(...int) int {
	return func(values ...int) int {
		initial := initial
		for _, v := range values {
			initial = reducer(initial, mapper(v))
		}
		return initial
	}
}
