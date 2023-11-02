package combo

import (
	"database/sql"
	"fmt"
	"strings"

	"farlogin/src/app"
	"farlogin/src/lib/context"
	"farlogin/src/lib/tax"
	"farlogin/src/lib/util"
)

type taxRec struct {
	name    int64
	parent  int64
	label   sql.NullString
	enum    int64
	isAdded bool
}

type taxComboItem struct {
	key   string
	label string
}

type TaxCombo struct {
	sqlStr       string
	defaultLabel string
	tax          *tax.Tax
}

func NewTaxCombo(sqlStr, defaultLabel string) *TaxCombo {
	rv := new(TaxCombo)
	rv.tax = tax.New()

	rv.sqlStr = sqlStr
	rv.defaultLabel = defaultLabel

	return rv
}

func (tc *TaxCombo) add(name, parent string, enum int64, isVisible bool, label string) {
	tc.tax.Add(name, parent, enum, isVisible, &taxComboItem{name, label})
}

func (tc *TaxCombo) IsEmpty() bool {
	return tc.tax.IsEmpty()
}

func (tc *TaxCombo) Set(ctx *context.Ctx) {
	rows, err := app.Db.Query(tc.sqlStr)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	list := make(map[int64]*taxRec)

	for rows.Next() {
		rec := new(taxRec)
		if err = rows.Scan(&rec.name, &rec.parent, &rec.label, &rec.enum); err != nil {
			panic(err)
		}

		list[rec.name] = rec
	}

	for k, _ := range list {
		tc.addItemToTax(ctx, list, k)
	}

	tc.tax.SortChildren()
}

func (tc *TaxCombo) addItemToTax(ctx *context.Ctx, list map[int64]*taxRec, item int64) {
	if item == 0 || list[item].isAdded {
		return
	}

	//if there is no parent, ad it first
	if !tc.tax.IsExists(tax.IntToName(list[item].parent)) {
		tc.addItemToTax(ctx, list, list[item].parent)
	}

	var label string
	if list[item].label.Valid {
		label = list[item].label.String
	} else {
		label = ""
	}

	if label == "root" {
		label = tc.defaultLabel
	}

	tc.add(tax.IntToName(list[item].name), tax.IntToName(list[item].parent), list[item].enum, true, label)
	list[item].isAdded = true
}

func (tc *TaxCombo) Format(curKey string) string {
	buf := util.NewBuf()
	isLast := make([]bool, 0, 10)
	isLast = append(isLast, true)

	tc.FormatItem("root", curKey, buf, 0)

	return *buf.String()
}

func (tc *TaxCombo) FormatItem(name, curKey string, buf *util.Buf, tab int) {
	item := tc.tax.GetItem(name)
	data := item.Data.(*taxComboItem)

	var selectedStr = ""
	if data.key == curKey {
		selectedStr = " selected"
	}

	buf.Add("<option value=\"%d\"%s>", tax.NameToInt(name), selectedStr)

	str := strings.Repeat("&nbsp;", tab)

	buf.Add(str + "&#9726; " + data.label)
	buf.Add("</option>")

	if tc.tax.IsParent(name) {
		children := tc.tax.GetChildren(name)

		for _, v := range children {
			tc.FormatItem(v, curKey, buf, tab+6)
		}
	}
}

func (tc *TaxCombo) FormatDir(name string) string {
	parentList := tc.tax.GetAllParents(name)

	var dirList []string
	for i := len(parentList) - 1; i >= 0; i-- {
		item := tc.tax.GetItem(parentList[i])
		data := item.Data.(*taxComboItem)
		dirList = append(dirList, data.label)
	}

	item := tc.tax.GetItem(name)
	data := item.Data.(*taxComboItem)

	return fmt.Sprintf("%s > <strong>%s</strong>", strings.Join(dirList, " > "), data.label)
}
