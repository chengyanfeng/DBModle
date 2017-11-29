package util

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"time"
)

var sessionMap map[string]*mgo.Session = map[string]*mgo.Session{}

type MongoModel struct {
	cfg    *P
	C      int
	Cname  string
	Query  *P     // find/query condition
	Start  int    // query start at
	Rows   int    // query max rows
	sort   string // sort
	Select *P     // select field
}

func (m *MongoModel) Session() (_ *mgo.Session, err error) {
	p := *m.cfg
	key := ToString(p)
	if sessionMap[key] == nil {
		session, err := mgo.DialWithInfo(&mgo.DialInfo{
			Addrs:    []string{ToString(p["host"])},
			Database: ToString(p["name"]),
			Username: ToString(p["username"]),
			Password: ToString(p["password"]),
			Timeout:  time.Duration(ToInt(p["timeout"], 10)) * time.Second,
		})
		if err != nil {
			Error(err)
			return nil, err
		} else {
			sessionMap[key] = session
			return session.Clone(), nil
		}
	} else {
		return sessionMap[key].Clone(), nil
	}
}

func (m *MongoModel) Run(collection string, f func(*mgo.Collection)) error {
	session, err := m.Session()
	if err != nil {
		return err
	}
	defer func() {
		if err := recover(); err != nil {
			Error("Mgo", err)
		}
		session.Close()
	}()
	p := *m.cfg
	c := session.DB(ToString(p["name"])).C(collection)
	f(c)
	return err
}

func (m *MongoModel) Like(v string) (result interface{}) {
	return &bson.RegEx{Pattern: v, Options: "i"}
}

func (m *MongoModel) Find(p P) *MongoModel {
	m.Query = &p
	return m
}

func (m *MongoModel) Or(ps ...P) *MongoModel {
	q := *m.Query
	tmp := q["$or"]
	or := []P{}
	if tmp != nil {
		or = tmp.([]P)
	}
	for _, p := range ps {
		or = append(or, p)
	}
	q["$or"] = or
	m.Query = &q
	return m
}

func (m *MongoModel) ToString() string {
	return JoinStr(m.Cname, m.Query, m.Start, m.Rows, m.sort, m.Select)
}

func (m *MongoModel) Cache(i ...int) *MongoModel {
	if len(i) > 0 {
		m.C = i[0]
	} else {
		m.C = 60
	}
	return m
}

func (m *MongoModel) Field(s ...string) *MongoModel {
	if m.Select == nil {
		m.Select = &P{}
	}
	for _, k := range s {
		(*m.Select)[k] = 1
	}
	return m
}

func (m *MongoModel) Skip(start int) *MongoModel {
	m.Start = start
	return m
}

func (m *MongoModel) Limit(rows int) *MongoModel {
	m.Rows = rows
	return m
}

func (m *MongoModel) Page(start int, rows int) *MongoModel {
	m.Start = start
	m.Rows = rows
	return m
}

func (m *MongoModel) Sort(s string) *MongoModel {
	m.sort = s
	return m
}

func (m *MongoModel) All() (r *[]P) {
	cacheKey := Md5(m.ToString(), "All")
	if m.C > 0 {
		tmp := S(cacheKey)
		if tmp != nil {
			Debug("all from cache", tmp)
			r = tmp.(*[]P)
		}
	}
	if r == nil {
		r = &[]P{}
		m.Run(m.Cname, func(c *mgo.Collection) {
			q := m.query(c)
			q.All(r)
		})
	}
	if len(*r) > 0 && m.C > 0 {
		S(cacheKey, r, m.C)
	}
	return
}

func (m *MongoModel) One() (r *P) {
	cacheKey := Md5(m.ToString(), "One")
	if m.C > 0 {
		tmp := S(cacheKey)
		if tmp != nil {
			Debug("one from cache", tmp)
			r = tmp.(*P)
		}
	}
	if r == nil {
		r = &P{}
		m.Run(m.Cname, func(c *mgo.Collection) {
			q := m.query(c)
			err := q.One(r)
			Debug("one from db", r)
			if err != nil {
				Error("One", err)
			}
		})
	}
	if len(*r) > 0 && m.C > 0 {
		S(cacheKey, r, m.C)
	}
	return
}

func (m *MongoModel) Count() (total int) {
	m.Run(m.Cname, func(c *mgo.Collection) {
		q := m.query(c)
		total, _ = q.Count()
	})
	return
}

func (m *MongoModel) Add(docs ...interface{}) (err error) {
	m.Run(m.Cname, func(c *mgo.Collection) {
		if len(docs) == 1 {
			c.Insert(docs[0])
		} else {
			err = c.Insert(docs)
		}
	})
	return
}

func (m *MongoModel) Upsert(selector interface{}, doc interface{}) (err error) {
	m.Run(m.Cname, func(c *mgo.Collection) {
		_, err = c.Upsert(selector, P{"$set": doc})
		if err != nil {
			Error(err)
		}
	})
	return err
}

func (m *MongoModel) Save(p *P) (err error) {
	m.Run(m.Cname, func(c *mgo.Collection) {
		id := (*p)["_id"]
		var oid bson.ObjectId
		switch id.(type) {
		case string:
			oid = bson.ObjectIdHex(id.(string))
		case bson.ObjectId:
			oid = id.(bson.ObjectId)
		}
		(*p)["_id"] = oid
		err = c.UpdateId(oid, P{"$set": p})
		if err != nil {
			Error(err)
		}
	})
	return
}

func (m *MongoModel) RemoveId(id string) {
	m.Run(m.Cname, func(c *mgo.Collection) {
		err := c.RemoveId(bson.ObjectIdHex(id))
		if err != nil {
			Error(err)
		}
	})
}

func (m *MongoModel) Remove(selector interface{}) (e error) {
	if selector == nil || IsEmpty(selector) {
		return
	}
	m.Run(m.Cname, func(c *mgo.Collection) {
		_, err := c.RemoveAll(selector)
		if err != nil {
			Error(err)
			e = err
		}
	})
	return
}

func (m *MongoModel) RemoveAll() (e error) {
	m.Run(m.Cname, func(c *mgo.Collection) {
		_, err := c.RemoveAll(nil)
		if err != nil {
			Error(err)
			e = err
		}
	})
	return
}

func (m *MongoModel) Explain() (result interface{}) {
	p := P{}
	m.Run(m.Cname, func(c *mgo.Collection) {
		q := m.query(c)
		q.Explain(p)
	})
	return p
}

func (m *MongoModel) query(c *mgo.Collection) *mgo.Query {
	q := c.Find(m.Query).Skip(m.Start)
	if m.Rows > 0 {
		q = q.Limit(m.Rows)
	}
	if len(m.sort) > 0 {
		q = q.Sort(m.sort)
	}
	if m.Select != nil {
		q = q.Select(m.Select)
	}
	return q
}

func D(name string, params ...P) (m *MongoModel) {
	dbhost := "127.0.0.1"
	db := "db"
	m = &MongoModel{Cname: name}
	if len(params) < 1 {
		p := P{"host": dbhost, "timeout": 10}
		p["name"] = db
		params = []P{p}
	}
	p := params[0]
	m.cfg = &p
	return
}

func (this *MongoModel) Import(tblname string, f func([]*P), page ...int) (e error) {
	pos := this.loadPos(tblname)
	data := []*P{}
	if IsEmpty(pos) {
		tmp := *D(tblname).Find(P{}).Field("_id").Sort("_id").One()
		pos = ToString(tmp["_id"])
	}
	if !IsEmpty(pos) {
		if len(page) == 0 {
			page = []int{1000}
		}
		p := P{"_id": P{"$gte": ToOid(pos)}}
		Debug("Import", JsonEncode(p))
		tmp := *D(tblname).Find(p).Limit(page[0]).Sort("_id").All()
		for _, v := range tmp {
			t := v
			data = append(data, &t)
		}
		f(data)
		if len(tmp) > 0 {
			lastRow := tmp[len(data)-1]
			pos = ToString(lastRow["_id"])
			this.savePos(tblname, pos)
		}
	}
	return
}

func ToOid(id string) (oid bson.ObjectId) {
	if bson.IsObjectIdHex(id) {
		return bson.ObjectIdHex(id)
	}
	return
}

func ToOids(ids interface{}) (oids []bson.ObjectId) {
	oids = []bson.ObjectId{}
	switch ids.(type) {
	case []string:
		for _, id := range ids.([]string) {
			if bson.IsObjectIdHex(id) {
				oids = append(oids, ToOid(id))
			}
		}
	case []interface{}:
		for _, id := range ids.([]interface{}) {
			if IsOid(ToString(id)) {
				oids = append(oids, ToOid(ToString(id)))
			}
		}
	}
	return
}

func NewId() bson.ObjectId {
	return bson.NewObjectId()
}

func IsOid(id string) bool {
	return bson.IsObjectIdHex(id)
}

func (this *MongoModel) loadPos(name string) (r string) {
	p := *D(DbPos).Find(P{"key": this.getStoreKey(name)}).One()
	if len(p) > 0 {
		r = ToString(p["pos"])
	}
	return
}

func (this *MongoModel) savePos(name string, pos string) {
	D(DbPos).Upsert(P{"key": this.getStoreKey(name)}, P{"key": this.getStoreKey(name), "pos": pos})
}

func (this *MongoModel) ClearPos(name string) {
	D(DbPos).Remove(P{"key": this.getStoreKey(name)})
}

func (this *MongoModel) getStoreKey(name string) string {
	return Md5(this.cfg, name)
}

func MgoLike(v string) (result interface{}) {
	return &bson.RegEx{Pattern: v, Options: "i"}
}
