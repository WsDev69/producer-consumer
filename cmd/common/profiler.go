package common

import (
	_ "net/http/pprof"

	"log"
	"net/http"
	"os"
	"runtime/pprof"
)

func ExposePprof(addr string) {
	go func() {
		log.Println(http.ListenAndServe(addr, nil)) // Exposes pprof on localhost:8080
	}()
}

func CPUProf() {
	f, err := os.Create("cpu.prof")
	if err != nil {
		log.Fatal(err)
	}
	err = pprof.StartCPUProfile(f)
	if err != nil {
		log.Println("could not start CPU profile: ", err)
	}
	defer pprof.StopCPUProfile()
}

func MEMProf() {
	f, err := os.Create("mem.prof")
	if err != nil {
		log.Fatal(err)
	}
	if err = pprof.WriteHeapProfile(f); err != nil {
		log.Println("could not write memory profile: ", err)
	}
}
