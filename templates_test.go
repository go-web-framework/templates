package templates

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestTemplates(t *testing.T) {

	Convey("templates", t, func() {
		Convey("parse", func() {
			s := Set{}
			So(s.Parse("testdata/templates"), ShouldBeNil)
		})
	})
}
