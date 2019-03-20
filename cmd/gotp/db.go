package main

import (
	"encoding/binary"
	"github.com/major1201/gotp"
	bolt "go.etcd.io/bbolt"
)

const defaultBucket = "otp"

// Store defines the database store
type Store struct {
	filename string
	db       *bolt.DB
}

// NewStore create a new store with a database file path
func NewStore(d string) (*Store, error) {
	s := &Store{filename: d}
	if err := s.open(); err != nil {
		return nil, err
	}

	// ensure bucket "otp" exists
	if err := s.db.Update(func(tx *bolt.Tx) error {
		_, err2 := tx.CreateBucketIfNotExists([]byte(defaultBucket))
		return err2
	}); err != nil {
		s.Close()
		return nil, err
	}

	return s, nil
}

// Open a database file
func (s *Store) open() error {
	db, err := bolt.Open(s.filename, 0644, nil)
	s.db = db
	return err
}

// Close the opened database file
func (s *Store) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// Add an OTP object
func (s *Store) Add(o gotp.Otp) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(defaultBucket))
		i, err := b.NextSequence()
		if err != nil {
			return err
		}
		return b.Put(gotp.Itob(i), []byte(o.URI()))
	})
}

// Delete an OTP object with id
func (s *Store) Delete(id uint64) error {
	return s.db.Update(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(defaultBucket))
		return b.Delete(gotp.Itob(id))
	})
}

// List all the OTP objects
func (s *Store) List() ([]gotp.Otp, error) {
	var otps []gotp.Otp
	err := s.db.View(func(tx *bolt.Tx) error {
		b := tx.Bucket([]byte(defaultBucket))
		return b.ForEach(func(k, v []byte) error {
			id := binary.BigEndian.Uint64(k)
			otp, err := gotp.NewOtpFromURI(string(v))
			if err != nil {
				return err
			}
			otp.SetID(id)
			otps = append(otps, otp)
			return nil
		})
	})

	return otps, err
}
