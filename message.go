package mail

import (
	"bytes"
	"text/template"
)

const contentTypeTextHTML = "text/html"

type Message struct {
	header  string
	mime    string
	subject string
	html    bytes.Buffer
}

type Option interface {
	get() string
}

type ContentTypeHTML struct{}

func (o ContentTypeHTML) get() string {
	return contentTypeTextHTML
}

func NewMessage() Message {
	return Message{}
}

func (m Message) SetHeader(from, to string) Message {
	m.header = "From: " + from + "\r\nTo: " + to + "\r\n"
	return m
}

func (m Message) SetMime(contentType Option) Message {
	m.mime = "MIME-version: 1.0;\r\nContent-Type: " + contentType.get() + "; charset=\"UTF-8\";\r\n"
	return m
}

func (m Message) SetSubject(subject string) Message {
	m.subject = subject
	return m
}

func (m Message) SetHTMLTemplate(filePath, fileName string, data interface{}) error {
	tmpl, err := template.New(fileName).ParseFiles(filePath)
	if err != nil {
		return err
	}

	var buff bytes.Buffer
	if err := tmpl.Execute(&buff, data); err != nil {
		return err
	}

	m.html = buff
	return nil
}

func (m *Message) Build() []byte {
	return []byte(m.header + m.subject + m.mime + m.html.String() + "\r\n")
}
