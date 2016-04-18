infant-info
========
A mobile application and website for compiling Wichita resources in an attempt to decrease the infant mortality rate.

Kansas' infant mortality rate is higher than the average for the rest of the country.
In part, this is due to a lack of knowledge about what resources are available.

The goal of this mobile application and website is to be a central repository for that info that nurses, midwives, and anyone else can use to find the resources that they need to help.

# To Run

* Install Git
  * https://help.github.com/articles/set-up-git/
* Install Go
  * https://golang.org/doc/install
* Clone the repository locally with
  * `go get github.com/openwichita/infant-info`

This should pull in the project with all of its dependencies.

Switch to the `$GOPATH/github.com/openwichita/infant-info` directory and run
`go build` to build the executable. You must be in the directory when you do in
order to find the templates and assets.

Run the project with `./infant-info` for silent mode or `./infant-info --dev`
for verbose console messages.

Navigate to `localhost:8080` in your web browser.

# To Contribute

* Install the project as defined above using `go get`.
* Fork the project on GitHub to your own account
* In the project folder add your fork as another git remote: `git remote add fork url-to-your-forked-repo`
* Make a branch off the current `master` branch for your feature. This is where
  you should perform your work.
* Push your branch to your remote `fork` repo then issue a pull request.

TODO
========
This is a (probably incomplete) list of things that still need to be done on this project

* Development
  * Admin
    * Resources probably need fields for more information, maybe something like:
      * Resource Title
      * Organization
      * Address
      * URL
      * Email
      * Phone
      * Hours
      * Fees
      * Languages
      * Description
      * Tags
    * Implement 'flash' messages to show the results of CRUD actions
    * Tag management system 
    * Make confirmation boxes not just use 'alert'

  * User facing
    * Basically everything here needs to be built
      * Currently the 'search' does nothing
      * Browse needs have handy tag filtering and generally be much more user friendly
      * Resources need to be links
    
  * Overall
    * Probably need better design... Everything is very bare-bones right now
    * i18n

* Non-Development
  * Figure out who is going to be keeping the site updated
    * This is probably a conversation with Rachel, who is basically N/A right now
  * Input resources
    * A sampling of the resources can be found at: http://mssconline.org/%5Cpdfs%5Cmihc-asset-map-resource-list.pdf
  * Marketing (end-game)
