package hashmap

import (
	"fmt"
	"math/rand"
	"testing"
)

// Helpers
func buildPersistent(size int) *HashMap[int, int] {
	m := NewHashMap[int, int]()
	for i := 0; i < size; i++ {
		m = m.Set(i, i)
	}
	return m
}

func buildNaive(size int) *NaiveHashMap[int, int] {
	m := NewNaiveHashMap[int, int]()
	for i := 0; i < size; i++ {
		m = m.Set(i, i)
	}
	return m
}

func buildGoMap(size int) map[int]int {
	m := make(map[int]int, size)
	for i := 0; i < size; i++ {
		m[i] = i
	}
	return m
}

func pregenKeys(size int, n int) []int {
	r := rand.New(rand.NewSource(1))
	keys := make([]int, n)
	for i := 0; i < n; i++ {
		keys[i] = r.Intn(size)
	}
	return keys
}

// Benchmarks
func BenchmarkBuild(b *testing.B) {
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		b.Run(fmt.Sprintf("HashMap/size_%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m := NewHashMap[int, int]()
				for j := 0; j < size; j++ {
					m = m.Set(j, j)
				}
			}
		})

		b.Run(fmt.Sprintf("NaiveHashMap/size_%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m := NewNaiveHashMap[int, int]()
				for j := 0; j < size; j++ {
					m = m.Set(j, j)
				}
			}
		})

		b.Run(fmt.Sprintf("GoMap/size_%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				m := make(map[int]int, size)
				for j := 0; j < size; j++ {
					m[j] = j
				}
			}
		})
	}
}

func BenchmarkGet(b *testing.B) {
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		hm := buildPersistent(size)
		naive := buildNaive(size)
		gm := buildGoMap(size)

		keysHit := pregenKeys(size, 1<<16)
		keysMiss := pregenKeys(size, 1<<16)

		b.Run(fmt.Sprintf("HashMap/size_%d/hit", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				k := keysHit[i&(len(keysHit)-1)]
				_, _ = hm.Get(k)
			}
		})

		b.Run(fmt.Sprintf("NaiveHashMap/size_%d/hit", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				k := keysHit[i&(len(keysHit)-1)]
				_, _ = naive.Get(k)
			}
		})

		b.Run(fmt.Sprintf("GoMap/size_%d/hit", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				k := keysHit[i&(len(keysHit)-1)]
				_, _ = gm[k]
			}
		})

		b.Run(fmt.Sprintf("HashMap/size_%d/miss", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				k := keysMiss[i&(len(keysMiss)-1)] + size + 1
				_, _ = hm.Get(k)
			}
		})

		b.Run(fmt.Sprintf("NaiveHashMap/size_%d/miss", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				k := keysMiss[i&(len(keysMiss)-1)] + size + 1
				_, _ = naive.Get(k)
			}
		})

		b.Run(fmt.Sprintf("GoMap/size_%d/miss", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				k := keysMiss[i&(len(keysMiss)-1)] + size + 1
				_, _ = gm[k]
			}
		})
	}
}

func BenchmarkSetUpdate(b *testing.B) {
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		hm := buildPersistent(size)
		naive := buildNaive(size)
		gm := buildGoMap(size)

		keys := pregenKeys(size, 1<<16)

		b.Run(fmt.Sprintf("HashMap/size_%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				k := keys[i&(len(keys)-1)]
				_ = hm.Set(k, 999)
			}
		})

		b.Run(fmt.Sprintf("NaiveHashMap/size_%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				k := keys[i&(len(keys)-1)]
				_ = naive.Set(k, 999)
			}
		})

		b.Run(fmt.Sprintf("GoMap/size_%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				k := keys[i&(len(keys)-1)]
				gm[k] = 999
			}
		})
	}
}

func BenchmarkSetInsert(b *testing.B) {
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		hm := buildPersistent(size)
		naive := buildNaive(size)

		keys := pregenKeys(size, 1<<16)

		b.Run(fmt.Sprintf("HashMap/size_%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				k := keys[i&(len(keys)-1)] + size + 1
				_ = hm.Set(k, 999)
			}
		})

		b.Run(fmt.Sprintf("NaiveHashMap/size_%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				k := keys[i&(len(keys)-1)] + size + 1
				_ = naive.Set(k, 999)
			}
		})

		b.Run(fmt.Sprintf("GoMap/size_%d", size), func(b *testing.B) {
			gm := buildGoMap(size)
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				k := keys[i&(len(keys)-1)] + size + 1
				gm[k] = 999
			}
		})
	}
}

func BenchmarkDelete(b *testing.B) {
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		hm := buildPersistent(size)
		naive := buildNaive(size)

		keys := pregenKeys(size, 1<<16)

		b.Run(fmt.Sprintf("HashMap/size_%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				k := keys[i&(len(keys)-1)]
				_ = hm.Delete(k)
			}
		})

		b.Run(fmt.Sprintf("NaiveHashMap/size_%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				k := keys[i&(len(keys)-1)]
				_ = naive.Delete(k)
			}
		})

		b.Run(fmt.Sprintf("GoMap/size_%d", size), func(b *testing.B) {
			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				gm := buildGoMap(size)
				k := keys[i&(len(keys)-1)]
				delete(gm, k)
			}
		})
	}
}

func BenchmarkIterate(b *testing.B) {
	sizes := []int{100, 1000, 10000}

	for _, size := range sizes {
		hm := buildPersistent(size)
		naive := buildNaive(size)
		gm := buildGoMap(size)

		b.Run(fmt.Sprintf("HashMap/size_%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for k, v := range hm.All() {
					_, _ = k, v
				}
			}
		})

		b.Run(fmt.Sprintf("NaiveHashMap/size_%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for k, v := range naive.data {
					_, _ = k, v
				}
			}
		})

		b.Run(fmt.Sprintf("GoMap/size_%d", size), func(b *testing.B) {
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				for k, v := range gm {
					_, _ = k, v
				}
			}
		})
	}
}

func BenchmarkVersionCreation(b *testing.B) {
	sizes := []int{1000, 10000}
	versionsList := []int{10, 100, 1000}

	for _, size := range sizes {
		base := buildPersistent(size)
		baseNaive := buildNaive(size)
		baseGo := buildGoMap(size)

		for _, versions := range versionsList {
			keys := pregenKeys(size, versions)

			b.Run(fmt.Sprintf("HashMap/size_%d/versions_%d", size, versions), func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					m := base
					for j := 0; j < versions; j++ {
						m = m.Set(keys[j], 999)
					}
				}
			})

			b.Run(fmt.Sprintf("NaiveHashMap/size_%d/versions_%d", size, versions), func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					m := baseNaive
					for j := 0; j < versions; j++ {
						m = m.Set(keys[j], 999)
					}
				}
			})

			b.Run(fmt.Sprintf("GoMapCopy/size_%d/versions_%d", size, versions), func(b *testing.B) {
				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					m := baseGo
					for j := 0; j < versions; j++ {
						next := make(map[int]int, len(m)+1)
						for k, v := range m {
							next[k] = v
						}
						next[keys[j]] = 999
						m = next
					}
				}
			})
		}
	}
}

func BenchmarkMixedOperations(b *testing.B) {
	size := 10000
	keys := pregenKeys(size, 2000)

	b.Run("HashMap", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			m := NewHashMap[int, int]()
			for j := 0; j < size; j++ {
				m = m.Set(j, j)
			}

			for j := 0; j < 1000; j++ {
				k := keys[j]
				_, _ = m.Get(k)
				m = m.Set(k, 999)
			}

			for j := 0; j < 500; j++ {
				m = m.Delete(keys[j])
			}
		}
	})

	b.Run("NaiveHashMap", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			m := NewNaiveHashMap[int, int]()
			for j := 0; j < size; j++ {
				m = m.Set(j, j)
			}

			for j := 0; j < 1000; j++ {
				k := keys[j]
				_, _ = m.Get(k)
				m = m.Set(k, 999)
			}

			for j := 0; j < 500; j++ {
				m = m.Delete(keys[j])
			}
		}
	})

	b.Run("GoMap", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()

		for i := 0; i < b.N; i++ {
			m := make(map[int]int, size)
			for j := 0; j < size; j++ {
				m[j] = j
			}

			for j := 0; j < 1000; j++ {
				k := keys[j]
				_, _ = m[k]
				m[k] = 999
			}

			for j := 0; j < 500; j++ {
				delete(m, keys[j])
			}
		}
	})
}

func BenchmarkMemoryAllocation(b *testing.B) {
	size := 100000
	keys := pregenKeys(size, 1<<16)

	b.Run("HashMap/multipleSet", func(b *testing.B) {
		m := buildPersistent(size)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			k := keys[i&(len(keys)-1)]
			_ = m.Set(k, 999)
		}
	})

	b.Run("NaiveHashMap/multipleSet", func(b *testing.B) {
		m := buildNaive(size)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			k := keys[i&(len(keys)-1)]
			_ = m.Set(k, 999)
		}
	})

	b.Run("GoMap/multipleSet", func(b *testing.B) {
		m := buildGoMap(size)
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			k := keys[i&(len(keys)-1)]
			m[k] = 999
		}
	})
}
