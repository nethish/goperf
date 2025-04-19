package main

import (
	"fmt"
	"runtime"
	"strconv"
)

type Data struct {
	name      string
	age       int
	gender    rune
	info      string
	encrypted [1024]byte
}

func (d *Data) GetFormattedString() string {
	return strconv.FormatInt(int64(d.age), 10) + d.name
}

// Trying to see for which value of GOGC this program is performant, and has less GC cycles
// I'll be using GOMEMLIMIT and GOGC flags while running the program with GOMAXPROCS=1
func main() {
	// Allocates almost 2GiB
	SIZE := 2_000_000

	var result string
	var retained []*Data
	for i := range SIZE {
		d := GetNewData(i)
		result = d.GetFormattedString()
		retained = append(retained, d)

		if i%100000 == 0 {
			var m runtime.MemStats
			runtime.ReadMemStats(&m)

			// Current Heap Allocation
			fmt.Printf("Iter = %d\n", i)
			fmt.Printf("Alloc = %v KB\n", m.Alloc/1024)
			fmt.Printf("TotalAlloc = %v KB\n", m.TotalAlloc/1024)
			fmt.Printf("Sys = %v KB\n", m.Sys/1024)
			fmt.Printf("NumGC = %v\n", m.NumGC)
			fmt.Printf("PauseTotalNs = %v\n", m.PauseTotalNs)
			fmt.Println("-----------------------")
			retained = []*Data{}
		}
	}

	fmt.Println(result)
}

func GetNewData(age int) *Data {
	return &Data{
		name:   "A",
		age:    age,
		gender: 0,
		info:   "A big brown fox ran away from here",
	}
}
