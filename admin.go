package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

type editUserData struct {
	Email      string
	Password   string
	FormAction string
}

type listData struct {
	List []string
}

// handleAdmin
// Handle entry into the Admin side of things
func handleAdmin(w http.ResponseWriter, req *http.Request) {
	initAdminRequest(w, req)

	vars := mux.Vars(req)
	adminCategory := vars["category"]
	printOutput("Admin Category: " + adminCategory + "\n")

	// First, check if we're logged in
	userEmail, _ := getSessionStringValue("email", w, req)

	// With a valid account
	printOutput("Checking for Valid Account: " + userEmail + "\n")
	validUser := adminIsUser(userEmail)

	if validUser != nil {
		// Not logged in, only allow access to the login page
		if adminCategory == "dologin" {
			handleAdminDoLogin(w, req)
			return
		}
		if adminCategory == "firstcreate" {
			if firstErr := adminCheckFirstRun(); firstErr != nil {
				handleAdminSaveUser(w, req)
			} else {
				// We already have an admin account... So...
				http.Redirect(w, req, "/", 302)
			}
			return
		}
		if adminCategory == "" {
			handleAdminLogin(w, req)
			return
		}
		http.Redirect(w, req, "/admin", 302)
		return
	}

	site.SubTitle = fmt.Sprintf("Logged in as %s", userEmail)

	setMenuItemActive("Admin")

	if adminCategory == "dologout" {
		handleAdminDoLogout(w, req)
		return
	}
	if adminCategory == "users" {
		handleAdminUsers(w, req)
		return
	}
	if adminCategory == "resources" {
		handleAdminResources(w, req)
		return
	}

	http.Redirect(w, req, "/admin/resources", 302)
}

func initAdminRequest(w http.ResponseWriter, req *http.Request) {
	printOutput(fmt.Sprintf("Admin Request: %s\n", req.URL))

	w.Header().Set("Cache-Control", "no-cache")

	// First, check if we're logged in
	userEmail, _ := getSessionStringValue("email", w, req)

	// With a valid account
	validUser := adminIsUser(userEmail)

	site.SubTitle = ""
	//site.Flash = new(flashMessage)
	site.Menu = make([]menuItem, 0, 0)
	site.BottomMenu = make([]menuItem, 0, 0)

	site.Stylesheets = make([]string, 0, 0)
	site.Stylesheets = append(site.Stylesheets, "/assets/css/pure-min.css")
	site.Stylesheets = append(site.Stylesheets, "https://maxcdn.bootstrapcdn.com/font-awesome/4.4.0/css/font-awesome.min.css")
	site.Stylesheets = append(site.Stylesheets, "/assets/css/ii.css")
	// TODO: Make Dynamic Scripts work
	site.Scripts = make([]string, 0, 0)
	site.Scripts = append(site.Scripts, "/assets/js/ii.js")
	site.Scripts = append(site.Scripts, "/assets/js/admin.js")

	if validUser == nil {
		site.Menu = append(site.Menu, menuItem{Text: "Users", Link: "/admin/users"})
		site.Menu = append(site.Menu, menuItem{Text: "Resources", Link: "/admin/resources"})

		site.BottomMenu = append(site.BottomMenu, menuItem{Text: "Logout", Link: "/admin/dologout"})
	}
	site.BottomMenu = append(site.BottomMenu, menuItem{Text: "Home", Link: "/"})
}

// handleAdminLogin
// Show the Login screen
func handleAdminLogin(w http.ResponseWriter, req *http.Request) {
	setMenuItemActive("Admin")
	if err := adminCheckFirstRun(); err != nil {
		site.SubTitle = "Create Admin Account"
		site.TemplateData = editUserData{Email: "", Password: "", FormAction: "/admin/firstcreate"}
		showPage("admin-edituser.html", site, w)
	} else {
		site.SubTitle = "Admin Login"
		showPage("admin-login.html", site, w)
	}
}

// handleAdminDoLogin
// Verify the provided credentials, set up a cookie (if requested)
// And redirect back to /admin
func handleAdminDoLogin(w http.ResponseWriter, req *http.Request) {
	// Fetch the login credentials
	email := req.FormValue("email")
	password := req.FormValue("password")
	// Remember functionality is not included (yet? ever?)
	// remember := req.FormValue("remember")
	if email != "" && password != "" {
		printOutput(fmt.Sprintf("  Login Request (%s)\n", email))
		if err := adminCheckCredentials(email, password); err != nil {
			// Couldn't find the credentials
			printOutput(fmt.Sprintf("		Failed!\n"))
		} else {
			printOutput(fmt.Sprintf("		Success!\n"))
			session, err := sessionStore.Get(req, site.SessionName)
			if err != nil {
				http.Error(w, err.Error(), 500)
				return
			}
			session.Values["email"] = email
			session.Save(req, w)
		}
	}
	// TODO: Show Flash Message
	//showFlashMessage(fmt.Sprintf("Logged in as %s", email), "warning")
	http.Redirect(w, req, "/admin", 302)
}

func handleAdminDoLogout(w http.ResponseWriter, req *http.Request) {
	session, err := sessionStore.Get(req, site.SessionName)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	session.Options.MaxAge = -1
	session.Save(req, w)

	site.SubTitle = "Login"
	setMenuItemActive("Admin")

	// TODO: Show Flash Message
	//showFlashMessage("You have been logged out.", "secondary")

	showPage("admin-login.html", site, w)

}

func handleAdminUsers(w http.ResponseWriter, req *http.Request) {
	site.SubTitle = "Admin User Management"
	setMenuItemActive("Users")

	vars := mux.Vars(req)
	userFunction := vars["action"]

	if userFunction == "edit" {
		handleAdminEditUser(w, req)
		return
	} else if userFunction == "save" {
		handleAdminSaveUser(w, req)
		return
	} else if userFunction == "delete" {
		handleAdminDeleteUser(w, req)
		return
	}

	// No action given, display users
	users, err := getAdminUsers()
	userList := make([]string, 0, 0)
	for i := range users {
		printOutput("Found User: " + users[i] + "\n")
		userList = append(userList, users[i])
	}
	site.TemplateData = listData{List: userList}
	if err == nil {
		showPage("admin-users.html", site, w)
	} else {
		printOutput(fmt.Sprintf("%s\n", err))
	}
}

func handleAdminEditUser(w http.ResponseWriter, req *http.Request) {
	site.SubTitle = "Edit Admin Account"
	vars := mux.Vars(req)
	userEmail := vars["item"]
	site.TemplateData = editUserData{Email: userEmail, Password: "", FormAction: "/admin/users/save/"}
	showPage("admin-edituser.html", site, w)
}

func handleAdminSaveUser(w http.ResponseWriter, req *http.Request) {
	// Fetch the login credentials
	email := req.FormValue("email")
	password := req.FormValue("password")
	repeatpw := req.FormValue("repeat")
	printOutput(fmt.Sprintf("  Starting Save User Process %s -> %s == %s\n", email, password, repeatpw))
	if email != "" && password != "" && password == repeatpw {
		printOutput(fmt.Sprintf("  Save User Request (%s)\n", email))
		if err := adminSaveUser(email, password); err != nil {
			printOutput(fmt.Sprintf("		Failed!\n"))
			// TODO: Set Flash Message for Failure
		} else {
			printOutput(fmt.Sprintf("		Success!\n"))
			// TODO: Set Flash Message for Success
		}
	}

	http.Redirect(w, req, "/admin/users", 302)
}

func handleAdminDeleteUser(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userItem := vars["item"]
	printOutput("Deleting User: " + userItem)
	if err := adminDeleteUser(userItem); err != nil {
		printOutput(fmt.Sprintf("		Failed!\n"))
		// TODO: Set Flash Message for Failure
	} else {
		printOutput(fmt.Sprintf("		Success!\n"))
		// TODO: Set Flash Message for Success
	}

	handleAdminUsers(w, req)
	//http.Redirect(w, req, "/admin/users", 302)
}

func handleAdminResources(w http.ResponseWriter, req *http.Request) {
	site.SubTitle = "Edit Resources"
	setMenuItemActive("Resources")

	vars := mux.Vars(req)
	resFunction := vars["action"]
	if resFunction == "save" {
		handleAdminSaveResource(w, req)
		return
	}

	// No action given, display users
	showPage("admin-resources.html", site, w)
}

func handleAdminSaveResource(w http.ResponseWriter, req *http.Request) {
	/* Create/Update Resource Example:
	if err := SaveResource(
		Resource{Title: "New Resource", Url: "http://www.google.com", Tags: make([]string, 0, 0)},
	); err != nil {
		Handle Error
	}
	*/

	http.Redirect(w, req, "/admin/resources", 302)
}
