package main

import (
	"fmt"
	"github.com/asnelzin/translate/yandex"
	"golang.org/x/text/language"
)

func main() {
	c := yandex.NewClient(nil)
	fmt.Print(c.TranslateString(language.English, language.Russian, "hello"))
}