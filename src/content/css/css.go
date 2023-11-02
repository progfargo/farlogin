package css

import (
	"farlogin/src/lib/util"
)

type Css struct {
	Data []string
}

func New() *Css {
	rv := new(Css)
	rv.Data = make([]string, 0, 10)

	return rv
}

func (css *Css) Add(fileName string) {
	css.Data = append(css.Data, fileName)
}

func (css *Css) Format() *string {
	buf := util.NewBuf()

	for _, v := range css.Data {
		buf.Add("<link rel=\"stylesheet\" href=\"%s\">", v)
	}

	return buf.StringSep("\n")
}
