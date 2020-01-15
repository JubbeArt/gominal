package main

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
