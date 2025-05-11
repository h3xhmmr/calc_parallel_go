package application

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthRequest struct {
	Password string `json:"password"`
	Login    string `json:"login"`
}

type LoginResponse struct {
	AcsToken string `json:"access_token"`
}

const (
	invalidHeader       = "error: invalid header"
	invalidToken        = "error: invalid token"
	invalidPass         = "error: invalid password"
	invalidBody         = "error: invalid body"
	userAlreadyExists   = "error: user already exists"
	userNotFound        = "error: user not found"
	internalServerError = "error: internal server error"
)

func (o *Orchestrator) MakeToken(id string) string {
	now := time.Now()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":  id,
		"nbf": now.Unix(),
		"exp": now.Add(time.Hour * 24).Unix(),
		"iat": now.Unix(),
	})
	tokenString, err := token.SignedString([]byte(o.jwt_key))
	if err != nil {
		panic(err)
	}
	return tokenString
}

func makeLoginResponse(token string, w http.ResponseWriter) {
	b, err := json.Marshal(LoginResponse{AcsToken: token})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(b)
}

func writeError(w http.ResponseWriter, text string, statusCode int) {
	w.WriteHeader(statusCode)
	fmt.Fprint(w, text)
}

func (o *Orchestrator) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	req := new(AuthRequest)
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		writeError(w, invalidBody, http.StatusUnprocessableEntity)
		return
	}
	if len(req.Password) < 5 {
		writeError(w, invalidPass, http.StatusUnprocessableEntity)
		return
	}
	_, ok := app.GetUser(req.Login, req.Password)
	if ok {
		writeError(w, userAlreadyExists, http.StatusUnprocessableEntity)
		return
	}
	_, err = app.AddUser(req.Login, req.Password)
	if err != nil {
		writeError(w, internalServerError, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (o *Orchestrator) LoginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
	defer r.Body.Close()
	req := new(AuthRequest)
	err := json.NewDecoder(r.Body).Decode(req)
	if err != nil {
		http.Error(w, invalidBody, http.StatusUnprocessableEntity)
		return
	}
	u, ok := o.GetUser(req.Login, req.Password)
	if !ok {
		http.Error(w, userNotFound, http.StatusUnauthorized)
		return
	}
	makeLoginResponse(o.MakeToken(u.ID), w)
}
