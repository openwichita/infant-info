package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

var siteTitle = "Infant Info"

// SiteData contains data needed for many templates
// Header/Footer/Menu, etc.
type SiteData struct {
	DevMode bool

	Title       string
	SubTitle    string
	AdminMode   bool
	Port        int
	SessionName string

	Stylesheets []string
	Scripts     []string

	Flash      flashMessage // Quick message at top of page
	Menu       []menuItem   // Top-aligned menu items
	BottomMenu []menuItem   // Bottom-aligned menu items

	// Any other template data
	TemplateData interface{}
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

// Set this to something else when in production
var sessionStore = sessions.NewCookieStore([]byte("webserver secret wahoo"))

var r *mux.Router

func main() {
	site.Title = siteTitle
	site.SubTitle = ""
	site.DevMode = false
	site.Port = 8080
	site.SessionName = "infant-info"

	if err := loadDatabase(); err != nil {
		log.Fatal("Error loading database", err)
	}
	defer closeDatabase()

	args := os.Args[1:]
	for i := range args {
		if args[i] == "--dev" {
			site.DevMode = true
		}
		if strings.HasPrefix(args[i], "--port=") {
			if newPort, err := strconv.Atoi(strings.Replace(args[i], "--port=", "", -1)); err == nil {
				site.Port = newPort
			}
		}
	}

	r = mux.NewRouter()
	r.StrictSlash(true)

	assetHandler := http.FileServer(http.Dir("./assets/"))
	http.Handle("/assets/", http.StripPrefix("/assets/", assetHandler))
	r.HandleFunc("/search/", handleSearch)
	r.HandleFunc("/browse/", handleBrowse)
	r.HandleFunc("/browse/{tags}", handleBrowse).Name("browse")
	r.HandleFunc("/about/", handleAbout).Name("about")

	// Admin Subrouter
	s := r.PathPrefix("/admin").Subrouter()
	s.HandleFunc("/", handleAdmin)
	s.HandleFunc("/{category}", handleAdmin)
	s.HandleFunc("/{category}/{action}", handleAdmin)
	s.HandleFunc("/{category}/{action}/{item}", handleAdmin)

	r.HandleFunc("/download", handleBackupData)

	r.HandleFunc("/", handleSearch)

	http.Handle("/", r)
	printOutput(fmt.Sprintf("Listening on port %d\n", site.Port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", site.Port), context.ClearHandler(http.DefaultServeMux)))
}

// showFlashMessage
// Will put text into the 'aside' in the header template
// Valid 'status' values include:
// - primary		(blue)
// - secondary (light blue)
// - success		(green)
// - error			(maroon)
// - warning		(orange)
/*
func showFlashMessage(msg, status string) {
	// TODO: Store Flash Message in Session Flash
	if status == "" {
		status = "primary"
	}
	site.Flash.Message = msg
	site.Flash.Status = status
}
*/

func (s *SiteData) addStylesheet(res string) {
	s.Stylesheets = append(site.Stylesheets, res)
}

func (s *SiteData) addScript(res string) {
	s.Scripts = append(site.Scripts, res)
}

func initRequest(w http.ResponseWriter, req *http.Request) {
	printOutput(fmt.Sprintf("Request: %s\n", req.URL))
	site.SubTitle = ""
	//site.Flash = new(flashMessage)

	site.Stylesheets = make([]string, 0, 0)
	site.Stylesheets = append(site.Stylesheets, "/assets/css/pure-min.css")
	site.Stylesheets = append(site.Stylesheets, "/assets/css/ii.css")
	site.Stylesheets = append(site.Stylesheets, "https://maxcdn.bootstrapcdn.com/font-awesome/4.4.0/css/font-awesome.min.css")
	site.Scripts = make([]string, 0, 0)
	site.Scripts = append(site.Scripts, "/assets/js/ii.js")

	site.Menu = make([]menuItem, 0, 0)
	site.BottomMenu = make([]menuItem, 0, 0)
	site.AdminMode = false
	site.Menu = append(site.Menu, menuItem{Text: "Search", Link: "/search/"})
	site.Menu = append(site.Menu, menuItem{Text: "Browse", Link: "/browse/"})
	site.Menu = append(site.Menu, menuItem{Text: "About", Link: "/about/"})

	site.BottomMenu = append(site.BottomMenu, menuItem{Text: "Admin", Link: "/admin/"})
}

// handleSearch
// The main handler for all 'search' functionality
func handleSearch(w http.ResponseWriter, req *http.Request) {
	initRequest(w, req)

	site.SubTitle = "Search Resources"
	setMenuItemActive("Search")
	// Was a search action requested?
	v := req.URL.Query()
	if qry := v.Get("q"); qry != "" {
		printOutput(fmt.Sprintf("  Query: %s\n", qry))
	}
	showPage("search.html", site, w)
}

// handleBrowse
// The main handler for all 'browse' functionality
func handleBrowse(w http.ResponseWriter, req *http.Request) {
	initRequest(w, req)
	type browseData struct {
		Tags      string
		Resources []resource
	}
	vars := mux.Vars(req)
	tags := vars["tags"]

	resources, err := getResources()
	if err != nil {
		// TODO: Show Flash Message
		//showFlashMessage("Error Loading Resources!", "error")
	}

	site.SubTitle = "Browse Resources"
	setMenuItemActive("Browse")

	site.TemplateData = browseData{
		Tags:      tags,
		Resources: resources,
	}
	showPage("browse.html", site, w)
}

// handleAbout
// Show the about screen
func handleAbout(w http.ResponseWriter, req *http.Request) {
	initRequest(w, req)

	site.SubTitle = "About"
	setMenuItemActive("About")

	showPage("about.html", site, w)
}

// handleBackupData
// Pushes a download of the resource database
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

// showPage
// Load a template and all of the surrounding templates
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

// outputTemplate
// Spit out a template
func outputTemplate(tmplName string, tmplData interface{}, w http.ResponseWriter) error {
	_, err := os.Stat("templates/" + tmplName)
	if err == nil {
		t := template.New(tmplName)
		t, _ = t.ParseFiles("templates/" + tmplName)
		return t.Execute(w, tmplData)
	}
	return fmt.Errorf("WebServer: Cannot load template (templates/%s): File not found", tmplName)
}

// setMenuItemActive
// Sets a menu item to active, all others to inactive
func setMenuItemActive(which string) {
	for i := range site.Menu {
		if site.Menu[i].Text == which {
			site.Menu[i].Active = true
		} else {
			site.Menu[i].Active = false
		}
	}
}

// printOutput
// Print something to the screen, if conditions are right
func printOutput(out string) {
	if site.DevMode {
		fmt.Printf(out)
	}
}
