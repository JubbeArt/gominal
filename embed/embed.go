package main

import (
	"fmt"
	"io/ioutil"
)

func main() {
	fontBytes, _ := ioutil.ReadFile("DejaVuSansMono.ttf")
	fontBoldBytes, _ := ioutil.ReadFile("DejaVuSansMono-Bold.ttf")

	ioutil.WriteFile("../font_data.go", []byte(fmt.Sprintf(`package gominal

var fontBytes = %#v

var fontBoldBytes = %#v
`, fontBytes, fontBoldBytes)), 0664)
}
