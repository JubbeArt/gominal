package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

var  windowStateFile string
func init() {
	cache, _ := os.UserCacheDir()
	windowStateFile = filepath.Join(cache, "gominal", "window.json")
}

func defaultState() WindowState {
	return WindowState{Width: 640, Height: 480, X: -1, Y: -1}
}

type WindowState struct {
	Width  int `json:"width"`
	Height int `json:"height"`
	X      int `json:"x"`
	Y      int `json:"y"`
}

// the moment i really wished go had try catch...
func getWindowState() WindowState {
	_, err := os.Stat(windowStateFile)

	if err != nil {
		return defaultState()
	}

	bytes, err := ioutil.ReadFile(windowStateFile)

	if err != nil {
		return defaultState()
	}

	var state WindowState
	err = json.Unmarshal(bytes, &state)

	if err != nil {
		return defaultState()
	}

	if state.Width < 150 {
		state.Width = 150
	}

	if state.Height < 100 {
		state.Height = 100
	}

	return state
}

func (state WindowState) save() {
	err := os.MkdirAll(filepath.Dir(windowStateFile), 0775)

	if err != nil {
		return
	}

	data, err := json.MarshalIndent(state, "", "\t")

	if err != nil {
		return
	}

	_ = ioutil.WriteFile(windowStateFile, data, 0664)
}
