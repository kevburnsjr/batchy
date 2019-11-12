package main

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kevburnsjr/batchy"
)


func main() {
	// Max batch size 100
	// Max wait time 100 milliseconds
	batcher := batchy.New(100, 100*time.Millisecond, func(items []interface{}) (errs []error) {
		errs = make([]error, len(items))
		var ids = make([]string, len(items))
		for i, v := range items {
			ids[i] = v.(string)
		}
		err := appendToFile(strings.Join(ids, "\n"))
		if err != nil {
			for i := range items {
				errs[i] = err
			}
		}
		return
	})
	http.HandleFunc("/batched", func(w http.ResponseWriter, r *http.Request) {
		batcher.Add(r.FormValue("id"))
	})
	http.HandleFunc("/unbatched", func(w http.ResponseWriter, r *http.Request) {
		appendToFile(r.FormValue("id"))
	})
	http.ListenAndServe(":8080", nil)
}

func appendToFile(str string) (err error) {
	f, err := os.OpenFile("items", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return
	}
	defer f.Close()

	_, err = f.Write([]byte(str + "\n"))
	if err != nil {
		return
	}
	return
}
