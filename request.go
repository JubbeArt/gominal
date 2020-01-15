package main

import (
	"encoding/base64"
	"encoding/json"
	"image"
	"image/color"
	_ "image/jpeg"
	_ "image/png"
	"strings"
	"unicode/utf8"

	"github.com/pkg/errors"
)

type request struct {
	Type *string `json:"type"`
}

type setCharRequest struct {
	Rune       *string     `json:"char"`
	Col        *int        `json:"col"`
	Row        *int        `json:"row"`
	Color      *color.RGBA `json:"color"`
	Background *color.RGBA `json:"background"`
	Style      *string     `json:"style"`
}

type imageRequest struct {
	Image *string `json:"image"`
	Col   *int    `json:"col"`
	Row   *int    `json:"row"`
}

type titleRequest struct {
	Title *string `json:"title"`
}

func handleRequest(line []byte) error {
	var request request

	err := json.Unmarshal(line, &request)

	if err != nil {
		return errors.WithMessage(err, "could not parse request")
	} else if request.Type == nil {
		return errors.New("request is missing \"type\" field")
	}

	switch *request.Type {
	case "char":
		var req setCharRequest
		err := json.Unmarshal(line, &req)

		if err != nil {
			return errors.WithMessage(err, "could not parse request")
		} else if req.Rune == nil {
			return errors.New("char request is missing \"char\" field")
		}

		char, width := utf8.DecodeRuneInString(*req.Rune)

		if char == utf8.RuneError && width == 0 {
			return errors.New("char request was sent with empty char")
		} else if char == utf8.RuneError && width == 1 {
			return errors.New("char request was sent with invalid utf8")
		} else if req.Col == nil {
			return errors.New("char request is missing \"col\" field")
		} else if req.Row == nil {
			return errors.New("char request is missing \"row\" field")
		}

		col := *req.Col
		row := *req.Row

		textColor := color.RGBA{R: 255, G: 255, B: 255, A: 255}
		if req.Color != nil {
			textColor = *req.Color
			textColor.A = 255
		}

		bg := color.RGBA{R: 0, G: 0, B: 0, A: 255}
		if req.Background != nil {
			bg = *req.Background
			bg.A = 255
		}

		style := styleNormal
		if req.Style != nil {
			if ok := styles[*req.Style]; !ok {
				return errors.Errorf("char request got invalid style: %q", *req.Style)
			}

			style = *req.Style
		}

		drawRequests <- charDrawRequest{char: char, col: col, row: row, textColor: textColor, bg: bg, style: style}
	case "image":
		var req imageRequest
		err := json.Unmarshal(line, &req)

		if err != nil {
			return errors.WithMessage(err, "could not parse request")
		} else if req.Col == nil {
			return errors.New("image request is missing \"col\" field")
		} else if req.Row == nil {
			return errors.New("image request is missing \"row\" field")
		} else if req.Image == nil {
			return errors.New("image request is missing \"image\" field")
		}

		dec := base64.NewDecoder(base64.StdEncoding, strings.NewReader(*req.Image))
		img, _, err := image.Decode(dec)

		if err != nil {
			return errors.WithMessage(err, "could not decode image")
		}

		drawRequests <- imageDrawRequest{img: img, col: *req.Col, row: *req.Row}
	case "clear":
		drawRequests <- clearDrawRequest{}
	case "title":
		var req titleRequest
		err := json.Unmarshal(line, &req)

		if err != nil {
			return errors.WithMessage(err, "could not parse request")
		} else if req.Title == nil {
			return errors.New("title request is missing \"title\" field")
		}

		window.SetTitle(*req.Title)
	case "close":
		window.SetShouldClose(true)
	default:
		return errors.Errorf("unknown request type %q", request.Type)
	}

	return nil
}
