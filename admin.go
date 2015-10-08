package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
)

// handleAdmin
// Handle entry into the Admin side of things
func handleAdmin(w http.ResponseWriter, req *http.Request) {
	initRequest(w, req)
	vars := mux.Vars(req)
	adminFunction := vars["function"]

	// First, check if we're logged in
	session, err := sessionStore.Get(req, site.SessionName)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	userEmail := session.Values["email"]
	if userEmail == "" || userEmail == nil {
		// Not logged in, only allow access to the login page
		if adminFunction == "dologin" {
			handleAdminDoLogin(w, req)
			return
		}
		handleAdminLogin(w, req)
		return
	}

	site.SubTitle = fmt.Sprintf("Logged in as %s", userEmail)

	setupMenu("admin")
	setMenuItemActive("Admin")
	printOutput("    Admin/" + adminFunction + "\n")

	if adminFunction == "firstcreate" {
		printOutput("    First time Create\n")
		handleAdminFirstCreate(w, req)
		return
	}
	if adminFunction == "logout" {
		printOutput("    Do Logout\n")
		handleAdminDoLogout(w, req)
		return
	}
	if adminFunction == "users" {
		printOutput("    Users\n")
		handleAdminUsers(w, req)
		return
	}
	if adminFunction == "resources" {
		printOutput("    Resources\n")
		handleAdminResources(w, req)
		return
	}

	/* Create/Update Resource Example:
	if err := SaveResource(
		Resource{Title: "New Resource", Url: "http://www.google.com", Tags: make([]string, 0, 0)},
	); err != nil {
		Handle Error
	}
	*/

	showPage("admin.html", site, w)
}

// handleAdminLogin
// Show the Login screen
func handleAdminLogin(w http.ResponseWriter, req *http.Request) {
	setupMenu("")
	setMenuItemActive("Admin")
	if err := adminCheckFirstRun(); err != nil {
		site.SubTitle = "Create Admin Account"
		showPage("admin-create.html", site, w)
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

	http.Redirect(w, req, "/admin", 301)
}

func handleAdminDoLogout(w http.ResponseWriter, req *http.Request) {
	printOutput("Do Logout\n")

	session, err := sessionStore.Get(req, site.SessionName)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}

	session.Options.MaxAge = -1
	session.Save(req, w)

	site.SubTitle = "Login"
	setupMenu("")
	setMenuItemActive("Admin")

	// TODO: Show Flash Message
	//showFlashMessage("You have been logged out.", "secondary")

	showPage("admin-login.html", site, w)

}

func handleAdminUsers(w http.ResponseWriter, req *http.Request) {
	site.SubTitle = "User Management"
	setupMenu("admin")
	setMenuItemActive("Users")

	vars := mux.Vars(req)
	userFunction := vars["subfunc"]
	if userFunction == "save" {
		handleAdminSaveUser(w, req)
		return
	}
	if userFunction == "create" {
	}

	// No subfunc given, display users
	if err := showPage("admin-users.html", site, w); err != nil {
		printOutput(fmt.Sprintf("%s", err))
	}
}

func handleAdminSaveUser(w http.ResponseWriter, req *http.Request) {
	// Fetch the login credentials
	email := req.FormValue("email")
	password := req.FormValue("password")
	if email != "" && password != "" {
		printOutput(fmt.Sprintf("  Save User Request (%s:%s)\n", email, password))
		if err := adminSaveUser(email, password); err != nil {
			printOutput(fmt.Sprintf("		Failed!\n"))
			// TODO: Set Flash Message for Failure
		} else {
			printOutput(fmt.Sprintf("		Success!\n"))
			// TODO: Set Flash Message for Success
		}
	}

	http.Redirect(w, req, "/admin/users", 301)
}

func handleAdminResources(w http.ResponseWriter, req *http.Request) {
	site.SubTitle = "Edit Resources"
	setupMenu("admin")
	setMenuItemActive("Resources")

	vars := mux.Vars(req)
	resFunction := vars["subfunc"]
	if resFunction == "save" {
		handleAdminSaveResource(w, req)
		return
	}

	// No subfunc given, display users
	showPage("admin-resources.html", site, w)
}

func handleAdminSaveResource(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, "/admin/resources", 301)
}

func handleAdminFirstCreate(w http.ResponseWriter, req *http.Request) {
	if err := adminCheckFirstRun(); err != nil {

	} else {
		// We already have an admin account... So...
		http.Redirect(w, req, "/", 301)
	}
}
