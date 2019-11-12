package main

import (
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/kevburnsjr/batchy"
)

var batcher = batchy.New(100, 100*time.Millisecond, func(items []interface{}) (errs []error) {
	var ids = make([]string, len(items))
	for i, v := range items {
		ids[i] = v.(string)
	}
	appendToFile(strings.Join(ids, "\n"))
	return
})

func main() {
	http.HandleFunc("/batched", func(w http.ResponseWriter, r *http.Request) {
		err := batcher.Add(r.FormValue("id"))
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	})
	http.HandleFunc("/unbatched", func(w http.ResponseWriter, r *http.Request) {
		err := appendToFile(r.FormValue("id"))
		if err != nil {
			http.Error(w, err.Error(), 500)
		}
	})
	http.ListenAndServe(":8080", nil)
}

func appendToFile(str string) (err error) {
	f, err := os.OpenFile("items", os.O_APPEND|os.O_RDWR|os.O_CREATE, 0600)
	if err != nil {
		return
	}
	_, err = f.Write([]byte(str + "\n"))
	f.Close()
	return
}
