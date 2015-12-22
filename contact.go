package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Contact struct {
	Title string
	Add   string

	Company           string
	Contact           string
	Phone             string
	Email             string
	Logo              string
	Body              []byte
	BodyRows          int
	ApplicationsTotal int

	Applications []Application
	Logos        []string
	Attention    string
}

func (p *Contact) save() error {
	filename := p.Title + ".txt"
	body := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s", p.Company, p.Contact, p.Phone, p.Email, p.Logo, p.Body)
	return ioutil.WriteFile(filename, []byte(body), 0600)
}

func loadPage(title string) (*Contact, error) {
	filename := title + ".txt"
	body2, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	// split body
	raw := strings.SplitN(string(body2), "\n", 6)
	company := raw[0]
	contact := raw[1]
	phone := raw[2]
	email := raw[3]
	logo := raw[4]
	body := []byte(raw[5])
	tmp := strings.Split(string(body), "\n")
	bodyRows := len(tmp) + 1
	now := time.Now()
	str := fmt.Sprintf("%s-%d%d%d%d%d%d", title, now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())

	c, _ := FindApplications("./", title)
	apps := c.Applications
	applicationsTotal := len(c.Applications)

	newContact := Contact{Title: title, Add: str, Company: company, Contact: contact, Phone: phone, Email: email, Logo: logo, Body: body, BodyRows: bodyRows, ApplicationsTotal: applicationsTotal, Applications: apps, Logos: logos}

	newContact.Attention = ContactRequiringAttention(&newContact)

	return &newContact, nil
	//return &Contact{Title: title, Add: str, Company: company, Contact: contact, Phone: phone, Email: email, Logo: logo, Body: body, BodyRows: bodyRows, ApplicationsTotal: applicationsTotal, Applications: apps, Logos: logos}, nil
}

func ContactRequiringAttention(c *Contact) string {
	layout := "02/01/06"
	result := ""

	for i := 0; i < len(c.Applications); i++ {
		app := c.Applications[i]
		t, _ := time.Parse(layout, app.Followup)

		duration := time.Since(t)
		//fmt.Println("hours = ", duration.Hours())

		if duration.Hours() > -24 && duration.Hours() < 24 {
			c.Applications[i].Attention = "highlight"
			result = "highlight"
		}
	}

	return result
}

func Walkdir(path string) (*Context, error) {

	files := []Contact{}

	var scan = func(filename string, info os.FileInfo, err error) error {
		if info.IsDir() == true {
			return nil
		}

		// ignore base file
		if strings.Contains(filename, "-") == true {
			return nil
		}

		if strings.HasSuffix(filename, ".txt") {
			name := strings.TrimSuffix(filename, filepath.Ext(filename))
			page, _ := loadPage(name)
			files = append(files, *page)
		}
		return nil
	}

	err := filepath.Walk(path, scan)

	if err != nil {
		return nil, err
	}

	now := time.Now()
	str := fmt.Sprintf("%d%d%d%d%d%d", now.Year(), now.Month(), now.Day(), now.Hour(), now.Minute(), now.Second())
	return &Context{Contacts: files, Title: "Application Tracker", Add: str}, nil //, Add: str}, nil
}
