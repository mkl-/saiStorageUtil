package saiStorageUtil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/saiset-co/saiUtil"
	"go.mongodb.org/mongo-driver/bson"
)

type Database struct {
	url string
	email string
	password string
}

func Storage(Url string, Email string, Password string) Database {
	return Database{
		url: Url,
		email: Email,
		password: Password,
	}
}

type Token struct {
	Id    string `json:"_id"`
	Token string `json:"token"`
}

type Login struct {
	email    string
	password string
}

type StorageRequest struct {
	login        Login
	token        string
	collection   string
	options      interface{}
	criteria     interface{}
	data         interface{}
}

var myToken Token

func (s StorageRequest) toJson() ([]byte, error) {
	if (Login{}) != s.login {
		return json.Marshal(bson.M{"email": s.login.email, "password": s.login.password})
	}

	jsonObj := bson.M{"collection": s.collection, "token": s.token}

	if s.data != nil {
		jsonObj["data"] = s.data
	}

	if s.criteria != nil {
		jsonObj["select"] = s.criteria
	}

	if s.options != nil {
		jsonObj["options"] = s.options
	}

	return json.Marshal(jsonObj)
}

func (db Database) login() {
	request := StorageRequest{login: Login{email: db.email, password: db.password}}
	err, token := db.makeRequest("login", request)

	if err != nil {
		fmt.Println("Login error: ", err)
	}

	err = json.Unmarshal(token, &myToken)

	if err != nil {
		fmt.Println("Login error: ", err)
	}
}

func (db Database) Get(collectionName string, criteria interface{}, options interface{}) (error, []byte) {
	request := StorageRequest{collection: collectionName, criteria: criteria, options: options}
	return db.makeRequest("get", request)
}

func (db Database) Put(collectionName string, data interface{}) (error, []byte) {
	request := StorageRequest{collection: collectionName, data: data}
	return db.makeRequest("save", request)
}

func (db Database) Update(collectionName string, criteria interface{}, data interface{}) (error, []byte)  {
	request := StorageRequest{collection: collectionName, criteria: criteria, data: data, options: "set"}
	return db.makeRequest("update", request)
}

func (db Database) makeRequest(method string, request StorageRequest) (error, []byte) {
	if method != "login" && myToken.Token == "" {
		db.login()
	}

	request.token = myToken.Token
	jsonStr, jsonErr := request.toJson()

	if jsonErr != nil {
		fmt.Println("Database request error: ", jsonErr)
		return jsonErr, []byte("")
	}

	return saiUtil.Send(db.url + "/" + method, bytes.NewBuffer(jsonStr))
}
