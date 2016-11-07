package templates

import (
	"log"
	"os"
	"path/filepath"
)

func Example() {
	s := Set{
		PartialsDir: "defines",
		DefaultArgs: Args{
			"Title": "Page title",
			"Name":  "someone",
		},
	}

	if err := s.Parse(filepath.Join("testdata", "example")); err != nil {
		log.Fatalln(err)
	}

	if err := s.Execute("example.html", os.Stdout, Args{
		"Name": "world",
	}); err != nil {
		log.Fatalln(err)
	}

	// Output:
	// <!doctype html>
	// <head>
	//     <title>Page title</title>
	// </head>
	// <h1>hello, world</h1>
}
