package main

import (
	"bufio"
	"fmt"
	"html/template"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

type Context struct {
	Title    string
	Add      string
	Contacts []Contact
}

func listHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := Walkdir("./")
	if err != nil {
		http.Redirect(w, r, "/list/", http.StatusFound)
		return
	}
	renderTemplateContext(w, "list", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	p, err := loadPage(title)
	if err != nil {
		p = &Contact{Title: title, BodyRows: 5}
	}
	renderTemplateContact(w, "edit", p)
}

func deleteHandler(w http.ResponseWriter, r *http.Request, title string) {
	err := os.Rename(title+".txt", title+".deleted")

	if err != nil {
		fmt.Println(err)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func applicationHandler(w http.ResponseWriter, r *http.Request, title string) {
	application, err := loadApplication(title)
	if err == nil {
		// load existing application
		renderTemplateApplication(w, "application", application)
		return
	}

	a := &Application{Title: title, Position: "", Company: "", DateApplied: "", Followup: "", Action: "", Notes: []byte(""), NotesRows: 5}
	renderTemplateApplication(w, "application", a)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	company := r.FormValue("company")
	contact := r.FormValue("contact")
	phone := r.FormValue("phone")
	email := r.FormValue("email")
	logo := r.FormValue("logo")
	body := r.FormValue("body")
	p := &Contact{Title: title, Company: company, Contact: contact, Phone: phone, Email: email, Logo: logo, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/", http.StatusFound)
}

func appsaveHandler(w http.ResponseWriter, r *http.Request, title string) {
	position := r.FormValue("position")
	company := r.FormValue("company")
	dateapplied := r.FormValue("dateapplied")
	followup := r.FormValue("followup")
	action := r.FormValue("action")
	notes := r.FormValue("notes")
	p := &Application{Title: title, Position: position, Company: company, DateApplied: dateapplied, Followup: followup, Action: action, Notes: []byte(notes)}
	err := p.saveApplication()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	editUrl := strings.Split(title, "-")
	http.Redirect(w, r, "/edit/"+editUrl[0], http.StatusFound)
}

var templates = populateTemplates()
var logos []string

func readLogos(path string) ([]string, error) {

	files := []string{}

	var scan = func(filename string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() == true {
			return nil
		}

		files = append(files, filepath.Base(filename))
		fmt.Println("Adding " + filepath.Base(filename))
		return nil
	}

	err := filepath.Walk(path, scan)

	if err != nil {
		return nil, err
	}

	return files, nil
}

func renderTemplateContact(w http.ResponseWriter, tmpl string, p *Contact) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, "renderContact "+err.Error(), http.StatusInternalServerError)
	}
}

func renderTemplateApplication(w http.ResponseWriter, tmpl string, p *Application) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, "renderApplication "+err.Error(), http.StatusInternalServerError)
	}
}

func renderTemplateContext(w http.ResponseWriter, tmpl string, p *Context) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, "renderContext "+err.Error(), http.StatusInternalServerError)
	}
}

var validPath = regexp.MustCompile("^/(edit|save|view|application)/([a-zA-Z0-9]+)$")

func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		rawFileName := strings.Split(r.URL.Path, "/")
		lastIdx := len(rawFileName) - 1
		for i := 0; i < lastIdx; i++ {
			fmt.Printf("%s", rawFileName[i])
		}
		fn(w, r, rawFileName[lastIdx])
	}
}

type MyHandler struct {
	http.Handler
}

func serveResource(w http.ResponseWriter, req *http.Request) {
	path := "public" + req.URL.Path
	//path := req.URL.Path
	var contentType string
	if strings.HasSuffix(path, ".css") {
		contentType = "text/css"
	} else if strings.HasSuffix(path, ".png") {
		contentType = "image/png"
	} else if strings.HasSuffix(path, ".jpg") {
		contentType = "image/jpg"
	} else {
		contentType = "text/plain"
	}

	f, err := os.Open(path)

	if err == nil {
		defer f.Close()
		w.Header().Add("Content Type", contentType)

		br := bufio.NewReader(f)
		br.WriteTo(w)
	} else {
		w.WriteHeader(404)
	}
}

func populateTemplates() *template.Template {
	result := template.New("templates")

	basePath := "templates"
	templateFolder, _ := os.Open(basePath)
	defer templateFolder.Close()

	templatePathsRaw, _ := templateFolder.Readdir(-1)

	templatePaths := new([]string)
	for _, pathInfo := range templatePathsRaw {
		if !pathInfo.IsDir() {
			*templatePaths = append(*templatePaths,
				basePath+"/"+pathInfo.Name())
		}
	}

	result.ParseFiles(*templatePaths...)

	return result
}

func main() {
	http.HandleFunc("/edit/", makeHandler(editHandler))
	http.HandleFunc("/delete/", makeHandler(deleteHandler))
	http.HandleFunc("/application/", makeHandler(applicationHandler))
	http.HandleFunc("/appsave/", makeHandler(appsaveHandler))
	http.HandleFunc("/save/", makeHandler(saveHandler))
	http.HandleFunc("/logos/", serveResource)
	http.HandleFunc("/img/", serveResource)
	http.HandleFunc("/css/", serveResource)
	http.HandleFunc("/", makeHandler(listHandler))

	logos, _ = readLogos("public/logos")

	http.ListenAndServe(":8080", nil)
}
