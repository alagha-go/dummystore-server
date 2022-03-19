package handler

import (
	u "dummystore/lib/user"
	"encoding/json"
	"errors"
	"net/http"
)



func CreateUser(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	if req.Method != "POST" {
		res.WriteHeader(405)
		return
	}
	var user u.User
	json.NewDecoder(req.Body).Decode(&user)
	token, err, status := u.CreateUser(user)
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


func Login(res http.ResponseWriter, req *http.Request) {
	res.Header().Add("content-type", "application/json")
	if req.Method != "POST" {
		res.WriteHeader(405)
		return
	}
	var user u.User
	json.NewDecoder(req.Body).Decode(&user)
	token, err, status := u.Login(user)
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
	user, err, status := VerifyUser(req)
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
	_, err, status := VerifyUser(req)
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
	user, err, status := VerifyUser(req)
	if err != nil && status != 401{
		res.WriteHeader(status)
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}
	deleted, err, status := u.DeleteUser(user, true)
	res.WriteHeader(status)
	if !deleted {
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}
	res.WriteHeader(http.StatusOK) 
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

	deleted, err, status := u.DeleteUser(user, false)
	res.WriteHeader(status)
	if !deleted {
		json.NewEncoder(res).Encode(Error{Error: err.Error()})
		return
	}
	res.WriteHeader(http.StatusOK) 
}


func VerifyUser(req *http.Request) (u.User, error, int) {
	var User u.User
	token, err := GetAuth(req)
	if err != nil {
		return User, errors.New("make sure to provide Authorization token"), 401
	}
	valid, user := u.ValidateToken(token)
	if !valid {
		return User, errors.New("provide a valid token"), 401
	}
	data, _ := json.Marshal(user)
	json.Unmarshal(data, &User)
	exist := u.CheckIfUserExists(User.ID)
	if !exist {
		return User, errors.New("User does not exist"), 404
	}
	return User, nil, 200
}