package main

import (
	"encoding/json"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"net/http"
	"strconv"
	"strings"
)

type Route struct {
	pool      *pgxpool.Pool
	errorBody *[]byte
}

type PathLengthError struct{}

func (PathLengthError) Error() string {
	return "the length of current array is less than requested parameter number"
}

func (route *Route) Handle(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		route.retrieveUserById(w, r)
	case http.MethodPost:
		route.createUser(w, r)
	case http.MethodPut:
		route.updateUser(w, r)
	case http.MethodDelete:
		route.deleteUser(w, r)
	}
}

func (route *Route) createUser(w http.ResponseWriter, r *http.Request) {
	var user User

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Fatal(err)
	}

	err = user.Create(route.pool)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte("successful operation"))
}

func (route *Route) retrieveUserById(w http.ResponseWriter, r *http.Request) {
	id, err := secondValueOfUrlPathAsNum(r)
	route.responseWithError(w, err)

	user, err := RetrieveById(id, route.pool)
	route.responseWithError(w, err)

	js, err := json.Marshal(user)
	route.responseWithError(w, err)

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}

func (route *Route) updateUser(w http.ResponseWriter, r *http.Request) {
	id, err := secondValueOfUrlPathAsNum(r)
	route.responseWithError(w, err)

	var user User

	err = json.NewDecoder(r.Body).Decode(&user)
	route.responseWithError(w, err)

	user.Update(id, route.pool)

	w.Header().Set("Content-Type", "application/json")
}

func (route *Route) deleteUser(w http.ResponseWriter, r *http.Request) {
	id, err := secondValueOfUrlPathAsNum(r)
	route.responseWithError(w, err)

	DeleteById(id, route.pool)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent)
	w.Write([]byte("user deleted"))
}

func (route *Route) responseWithError(w http.ResponseWriter, err error) {
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.Write(*route.errorBody)
		log.Fatal(err)
	}
}

func secondValueOfUrlPathAsNum(r *http.Request) (int, error) {
	id, err := secondValueOfUrlPath(r)
	if err != nil {
		return 0, err
	}

	numId, err := strconv.Atoi(id)
	if err != nil {
		return 0, err
	}

	return numId, nil
}

func secondValueOfUrlPath(r *http.Request) (string, error) {
	values := strings.Split(r.URL.Path, "/")
	if len(values) < 3 {
		return "", PathLengthError{}
	}

	return values[2], nil
}
