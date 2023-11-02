package combo

import (
	"farlogin/src/lib/util"
)

type enumComboItem struct {
	key   string
	label string
}

type enumCombo struct {
	defaultOption string
	list          []*enumComboItem
}

func NewEnumCombo() *enumCombo {
	rv := new(enumCombo)
	rv.list = make([]*enumComboItem, 0, 50)

	return rv
}

func (ec *enumCombo) Add(key, val string) {
	ec.list = append(ec.list, &enumComboItem{key, val})
}

func (ec *enumCombo) Format(curKey string) string {
	buf := util.NewBuf()

	var selectedStr string
	for _, item := range ec.list {
		selectedStr = ""
		if curKey == item.key {
			selectedStr = " selected"
		}

		buf.Add("<option value=\"%s\"%s>%s</option>", item.key, selectedStr, item.label)
	}

	return *buf.String()
}
