package main

import (
	"log"
	"fmt"
	"time"
	"crypto/rand"
	"runtime"
	_ "runtime/debug"
	"net/http"
	_ "net/http/pprof"
) 

type Stocks interface {
	FillWithValue(val string)
	ShowAll()
	Clear()
	CopyTo(to map[int]string)
	Realloc()
}

type stocks struct {
	storeMap map[int]string
}

func NewStocks() *stocks {
	var s stocks
	s.storeMap = make(map[int]string)
	return &s
} 

func main() {


	http.HandleFunc("/", DoStuff)
	http.HandleFunc("/view", View)
	fmt.Println("Server Started")
	log.Fatal(http.ListenAndServe("localhost:8080", nil))
}

func View(w http.ResponseWriter, r *http.Request) {

}

func DoStuff(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("All Stuff Started"))

	chanUpdate := make(chan bool)
	chanStats := make(chan bool)

	go timerStats(chanStats)
	go timer(chanUpdate)

	go func(){
		for <-chanStats {
			PrintMemUsage()
		}
	}()

	stocksInitial := NewStocks()
	stocksInitial.FillWithValue("one")

	for <-chanUpdate {
		PrintMemUsage()
		fmt.Println("Start realloc")
		// stocksUpdated := NewStocks()
		// stocksUpdated.FillWithValue(genRand())
		// stocksInitial = stocksUpdated
		stocksInitial.Realloc()
		// PrintMemUsage()
		fmt.Println("end realloc")
	}
}

func (s *stocks) Realloc() {
	stocksUpdated := NewStocks()
	stocksUpdated.FillWithValue(genRand())
	s.storeMap = stocksUpdated.storeMap
	fmt.Print("a:")
	PrintMemUsage()
}

func (s *stocks) FillWithValue(val string) {
	for i := 0; i < 10000000; i++ {
		s.storeMap[i] = val
	}
	time.Sleep(15 * time.Second)
}

func (s *stocks) ShowAll() {
	for _,v := range s.storeMap {
		fmt.Printf("%s ", v)
	}
	fmt.Println()
}

func (s *stocks) CopyTo(to map[int]string) {
	for k,v := range s.storeMap {
		to[k] = v
	}
}

func timer(chanUpdate chan bool) {
	for {
		chanUpdate<-true
		time.Sleep(8 * time.Minute)
	}
}

func timerStats(chanStats chan bool) {
	for {
		chanStats<-true
		time.Sleep(time.Minute)
	}
}

func (s *stocks) Clear() {
	fmt.Println("before clear: ", len(s.storeMap))
	for v := range s.storeMap {
		delete(s.storeMap,v)
	}
	fmt.Println("after clear: ", len(s.storeMap))
}

func genRand() string {
	b := make([]byte, 2)
    if _, err := rand.Read(b); err != nil {
        panic(err)
    }

	return fmt.Sprintf("%X", b)
}

func PrintMemUsage() {
        var m runtime.MemStats
        runtime.ReadMemStats(&m)
        // For info on each, see: https://golang.org/pkg/runtime/#MemStats
        fmt.Printf(" LastGC=%v", m.LastGC)
        fmt.Printf(" NumGC=%v", m.NumGC)
        fmt.Printf(" NextGC=%vmB", bToMb(m.NextGC))
        fmt.Printf(" HeapAlloc=%vmB", bToMb(m.HeapAlloc))
        fmt.Printf(" HeapIdle=%vmB", bToMb(m.HeapIdle))
        fmt.Printf(" HeapReleased=%vmB", bToMb(m.HeapReleased))
        fmt.Printf(" HeapObjects=%v\n", m.HeapObjects)
        // fmt.Printf(" Mallocs = %v", m.Mallocs)
        // fmt.Printf(" Frees = %v\n", m.Frees)
}

func bToMb(b uint64) uint64 {
    return b / 1024 / 1024
}
