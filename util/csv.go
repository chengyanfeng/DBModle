package util

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type Csv struct {
	Body     string
	Split    string
	File     string
	Tmp      string
	Filter   []string
	Fh       bool
	Limit    int
	Offset   int
	Head     []P
	Data     []P
	Err      error
	LockHead bool
}

func (this *Csv) Scan(head []P) (count int) {
	if this.Limit < 1 {
		this.Limit = 1
	}
	this.Data = nil
	this.Err = nil
	this.Head = []P{}
	if !IsEmpty(head) {
		this.LockHead = true
		for _, v := range head {
			if !IsEmpty(v["o"]) {
				this.Head = append(this.Head, v)
			}
		}
		if !IsEmpty(this.Tmp) {
			this.saveTmpHead(head)
		}
	}
	if !IsEmpty(this.File) {
		f, err := os.Open(this.File)
		defer f.Close()
		if err != nil {
			this.Err = err
			return
		}
		buf := bufio.NewReader(f)
		cols := []string{}
		half := false
		line := ""
		for {
			seg, err := buf.ReadString('\n')
			if err != nil {
				if err == io.EOF {
					return
				}
				this.Err = err
				return
			}
			if !IsCsvEnd(seg, half) {
				half = true
				line += Trim(seg)
				continue
			} else {
				half = false
				line += Trim(seg)
				if len(cols) == 0 {
					cols = ToFields(Trim(line), this.Split)
					if IsEmpty(head) {
						for i, col := range cols {
							if this.Filter != nil {
								col = Replace(col, this.Filter, "")
							}
							p := P{"o": JoinStr("c", i), "n": col}
							if !this.Fh {
								p["n"] = p["o"]
							}
							this.Head = append(this.Head, p)
						}
					}
				}
				if count == 0 {
					if !this.Fh {
						this.scanData(count, line)
					}
				} else {
					this.scanData(count, line)
				}
				line = ""
				count++
			}
		}
	} else {
		md5 := Md5(this.Body)
		tmpdir := "/data/tmp/"
		tmp := tmpdir + md5
		if !FileExists(tmp) {
			Mkdir(tmpdir)
			WriteFile(tmp, []byte(this.Body))
		}
		this.File = tmp
		this.Scan(head)
	}
	return
}

func (this *Csv) scanData(count int, line string) {
	data := ToFields(Trim(line), this.Split)
	if len(data) == len(this.Head) {
		p := P{}
		for i, v := range this.Head {
			data[i] = Trim(data[i])
			k := ToString(v["o"])
			v[k] = data[i]
			this.setType(data, i, v)
			p[k] = data[i]
		}
		if count >= this.Offset && len(this.Data) < this.Limit && len(this.Data) < ROW_LIMIT_MAX {
			this.Data = append(this.Data, p)
		}
	} else {
		Error("scanData", count, len(data), len(this.Head), JsonEncode(data))
	}
	if !IsEmpty(this.Tmp) {
		AppendFile(this.Tmp, this.ToLine(data)+"\n")
	}
}

func (this *Csv) ToLine(row []string) string {
	line := ""
	for _, v := range row {
		line += v + ","
	}
	if len(line) > 1 {
		line = line[0 : len(line)-1]
	}
	return line
}

func (this *Csv) saveTmpHead(head []P) {
	row := []string{}
	for _, h := range head {
		v := ToString(h["n"], ToString(h["o"]))
		row = append(row, v)
	}
	if this.Fh {
		WriteFile(this.Tmp, []byte(this.ToLine(row)+"\n"))
	} else {
		WriteFile(this.Tmp, []byte(""))
	}
}

func (this *Csv) setType(row []string, i int, p P) {
	if p["type"] == "string" {
		return
	}
	v := row[i]
	if !this.LockHead {
		if IsInt(v) {
			if p["type"] == nil {
				p["type"] = "number"
			}
		} else if IsFloat(v) {
			if p["type"] == nil {
				p["type"] = "number"
			}
		} else if IsDate(v) {
			if p["type"] == nil {
				p["type"] = "date"
			}
		} else {
			if IsEmpty(v) && p["type"] != nil {
				return
			}
			p["type"] = "string"
		}
	}
	switch ToString(p["type"]) {
	case "number":
		row[i] = ToString(ToFloat(v))
	case "date":
		row[i], _ = ToDate(v)
	}
}

func (this *Csv) Cut(cols []string) (dst string) {
	str := ""
	for _, v := range cols {
		str = JoinStr(str, ToInt(Replace(v, []string{"C", "c"}, ""))+1, ",")
	}
	if len(str) > 0 {
		str = str[0 : len(str)-1]
	}
	dst = Replace(this.File, []string{".csv"}, ".cut.csv")
	cmd := fmt.Sprintf("cut -d, -f%v %v > %v", str, this.File, dst)
	Exec(cmd)
	return
}
