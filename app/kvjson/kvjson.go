// USAGE:

// package main

// import (
// 	"fmt"

// 	"github.com/ClimenteA/kvjson"
// )

// func main() {

// 	type Person struct {
// 		Name string
// 		Age  int
// 	}

// 	type Car struct {
// 		Model string
// 		Year  int
// 	}

// 	// Next, in the main func initialize the "db" with a folder path
// 	// Use `kvjson.DB` to get the types. Ex: `func someFunc(db kvjson.DB) {etc}`
// 	db := kvjson.New("./db")

// 	// Set a `key` with an initialized `struct`
// 	// In the `./db` folder a `key.json` will be saved
// 	db.Set("person", Person{Name: "Alin", Age: 30})
// 	db.Set("car", Car{Model: "Dacia", Year: 2020})

// 	// Get value of a `key`
// 	// This will just unmarshal the data into the struct
// 	var car Car
// 	err := db.Get("car", &car)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(car)

// 	var person Person
// 	err = db.Get("person", &person)
// 	if err != nil {
// 		panic(err)
// 	}
// 	fmt.Println(person)

// 	// Delete a `key`
// 	err = db.Del("person")
// 	if err != nil {
// 		panic(err)
// 	}
// 	err = db.Del("car")
// 	if err != nil {
// 		panic(err)
// 	}
// }

package kvjson

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

type DB struct {
	path string
}

func New(dbPath string) DB {
	err := os.MkdirAll(dbPath, os.ModePerm)
	if err != nil {
		panic(err)
	}
	return DB{path: dbPath}
}

func (db DB) Path() string {
	return db.path
}

func (db DB) UUID4() string {
	newUUID4, err := exec.Command("uuidgen").Output()
	if err != nil {
		panic(err)
	}
	return string(newUUID4)
}

func (db DB) SetWithKey(key string, val interface{}) error {
	entry := filepath.Join(db.path, key+".json")
	file, err := os.Create(entry)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(val)
	if err != nil {
		return err
	}

	return nil
}

func (db DB) Set(val interface{}) string {
	valId := db.UUID4()
	err := db.SetWithKey(valId, val)
	if err != nil {
		panic(err)
	}
	return valId
}

func (db DB) SetWithKeyIfNew(key string, val interface{}) error {
	entry := filepath.Join(db.path, key+".json")
	if _, err := os.Stat(entry); os.IsNotExist(err) {
		return db.SetWithKey(key, val)
	}
	return nil
}

func (db DB) Get(key string, val interface{}) error {
	entry := filepath.Join(db.path, key+".json")
	file, err := os.Open(entry)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	err = decoder.Decode(&val)
	if err != nil {
		fmt.Println(err)
		return err
	}

	return nil
}

// func (db DB) GetAll(values interface{}) error {

// 	entries, _ := os.ReadDir(db.path)

// 	for _, e := range entries {
// 		if e.IsDir() {
// 			continue
// 		}
// 		key := strings.TrimSuffix(e.Name(), filepath.Ext(e.Name()))
// 		var val interface{}
// 		err := db.Get(key, &val)
// 		if err != nil {
// 			continue
// 		}
// 		values = append(values, val)
// 	}

// 	return nil
// }

func (db DB) Del(key string) error {
	entry := filepath.Join(db.path, key+".json")
	err := os.Remove(entry)
	if err != nil {
		return err
	}
	return nil
}
