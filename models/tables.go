package models



type Studentth struct {
	Id        int    `json:"id";"size:32;column:id;auto_increment"`
	Username  string `json:"username";"size:512;column:username"`
	Companyid string `json:"companyid";"size:512;column:companyid"`
	Date      string `json:"date";"size:512;column:date"`
	Text      string `json:"text";"size:1024;column:text"`
	Messageid      int `json:"messageid";"size:64;column:messageid"`
	Credit Credit


}

type Credit struct {

	Credit  string `json:"size:512;column:credit"`
	Cred string `json:"size:512;column:cred"`




}


