package main

import (
	"encoding/json"
	"fmt"

	"github.com/go-gl/glfw/v3.3/glfw"
)

var (
	mouseCol = -1
	mouseRow = -1
)

func charCallback(char rune) {
	sendResponse(charResponse{Type: "char", Rune: string(char)})
}

func keyCallback(key glfw.Key, scanCode int, action glfw.Action, mods glfw.ModifierKey) {
	if action == glfw.Repeat {
		return
	}

	keyName := glfw.GetKeyName(key, scanCode)

	if keyName == "" {
		keyName = keyLookup[key]
	}

	if keyName == "" {
		keyName = "unknown"
	}

	sendResponse(keyResponse{
		Type:  "key",
		Key:   keyName,
		State: actionLookup[action],
		Ctrl:  mods&glfw.ModControl != 0,
		Shift: mods&glfw.ModShift != 0,
		Alt:   mods&glfw.ModAlt != 0,
		Super: mods&glfw.ModSuper != 0,
	})
}

func mouseClickCallback(button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
	if buttonText, ok := mouseLookup[button]; ok {
		mouseX, mouseY := window.GetCursorPos()
		sendResponse(mouseClickResponse{
			Type:   "mouseClick",
			Button: buttonText,
			State:  actionLookup[action],
			Col:    int(mouseX) / colWidth,
			Row:    int(mouseY) / rowHeight,
			Ctrl:   mods&glfw.ModControl != 0,
			Shift:  mods&glfw.ModShift != 0,
			Alt:    mods&glfw.ModAlt != 0,
			Super:  mods&glfw.ModSuper != 0,
		})
	}
}

func mouseMoveCallback(x, y int) {
	newMouseCol := x / colWidth
	newMouseRow := y / rowHeight

	if mouseCol != newMouseCol || mouseRow != newMouseRow {
		mouseCol = newMouseCol
		mouseRow = newMouseRow
		sendResponse(mouseMoveResponse{
			Type: "mouseMove",
			Col:  mouseCol,
			Row:  mouseRow,
		})
	}
}

func sizeCallback(width int, height int) {
	newCols := width / colWidth
	newRows := height / rowHeight

	if newCols == cols && newRows == rows {
		return
	}

	cols = newCols
	rows = newRows

	sendResponse(sizeResponse{Type: "size", Rows: rows, Cols: cols, ColWidth: colWidth, RowHeight: rowHeight})
}

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
	State string `json:"state"`
	Ctrl  bool   `json:"ctrl"`
	Shift bool   `json:"shift"`
	Alt   bool   `json:"alt"`
	Super bool   `json:"super"`
}

type charResponse struct {
	Type string `json:"type"`
	Rune string `json:"char"`
}

type mouseClickResponse struct {
	Type   string `json:"type"`
	Button string `json:"button"`
	State  string `json:"state"`
	Col    int    `json:"col"`
	Row    int    `json:"row"`
	Ctrl   bool   `json:"ctrl"`
	Shift  bool   `json:"shift"`
	Alt    bool   `json:"alt"`
	Super  bool   `json:"super"`
}

type mouseMoveResponse struct {
	Type string `json:"type"`
	Col  int    `json:"col"`
	Row  int    `json:"row"`
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

var actionLookup = map[glfw.Action]string{
	glfw.Press:   "press",
	glfw.Release: "release",
}

var mouseLookup = map[glfw.MouseButton]string{
	glfw.MouseButtonLeft:   "left",
	glfw.MouseButtonMiddle: "middle",
	glfw.MouseButtonRight:  "right",
}

var keyLookup = map[glfw.Key]string{
	glfw.KeyLeftControl:  "ctrl",
	glfw.KeyRightControl: "ctrl",
	glfw.KeyLeftShift:    "shift",
	glfw.KeyRightShift:   "shift",
	glfw.KeyLeftAlt:      "alt",
	glfw.KeyRightAlt:     "alt",
	glfw.KeyLeftSuper:    "super",
	glfw.KeyRightSuper:   "super",
	glfw.KeyTab:          "tab",
	glfw.KeyEnter:        "enter",
	glfw.KeySpace:        "space",
	glfw.KeyBackspace:    "backspace",
	glfw.KeyEscape:       "escape",
	glfw.KeyLeft:         "left",
	glfw.KeyRight:        "right",
	glfw.KeyUp:           "up",
	glfw.KeyDown:         "down",
	glfw.KeyCapsLock:     "capsLock",
	glfw.KeyDelete:       "delete",
	glfw.KeyInsert:       "insert",
	glfw.KeyHome:         "home",
	glfw.KeyPageUp:       "pageUp",
	glfw.KeyPageDown:     "pageDown",
	glfw.KeyEnd:          "end",
	glfw.KeyNumLock:      "numLock",

	glfw.KeyF1:  "f1",
	glfw.KeyF2:  "f2",
	glfw.KeyF3:  "f3",
	glfw.KeyF4:  "f4",
	glfw.KeyF5:  "f5",
	glfw.KeyF6:  "f6",
	glfw.KeyF7:  "f7",
	glfw.KeyF8:  "f8",
	glfw.KeyF9:  "f9",
	glfw.KeyF10: "f10",
	glfw.KeyF11: "f11",
	glfw.KeyF12: "f12",
	glfw.KeyF13: "f13",
	glfw.KeyF14: "f14",
	glfw.KeyF15: "f15",
	glfw.KeyF16: "f16",
	glfw.KeyF17: "f17",
	glfw.KeyF18: "f18",
	glfw.KeyF19: "f19",
	glfw.KeyF20: "f20",
	glfw.KeyF21: "f21",
	glfw.KeyF22: "f22",
	glfw.KeyF23: "f23",
	glfw.KeyF24: "f24",
	glfw.KeyF25: "f25",
}
