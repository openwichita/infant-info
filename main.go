package main

import (
	"bytes"
	"fmt"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"html/template"
	"net/http"
	"os"
	"strconv"
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
	loadDatabase()
	defer closeDatabase()

	args := os.Args[1:]
	for i := range args {
		if args[i] == "--dev" {
			site.DevMode = true
		}
	}

	r = mux.NewRouter()
	assetHandler := http.FileServer(http.Dir("./assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", assetHandler))
	r.HandleFunc("/", handleSearch)
	r.HandleFunc("/search", handleSearch)
	r.HandleFunc("/browse", handleBrowse)
	r.HandleFunc("/browse/{tags}", handleBrowse).Name("browse")
	r.HandleFunc("/about", handleAbout).Name("about")
	r.HandleFunc("/admin", handleAdmin)
	r.HandleFunc("/admin/{category}/{id}", handleAdmin)

	r.HandleFunc("/download", handleBackupData)

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
	printOutput(fmt.Sprintf("Request: %s\n", req.URL))

	site.SubTitle = "Search Resources"
	setupMenu("")
	setMenuItemActive("Search")
	// Was a search action requested?
	v := req.URL.Query()
	if qry := v.Get("q"); qry != "" {
		printOutput(fmt.Sprintf("  Query: %s\n", qry))
	}
	showPage("search.html", site, w)
}

/*
handleBrowse
The main handler for all 'browse' functionality
*/
func handleBrowse(w http.ResponseWriter, req *http.Request) {
	type browseData struct {
		Tags      string
		Resources []resource
	}
	vars := mux.Vars(req)
	tags := vars["tags"]

	printOutput(fmt.Sprintf("Request: %s\n", req.URL))
	resources, err := getResources()
	if err != nil {
		showFlashMessage("Error Loading Resources!", "error")
	}

	site.SubTitle = "Browse Resources"
	setupMenu("")
	setMenuItemActive("Browse")

	site.TemplateData = browseData{
		Tags:      tags,
		Resources: resources,
	}
	showPage("browse.html", site, w)
}

/*
handleAbout
Show the about screen
*/
func handleAbout(w http.ResponseWriter, req *http.Request) {
	printOutput(fmt.Sprintf("Request: %s\n", req.URL))

	site.SubTitle = "About"
	setupMenu("")
	setMenuItemActive("About")

	showPage("about.html", site, w)
}

/*
handleAdmin
Handle entry into the Admin side of things
*/
func handleAdmin(w http.ResponseWriter, req *http.Request) {
	printOutput(fmt.Sprintf("Request: %s\n", req.URL))

	site.SubTitle = ""
	setupMenu("admin")
	setMenuItemActive("Admin")

	/* Create/Update Resource Example:
	if err := SaveResource(
		Resource{Title: "New Resource", Url: "http://www.google.com", Tags: make([]string, 0, 0)},
	); err != nil {
		Handle Error
	}
	*/

	showPage("admin.html", site, w)
}

/*
handleBackupData
Pushes a download of the resource database
*/
func handleBackupData(w http.ResponseWriter, req *http.Request) {
	var b *bytes.Buffer
	err := backupDatabase(b)
	fmt.Println("DB Backup Requested")
	fmt.Printf("DB Size: %d\n", b.Len())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Disposition", `attachment; filename="infant-info.db"`)
	w.Header().Set("Content-Length", strconv.Itoa(int(b.Len())))
	_, err = b.WriteTo(w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
		tmplName,
		"footer.html",
		"htmlfooter.html",
	} {
		if err := outputTemplate(tmpl, tmplData, w); err != nil {
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
	_, err := os.Stat("templates/" + tmplName)
	if err == nil {
		t := template.New(tmplName)
		t, _ = t.ParseFiles("templates/" + tmplName)
		return t.Execute(w, tmplData)
	}
	return fmt.Errorf("WebServer: Cannot load template (templates/%s): File not found", tmplName)
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
