package templates

import (
	"bytes"
	"html/template"
	"io/ioutil"
	"path/filepath"
	"sort"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTemplates(t *testing.T) {
	var templateFilenames = []string{
		"hello.html",
		"foo.html",
		"func.html",
	}

	Convey("templates", t, func() {
		set := Set{
			PartialsDir: "partials",
			Funcs: template.FuncMap{
				"Repeat": strings.Repeat,
			},
		}

		Convey("Parse", func() {
			So(set.Parse("testdata/templates"), ShouldBeNil)

			Convey("parsed all templates", func() {
				var keys []string
				for k := range set.Templates {
					keys = append(keys, k)
				}
				sort.Sort(sort.StringSlice(keys))
				sort.Sort(sort.StringSlice(templateFilenames))
				So(keys, ShouldResemble, templateFilenames)
			})

			Convey("delimiter", func() {
				set.LDelim = "[["
				set.RDelim = "]]"
				So(set.Parse("testdata/templates-delim-brackets"), ShouldBeNil)
			})
		})

		Convey("Execute", func() {
			set := Set{
				PartialsDir: "partials",
				Funcs: template.FuncMap{
					"Repeat": strings.Repeat,
				},
			}
			So(set.Parse("testdata/templates"), ShouldBeNil)
			buf := bytes.Buffer{}

			Convey("normal", func() {
				So(set.Execute("hello.html", &buf, map[string]int{"V": 42}), ShouldBeNil)
				b, err := ioutil.ReadFile(filepath.Join("testdata", "expected", "hello0"))
				So(err, ShouldBeNil)
				So(buf.String(), ShouldResemble, string(b))
			})

			Convey("template does not exist", func() {
				So(set.Execute("doesnotexist", nil, nil), ShouldEqual, ErrNoSuchTemplate)
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

			Convey("funcs", func() {
				So(set.Execute("func.html", &buf, "repeatme"), ShouldBeNil)
				b, err := ioutil.ReadFile(filepath.Join("testdata", "expected", "func0"))
				So(err, ShouldBeNil)
				So(buf.String(), ShouldResemble, string(b))
			})
		})
	})
}
