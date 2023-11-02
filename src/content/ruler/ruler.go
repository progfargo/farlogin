package ruler

import (
	"strconv"

	"farlogin/src/lib/context"
	"farlogin/src/lib/util"
)

type ruler struct {
	pageNo    int64
	totalPage int64
	initUrl   string
	urlStr    string
	list      []*rulerNode
}

type rulerNode struct {
	label  string
	urlStr string
	active bool
}

func NewRuler(totalPage, pageNo int64, initUrl string) *ruler {
	rv := new(ruler)
	rv.initUrl = initUrl
	rv.pageNo = pageNo
	rv.totalPage = totalPage

	rv.list = make([]*rulerNode, 0, 20)
	return rv
}

func (rl *ruler) IsEmpty() bool {
	if rl.totalPage < 2 {
		return true
	}

	return false
}

func (rl *ruler) Set(ctx *context.Ctx) {

	if rl.totalPage < 2 {
		return
	}

	rulerLen := ctx.Config.Int("rulerLen")

	start := rl.pageNo - rulerLen
	if start < 1 {
		start = 1
	}

	end := start + rulerLen*2
	if end > rl.totalPage {
		start -= end - rl.totalPage
		if start < 1 {
			start = 1
		}
		end = rl.totalPage
	}

	if start > 1 {
		node := new(rulerNode)
		node.label = "First"
		ctx.Cargo.SetInt("pn", 1)
		node.urlStr = ctx.U(rl.initUrl, "pn")
		rl.list = append(rl.list, node)
	}

	if start > 1 {
		node := new(rulerNode)
		node.label = "&lt;&lt;"
		ctx.Cargo.SetInt("pn", rl.pageNo-rulerLen)
		node.urlStr = ctx.U(rl.initUrl, "pn")
		rl.list = append(rl.list, node)
	}

	for i := start; i <= end; i++ {
		node := new(rulerNode)
		if i == rl.pageNo {
			node.active = true
		}

		node.label = strconv.FormatInt(i, 10)
		ctx.Cargo.SetInt("pn", i)
		node.urlStr = ctx.U(rl.initUrl, "pn")
		rl.list = append(rl.list, node)
	}

	if end < rl.totalPage {
		node := new(rulerNode)
		node.label = "&gt;&gt;"
		ctx.Cargo.SetInt("pn", rl.pageNo+rulerLen)
		node.urlStr = ctx.U(rl.initUrl, "pn")
		rl.list = append(rl.list, node)
	}

	if end < rl.totalPage {
		node := new(rulerNode)
		node.label = "Last"
		ctx.Cargo.SetInt("pn", rl.totalPage)
		node.urlStr = ctx.U(rl.initUrl, "pn")
		rl.list = append(rl.list, node)
	}

	ctx.Cargo.SetInt("pn", rl.pageNo)
}

func (rl *ruler) Format() string {
	rv := util.NewBuf()

	rv.Add("<div class=\"pagination\">")
	rv.Add("<ul class=\"clear\">")

	var classStr string
	for _, val := range rl.list {
		classStr = ""
		if val.active {
			classStr = " class=\"active\""
		}

		rv.Add("<li%s><a href=\"%s\">%s</a></li>", classStr, val.urlStr, val.label)
	}

	rv.Add("</ul>")
	rv.Add("</div>")

	return *rv.String()
}
