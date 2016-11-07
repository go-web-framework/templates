package templates

import (
	"bytes"
	"io/ioutil"
	"path/filepath"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTemplates(t *testing.T) {
	var templateFilenames = []string{
		"hello.html",
	}

	Convey("templates", t, func() {
		Convey("parse", func() {
			set := Set{PartialsDir: "partials"}
			So(set.Parse("testdata/templates"), ShouldBeNil)

			Convey("parsed all templates", func() {
				for _, s := range templateFilenames {
					So(set.Templates[s], ShouldNotBeNil)
				}
			})

			Convey("execute", func() {
				buf := bytes.Buffer{}

				Convey("normal", func() {
					So(set.Execute("hello.html", &buf, map[string]int{"V": 42}), ShouldBeNil)
					b, err := ioutil.ReadFile(filepath.Join("testdata", "expected", "hello0"))
					So(err, ShouldBeNil)
					So(buf.String(), ShouldResemble, string(b))
				})

				Convey("default args", func() {
					set.DefaultArgs = Args{"V": "default"}

					Convey("no override", func() {
						So(set.Execute("hello.html", &buf, nil), ShouldBeNil)
						b, err := ioutil.ReadFile(filepath.Join("testdata", "expected", "hello1"))
						So(err, ShouldBeNil)
						So(buf.String(), ShouldResemble, string(b))
					})

					Convey("override args", func() {
						So(set.Execute("hello.html", &buf, Args{"V": "overriden"}), ShouldBeNil)
						b, err := ioutil.ReadFile(filepath.Join("testdata", "expected", "hello2"))
						So(err, ShouldBeNil)
						So(buf.String(), ShouldResemble, string(b))
					})
				})

			})
		})
	})
}
