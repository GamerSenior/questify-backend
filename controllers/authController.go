package controllers

import (
	"encoding/json"
	"github.com/GamerSenior/questify-backend/models"
	u "github.com/GamerSenior/questify-backend/utils"
	"net/http"
)

// CreateAccount explanation
var CreateAccount = func(w http.ResponseWriter, r *http.Request) {
	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return 
	}
	resp, _ := account.Create()
	u.Respond(w, resp)
}
// Authenticate explanation
var Authenticate = func(w http.ResponseWriter, r *http.Request) {
	account := &models.Account{}
	err := json.NewDecoder(r.Body).Decode(account)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid request"))
		return
	}
	resp := models.Login(account.Email, account.Password)
	u.Respond(w, resp)
}