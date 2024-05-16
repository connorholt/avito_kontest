package main

var secretKey = []byte("secret_token")

type RegisterRequest struct {
	Name     string `json:"name" example:"John"`
	Password string `json:"password" example:"123"`
	Role     string `json:"role" example:"admin"`
}

type LoginRequest struct {
	Name     string `json:"name" example:"admin"`
	Password string `json:"password" example:"123"`
}

type LoginToken struct {
	AccessToken string `json:"access_token"`
}

type contextKey string

var roleKey = contextKey("role")
