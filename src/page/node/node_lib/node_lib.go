package node_lib

import (
	"fmt"
	"regexp"

	"farlogin/src/app"
	"farlogin/src/lib/context"
	"farlogin/src/lib/util"
)

type NodeRec struct {
	NodeId   int64
	Name     string
	Exp      string
	LastSeen int64
}

func GetNodeRec(nodeId int64) (*NodeRec, error) {
	sqlBuf := util.NewBuf()
	sqlBuf.Add(`select
					node.nodeId,
					node.name,
					node.exp,
					node.lastSeen
				from
					node`)

	conBuf := util.NewBuf()
	conBuf.Add("(node.nodeId = %d)", nodeId)

	sqlBuf.Add("where")
	sqlBuf.Add(*conBuf.StringSep("and"))

	row := app.Db.QueryRow(*sqlBuf.String())

	rec := new(NodeRec)
	err := row.Scan(&rec.NodeId, &rec.Name, &rec.Exp, &rec.LastSeen)

	if err != nil {
		return nil, err
	}

	return rec, nil
}

func CountNode(key string) int64 {
	sqlBuf := util.NewBuf()
	sqlBuf.Add(`select
					count(*)`)

	fromBuf := util.NewBuf()
	fromBuf.Add("node")

	conBuf := util.NewBuf()

	if key != "" {
		key = util.DbStr(key)
		conBuf.Add(`(node.name like('%%%s%%'))`, key)
	}

	sqlBuf.Add("from")
	sqlBuf.Add(*fromBuf.StringSep(", "))

	if !conBuf.IsEmpty() {
		sqlBuf.Add("where")
		sqlBuf.Add(*conBuf.StringSep("and"))
	}

	row := app.Db.QueryRow(*sqlBuf.String())

	var rv int64
	err := row.Scan(&rv)
	if err != nil {
		panic(err)
	}

	return rv
}

func GetNodePage(ctx *context.Ctx, key string, pageNo int64) []*NodeRec {
	sqlBuf := util.NewBuf()
	sqlBuf.Add(`select
					node.nodeId,
					node.name,
					node.exp,
					node.lastSeen`)

	fromBuf := util.NewBuf()
	fromBuf.Add("node")

	conBuf := util.NewBuf()

	if key != "" {
		key = util.DbStr(key)
		conBuf.Add(`(node.name like('%%%s%%'))`, key)
	}

	sqlBuf.Add("from")
	sqlBuf.Add(*fromBuf.StringSep(", "))

	if !conBuf.IsEmpty() {
		sqlBuf.Add("where")
		sqlBuf.Add(*conBuf.StringSep(" and "))
	}

	sqlBuf.Add("order by nodeId")

	pageLen := ctx.Config.Int("pageLen")
	start := (pageNo - 1) * pageLen
	sqlBuf.Add("limit %d, %d", start, pageLen)

	rows, err := app.Db.Query(*sqlBuf.String())
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	rv := make([]*NodeRec, 0, 100)
	for rows.Next() {
		rec := new(NodeRec)
		err = rows.Scan(&rec.NodeId, &rec.Name, &rec.Exp, &rec.LastSeen)
		if err != nil {
			panic(err)
		}

		rv = append(rv, rec)
	}

	return rv
}

type NodeSessionRec struct {
	NodeSessionId int64
	RecordTime    int64
	SessionHash   string
	Status        string
}

func CountNodeSessionList(nodeId int64) int64 {
	sqlStr := `select
					count(*)
				from
					nodeSession
				where
					nodeSession.nodeId = ?`

	row := app.Db.QueryRow(sqlStr, nodeId)

	var rv int64
	err := row.Scan(&rv)
	if err != nil {
		panic(err)
	}

	return rv
}

func GetNodeSessionList(nodeId int64) []*NodeSessionRec {
	sqlStr := `select
					nodeSession.nodeSessionId,
					nodeSession.recordTime,
					nodeSession.sessionHash,
					nodeSession.status
				from
					nodeSession
				where
					nodeSession.nodeId = ?
				order by nodeSession.recordTime`

	rows, err := app.Db.Query(sqlStr, nodeId)
	if err != nil {
		panic(err)
	}

	defer rows.Close()

	rv := make([]*NodeSessionRec, 0, 100)
	for rows.Next() {
		rec := new(NodeSessionRec)
		err = rows.Scan(&rec.NodeSessionId, &rec.RecordTime, &rec.SessionHash, &rec.Status)
		if err != nil {
			panic(err)
		}

		rv = append(rv, rec)
	}

	return rv
}

func GetSessionByNodeName(nodeName string) (*NodeSessionRec, error) {
	sqlStr := `select
					nodeSession.nodeSessionId,
					nodeSession.recordTime,
					nodeSession.sessionHash,
					nodeSession.status
				from
					nodeSession,
					node
				where
					node.nodeId = nodeSession.nodeId and
					node.name = ?`

	row := app.Db.QueryRow(sqlStr, nodeName)

	rec := new(NodeSessionRec)
	err := row.Scan(&rec.NodeSessionId, &rec.RecordTime, &rec.SessionHash, &rec.Status)

	if err != nil {
		return nil, err
	}

	return rec, nil
}

func GetSessionByHash(sessionHash string) (*NodeSessionRec, error) {
	sqlStr := `select
					nodeSession.nodeSessionId,
					nodeSession.recordTime,
					nodeSession.sessionHash,
					nodeSession.status
				from
					nodeSession
				where
					nodeSession.sessionHash = ?`

	row := app.Db.QueryRow(sqlStr, sessionHash)

	rec := new(NodeSessionRec)
	err := row.Scan(&rec.NodeSessionId, &rec.RecordTime, &rec.SessionHash, &rec.Status)
	if err != nil {
		return nil, err
	}

	return rec, nil
}

func CheckName(name string) bool {
	r, _ := regexp.Compile("^[a-z0-9][a-z0-9.]+[a-z0-9]$")

	return r.MatchString(name)
}

func IsNodeOn(lastSeen int64) bool {
	now := util.Now()

	return now-lastSeen < 20
}

func NodeStatusToLabel(lastSeen int64) string {

	if IsNodeOn(lastSeen) {
		return fmt.Sprintf("<span class=\"label labelSuccess labelXs\">on</span>")
	} else {
		return fmt.Sprintf("<span class=\"label labelWarning labelXs\">off</span>")
	}
}

func SessionStatusToLabel(str string) string {
	switch str {
	case "new":
		return fmt.Sprintf("<span class=\"label labelInfo labelXs\">new</span>")
	case "used":
		return fmt.Sprintf("<span class=\"label labelWarning labelXs\">used</span>")
	default:
		return fmt.Sprintf("<span class=\"label labelDefault labelXs\">unknown</span>")
	}
}
