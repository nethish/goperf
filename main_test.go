package main

import (
	"sync"
	"testing"
)

// It has to allocate that much memory
type Data struct {
	Values [1024]int
}

// BenchmarkWithoutPooling measures the performance of direct heap allocations.
func BenchmarkWithoutPooling(b *testing.B) {
	// This is needed. go optimizes unused vars and without this the benchmark will tell you this is faster
	var a *Data
	for b.Loop() {
		data := &Data{}     // Allocating a new object each time
		data.Values[0] = 42 // Simulating some memory activity
		a = data
	}
	_ = a
}

// dataPool is a sync.Pool that reuses instances of Data to reduce memory allocations.
var dataPool = sync.Pool{
	New: func() any {
		return &Data{}
	},
}

// BenchmarkWithPooling measures the performance of using sync.Pool to reuse objects.
func BenchmarkWithPooling(b *testing.B) {
	var a *Data
	for b.Loop() {
		obj := dataPool.Get().(*Data) // Retrieve from pool
		obj.Values[0] = 42            // Simulate memory usage
		dataPool.Put(obj)             // Return object to pool for reuse
		a = obj
	}
	_ = a
}

type A struct{}

func (a A) Do() {}

type I interface {
	Do()
}

var a I = &A{}

func B() {
	a.Do()
}
