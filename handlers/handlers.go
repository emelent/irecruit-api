package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/gorilla/mux"

	db "../database"
	models "../models"
	"fmt"
	"gopkg.in/mgo.v2/bson"
)

const productsCollection = "products"

func jsonEncode(w http.ResponseWriter, v interface{}, statusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if err := json.NewEncoder(w).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func jsonDecode(r *http.Request, v models.Validator) error {
	if err := json.NewDecoder(r.Body).Decode(v); err != nil {
		return nil
	}
	return v.OK()
}

//NewProductHandler endpoint
func NewProductHandler(crud *db.CRUD) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer crud.CloseCopy()
		fmt.Println(crud.TempStorage[productsCollection])
		var p models.Product
		if err := jsonDecode(r, &p); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if err := p.OK(); err != nil {
			http.Error(w, err.Error(), http.StatusUnprocessableEntity)
			return
		}
		p.ID = bson.NewObjectId()

		if err := crud.Insert(productsCollection, p); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		fmt.Println(crud.TempStorage[productsCollection])
		jsonEncode(w, p, http.StatusCreated)
	}
}

//PreflightSyncHandler endpoint
func PreflightSyncHandler(crud *db.CRUD) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer crud.CloseCopy()
		w.WriteHeader(http.StatusOK)
		return
	}
}

//SyncHandler endpoint
func SyncHandler(crud *db.CRUD) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer crud.CloseCopy()

		var rawProducts []interface{}
		if err := json.NewDecoder(r.Body).Decode(&rawProducts); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if lenProd := len(rawProducts); lenProd > 0 {
			products := make([]interface{}, lenProd)
			for i, raw := range rawProducts {
				p := mapToProd(raw.(map[string]interface{}))
				if err := p.OK(); err != nil {
					http.Error(w, err.Error(), http.StatusUnprocessableEntity)
					return
				}

				p.ID = bson.NewObjectId()
				products[i] = p
			}

			if err := crud.Insert(productsCollection, products...); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}
		}
		allRaw, err := crud.FindAll(productsCollection, nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		allProducts := make([]models.Product, len(allRaw))
		for i, rawProd := range allRaw {
			switch prod := rawProd.(type) {
			case bson.M:
				allProducts[i] = bsonToProd(prod)
			case models.Product:
				allProducts[i] = prod
			}
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(&allProducts); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}
}

//InitHandler just dummy db inserts
func InitHandler(crud *db.CRUD) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer crud.CloseCopy()

		products := []interface{}{
			map[string]interface{}{
				"id":       bson.NewObjectId(),
				"name":     "Rusher",
				"brand":    "Coala",
				"category": "Misc",
				"quantity": 400,
				"sell":     50.5,
				"buy":      30.0,
			},
			models.Product{bson.NewObjectId(), "Fouton Red", "Seaters", "Furniture", 30, 500.99, 100.0},
			models.Product{bson.NewObjectId(), "Fouton Blue", "Seaters", "Furniture", 30, 500.99, 100.0},
			models.Product{bson.NewObjectId(), "Fouton White", "Seaters", "Furniture", 30, 500.99, 100.0},
		}

		if err := crud.Insert(productsCollection, products...); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		jsonEncode(w, "Initialization successful", http.StatusOK)
	}
}

//AllProductsHandler endpoint
func AllProductsHandler(crud *db.CRUD) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer crud.CloseCopy()

		var results []interface{}
		results, err := crud.FindAll(productsCollection, nil)
		products := make([]models.Product, len(results))
		for index, rawProd := range results {
			switch prod := rawProd.(type) {
			case bson.M:
				products[index] = bsonToProd(prod)
			case models.Product:
				products[index] = prod
			}
		}
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		jsonEncode(w, products, http.StatusOK)
	}
}

func bsonToProd(bsonProd map[string]interface{}) models.Product {
	return models.Product{
		ID:       bsonProd["_id"].(bson.ObjectId),
		Name:     bsonProd["name"].(string),
		Brand:    bsonProd["brand"].(string),
		Category: bsonProd["category"].(string),
		Quantity: bsonProd["quantity"].(int),
		Sell:     bsonProd["sell"].(float64),
		Buy:      bsonProd["buy"].(float64),
	}
}

func mapToProd(mapProd map[string]interface{}) models.Product {
	return models.Product{
		ID:       "",
		Name:     mapProd["name"].(string),
		Brand:    mapProd["brand"].(string),
		Category: mapProd["category"].(string),
		Quantity: int(mapProd["quantity"].(float64)),
		Sell:     mapProd["sell"].(float64),
		Buy:      mapProd["buy"].(float64),
	}
}

//UpdateProductHandler endpoint
func UpdateProductHandler(crud *db.CRUD) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer crud.CloseCopy()

		params := mux.Vars(r)
		productID := bson.ObjectIdHex(params["productID"])
		var p models.Product
		if err := json.NewDecoder(r.Body).Decode(&p); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		p.ID = productID
		if err := crud.UpdateID(productsCollection, productID, p); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonEncode(w, p, http.StatusOK)
	}
}

//DeleteProductHandler endpoint
func DeleteProductHandler(crud *db.CRUD) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer crud.CloseCopy()

		params := mux.Vars(r)
		productID := bson.ObjectIdHex(params["productID"])

		if err := crud.DeleteID(productsCollection, productID); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		jsonEncode(w, "Product successfully deleted.", http.StatusOK)
	}
}
