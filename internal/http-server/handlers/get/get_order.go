package get

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Request struct {
	OrderUID string `json:"order_uid"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.32 --name=Getter
type Getter interface {
	Get(key string) ([]byte, error)
}

// New hendler for GET and POST
func New(cache Getter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			newGet(w, r, cache)
		case http.MethodPost:
			newPost(w, r, cache)
		}
	}
}

func newGet(w http.ResponseWriter, r *http.Request, cache Getter) {

	var uid Request
	vars := mux.Vars(r)
	uid.OrderUID = vars["order_uid"]

	get(w, cache, &uid)
}

func newPost(w http.ResponseWriter, r *http.Request, cache Getter) {

	var uid Request

	defer func() {
		err := r.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	err := json.NewDecoder(r.Body).Decode(&uid)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	get(w, cache, &uid)
}

func get(w http.ResponseWriter, cache Getter, req *Request) {

	log.Println(req.OrderUID)

	order, err := cache.Get(req.OrderUID)
	if err != nil {
		log.Println(err)
		w.WriteHeader(http.StatusNotFound)
		_, err := w.Write([]byte(`{"order": "Not_found"}`))
		if err != nil {
			log.Println(err)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusOK)
	_, err = w.Write(order)
	if err != nil {
		log.Println(err)
	}
}
