package text_content_menu

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

type textMenu struct {
	tax       *tax.Tax
	cargoList []string
}

func New(cargoVar ...string) *textMenu {
	rv := new(textMenu)
	rv.tax = tax.New()
	rv.Add("root", "end", 0, true, "", "", "")
	rv.cargoList = cargoVar

	return rv
}

func (tcm *textMenu) Add(name, parent string, enum int64, isVisible bool, url, label, icon string) {
	tcm.tax.Add(name, parent, enum, isVisible, &menuItem{url: url, label: label, icon: icon})
}

func (tcm *textMenu) Set(ctx *context.Ctx, name ...string) {
	tcm.Add("text_content_display", "root", 10, ctx.IsRight("text_content", "browse"), ctx.U("/text_content_display", tcm.cargoList...), "Display", "fas fa-eye")
	tcm.Add("text_content_set", "root", 20, ctx.IsRight("text_content", "set"), ctx.U("/text_content_set", tcm.cargoList...), "Set", "fas fa-edit")
	tcm.Add("text_content_image", "root", 30, ctx.IsRight("text_content", "update"), ctx.U("/text_content_image", tcm.cargoList...), "Images", "fas fa-images")
	tcm.Add("text_content_update", "root", 40, ctx.IsRight("text_content", "update"), ctx.U("/text_content_update", tcm.cargoList...), "Update", "fas fa-edit")
	tcm.Add("text_content_delete", "root", 50, ctx.IsRight("text_content", "delete"), ctx.U("/text_content_delete", tcm.cargoList...), "Delete", "fas fa-trash-can")

	if len(name) > 1 {
		panic("wrong number of parameters.")
	}

	if len(name) == 1 {
		tcm.setActive(name[0])
	}

	tcm.tax.SortChildren()
	tcm.reduce("root")
}

func (tcm *textMenu) setActive(name string) {
	item := tcm.tax.GetItem(name)
	data := item.Data.(*menuItem)
	data.isActive = true
}

func (tcm *textMenu) reduce(name string) {
	item := tcm.tax.GetItem(name)

	if tcm.tax.IsParent(name) {
		children := tcm.tax.GetChildren(name)
		for _, val := range children {
			tcm.reduce(val)
		}

		if item.IsVisible() {
			return
		}

		if !tcm.tax.IsParent(name) {
			tcm.tax.Delete(name)

			return
		}

		//not visible and still parent
		children = tcm.tax.GetChildren(name)
		firstChild := children[0]
		firstChildItem := tcm.tax.GetItem(firstChild)

		nameData := item.Data.(*menuItem)
		firstChildData := firstChildItem.Data.(*menuItem)

		nameData.url = firstChildData.url
		return
	}

	if !item.IsVisible() {
		tcm.tax.Delete(name)
	}
}

func (tcm *textMenu) Format(ctx *context.Ctx) string {
	rv := util.NewBuf()

	children := tcm.tax.GetChildren("root")

	if len(children) == 0 {
		return ""
	}

	rv.Add("<div class=\"buttonBar textMenu\">")

	for _, v := range children {
		item := tcm.tax.GetItem(v)
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
