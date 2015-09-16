package main

import (
	"github.com/boltdb/bolt"
	"strings"
)

type resource struct {
	Title string
	URL   string
	Tags  []string
}

var databaseFile string
var db *bolt.DB

/*
loadDatabase Opens the database file and makes sure that the
initial 'resources' bucket exists
*/
func loadDatabase() error {
	var err error
	db, err = bolt.Open("ii.db", 0600, nil)
	if err != nil {
		return err
	}
	// Make sure that the 'resources' bucket exists
	db.Update(func(tx *bolt.Tx) error {
		_, err := tx.CreateBucketIfNotExists([]byte("resources"))
		if err != nil {
			return err
		}
		return nil
	})
	return nil
}

/*
All resources are saved in the boltdb like so:
resources		(bucket)
|- Title 1	(bucket)
|	 |-url		(pair)
|	 \-tags		(pair) (csv)
|
\- Title 2	(bucket)
   |-url		(pair)
   \-tags		(pair) (csv)
*/

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
		b.ForEach(func(k, v []byte) error {
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
		return nil
	})
	return ret, err
}
