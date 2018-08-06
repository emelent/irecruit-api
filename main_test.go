package main

import (
	"fmt"
	"os"
	"testing"

	db "./database"
	"github.com/fatih/structs"
	"github.com/stretchr/testify/assert"
	"gopkg.in/mgo.v2/bson"
)

type person struct {
	ID   bson.ObjectId `bson:"_id"`
	Name string
	Age  int
	Food string
}

const collection = "collection"

var people = []interface{}{
	person{bson.NewObjectId(), "John", 23, "Ice-Cream"},
	person{bson.NewObjectId(), "Lisa", 33, "Cake"},
	person{bson.NewObjectId(), "Mark", 42, "Yoghurt"},
	person{bson.NewObjectId(), "Martha", 23, "Ice-Cream"},
}

func TestMain(m *testing.M) {
	retCode := m.Run()
	os.Exit(retCode)
}

func loadedCRUD() *db.CRUD {
	crud := db.NewCRUD(nil)
	crud.Insert(collection, people...)
	return crud
}

func Test_CrudFindAll(t *testing.T) {
	crud := loadedCRUD()

	// prepare results
	r1, _ := crud.FindAll(collection, nil)                          // expect all
	r2, _ := crud.FindAll(collection, &bson.M{"Food": "Ice-Cream"}) // expect some
	r3, _ := crud.FindAll(collection, &bson.M{"Name": "Jake"})      // expect none

	// make assertions
	assert := assert.New(t)
	assert.Equal(len(people), len(r1), "crud.FindAll with nil query does not return all results")
	assert.Equal(2, len(r2), "crud.FindAll with query does not return the right results")
	assert.Equal(0, len(r3), "crud.FindAll with query does not return 0 results")
}

func Test_CrudFindOne(t *testing.T) {
	crud := loadedCRUD()
	p0 := bson.M(structs.Map(people[0]))

	// prepare results
	i1, _ := crud.FindOne(collection, nil)                          // expect one
	i2, _ := crud.FindOne(collection, &bson.M{"Food": "Ice-Cream"}) // expect one
	r3, _ := crud.FindOne(collection, &bson.M{"Name": "Jake"})      // expect none

	r1 := i1.(bson.M)
	fmt.Println(r1)
	r2 := i2.(bson.M)

	//make assertions
	assert := assert.New(t)
	assert.Equal(p0["Name"], r1["Name"], "crud.FindOne with nil query does not return expected result")
	assert.Equal(p0["Name"], r2["Name"], "crud.FindOne with matching query does not return expected result")
	assert.Nil(r3, "crud.FindOne with non-matching query does not return expected result")
}

func Test_CrudUpdateID(t *testing.T) {
	crud := loadedCRUD()
	p0 := people[0].(person)
	oldName := p0.Name
	newName := "New Monicker"

	// prepare results
	crud.UpdateID(collection, p0.ID, bson.M{"Name": newName})
	i1, _ := crud.FindOne(collection, &bson.M{"_id": p0.ID})
	r1 := i1.(bson.M)

	// make asserts
	assert := assert.New(t)
	assert.NotEqual(oldName, newName, "New name and old name are the same.")
	assert.NotEqual(oldName, r1["Name"], "crud.UpdateID did not update the old value.")
	assert.Equal(newName, r1["Name"], "crud.UpdateID did not update to the correct value.")

}

func Test_CrudDeleteID(t *testing.T) {
	crud := loadedCRUD()
	p0 := people[0].(person)

	// prepare results
	b1, _ := crud.FindOne(collection, &bson.M{"_id": p0.ID})
	crud.DeleteID(collection, p0.ID)
	r1, _ := crud.FindOne(collection, &bson.M{"_id": p0.ID})

	// make asserts
	assert := assert.New(t)
	assert.NotNil(b1, "An expected value does not exist in the db.")
	assert.Nil(r1, "crud.DeleteID did not remove value from the db.")
}
