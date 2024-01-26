package main

import (
	"fmt"
	cache "github.com/aliforever/go-gcache"
	"sync"
	"time"
)

func main() {
	gCache := cache.New[string]()

	wg := &sync.WaitGroup{}
	wg.Add(100)

	for i := 0; i < 10; i++ {
		go fetchVal(wg, gCache)
		go fetchVal(wg, gCache)
		go fetchVal(wg, gCache)
		go fetchVal(wg, gCache)
		go fetchVal(wg, gCache)
		go fetchVal(wg, gCache)
		go fetchVal(wg, gCache)
		go fetchVal(wg, gCache)
		go fetchVal(wg, gCache)
		go fetchVal(wg, gCache)

		time.Sleep(time.Second * 2)
	}

	wg.Wait()
}

func loadVal() (*string, error) {
	// fmt.Println("loaded from db")

	val := "value"

	return &val, nil
}

func fetchVal(wg *sync.WaitGroup, gCache *cache.Cache[string]) {
	defer wg.Done()

	val, loaded, _ := gCache.LoadOrStoreFunc("key", loadVal, time.Second*5)

	fmt.Println("val: ", val, " loaded: ", loaded)
}
