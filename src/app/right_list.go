package app

import (
	"fmt"
)

type Right struct {
	PageName string
	FuncName string
	Exp      string
}

func newRight(pageName, funcName, exp string) *Right {
	rv := new(Right)
	rv.PageName = pageName
	rv.FuncName = funcName
	rv.Exp = exp

	return rv
}

type RightTab struct {
	Title string
	List  []*Right
	rMap  map[string]*Right
}

func newRightTab(title string) *RightTab {
	rv := new(RightTab)

	rv.Title = title
	rv.List = make([]*Right, 0, 100)
	rv.rMap = make(map[string]*Right)

	return rv
}

func (rt *RightTab) Add(r *Right) {
	if _, ok := rt.rMap[MakeKey(r.PageName, r.FuncName)]; ok {
		panic("This item already exists. name: " + MakeKey(r.PageName, r.FuncName))
	}

	rt.List = append(rt.List, r)
	rt.rMap[MakeKey(r.PageName, r.FuncName)] = r
}

func (rt *RightTab) GetRight(pageName, funcName string) *Right {
	if _, ok := rt.rMap[MakeKey(pageName, funcName)]; !ok {
		panic(fmt.Sprintf("Can not find right. page name: %s function name: %s", pageName, funcName))
	}

	return rt.rMap[MakeKey(pageName, funcName)]
}

type rightList struct {
	List []*RightTab
	rMap map[string]bool
}

func newRightList() *rightList {
	rv := new(rightList)

	rv.List = make([]*RightTab, 0, 100)
	rv.rMap = make(map[string]bool)
	return rv
}

func (rl *rightList) Add(name string, rt *RightTab) {
	if _, ok := rl.rMap[name]; ok {
		panic("This right tab already exists.")
	}

	rl.List = append(rl.List, rt)
	rl.rMap[name] = true
}

func (rl *rightList) GetRightMap() map[string]*Right {
	rightMap := make(map[string]*Right, 100)

	var key string
	for _, rightTab := range rl.List {
		for _, right := range rightTab.List {
			key = MakeKey(right.PageName, right.FuncName)
			rightMap[key] = right
		}
	}

	return rightMap
}

func MakeKey(pageName, funcName string) string {
	return pageName + ":" + funcName
}

var UserRightList *rightList
var UserRightMap map[string]*Right

func SetRoleRightList() {
	UserRightList = newRightList()
	var rl = UserRightList
	var tab *RightTab

	//user
	tab = newRightTab("Users")
	tab.Add(newRight("user", "browse", "Can browse user records."))
	tab.Add(newRight("user", "select", "Can select user as effective user."))
	tab.Add(newRight("user", "insert", "Can insert user record."))
	tab.Add(newRight("user", "update", "Can update user records."))
	tab.Add(newRight("user", "update_pass", "Can update user passowrd."))
	tab.Add(newRight("user", "delete", "Can delete user records."))
	tab.Add(newRight("user", "role_browse", "Can browse user roles."))
	tab.Add(newRight("user", "role_revoke", "Can revoke user roles."))
	tab.Add(newRight("user", "role_grant", "Can grant user roles."))
	rl.Add("user", tab)

	//role
	tab = newRightTab("Roles")
	tab.Add(newRight("role", "browse", "Can browse role records."))
	tab.Add(newRight("role", "insert", "Can insert role record."))
	tab.Add(newRight("role", "update", "Can update role records."))
	tab.Add(newRight("role", "delete", "Can delete role records."))
	tab.Add(newRight("role", "role_right", "Can update role rights."))
	rl.Add("role", tab)

	//config
	tab = newRightTab("Configuration")
	tab.Add(newRight("config", "browse", "Can browse configuration records."))
	tab.Add(newRight("config", "set", "Can set configuration record value."))
	rl.Add("config", tab)

	//profile
	tab = newRightTab("Profile")
	tab.Add(newRight("profile", "browse", "Can browse own profile record."))
	tab.Add(newRight("profile", "update", "Can update own profile record."))
	tab.Add(newRight("profile", "update_password", "Can update own password."))
	rl.Add("profile", tab)

	//node
	tab = newRightTab("Node")
	tab.Add(newRight("node", "browse", "Can browse own node records."))
	tab.Add(newRight("node", "insert", "Can insert node records."))
	tab.Add(newRight("node", "update", "can update node records."))
	tab.Add(newRight("node", "delete", "Can delete node records."))
	rl.Add("node", tab)

	UserRightMap = UserRightList.GetRightMap() //to set AppRightMap
}

func makeSet(list ...string) map[string]bool {
	rv := make(map[string]bool)

	for _, v := range list {
		rv[v] = true
	}

	return rv
}
