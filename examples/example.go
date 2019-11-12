package main

import (
	"errors"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kevburnsjr/batchy"
)

type unbatchedHandler struct{}

func (h unbatchedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := appendIdToFile(r.FormValue("id"))
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

type batchedHandler struct {
	batcher batchy.Batcher
}

func (h batchedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := h.batcher.Add(r.FormValue("id"))
	if err != nil {
		http.Error(w, err.Error(), 500)
	}
}

func main() {
	http.Handle("/unbatched", unbatchedHandler{})
	http.Handle("/batched", batchedHandler{
		// Max batch size 100
		// Max wait time 100 milliseconds
		batchy.New(100, 100*time.Millisecond, func(items []interface{}) (errs []error) {
			var ok bool
			errs = make([]error, len(items))
			var ids = make([]string, len(items))
			for i, v := range items {
				ids[i], ok = v.(string)
				if !ok {
					errs[i] = errors.New("Wrong data type")
				}
			}
			err := appendIdsToFile(ids)
			if err != nil {
				for i := range items {
					errs[i] = err
				}
			}
			return
		}),
	})
	http.ListenAndServe(":8080", nil)
}

func appendIdsToFile(ids []string) (err error) {
	if len(ids) == 0 {
		return
	}
	f, err := os.OpenFile("items", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return
	}
	defer f.Close()

	_, err = f.Write([]byte(strings.Join(ids, "\n") + "\n"))
	if err != nil {
		return
	}
	return
}

func appendIdToFile(id string) (err error) {
	f, err := os.OpenFile("items", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return
	}
	defer f.Close()

	_, err = f.Write([]byte(id + "\n"))
	if err != nil {
		return
	}
	return
}
