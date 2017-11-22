package main

import (
	"fmt"
	"github.com/asnelzin/translate/yandex"
	"golang.org/x/text/language"
	"strings"
	"os"
	"github.com/jessevdk/go-flags"
)

var revision string
var yandexAPIKey = getenv("YANDEX_API_KEY")

var opts struct {
	From    string `short:"f" long:"from" description:"From language"`
	To      string `short:"t" long:"to" description:"To language" default:"ru" required:"true"`
	Version bool   `short:"v" long:"version" description:"Print the version information and exit"`

	Positional struct {
		Text []string `positional-arg-name:"text" description:"Text to translate"`
	} `positional-args:"yes" required:"yes"`
}
var parser = flags.NewParser(&opts, flags.Default)

func getenv(name string) string {
	v := os.Getenv(name)
	if v == "" {
		panic("missing required environment variable " + name)
	}
	return v
}

func main() {
	_, err := parser.Parse()
	if err != nil {
		os.Exit(1)
	}

	if opts.Version {
		showVersion()
		os.Exit(0)
	}

	if len(opts.Positional.Text) == 0 {
		errorf("the required argument `text (at least 1 argument)` was not provided")
	}

	to, err := language.Parse(opts.To)
	if err != nil {
		errorf("Failed to parse `to language`: %v", err)
	}

	var from language.Tag
	if opts.From != "" {
		from, err = language.Parse(opts.From)
		if err != nil {
			errorf("Failed to parse `from language`: %v", err)
		}
	}

	c := yandex.NewClient(nil, yandexAPIKey)
	result, err := c.TranslateString(from, to, strings.Join(opts.Positional.Text, " "))
	if err != nil {
		errorf("Failed to translate: %v", err)
	}

	fmt.Println(result)
	fmt.Printf("\nPowered by Yandex.Translate (http://translate.yandex.ru/)\n")
}

func showVersion() {
	fmt.Printf("Translate version %s\n", revision)
}

func errorf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format+"\n", args...)
	os.Exit(2)
}
