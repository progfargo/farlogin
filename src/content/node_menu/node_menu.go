package node_menu

import (
	"fmt"

	"farlogin/src/lib/context"
	"farlogin/src/lib/tax"
	"farlogin/src/lib/util"
)

type menuItem struct {
	url      string
	label    string
	icon     string
	isActive bool
}

type nodeMenu struct {
	tax       *tax.Tax
	cargoList []string
}

func New(cargoVar ...string) *nodeMenu {
	rv := new(nodeMenu)
	rv.tax = tax.New()
	rv.Add("root", "end", 0, true, "", "", "")
	rv.cargoList = cargoVar

	return rv
}

func (nm *nodeMenu) Add(name, parent string, enum int64, isVisible bool, url, label, icon string) {
	nm.tax.Add(name, parent, enum, isVisible, &menuItem{url: url, label: label, icon: icon})
}

func (nm *nodeMenu) Set(ctx *context.Ctx, name ...string) {
	nm.Add("node_display", "root", 10, true, ctx.U("/node_display", nm.cargoList...), "Display", "fas fa-info-circle")
	nm.Add("node_update", "root", 20, true, ctx.U("/node_update", nm.cargoList...), "Edit", "fas fa-edit")
	nm.Add("node_new_session", "root", 30, true, ctx.U("/node_new_session", nm.cargoList...), "New Session", "fas fa-plug")
	nm.Add("node_delete_all_session", "root", 40, true, ctx.U("/node_delete_all_session", nm.cargoList...), "Delete Sessions", "fas fa-trash-alt")

	if len(name) > 1 {
		panic("wrong number of parameters.")
	}

	if len(name) == 1 {
		nm.setActive(name[0])
	}

	nm.tax.SortChildren()
	nm.reduce("root")
}

func (nm *nodeMenu) setActive(name string) {
	item := nm.tax.GetItem(name)
	data := item.Data.(*menuItem)
	data.isActive = true
}

func (nm *nodeMenu) reduce(name string) {
	item := nm.tax.GetItem(name)

	if nm.tax.IsParent(name) {
		children := nm.tax.GetChildren(name)
		for _, val := range children {
			nm.reduce(val)
		}

		if item.IsVisible() {
			return
		}

		if !nm.tax.IsParent(name) {
			nm.tax.Delete(name)

			return
		}

		//not visible and still parent
		children = nm.tax.GetChildren(name)
		firstChild := children[0]
		firstChildItem := nm.tax.GetItem(firstChild)

		nameData := item.Data.(*menuItem)
		firstChildData := firstChildItem.Data.(*menuItem)

		nameData.url = firstChildData.url
		return
	}

	if !item.IsVisible() {
		nm.tax.Delete(name)
	}
}

func (nm *nodeMenu) Format(ctx *context.Ctx) string {
	rv := util.NewBuf()

	children := nm.tax.GetChildren("root")

	if len(children) == 0 {
		return ""
	}

	rv.Add("<div class=\"buttonBar nodeMenu\">")

	for _, v := range children {
		item := nm.tax.GetItem(v)
		data := item.Data.(*menuItem)

		class := ""
		if data.isActive {
			class = " selected"
		}

		icon := ""
		if data.icon != "" {
			icon = fmt.Sprintf("<i class=\"%s fa-fw left\"></i>", data.icon)
		}

		rv.Add("<a href=\"%s\" class=\"button buttonDefault buttonSm%s\">%s%s</a>", data.url, class, icon, data.label)
	}

	rv.Add("</div>")

	return *rv.String()
}
