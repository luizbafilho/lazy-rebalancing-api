package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func main() {

	http.HandleFunc("/rebalancing", rebalancingHandler)

	fmt.Println("Running on :3000")
	http.ListenAndServe(":3000", nil)
}

type rebalancingRequest struct {
	AmountToContribute float64   `json:"amount_to_contribute,omitempty"`
	Portfolio          Portfolio `json:"portfolio,omitempty"`
}

func rebalancingHandler(w http.ResponseWriter, r *http.Request) {
	rq := rebalancingRequest{}

	err := json.NewDecoder(r.Body).Decode(&rq)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("invalid json"))
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