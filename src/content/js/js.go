package js

import (
	"farlogin/src/lib/util"
)

type jsItem struct {
	Name    string
	IsDefer bool
}

type Js []*jsItem

func New() Js {
	rv := make([]*jsItem, 0, 10)

	return rv
}

func (js *Js) Add(fileName string) {
	item := new(jsItem)
	item.Name = fileName
	item.IsDefer = false

	*js = append(*js, item)
}

func (js *Js) AddDefer(fileName string) {
	item := new(jsItem)
	item.Name = fileName
	item.IsDefer = true

	*js = append(*js, item)
}

func (js Js) Format() *string {
	buf := util.NewBuf()

	var str string
	for _, v := range js {
		if v.IsDefer {
			str = "async=\"false\""
		} else {
			str = ""
		}

		buf.Add("<script %s src=\"%s\"></script>", str, v.Name)
	}

	return buf.StringSep("\n")
}
