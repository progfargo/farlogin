package tax

import (
	"fmt"
	"sort"
	"strconv"
)

type child struct {
	name string
	enum int64
}

type childList []*child

func (cl *childList) Add(name string, enum int64) {
	*cl = append(*cl, &child{name: name, enum: enum})
}

func (cl *childList) delete(name string) {
	pos := -1
	for key, val := range *cl {
		if val.name == name {
			pos = key
			break
		}
	}

	if pos != -1 {
		*cl = append((*cl)[0:pos], (*cl)[pos+1:]...)
	}
}

func (cl childList) Len() int {
	return len(cl)
}

func (cl childList) Swap(i, j int) {
	cl[i], cl[j] = cl[j], cl[i]
}

func (cl childList) Less(i, j int) bool {
	if cl[i].enum != cl[j].enum {
		return cl[i].enum < cl[j].enum
	} else {
		return cl[i].name > cl[j].name
	}
}

type TaxItem struct {
	name      string
	parent    string
	children  childList
	enum      int64
	isVisible bool
	Data      interface{}
}

func (ti *TaxItem) Name() string {
	return ti.name
}

func (ti *TaxItem) Enum() int64 {
	return ti.enum
}

func (ti *TaxItem) IsVisible() bool {
	return ti.isVisible
}

func newTaxItem(name, parent string, enum int64, isVisible bool, data interface{}) *TaxItem {
	rv := new(TaxItem)
	rv.name = name
	rv.parent = parent
	rv.isVisible = isVisible
	rv.Data = data
	rv.enum = enum
	rv.children = make(childList, 0, 10)

	return rv
}

type Tax struct {
	list map[string]*TaxItem
}

func New() *Tax {
	rv := new(Tax)
	rv.list = make(map[string]*TaxItem, 50)

	return rv
}

func (tax *Tax) Add(name, parent string, enum int64, isVisible bool, data interface{}) {
	if _, ok := tax.list[name]; ok {
		panic(fmt.Sprintf("Tax item already exists: %s", name))
	}

	if _, ok := tax.list[parent]; !ok && parent != "end" {
		panic(fmt.Sprintf("Parent item does not exist: %s", parent))
	}

	item := newTaxItem(name, parent, enum, isVisible, data)
	if parent != "end" {
		tax.list[parent].children.Add(name, enum)
	}

	tax.list[name] = item
}

func (tax *Tax) IsExists(name string) bool {
	_, ok := tax.list[name]

	return ok
}

func (tax *Tax) IsEmpty() bool {
	return len(tax.list) == 0
}

func (tax *Tax) Delete(name string) {
	if _, ok := tax.list[name]; !ok {
		panic(fmt.Sprintf("Item does not exist: %s", name))
	}

	if tax.IsParent(name) {
		panic("Parent taxonomy item can not be deleted.")
	}

	parent := tax.list[name].parent
	tax.list[parent].children.delete(name)
	delete(tax.list, name)
}

func (tax *Tax) GetItem(name string) *TaxItem {
	item, ok := tax.list[name]
	if !ok {
		panic(fmt.Sprintf("Tax item does not exists: %s", name))
	}

	return item
}

func (tax *Tax) GetAllParents(name string) []string {
	if _, ok := tax.list[name]; !ok {
		panic(fmt.Sprintf("Tax item does not exists: %s", name))
	}

	rv := make([]string, 0, 10)
	for parent := tax.list[name].parent; parent != "end"; parent = tax.list[parent].parent {
		rv = append(rv, parent)
	}

	return rv
}

func (tax *Tax) IsParent(name string) bool {
	item, ok := tax.list[name]
	if !ok {
		panic(fmt.Sprintf("Tax item does not exists: %s", name))
	}

	if len(item.children) > 0 {
		return true
	}

	return false
}

func (tax *Tax) GetChildren(name string) []string {
	item, ok := tax.list[name]
	if !ok {
		panic(fmt.Sprintf("Tax item does not exists: %s", name))
	}

	rv := make([]string, 0, 10)
	for _, val := range item.children {
		rv = append(rv, val.name)
	}

	return rv
}

func (tax *Tax) GetAllChildren(name string) []string {
	rv := make([]string, 0, 10)

	_, ok := tax.list[name]
	if !ok {
		return rv
	}

	tax.getAllChildrenR(name, &rv)

	return rv
}

func (tax *Tax) getAllChildrenR(name string, rv *[]string) {
	*rv = append(*rv, name)

	for _, val := range tax.list[name].children {
		tax.getAllChildrenR(val.name, rv)
	}
}

func (tax *Tax) SortChildren() {
	for key, val := range tax.list {
		if tax.IsParent(key) {
			sort.Sort(val.children)
		}
	}
}

func IntToName(id int64) string {
	if id == 0 {
		return "end"
	}

	if id == 1 {
		return "root"
	}

	return fmt.Sprintf("%d", id)
}

func NameToInt(id string) int64 {
	if id == "end" {
		return 0
	}

	if id == "root" {
		return 1
	}

	rv, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		rv = 0
	}

	return rv
}
