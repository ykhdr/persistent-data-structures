package array

import (
	"fmt"
	"testing"
)

func BenchmarkAppend(b *testing.B) {
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Vector/size_%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				v := NewVector[int]()
				for j := 0; j < size; j++ {
					v = v.Append(j)
				}
			}
		})

		b.Run(fmt.Sprintf("NaiveArray/size_%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				a := NewNaiveArray[int]()
				for j := 0; j < size; j++ {
					a = a.Append(j)
				}
			}
		})

		b.Run(fmt.Sprintf("Slice/size_%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				s := make([]int, 0)
				for j := 0; j < size; j++ {
					s = append(s, j)
				}
			}
		})

		b.Run(fmt.Sprintf("SlicePrealloc/size_%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				s := make([]int, 0, size)
				for j := 0; j < size; j++ {
					s = append(s, j)
				}
			}
		})
	}
}

func BenchmarkGet(b *testing.B) {
	sizes := []int{100, 1000, 10000, 100000}

	for _, size := range sizes {
		vector := NewVector[int]()
		naive := NewNaiveArray[int]()
		slice := make([]int, size)

		for i := 0; i < size; i++ {
			vector = vector.Append(i)
			naive = naive.Append(i)
			slice[i] = i
		}

		mid := size / 2

		b.Run(fmt.Sprintf("Vector/size_%d/middle", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = vector.Get(mid)
			}
		})

		b.Run(fmt.Sprintf("NaiveArray/size_%d/middle", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_, _ = naive.Get(mid)
			}
		})

		b.Run(fmt.Sprintf("Slice/size_%d/middle", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = slice[mid]
			}
		})
	}
}

func BenchmarkSet(b *testing.B) {
	sizes := []int{100, 1000, 10000, 100000}

	for _, size := range sizes {
		vector := NewVector[int]()
		naive := NewNaiveArray[int]()
		slice := make([]int, size)

		for i := 0; i < size; i++ {
			vector = vector.Append(i)
			naive = naive.Append(i)
			slice[i] = i
		}

		mid := size / 2

		b.Run(fmt.Sprintf("Vector/size_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = vector.Set(mid, 999)
			}
		})

		b.Run(fmt.Sprintf("NaiveArray/size_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				_ = naive.Set(mid, 999)
			}
		})

		b.Run(fmt.Sprintf("SliceCopy/size_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				newSlice := make([]int, size)
				copy(newSlice, slice)
				newSlice[mid] = 999
			}
		})
	}
}

func BenchmarkPop(b *testing.B) {
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("Vector/size_%d/single", size), func(b *testing.B) {
			vector := NewVector[int]()
			for i := 0; i < size; i++ {
				vector = vector.Append(i)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = vector.Pop()
			}
		})

		b.Run(fmt.Sprintf("NaiveArray/size_%d/single", size), func(b *testing.B) {
			naive := NewNaiveArray[int]()
			for i := 0; i < size; i++ {
				naive = naive.Append(i)
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				_, _, _ = naive.Pop()
			}
		})

		b.Run(fmt.Sprintf("SliceCopy/size_%d/single", size), func(b *testing.B) {
			slice := make([]int, size)
			for i := 0; i < size; i++ {
				slice[i] = i
			}

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				newSlice := make([]int, len(slice)-1)
				copy(newSlice, slice[:len(slice)-1])
			}
		})
	}
}

func BenchmarkIterate(b *testing.B) {
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		vector := NewVector[int]()
		naive := NewNaiveArray[int]()
		slice := make([]int, size)

		for i := 0; i < size; i++ {
			vector = vector.Append(i)
			naive = naive.Append(i)
			slice[i] = i
		}

		b.Run(fmt.Sprintf("Vector/size_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, val := range vector.All() {
					_ = val
				}
			}
		})

		b.Run(fmt.Sprintf("NaiveArray/size_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, val := range naive.All() {
					_ = val
				}
			}
		})

		b.Run(fmt.Sprintf("Slice/size_%d", size), func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				for _, val := range slice {
					_ = val
				}
			}
		})
	}
}

func BenchmarkVersionCreation(b *testing.B) {
	sizes := []int{1000, 10000}
	numVersions := []int{10, 100}

	for _, size := range sizes {
		vector := NewVector[int]()
		naive := NewNaiveArray[int]()

		for i := 0; i < size; i++ {
			vector = vector.Append(i)
			naive = naive.Append(i)
		}

		for _, versions := range numVersions {
			b.Run(fmt.Sprintf("Vector/size_%d/versions_%d", size, versions), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					v := vector
					for j := 0; j < versions; j++ {
						v = v.Set((j*100)%size, 999)
					}
				}
			})

			b.Run(fmt.Sprintf("NaiveArray/size_%d/versions_%d", size, versions), func(b *testing.B) {
				for i := 0; i < b.N; i++ {
					a := naive
					for j := 0; j < versions; j++ {
						a = a.Set((j*100)%size, 999)
					}
				}
			})
		}
	}
}

func BenchmarkMixedOperations(b *testing.B) {
	size := 1000

	b.Run("Vector", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			v := NewVector[int]()

			for j := 0; j < size; j++ {
				v = v.Append(j)
			}

			for j := 0; j < 100; j++ {
				_, _ = v.Get(j * 10)
				v = v.Set(j*10, 999)
			}

			for j := 0; j < 50; j++ {
				v, _, _ = v.Pop()
			}
		}
	})

	b.Run("NaiveArray", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			a := NewNaiveArray[int]()

			for j := 0; j < size; j++ {
				a = a.Append(j)
			}

			for j := 0; j < 100; j++ {
				_, _ = a.Get(j * 10)
				a = a.Set(j*10, 999)
			}

			for j := 0; j < 50; j++ {
				a, _, _ = a.Pop()
			}
		}
	})
}

func BenchmarkMemoryAllocation(b *testing.B) {
	size := 10000

	b.Run("Vector/multipleSet", func(b *testing.B) {
		vector := NewVector[int]()
		for i := 0; i < size; i++ {
			vector = vector.Append(i)
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = vector.Set(size/2, 999)
		}
	})

	b.Run("NaiveArray/multipleSet", func(b *testing.B) {
		naive := NewNaiveArray[int]()
		for i := 0; i < size; i++ {
			naive = naive.Append(i)
		}

		b.ResetTimer()
		b.ReportAllocs()
		for i := 0; i < b.N; i++ {
			_ = naive.Set(size/2, 999)
		}
	})
}
