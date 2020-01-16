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
	Alt   bool   `json:"alt"`
	Super bool   `json:"super"`
}

type charResponse struct {
	Type string `json:"type"`
	Rune string `json:"char"`
}

type mouseResponse struct {
	Type   string `json:"type"`
	Button string `json:"button"`
	Col    int    `json:"col"`
	Row    int    `json:"row"`
	Ctrl   bool   `json:"ctrl"`
	Shift  bool   `json:"shift"`
	Alt    bool   `json:"alt"`
	Super  bool   `json:"super"`
}

type sizeResponse struct {
	Type      string `json:"type"`
	Rows      int    `json:"rows"`
	Cols      int    `json:"cols"`
	ColWidth  int    `json:"colWidth"`
	RowHeight int    `json:"rowHeight"`
}

type errorResponse struct {
	Type  string `json:"type"`
	Error string `json:"error"`
}
