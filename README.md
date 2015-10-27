infant-info
========
A mobile application and website for compiling Wichita resources in an attempt to decrease the infant mortality rate.

Kansas' infant mortality rate is higher than the average for the rest of the country.
In part, this is due to a lack of knowledge about what resources are available.

The goal of this mobile application and website is to be a central repository for that info that nurses, midwives, and anyone else can use to find the resources that they need to help.

To run:
* Install Git
  * https://help.github.com/articles/set-up-git/
* Install Go
  * https://golang.org/doc/install

* Clone the repository locally with
  * 'go get github.com/openwichita/infant-info'
* or clone your fork locally with
  * 'go get github.com/[your github username]/infant-info'

This should pull in the project with all of its dependencies.

'go build' to build the executable.

'./infant-info' for silent mode or './infant-info --dev' for verbose console messages.

Navigate to 'localhost:8080' in your web browser.
