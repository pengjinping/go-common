package strings

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestIsBlank(t *testing.T) {
	Convey("空白字符串测试", t, func() {
		So(IsBlank("你好 "), ShouldEqual, false)
		So(IsBlank("null"), ShouldEqual, false)
		So(IsBlank("0"), ShouldEqual, false)
		So(IsBlank(fmt.Sprint(0)), ShouldEqual, false)
		var p string
		So(IsBlank(p), ShouldEqual, true)
		So(IsBlank(""), ShouldEqual, true)
		So(IsBlank(" "), ShouldEqual, true)
	})
}

func TestIsNotBlank(t *testing.T) {
	Convey("空白字符串测试", t, func() {
		So(IsNotBlank("你好 "), ShouldEqual, true)
		So(IsNotBlank("null"), ShouldEqual, true)
		So(IsNotBlank("0"), ShouldEqual, true)
		So(IsNotBlank(fmt.Sprint(0)), ShouldEqual, true)
		var p string
		So(IsNotBlank(p), ShouldEqual, false)
		So(IsNotBlank(""), ShouldEqual, false)
		So(IsNotBlank(" "), ShouldEqual, false)
	})
}
