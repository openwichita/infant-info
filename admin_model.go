package main

import (
	"fmt"

	"github.com/br0xen/bolt"
	"golang.org/x/crypto/bcrypt"
)

var dbAdmin *bolt.DB

// Admin Model Functions
// All admin accounts are stored in the admin boltdb like so
// users		(bucket)
// |- <email address 1> (bucket)
// | \-password		(pair)
// |
// |- <email address 2> (bucket)
//   \-password		(pair)
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

// getAdminUsers
// Returns a slice of all of the admin email addresses
func getAdminUsers() ([]string, error) {
	u := make([]string, 0, 0)
	if err := loadAdminDatabase(); err != nil {
		return u, err
	}
	err := dbAdmin.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		err := b.ForEach(func(k, v []byte) error {
			if v == nil { // Nested Bucket
				u = append(u, string(k))
			}
			return nil
		})
		return err
	})
	closeAdminDatabase()
	return u, err
}

func adminIsUser(email string) error {
	if err := loadAdminDatabase(); err != nil {
		return err
	}
	err := dbAdmin.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		if userBucket := b.Bucket([]byte(email)); userBucket != nil {
			return nil
		}
		return fmt.Errorf("Invalid User")
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

func adminDeleteUser(email string) error {
	err := dbAdmin.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("users"))
		return b.DeleteBucket([]byte(email))
	})
	return err
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
			if v == nil {
				if userBucket := b.Bucket(k); userBucket != nil {
					// It's a user bucket
					// Does the bucket have a 'password' key
					if pw := userBucket.Get([]byte("password")); pw != nil {
						// We have a 'password' key, call it good
						foundOne = true
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
