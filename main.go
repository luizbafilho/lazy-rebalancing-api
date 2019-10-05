package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	http.HandleFunc("/rebalancing", rebalancingHandler)

	fmt.Println("Running on :" + port)
	http.ListenAndServe(":"+port, nil)
}

type rebalancingRequest struct {
	AmountToContribute float64   `json:"amount_to_contribute,omitempty"`
	Portfolio          Portfolio `json:"portfolio,omitempty"`
}

func rebalancingHandler(w http.ResponseWriter, r *http.Request) {
	rq := rebalancingRequest{}

	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	err = json.Unmarshal(b, &rq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid json: " + err.Error()))
		return
	}

	balancedPortfolio := lazyRebalance(rq.AmountToContribute, rq.Portfolio)

	js, err := json.Marshal(balancedPortfolio)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(js)
}
