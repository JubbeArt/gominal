package main

import (
	"encoding/json"
	"fmt"
)

func sendError(err error) {
	sendErrorStr(err.Error())
}

func sendErrorStr(err string) {
	sendResponse(errorResponse{Type: "error", Error: err})
}

func sendResponse(response interface{}) {
	bytes, err := json.Marshal(response)

	if err != nil {
		panic(err)
	}

	fmt.Println(string(bytes))
}

type keyResponse struct {
	Type  string `json:"type"`
	Key   string `json:"key"`
	Ctrl  bool   `json:"ctrl"`
	Shift bool   `json:"shift"`
	Super bool   `json:"meta"`
	Alt   bool   `json:"alt"`
}

type charResponse struct {
	Type string `json:"type"`
	Rune string `json:"char"`
}

type errorResponse struct {
	Type  string `json:"type"`
	Error string `json:"error"`
}

type sizeResponse struct {
	Type      string `json:"type"`
	Rows      int    `json:"rows"`
	Cols      int    `json:"cols"`
	ColWidth  int    `json:"colWidth"`
	RowHeight int    `json:"rowHeight"`
}
