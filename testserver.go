package main

import (
	"context"
	"encoding/json"
	"flag"
	"log"
	"math"
	"math/rand"
	"net/http"
	"strconv"
	"sync/atomic"
	"time"
)

var (
	limit     int64
	threshold int64 = 20
)

type GeneratorFunc func(context.Context, int) ([]int, error)

func main() {
	listenAddr := flag.String("http.addr", ":8090", "http listen address")
	flag.Parse()

	http.HandleFunc("/primes", handler(primes))
	http.HandleFunc("/fibo", handler(fibo))
	http.HandleFunc("/odd", handler(odd))
	http.HandleFunc("/rand", handler(random))

	log.Fatal(http.ListenAndServe(*listenAddr, nil))
}

func handler(fn GeneratorFunc) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		atomic.AddInt64(&limit, 1)
		defer atomic.AddInt64(&limit, -1)

		if atomic.LoadInt64(&limit) > threshold {
			<-time.After(100 * time.Millisecond)

			http.Error(w, "service unavailable", http.StatusServiceUnavailable)
			return
		}

		var (
			factor int64
			err    error
		)
		// Factor header is only for testing purposes.
		if factor, err = strconv.ParseInt(r.Header.Get("Factor"), 10, 64); err != nil {
			switch time.Now().Unix() % 6 {
			case 0:
				factor = 1000
			case 1:
				factor = 10000
			case 2:
				factor = 100000
			case 3:
				factor = 1000000
			case 4:
				factor = 10000000
			case 5:
				factor = 100000000
			}
		}

		dl, _ := r.Context().Deadline()
		log.Printf("%s: request received (factor=%d, deadline=%s, goroutines=%d)", r.URL.String(), factor, dl.String(), atomic.LoadInt64(&limit))

		numbers, err := fn(r.Context(), int(factor))

		if err != nil {
			switch err {
			case context.DeadlineExceeded:
				http.Error(w, "timeout", http.StatusRequestTimeout)
			case context.Canceled:
				http.Error(w, "request canceled by the client", http.StatusNoContent)
			default:
				http.Error(w, "internal server error", http.StatusInternalServerError)
			}
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		json.NewEncoder(w).Encode(map[string]interface{}{"numbers": numbers})
	}
}

func primes(ctx context.Context, max int) ([]int, error) {
	r := make([]int, 0)
	b := make([]bool, max)

	for i := 2; i < max; i++ {
		if b[i] == true {
			continue
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			r = append(r, i)
		}
		for k := i * i; k < max; k += i {
			b[k] = true
		}
	}
	return r, nil
}

func fibo(ctx context.Context, max int) ([]int, error) {
	var res []int

	for i := 0; i < max; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			Phi := (1 + math.Sqrt(5)) / 2
			phi := (1 - math.Sqrt(5)) / 2

			res = append(res, int((math.Pow(Phi, float64(i))-math.Pow(phi, float64(i)))/math.Sqrt(5)))
		}
	}
	return res, nil
}

func odd(ctx context.Context, max int) ([]int, error) {
	var res []int
	for i := 0; i < max; i++ {
		if i%2 == 0 {
			continue
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			res = append(res, i)
		}
	}
	return res, nil
}

func random(ctx context.Context, max int) ([]int, error) {
	var res []int
	for i := 0; i < max; i++ {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
			res = append(res, rand.Int())
		}
	}
	return res, nil
}
