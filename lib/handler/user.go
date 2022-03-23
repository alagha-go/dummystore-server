package handler

import (
	u "dummystore/lib/user"
	v "dummystore/lib/variables"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)



func CreateUser(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	if req.Method != "POST" {
		res.WriteHeader(405)
		return
	}
	req.ParseMultipartForm(0)
	data := req.FormValue("data")
	var user u.User
	err := json.Unmarshal([]byte(data), &user)
	if err != nil {
		res.WriteHeader(400)
		json.NewEncoder(res).Encode(Error{Error: "invalid json data"})
		return
	}
	user.ImageFile, _,  _ = req.FormFile("image")
	token, status, err := u.CreateUser(user)
	if err != nil {
		res.WriteHeader(status)
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}
	res.Header().Add("Authorization", token.Token)
	http.SetCookie(res,
		&http.Cookie{
			Name: "Authorization",
			Value: token.Token,
			Expires: token.Expires,
			Path: "/",
		},
	)
	res.WriteHeader(status)
	json.NewEncoder(res).Encode(token)
}


func Login(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	if req.Method != "POST" {
		res.WriteHeader(405)
		return
	}
	var user u.User
	json.NewDecoder(req.Body).Decode(&user)
	token, status, err := u.Login(user)
	if err != nil {
		res.WriteHeader(status)
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}
	res.Header().Add("Authorization", token.Token)
	http.SetCookie(res,
		&http.Cookie{
			Name: "Authorization",
			Value: token.Token,
			Expires: token.Expires,
			Path: "/",
		},
	)
	json.NewEncoder(res).Encode(token)
}

func GetUser(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	if req.Method != "GET" {
		res.WriteHeader(405)
		return
	}
	user, status, err := VerifyUser(req)
	res.WriteHeader(status)
	if err != nil {
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}
	json.NewEncoder(res).Encode(user)
}


func RefreshToken(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	if req.Method != "GET" {
		res.WriteHeader(405)
		return
	}
	_, status, err := VerifyUser(req)
	if err != nil {
		res.WriteHeader(status)
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}

	token, err := GetAuth(req)
	if err != nil {
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}
	newToken, err := u.RefreshToken(token)
	if err != nil {
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}

	http.SetCookie(res,
		&http.Cookie{
			Name: "Authorization",
			Value: newToken.Token,
			Expires: newToken.Expires,
		},
	)
	json.NewEncoder(res).Encode(newToken)
}


func DeleteAccount(res http.ResponseWriter, req *http.Request)   {
	res.Header().Add("content-type", "application/json")
	if req.Method != "DELETE" {
		res.WriteHeader(405)
		return
	}
	user, status, err := VerifyUser(req)
	if err != nil && status != 401{
		res.WriteHeader(status)
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}
	deleted, status, err := u.DeleteUser(user, true)
	res.WriteHeader(status)
	if !deleted {
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}
}

func DeleteUser(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	if req.Method != "DELETE" {
		res.WriteHeader(405)
		return
	}
	var user u.User
	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		res.WriteHeader(400)
		json.NewEncoder(res).Encode(Error{Error: "invalid json"})
	}

	deleted, status, err := u.DeleteUser(user, false)
	res.WriteHeader(status)
	if !deleted {
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}
	res.WriteHeader(http.StatusOK) 
}

func UpdateUser(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	if req.Method != "PUT" && req.Method != "UPDATE" {
		res.WriteHeader(http.StatusBadGateway)
		return
	}
	_, status, err := VerifyUser(req)
	if err != nil {
		res.WriteHeader(status)
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}
	var user u.User
	json.NewDecoder(req.Body).Decode(&user)
	user, status, err = u.UpdateUser(user)
	if err != nil {
		res.WriteHeader(status)
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}
	res.WriteHeader(status)
	json.NewEncoder(res).Encode(user)
}

func ProfileImage(res http.ResponseWriter, req *http.Request) {
	if req.Method != "GET" {
		res.WriteHeader(405)
		return
	}
	user, status, err := VerifyUser(req)
	if err != nil {
		res.WriteHeader(status)
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}

	http.ServeFile(res, req, fmt.Sprintf("./profiles/%s.png", user.ID.Hex()))
}


func VerifyUser(req *http.Request) (u.User, int, error) {
	var User u.User
	token, err := GetAuth(req)
	if err != nil {
		return User, 401, errors.New("make sure to provide Authorization token")
	}
	valid, user := u.ValidateToken(token)
	if !valid {
		return User, 401, errors.New("provide a valid token")
	}
	data, _ := json.Marshal(user)
	json.Unmarshal(data, &User)
	exist := u.CheckIfUserExists(User.ID)
	if !exist {
		return User, 404, v.UserDoesNotExist
	}
	return User, 200, nil
}