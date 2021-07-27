package model

type DocuFirst struct {
	Id    int    `xorm: "pk int" form: "id" json: "id"`
	Title string `xorm: "html title string" form: "title" json: "title"`
	Body  string `xorm: "html body string" form: "body" json: "body"`
}
