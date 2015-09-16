package main

import (
	"fmt"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"
	"os"
)

var siteTitle = "Infant Info"

/*
SiteData contains data needed for many templates
Header/Footer/Menu, etc.
*/
type SiteData struct {
	Title        string
	SubTitle     string
	AdminEmail   string
	DevMode      bool
	Menu         []menuItem
	AdminMode    bool
	TemplateData interface{}
	Flash        flashMessage
}

type flashMessage struct {
	Message string
	Status  string
}

type menuItem struct {
	Text   string
	Link   string
	Active bool
}

var site SiteData

/* Set this to something else when in production */
var sessionStore = sessions.NewCookieStore([]byte("webserver secret wahoo"))

var r *mux.Router

func main() {
	site.Title = siteTitle
	site.SubTitle = ""
	site.DevMode = false
	LoadDatabase()

	args := os.Args[1:]
	for i := range args {
		if args[i] == "--dev" {
			site.DevMode = true
		}
	}

	r = mux.NewRouter()
	assetHandler := http.FileServer(http.Dir("./assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", assetHandler))
	r.HandleFunc("/", HandleSearch)
	r.HandleFunc("/search", HandleSearch)
	r.HandleFunc("/browse", HandleBrowse)
	r.HandleFunc("/browse/{tags}", HandleBrowse).Name("browse")
	r.HandleFunc("/about", HandleAbout).Name("about")
	r.HandleFunc("/admin", HandleAdmin)
	r.HandleFunc("/admin/{category}/{id}", HandleAdmin)

	http.Handle("/", r)
	http.ListenAndServe(":8080", context.ClearHandler(http.DefaultServeMux))
}

/*
showFlashMessage
Will put text into the 'aside' in the header template
Valid 'status' values include:
- primary		(blue)
- secondary (light blue)
- success		(green)
- error			(maroon)
- warning		(orange)
*/
func showFlashMessage(msg, status string) {
	if status == "" {
		status = "primary"
	}
	site.Flash.Message = msg
	site.Flash.Status = status
}

/*
Maybe we want a different menu for the 'admin' stuff?
Probably.
*/
func setupMenu(which string) {
	if which == "admin" {
		site.AdminMode = true
		site.Menu = make([]menuItem, 0, 0)
		site.Menu = append(site.Menu, menuItem{Text: "Users", Link: "/admin/users"})
		site.Menu = append(site.Menu, menuItem{Text: "Resources", Link: "/admin/resources"})
	} else {
		site.AdminMode = false
		site.Menu = make([]menuItem, 0, 0)
		site.Menu = append(site.Menu, menuItem{Text: "Search", Link: "/search"})
		site.Menu = append(site.Menu, menuItem{Text: "Browse", Link: "/browse"})
		site.Menu = append(site.Menu, menuItem{Text: "About", Link: "/about"})
	}
}

/*
handleSearch
The main handler for all 'search' functionality
*/
func handleSearch(w http.ResponseWriter, req *http.Request) {
	PrintOutput(fmt.Sprintf("Request: %s\n", req.URL))

	site.SubTitle = "Search Resources"
	SetupMenu("")
	SetMenuItemActive("Search")
	// Was a search action requested?
	v := req.URL.Query()
	if qry := v.Get("q"); qry != "" {
		PrintOutput(fmt.Sprintf("  Query: %s\n", qry))
	}
	ShowPage("search.html", site, w)
}

/*
handleBrowse
The main handler for all 'browse' functionality
*/
func handleBrowse(w http.ResponseWriter, req *http.Request) {
	type browseData struct {
		Tags      string
		Resources []Resource
	}
	vars := mux.Vars(req)
	tags := vars["tags"]

	PrintOutput(fmt.Sprintf("Request: %s\n", req.URL))
	resources, err := GetResources()
	if err != nil {
		showFlashMessage("Error Loading Resources!", "error")
	}

	site.SubTitle = "Browse Resources"
	SetupMenu("")
	SetMenuItemActive("Browse")

	site.TemplateData = browseData{
		Tags:      tags,
		Resources: resources,
	}
	ShowPage("browse.html", site, w)
}

/*
handleAbout
Show the about screen
*/
func handleAbout(w http.ResponseWriter, req *http.Request) {
	PrintOutput(fmt.Sprintf("Request: %s\n", req.URL))

	site.SubTitle = "About"
	SetupMenu("")
	SetMenuItemActive("About")

	ShowPage("about.html", site, w)
}

/*
handleAdmin
Handle entry into the Admin side of things
*/
func handleAdmin(w http.ResponseWriter, req *http.Request) {
	PrintOutput(fmt.Sprintf("Request: %s\n", req.URL))

	site.SubTitle = ""
	SetupMenu("admin")
	SetMenuItemActive("Admin")

	/* Create/Update Resource Example:
	if err := SaveResource(
		Resource{Title: "New Resource", Url: "http://www.google.com", Tags: make([]string, 0, 0)},
	); err != nil {
		Handle Error
	}
	*/

	ShowPage("admin.html", site, w)
}

/*
showPage
Load a template and all of the surrounding templates
*/
func showPage(tmplName string, tmplData interface{}, w http.ResponseWriter) error {
	for _, tmpl := range []string{
		"htmlheader.html",
		"menu.html",
		"header.html",
		tmpl_name,
		"footer.html",
		"htmlfooter.html",
	} {
		if err := OutputTemplate(tmpl, tmpl_data, w); err != nil {
			return err
		}
	}
	return nil
}

/*
outputTemplate
Spit out a template
*/
func outputTemplate(tmplName string, tmplData interface{}, w http.ResponseWriter) error {
	_, err := os.Stat("templates/" + tmpl_name)
	if err == nil {
		t := template.New(tmpl_name)
		t, _ = t.ParseFiles("templates/" + tmpl_name)
		return t.Execute(w, tmpl_data)
	}
	return fmt.Errorf("WebServer: Cannot load template (templates/%s): File not found", tmpl_name)
}

/*
setMenuItemActive
Sets a menu item to active, all others to false
*/
func setMenuItemActive(which string) {
	for i := range site.Menu {
		if site.Menu[i].Text == which {
			site.Menu[i].Active = true
		} else {
			site.Menu[i].Active = false
		}
	}
}

/*
printOutput
Print something to the screen, if conditions are right
*/
func printOutput(out string) {
	if site.DevMode {
		fmt.Printf(out)
	}
}
