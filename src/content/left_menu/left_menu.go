package left_menu

import (
	"fmt"
	"strings"

	"farlogin/src/lib/context"
	"farlogin/src/lib/tax"
	"farlogin/src/lib/util"
)

type leftMenuItem struct {
	url           string
	label         string
	icon          string
	isActive      bool
	isSubmenuOpen bool
}

type leftMenu struct {
	tax *tax.Tax
}

func New() *leftMenu {
	rv := new(leftMenu)
	rv.tax = tax.New()
	rv.Add("root", "end", 0, true, "", "", "")

	return rv
}

func (lm *leftMenu) Add(name, parent string, enum int64, isVisible bool, url, label, icon string) {
	lm.tax.Add(name, parent, enum, isVisible, &leftMenuItem{url, label, icon, false, false})
}

func (lm *leftMenu) Set(ctx *context.Ctx, name ...string) {
	lm.Add("admin", "root", 10, ctx.IsRight("user", "browse"), ctx.U("/user"), "Admin", "fas fa-gear")

	lm.Add("user_", "admin", 10, ctx.IsRight("user", "browse"), ctx.U("/user"), "User Management", "fas fa-user")
	lm.Add("user", "user_", 10, ctx.IsRight("user", "browse"), ctx.U("/user"), "Users", "fas fa-user-cog")
	lm.Add("role", "user_", 20, ctx.IsRight("role", "browse"), ctx.U("/role"), "Roles", "fas fa-users-gear")

	lm.Add("config", "admin", 30, ctx.IsRight("config", "browse"), ctx.U("/config"), "App Configuration", "fas fa-gauge")

	lm.Add("node", "root", 30, ctx.IsRight("node", "browse"), ctx.U("/node"), "Node List", "fas fa-microchip")

	if len(name) > 1 {
		panic("wrong number of parameters.")
	}

	if len(name) == 1 {
		lm.setActive(name[0])
	}

	lm.tax.SortChildren()
	lm.reduce("root")
	lm.format(ctx)
}

func (lm *leftMenu) setActive(name string) {
	allParents := lm.tax.GetAllParents(name)
	allParents = append(allParents, name)

	for _, val := range allParents {
		item := lm.tax.GetItem(val)
		data := item.Data.(*leftMenuItem)
		data.isSubmenuOpen = true

		if val == name {
			data.isActive = true
		}
	}
}

func (lm *leftMenu) reduce(name string) {
	item := lm.tax.GetItem(name)

	if lm.tax.IsParent(name) {
		children := lm.tax.GetChildren(name)
		for _, val := range children {
			lm.reduce(val)
		}

		if item.IsVisible() {
			return
		}

		if !lm.tax.IsParent(name) {
			lm.tax.Delete(name)

			return
		}

		//not visible and still parent
		children = lm.tax.GetChildren(name)
		firstChild := children[0]
		firstChildItem := lm.tax.GetItem(firstChild)

		nameData := item.Data.(*leftMenuItem)
		firstChildData := firstChildItem.Data.(*leftMenuItem)

		nameData.url = firstChildData.url
		return
	}

	if !item.IsVisible() {
		lm.tax.Delete(name)
	}
}

func (lm *leftMenu) format(ctx *context.Ctx) {
	rv := util.NewBuf()

	children := lm.tax.GetChildren("root")

	if len(children) == 0 {
		return
	}

	rv.Add("<nav>")
	rv.Add("<ul>")

	for _, v := range children {
		lm.FormatItem(v, rv)
	}

	rv.Add("</ul>")
	rv.Add("</nav>")

	ctx.AddHtml("leftMenu", rv.String())
}

func (lm *leftMenu) FormatItem(name string, rv *util.Buf) {
	item := lm.tax.GetItem(name)
	data := item.Data.(*leftMenuItem)

	classList := make([]string, 0, 5)
	if data.isActive {
		classList = append(classList, "active")
	}

	icon := ""
	if data.icon != "" {
		icon = fmt.Sprintf("<i class=\"%s fa-fw menuIcon left\"></i>", data.icon)
	}

	if lm.tax.IsParent(name) {
		classList = append(classList, "subMenu")

		var subMenuShow, subMenuIcon string
		if data.isSubmenuOpen {
			subMenuShow = " class=\"subMenuShow\""
			subMenuIcon = "<i class=\"subMenuButton fas fa-minus fa-fw\">"
		} else {
			subMenuShow = ""
			subMenuIcon = "<i class=\"subMenuButton fas fa-plus fa-fw\">"
		}

		classStr := ""
		if len(classList) > 0 {
			classStr = fmt.Sprintf(" class=\"%s\"", strings.Join(classList, " "))

		}

		rv.Add("<li%s>", classStr)
		rv.Add("<a href=\"%s\">%s%s%s"+
			"</i></a>", data.url, icon, data.label, subMenuIcon)
		rv.Add("<ul%s>", subMenuShow)

		children := lm.tax.GetChildren(name)
		for _, v := range children {
			lm.FormatItem(v, rv)
		}

		rv.Add("</ul>")
		rv.Add("</li>")
	} else {
		classStr := ""
		if len(classList) > 0 {
			classStr = fmt.Sprintf(" class=\"%s\"", strings.Join(classList, " "))
		}

		rv.Add("<li%s>", classStr)
		rv.Add("<a href=\"%s\">%s%s</a>", data.url, icon, data.label)
		rv.Add("</li>")
	}
}
