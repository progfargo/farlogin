package profile_menu

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

func (pm *textMenu) Add(name, parent string, enum int64, isVisible bool, url, label, icon string) {
	pm.tax.Add(name, parent, enum, isVisible, &menuItem{url: url, label: label, icon: icon})
}

func (pm *textMenu) Set(ctx *context.Ctx, name ...string) {
	pm.Add("profile", "root", 10, ctx.IsRight("profile", "browse"), ctx.U("/profile", pm.cargoList...), "Profile", "fas fa-user")
	pm.Add("update_password", "root", 40, ctx.IsRight("profile", "update_password"), ctx.U("/profile_update_pass", pm.cargoList...), "Change Password", "fas fa-key")

	if len(name) > 1 {
		panic("wrong number of parameters.")
	}

	if len(name) == 1 {
		pm.setActive(name[0])
	}

	pm.tax.SortChildren()
	pm.reduce("root")
}

func (pm *textMenu) setActive(name string) {
	item := pm.tax.GetItem(name)
	data := item.Data.(*menuItem)
	data.isActive = true
}

func (pm *textMenu) reduce(name string) {
	item := pm.tax.GetItem(name)

	if pm.tax.IsParent(name) {
		children := pm.tax.GetChildren(name)
		for _, val := range children {
			pm.reduce(val)
		}

		if item.IsVisible() {
			return
		}

		if !pm.tax.IsParent(name) {
			pm.tax.Delete(name)

			return
		}

		//not visible and still parent
		children = pm.tax.GetChildren(name)
		firstChild := children[0]
		firstChildItem := pm.tax.GetItem(firstChild)

		nameData := item.Data.(*menuItem)
		firstChildData := firstChildItem.Data.(*menuItem)

		nameData.url = firstChildData.url
		return
	}

	if !item.IsVisible() {
		pm.tax.Delete(name)
	}
}

func (pm *textMenu) Format(ctx *context.Ctx) string {
	rv := util.NewBuf()

	children := pm.tax.GetChildren("root")

	if len(children) == 0 {
		return ""
	}

	rv.Add("<div class=\"buttonBar textMenu\">")

	for _, v := range children {
		item := pm.tax.GetItem(v)
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
