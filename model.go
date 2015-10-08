package main

import (
	"bytes"
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/boltdb/bolt"
)

type resource struct {
	Title string
	URL   string
	Tags  []string
}

var databaseFile string
var db *bolt.DB
var dbAdmin *bolt.DB

// loadDatabase Opens the database file and makes sure that the
// initial 'resources' bucket exists
func loadDatabase() error {
	var err error
	db, err = bolt.Open("ii.db", 0600, nil)
	if err != nil {
		return err
	}

	// Make sure that the 'resources' bucket exists
	err = db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("resources"))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func closeDatabase() error {
	return db.Close()
}

// All resources are saved in the boltdb like so:
// resources		(bucket)
// |- Title 1	(bucket)
// |	 |-url		(pair)
// |	 \-tags		(pair) (csv)
// |
// \- Title 2	(bucket)
//    |-url		(pair)
//    \-tags		(pair) (csv)

func saveResource(res resource) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("resources"))
		var newB *bolt.Bucket
		var err error
		if newB, err = b.CreateBucketIfNotExists([]byte(res.Title)); err != nil {
			return err
		}
		if err := newB.Put([]byte("url"), []byte(res.URL)); err != nil {
			return err
		}
		if err := newB.Put([]byte("tags"), []byte(strings.Join(res.Tags, ","))); err != nil {
			return err
		}
		return nil
	})
	return err
}

func getResources() ([]resource, error) {
	var ret []resource
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("resources"))
		err := b.ForEach(func(k, v []byte) error {
			if v == nil {
				// Nested Bucket
				rB := b.Bucket(k)
				var retRes resource
				retRes.Title = string(k)
				if rVal := rB.Get([]byte("tags")); rVal != nil {
					retRes.Tags = strings.Split(string(rVal), ",")
				}
				if rVal := rB.Get([]byte("url")); rVal != nil {
					retRes.URL = string(rVal)
				}
				ret = append(ret, retRes)
			}
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	return ret, err
}

func backupDatabase(b *bytes.Buffer) error {
	closeDatabase()
	err := db.View(func(tx *bolt.Tx) error {
		_, err := tx.WriteTo(b)
		return err
	})
	loadDatabase()
	return err
	/*
		err := db.View(func(tx *bolt.Tx) error {
			w.Header().Set("Content-Type", "application/octet-stream")
			w.Header().Set("Content-Disposition", `attachment; filename="infant-info.db"`)
			w.Header().Set("Content-Length", strconv.Itoa(int(tx.Size())))
			_, err := tx.WriteTo(w)
			return err
		})
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	*/
}

// Admin Functions
// All admin accounts are stored in the admin boltdb like so
// users		(bucket)
// |- <email address 1> (bucket)
// |	 \-password		(pair)
// |
// |- <email address 2> (bucket)
//  	 \-password		(pair)
func loadAdminDatabase() error {
	var err error
	dbAdmin, err = bolt.Open("iiAdmin.db", 0600, nil)
	if err != nil {
		return err
	}

	// Make sure that the 'users' bucket exists
	err = dbAdmin.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("users"))
		if err != nil {
			return err
		}
		return nil
	})

	if err != nil {
		return err
	}
	return nil
}

func closeAdminDatabase() error {
	return dbAdmin.Close()
}

// adminCheckFirstRun
// Check if there is an admin account.
func adminCheckFirstRun() error {
	if err := loadAdminDatabase(); err != nil {
		return err
	}
	err := dbAdmin.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		// Make sure that we have a bucket in users
		foundOne := false
		err := b.ForEach(func(k, v []byte) error {
			foundOne = true
			if v == nil {
				if userBucket := b.Bucket(k); userBucket != nil {
					// It's a user bucket
					// Does the bucket have a 'password' key
					if pw := userBucket.Get([]byte("password")); pw != nil {
						// We have a 'password' key, call it good
						return nil
					}
				}
			}
			return fmt.Errorf("Couldn't find an Admin User")
		})
		if !foundOne {
			return fmt.Errorf("Couldn't find an Admin User")
		}
		return err
	})
	closeAdminDatabase()
	return err
}

func adminCheckCredentials(email, password string) error {
	if err := loadAdminDatabase(); err != nil {
		return err
	}
	err := dbAdmin.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		if userBucket := b.Bucket([]byte(email)); userBucket != nil {
			if pw := userBucket.Get([]byte("password")); pw != nil {
				return bcrypt.CompareHashAndPassword(pw, []byte(password))
			}
		}
		return fmt.Errorf("Invalid User")
	})
	closeAdminDatabase()
	return err
}

func adminSaveUser(email, password string) error {
	cryptPW, cryptError := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if cryptError != nil {
		return cryptError
	}
	if err := loadAdminDatabase(); err != nil {
		return err
	}
	err := dbAdmin.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		var newB *bolt.Bucket
		var err error
		if newB, err = b.CreateBucketIfNotExists([]byte(email)); err != nil {
			return err
		}
		if err := newB.Put([]byte("password"), cryptPW); err != nil {
			return err
		}
		return nil
	})
	closeAdminDatabase()
	return err
}
