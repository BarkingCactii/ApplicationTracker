package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Application struct {
	Parent      string
	Title       string
	Position    string
	Company     string
	DateApplied string
	Followup    string
	Action      string // followup action
	Notes       []byte
	NotesRows   int
	Attention   string
}

func (p *Application) saveApplication() error {
	filename := p.Title + ".txt"
	body := fmt.Sprintf("%s\n%s\n%s\n%s\n%s\n%s", p.Position, p.Company, p.DateApplied, p.Followup, p.Action, p.Notes)
	return ioutil.WriteFile(filename, []byte(body), 0600)
}

func loadApplication(title string) (*Application, error) {
	filename := title + ".txt"
	raw := strings.SplitN(string(title), "-", 2)
	parent := raw[0]
	body2, err := ioutil.ReadFile(filename)

	if err != nil {
		return nil, err
	}
	// split body
	raw = strings.SplitN(string(body2), "\n", 6)
	position := raw[0]
	company := raw[1]
	dateapplied := raw[2]
	followup := raw[3]
	action := raw[4]
	notes := []byte(raw[5])
	tmp := strings.Split(string(notes), "\n")
	notesRows := len(tmp) + 1

	newApplication := Application{Title: title, Parent: parent, Position: position, Company: company, DateApplied: dateapplied, Followup: followup, Action: action, Notes: notes, NotesRows: notesRows}

	newApplication.Attention = ""
	//	newApplication.Attention = ContactRequiringAttention(&newContact)

	//	return &Application{Title: title, Parent: parent, Position: position, Company: company, DateApplied: dateapplied, Followup: followup, Action: action, Notes: notes, NotesRows: notesRows}, nil
	return &newApplication, nil
}

func RequiringAttention(a *Application) bool {
	layout := "02/01/06"
	t, _ := time.Parse(layout, a.Followup)

	//if t > time.Now {
	duration := time.Since(t)
	if duration.Hours() < 24 && duration.Hours() > 0 {
		return true
	}

	return false
}

func FindApplications(path string, fn string) (*Contact, error) {

	files := []Application{}

	var scan = func(filename string, info os.FileInfo, err error) error {
		if info.IsDir() == true {
			return nil
		}

		// ignore base file
		if strings.Contains(filename, "-") == false {
			return nil
		}

		if strings.HasPrefix(filename, fn) && strings.HasSuffix(filename, ".txt") {
			name := strings.TrimSuffix(filename, filepath.Ext(filename))
			page, _ := loadApplication(name)
			files = append(files, *page)
		}
		return nil
	}

	err := filepath.Walk(path, scan)

	if err != nil {
		return nil, err
	}

	return &Contact{Applications: files, Title: fn, Add: fn}, nil
}
