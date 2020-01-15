package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

const cacheFolder = "gominal"
const windowFile = "window.json"

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
func savedWindowState() WindowState {
	cacheDir, err := os.UserConfigDir()

	if err != nil {
		return defaultState()
	}

	filePath := filepath.Join(cacheDir, cacheFolder, windowFile)
	_, err = os.Stat(filePath)

	if err != nil {
		return defaultState()
	}

	bytes, err := ioutil.ReadFile(filePath)

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
	cacheDir, err := os.UserConfigDir()

	if err != nil {
		return
	}

	err = os.MkdirAll(filepath.Join(cacheDir, cacheFolder), 0775)

	if err != nil {
		return
	}

	filePath := filepath.Join(cacheDir, cacheFolder, windowFile)
	data, err := json.MarshalIndent(state, "", "\t")

	if err != nil {
		return
	}

	_ = ioutil.WriteFile(filePath, data, 0664)
}
