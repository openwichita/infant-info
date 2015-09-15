package main

import (
	"github.com/boltdb/bolt"
	"strings"
)

type Resource struct {
	Title string
	Url   string
	Tags  []string
}

var database_file string
var db *bolt.DB

func LoadDatabase() error {
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

/**
 * All resources are saved in the boltdb like so:
 *	resources		(bucket)
 *	|- Title 1	(bucket)
 *	|	 |-url		(pair)
 *	|	 \-tags		(pair) (csv)
 *	|
 *	\- Title 2	(bucket)
 *		 |-url		(pair)
 *		 \-tags		(pair) (csv)
 */

func SaveResource(res Resource) error {
	err := db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("resources"))
		if new_b, err := b.CreateBucketIfNotExists([]byte(res.Title)); err != nil {
			return err
		} else {
			if err := new_b.Put([]byte("url"), []byte(res.Url)); err != nil {
				return err
			}
			if err := new_b.Put([]byte("tags"), []byte(strings.Join(res.Tags, ","))); err != nil {
				return err
			}
		}
		return nil
	})
	return err
}

func GetResources() ([]Resource, error) {
	var ret []Resource
	err := db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte("resources"))
		b.ForEach(func(k, v []byte) error {
			if v == nil {
				// Nested Bucket
				r_b := b.Bucket(k)
				var ret_res Resource
				ret_res.Title = string(k)
				if r_val := r_b.Get([]byte("tags")); r_val != nil {
					ret_res.Tags = strings.Split(string(r_val), ",")
				}
				if r_val := r_b.Get([]byte("url")); r_val != nil {
					ret_res.Url = string(r_val)
				}
				ret = append(ret, ret_res)
			}
			return nil
		})
		return nil
	})
	return ret, err
}
