package top_menu

import (
	"fmt"

	"farlogin/src/lib/context"
	"farlogin/src/lib/tax"
	"farlogin/src/lib/util"
)

type topMenuItem struct {
	url      string
	label    string
	icon     string
	isActive bool
}

type topMenu struct {
	tax *tax.Tax
}

func New() *topMenu {
	rv := new(topMenu)
	rv.tax = tax.New()
	rv.Add("root", "end", 0, true, "", "", "", false)

	return rv
}

func (tm *topMenu) Add(name, parent string, enum int64, isVisible bool, url, label, icon string, isActive bool) {
	tm.tax.Add(name, parent, enum, isVisible, &topMenuItem{url, label, icon, isActive})
}

func (tm *topMenu) Set(ctx *context.Ctx, name ...string) {
	tm.Add("profile", "root", 10, ctx.IsRight("profile", "browse"), ctx.U("/profile"), "Profile", "fas fa-user", false)
	tm.Add("logout", "root", 30, ctx.IsLoggedIn(), ctx.U("/logout"), "Logout", "fas fa-sign-out-alt", false)

	if len(name) > 1 {
		panic("wrong number of parameters.")
	}

	if len(name) == 1 {
		tm.setActive(name[0])
	}

	tm.tax.SortChildren()
	tm.reduce("root")
	tm.format(ctx)
}

func (tm *topMenu) setActive(name string) {
	item := tm.tax.GetItem(name)
	data := item.Data.(*topMenuItem)
	data.isActive = true
}

func (tm *topMenu) reduce(name string) {
	item := tm.tax.GetItem(name)

	if tm.tax.IsParent(name) {
		children := tm.tax.GetChildren(name)
		for _, val := range children {
			tm.reduce(val)
		}

		if item.IsVisible() {
			return
		}

		if !tm.tax.IsParent(name) {
			tm.tax.Delete(name)

			return
		}

		//not visible and still parent
		children = tm.tax.GetChildren(name)
		firstChild := children[0]
		firstChildItem := tm.tax.GetItem(firstChild)

		nameData := item.Data.(*topMenuItem)
		firstChildData := firstChildItem.Data.(*topMenuItem)

		nameData.url = firstChildData.url
		return
	}

	if !item.IsVisible() {
		tm.tax.Delete(name)
	}
}

func (tm *topMenu) format(ctx *context.Ctx) {
	rv := util.NewBuf()

	children := tm.tax.GetChildren("root")

	if len(children) == 0 {
		return
	}

	rv.Add("<div class=\"topMenu\">")

	first := true
	sep := ""
	for _, v := range children {
		item := tm.tax.GetItem(v)
		data := item.Data.(*topMenuItem)

		if first {
			first = false
		} else {
			sep = " |"
		}

		class := ""
		if data.isActive {
			class = " class=\"active\""
		}

		icon := ""
		if data.icon != "" {
			icon = fmt.Sprintf("<i class=\"%s fa-fw left\"></i>", data.icon)
		}

		rv.Add("%s<a href=\"%s\"%s>%s%s</a>", sep, data.url, class, icon, data.label)
	}

	rv.Add("</div>")

	ctx.AddHtml("topMenu", rv.String())
}
