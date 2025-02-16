package main

import (
	"errors"
	"html/template"
	"log"
	"net/http"
	"os"
	"regexp"
)

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := "data/" + p.Title + ".txt"
	return os.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := "data/" + title + ".txt"
	body, err := os.ReadFile(filename)
	//os.ReadFile returns []byte and error
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

var templates = template.Must(template.ParseFiles("tmpl/edit.html", "tmpl/view.html"))

// remove this duplication by moving the templating code to its own function
func renderTemplate(w http.ResponseWriter, tmpl string, p *Page) {
	err := templates.ExecuteTemplate(w, tmpl+".html", p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

func viewHandler(w http.ResponseWriter, r *http.Request, title string) {
	//version 0.0.4 validation
	//title := r.URL.Path[len("/view/"):]
	/*version 0.0.5 closures
	title, err := getTitle(w, r)
	if err != nil {
		return
	}
	*/
	p, err := loadPage(title)
	if err != nil {
		http.Redirect(w, r, "/data/edit/"+title, http.StatusFound)
		return
	}
	renderTemplate(w, "view", p)
}

func editHandler(w http.ResponseWriter, r *http.Request, title string) {
	//title := r.URL.Path[len("/edit/"):]
	/*
		title, err := getTitle(w, r)
		if err != nil {
			return
		}
	*/
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	//version 0.0.2
	/*	fmt.Fprint(w, "<h1>Editing %s</h1>"+
		"<form action=\"/save/%s\" method=\"POST\">"+
		"<textarea name=\"body\">%s<textarea><br>"+
		"<input type=\"submit\" value=\"Save\">"+
		"</form>",
		p.Title, p.Title, p.Body)
	*/
	//version 0.0.3
	/*
		t, _ := template.ParseFiles("edit.html")
		t.Execute(w, p)
	*/
	renderTemplate(w, "edit", p)
}

func saveHandler(w http.ResponseWriter, r *http.Request, title string) {
	//title := r.URL.Path[len("/save/"):]
	/*
		title, err := getTitle(w, r)
		if err != nil {
			return
		}
	*/
	body := r.FormValue("body")
	p := &Page{Title: title, Body: []byte(body)}
	err := p.save()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/data/view/"+title, http.StatusFound)

}

func orgiHandler(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/data/view/FrontPage", http.StatusFound)
}

var validPath = regexp.MustCompile("^/data/(edit|save|view)/([a-zA-Z0-9]+)$")

func getTitle(w http.ResponseWriter, r *http.Request) (string, error) {
	m := validPath.FindStringSubmatch(r.URL.Path)
	if m == nil {
		http.NotFound(w, r)
		return "", errors.New("invaild Page Title")
	}
	return m[2], nil // The title is the second subexpression.
}

/*
version 0.0.5
func viewHandler(w http.ResponseWriter, r *http.Request, title string)
func editHandler(w http.ResponseWriter, r *http.Request, title string)
func saveHandler(w http.ResponseWriter, r *http.Request, title string)
*/
func makeHandler(fn func(http.ResponseWriter, *http.Request, string)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Here we will extract the page title from the Request,
		// and call the provided handler 'fn'
		m := validPath.FindStringSubmatch(r.URL.Path)
		if m == nil {
			http.NotFound(w, r)
			return
		}
		fn(w, r, m[2])
	}
}

/*
version 0.0.1

	func main() {
		p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
		p1.save()
		p2, _ := loadPage("TestPage")
		fmt.Println(string(p2.Body))
	}
*/
func main() {
	if _, err := os.Stat("data"); os.IsNotExist(err) {
		err := os.Mkdir("data", 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	http.HandleFunc("/", orgiHandler)
	http.HandleFunc("/data/view/", makeHandler(viewHandler))
	http.HandleFunc("/data/edit/", makeHandler(editHandler))
	http.HandleFunc("/data/save/", makeHandler(saveHandler))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
