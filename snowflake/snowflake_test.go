package main

import "testing"

func BenchmarkSnowFlake(b *testing.B) {
	b.ResetTimer()

	generate := Constructor(1, 1)
	for i := 0; i < b.N; i++ {
		generate.nextId()
	}
}

func BenchmarkCurTime(b *testing.B) {
	b.ResetTimer()

	generate := Constructor(1, 1)
	for i := 0; i < b.N; i++ {
		generate.nextId()
	}
}
