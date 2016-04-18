package main

import (
	"bytes"
	"strings"

	"github.com/boltdb/bolt"
)

type resource struct {
	Title string
	URL   string
	Tags  []string
}

var db *bolt.DB

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
// Likely there will be changes here when we actually get resources
// resources		(bucket)
// |- Title 1	(bucket)
// | |-url		(pair)
// | \-tags		(pair) (csv)
// |
// \- Title 2	(bucket)
//   |-url		(pair)
//   \-tags		(pair) (csv)

func saveResource(res resource) error {
	if err := loadDatabase(); err != nil {
		return err
	}
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
	closeDatabase()
	return err
}

func getResources() ([]resource, error) {
	ret := make([]resource, 0, 0)
	if err := loadDatabase(); err != nil {
		return ret, err
	}
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
	closeDatabase()
	return ret, err
}

func getResource(title string) (resource, error) {
	var ret resource
	if err := loadDatabase(); err != nil {
		return ret, err
	}
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("resources"))
		rB := b.Bucket([]byte(title))
		ret.Title = title
		if rVal := rB.Get([]byte("tags")); rVal != nil {
			ret.Tags = strings.Split(string(rVal), ",")
		}
		if rVal := rB.Get([]byte("url")); rVal != nil {
			ret.URL = string(rVal)
		}
		return nil
	})
	closeDatabase()
	return ret, err
}

func deleteResource(title string) error {
	if err := loadDatabase(); err != nil {
		return err
	}
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("resources"))
		return b.DeleteBucket([]byte(title))
	})
	closeDatabase()
	return err
}

func backupDatabase(b *bytes.Buffer) error {
	if err := loadDatabase(); err != nil {
		return err
	}
	err := db.View(func(tx *bolt.Tx) error {
		_, err := tx.WriteTo(b)
		return err
	})
	closeDatabase()
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
