package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	fontBytes, _ := ioutil.ReadFile("DejaVuSansMono.ttf")
	fontBoldBytes, _ := ioutil.ReadFile("DejaVuSansMono-Bold.ttf")

	file := fmt.Sprintf(`package gominal

var fontBytes = %#v

var fontBoldBytes = %#v
`, fontBytes, fontBoldBytes)

	_ = ioutil.WriteFile("../font_data.go", []byte(file), 0664)

}
