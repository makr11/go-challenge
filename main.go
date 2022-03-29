package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-challenge/math_cache"
	"net/http"
	"strconv"
	"strings"
)

type MathResponse struct {
	Action string `json:"action"`
	InputX int    `json:"x"`
	InputY int    `json:"y"`
	Answer int    `json:"answer"`
	Cached bool   `json:"cached"`
}

type operation func(x int, y int) (int, error)

func add(w http.ResponseWriter, req *http.Request) {
	op := func(x int, y int) (int, error) {
		answer := x + y
		return answer, nil
	}
	simple_math(w, req, op)
}

func subtract(w http.ResponseWriter, req *http.Request) {
	op := func(x int, y int) (int, error) {
		answer := x - y
		return answer, nil
	}
	simple_math(w, req, op)
}

func multiply(w http.ResponseWriter, req *http.Request) {
	op := func(x int, y int) (int, error) {
		answer := x * y
		return answer, nil
	}
	simple_math(w, req, op)
}

func divide(w http.ResponseWriter, req *http.Request) {
	op := func(x int, y int) (int, error) {
		if y == 0 {
			return 0, errors.New("Can not divide with '0'")
		}
		answer := x / y
		return answer, nil
	}
	simple_math(w, req, op)
}

func parse_operands(xs string, ys string) (xi int, yi int, err error) {
	var msg string

	xi, x_err := strconv.Atoi(xs)
	if x_err != nil {
		msg = fmt.Sprintf("X parameter value '%s' can not be parsed to integer.\n", xs)
	}

	yi, y_err := strconv.Atoi(ys)
	if y_err != nil {
		msg = msg + fmt.Sprintf("Y parameter value '%s' can not be parsed to integer.", ys)
	}

	if x_err != nil || y_err != nil {
		return xi, yi, errors.New(msg)
	}
	return xi, yi, nil
}

func simple_math(w http.ResponseWriter, req *http.Request, op operation) {
	x, y, parse_err := parse_operands(req.URL.Query().Get("x"), req.URL.Query().Get("y"))
	if parse_err != nil {
		http.Error(w, parse_err.Error(), http.StatusBadRequest)
		return
	}
	operation := strings.Trim(req.URL.Path, "/")

	lc := math_cache.NewLocalCache()
	answer, cache_hit := lc.Read(operation, x, y)
	if !cache_hit {
		var op_err error
		answer, op_err = op(x, y)
		if op_err != nil {
			http.Error(w, op_err.Error(), http.StatusBadRequest)
			return
		}
	}
	lc.Update(operation, x, y, answer)

	payload := &MathResponse{
		Action: operation,
		InputX: x,
		InputY: y,
		Answer: answer,
		Cached: cache_hit,
	}

	p, err := json.Marshal(payload)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	fmt.Fprintf(w, string(p))
}

func main() {
	http.HandleFunc("/add", add)
	http.HandleFunc("/subtract", subtract)
	http.HandleFunc("/multiply", multiply)
	http.HandleFunc("/divide", divide)

	http.ListenAndServe(":8090", nil)
}
