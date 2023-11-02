package context

import (
	"farlogin/src/lib/util"
)

type messageType int

const (
	ERROR messageType = iota
	WARNING
	INFO
	SUCCESS
)

func MesageTypeToStr(msgType messageType) string {
	switch msgType {
	case ERROR:
		return "Error:"
	case WARNING:
		return "Warning:"
	case INFO:
		return "Info:"
	case SUCCESS:
		return "Success:"
	default:
		return "To Title:"
	}
}

//messages
type message struct {
	List        []string
	MsgType     messageType
	IsAutoClear bool
	Title       string
}

func newMessage(msgType messageType) *message {
	rv := new(message)
	rv.List = make([]string, 0, 10)
	rv.MsgType = msgType
	rv.IsAutoClear = true

	return rv
}

func (m *message) Add(str string) *message {
	m.List = append(m.List, str)

	return m
}

func (m *message) SetTitle(str string) *message {
	m.Title = str

	return m
}

func (m *message) Stick() *message {
	m.IsAutoClear = false

	return m
}

func (m *message) format(ctx *Ctx) *string {
	buf := util.NewBuf()

	var classStr string
	if m.IsAutoClear {
		classStr = " autoClear"
	}

	var title string
	if m.Title == "" {
		title = MesageTypeToStr(m.MsgType)
	} else {
		title = m.Title
	}

	switch m.MsgType {
	case ERROR:
		buf.Add("<div class=\"alert alertError%s\">", classStr)
		buf.Add("<i class=\"fa fa-times buttonClose\"></i>")
		buf.Add("<p><strong>%s</strong></p>", title)
	case WARNING:
		buf.Add("<div class=\"alert alertWarning%s\">", classStr)
		buf.Add("<i class=\"fa fa-times buttonClose\"></i>")
		buf.Add("<p><strong>%s</strong></p>", title)
	case INFO:
		buf.Add("<div class=\"alert alertInfo%s\">", classStr)
		buf.Add("<i class=\"fa fa-times buttonClose\"></i>")
		buf.Add("<p><strong>%s</strong></p>", title)
	case SUCCESS:
		buf.Add("<div class=\"alert alertSuccess%s\">", classStr)
		buf.Add("<i class=\"fa fa-times buttonClose\"></i>")
		buf.Add("<p><strong>%s</strong></p>", title)
	}

	for _, val := range m.List {
		buf.Add("<p>%s</p>", val)
	}

	buf.Add("</div>")

	return buf.String()
}

type MessageList struct {
	List []*message
}

func NewMessageList() *MessageList {
	rv := new(MessageList)
	rv.List = make([]*message, 0, 10)

	return rv
}

func (ml *MessageList) add(m *message) {
	ml.List = append(ml.List, m)
}

func (ml *MessageList) Format(ctx *Ctx) *string {
	buf := util.NewBuf()

	for _, v := range ml.List {
		buf.Add(*v.format(ctx))
	}

	return buf.String()
}

func (ml *MessageList) Error(str string) *message {
	rv := newMessage(ERROR)
	rv.List = append(rv.List, str)

	ml.add(rv)

	return rv
}

func (ml *MessageList) Warning(str string) *message {
	rv := newMessage(WARNING)
	rv.List = append(rv.List, str)

	ml.add(rv)

	return rv
}

func (ml *MessageList) Info(str string) *message {
	rv := newMessage(INFO)
	rv.List = append(rv.List, str)

	ml.add(rv)

	return rv
}

func (ml *MessageList) Success(str string) *message {
	rv := newMessage(SUCCESS)
	rv.List = append(rv.List, str)

	ml.add(rv)

	return rv
}
