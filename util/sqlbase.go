package util

import (
	"errors"
)

type SqlBase struct {
	P      P
	OutFmt string
}

func (this *SqlBase) Sql(sql string) (r string, err error) {
	p := &P{"sql": sql, "db": JsonEncode(this.P), "o": this.OutFmt}
	Debug("Sql:", sql, Jdbc_proxy_url)
	r, err = HttpPost(Jdbc_proxy_url, nil, p)
	return
}

func (this *SqlBase) Error(err error) {
	Error(err)
}

func (this *SqlBase) loadPos(name string) (r string) {
	p := *D(DbPos).Find(P{"key": this.getStoreKey(name)}).One()
	if !IsEmpty(p["_id"]) {
		r = ToString(p["pos"])
	}
	Debug("loadPos", name, r)
	return
}

func (this *SqlBase) savePos(name string, pos string) {
	Debug("savePos", name, pos)
	D(DbPos).Upsert(P{"key": this.getStoreKey(name)}, P{"key": this.getStoreKey(name), "pos": pos})
}

func (this *SqlBase) ClearPos(name string) {
	D(DbPos).Remove(P{"key": this.getStoreKey(name)})
}

func (this *SqlBase) getStoreKey(name string) string {
	return Md5(this.P, name)
}

func (this *SqlBase) getPos(data string, pk string) (string, error) {
	csv := Csv{Body: data, Fh: true, Split: ","}
	csv.Scan(nil)
	json := csv.Data
	if len(json) < 1 {
		return "", errors.New("数据表记录数为0")
	}
	head := csv.Head
	last := json[len(json)-1]
	col := ""
	for _, v := range head {
		if pk == ToString(v["n"]) {
			col = ToString(v["o"])
			break
		}
	}
	return ToString(last[col]), nil
}
