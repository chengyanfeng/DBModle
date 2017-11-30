package util

import (
	//"errors"
	"errors"
	"fmt"
	"strings"
	."testthree/models"
)

type Oracle struct {
	SqlBase
}

func (this *Oracle) Init(p P) {
	this.P = p
	this.P["fmt"] = "oracle"
	return
}

func (this *Oracle) Tables(owner string) (r string, e error) {
	sql := ""
	if IsEmpty(owner) {
		sql = "SELECT table_name FROM user_tables"
	} else {
		sql = fmt.Sprintf("select (owner||'.'||table_name) as table_name from all_tables where owner='%v'", owner)
	}
	r, e = HttpPost(Jdbc_proxy_url, nil, &P{"sql": sql, "db": JsonEncode(Orclconn)})
	return
}

func (this *Oracle) Select(tableName string,field string,where string) (r string, e error) {
	sql := ""

	sql = fmt.Sprintf("select *  from %v where %v='%v'", tableName,field,where)

	r, e = HttpPost(Jdbc_proxy_url, nil, &P{"sql": sql, "db": JsonEncode(Orclconn)})
	return
}
func (this *Oracle) SelectALL(tableName string) (r string, e error) {
	sql := ""
	sql = fmt.Sprintf("select *  from %v", tableName)
	r, e = HttpPost(Jdbc_proxy_url, nil, &P{"sql": sql, "db": JsonEncode(Orclconn)})
	return
}
func (this *Oracle) Add(tableName string, p P){

}


func (this *Oracle) TableInfo(table string) (info string, sample string, e error) {
	table = strings.ToUpper(table)
	index := strings.Index(table, ".")
	sql := ""
	if index > -1 {
		owner := table[:index]
		t := table[index+1:]
		sql = fmt.Sprintf("SELECT all_tab_cols.column_name,all_tab_cols.DATA_TYPE,(select COMMENTS from all_col_comments where owner='%v' and table_name = '%v' and COLUMN_NAME=all_tab_cols.column_name) COMMENTS FROM all_tab_cols where owner='%v' and table_name='%v'", owner, t, owner, t)
	} else {
		sql = fmt.Sprintf("SELECT user_tab_cols.column_name,user_tab_cols.DATA_TYPE,(select COMMENTS from user_col_comments where table_name = '%v' and COLUMN_NAME=user_tab_cols.column_name) COMMENTS FROM user_tab_cols where table_name='%v'", table, table)
	}
	body, err := HttpPost(Jdbc_proxy_url, nil, &P{"sql": sql, "db": JsonEncode(this.P)})
	info = body
	if err == nil && IsJson([]byte(body)) {
		tmp := []P{}
		tmp, err = JsonDecodeArray([]byte(body))
		if err != nil {
			e = err
			return
		}
		r := []P{}
		for _, v := range tmp {
			p := P{}
			p["o"] = v["COLUMN_NAME"]
			p["type"] = v["DATA_TYPE"]
			cmt := Trim(ToString(v["COMMENTS"]))
			if !IsEmpty(cmt) {
				cmt = Replace(cmt, []string{"("}, "（")
				cmt = Replace(cmt, []string{")"}, "）")
				cmt = Replace(cmt, []string{",", ";", "#", "?", "&", "=", "%", "＃"}, "")
				p["n"] = cmt
			}
			r = append(r, p)
		}
		info = JsonEncode(r)
		sql = fmt.Sprintf("select * from %v where rownum < %v", table, ToInt(this.P["limit"], 50))
		sample, _ = HttpPost(Jdbc_proxy_url, nil, &P{"sql": sql, "db": JsonEncode(this.P)})
		return
	} else {
		return body, "", err
	}
}

func (this *Oracle) Import(ds P) (body string, err error) {
	mode := ToInt(ds["mode"])
	sql := ToString(ds["sql"])
	if mode == 1 {
		table := ds["table"]
		sql = fmt.Sprintf("select * from (select %v.*,rownum from %v) where rownum<1000", table, table)
	}
	if !IsEmpty(sql) {
		this.OutFmt = "csv"
		body, err = this.Sql(sql)
		if err == nil {
			file := JoinStr("upload/", ds["_id"], ".csv")
			WriteFile(file, []byte(body))
			ds["url"] = file
		} else {
			err = errors.New(JoinStr(err, body))
		}
	} else {
		Debug("Import", ds, "append mode")
		tblname := ToString(ds["table"])
		pk := ToString(ds["pk"])
		if IsEmpty(pk) {
			err = errors.New("需要指定pk才可以增量导入")
			return
		}
		this.OutFmt = FMT_CSV
		pos := this.loadPos(tblname)
		if IsEmpty(pos) {
			sql = fmt.Sprintf("select min(%v) as %v from %v", pk, pk, tblname)
			body, err = this.Sql(sql)
			if err == nil {
				pos, err = this.getPos(body, pk)
			} else {
				return
			}
		}
		if !IsEmpty(pos) {
			where := fmt.Sprintf("where %v >= '%v' and rownum<1000", pk, pos)
			sql = fmt.Sprintf("select * from (select %v.*,rownum from %v) %v order by %v",
				tblname, tblname, where, pk)
			body, err = this.Sql(sql)
			if err == nil {
				file := JoinStr("upload/", ds["_id"], ".csv")
				WriteFile(file, []byte(body))
				pos, err = this.getPos(body, pk)
				this.savePos(tblname, pos)
			} else {
				return
			}
		} else {
			err = errors.New(fmt.Sprintf("Can't find min(%v) ", pk))
		}
	}
	return
}
