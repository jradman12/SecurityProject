package handlers

import (
	"common/module/interceptor"
	"context"
	"encoding/json"
	"fmt"
	"gateway/module/application/helpers"
	"gateway/module/application/services"
	"gateway/module/auth"
	"gateway/module/domain/dto"
	modelGateway "gateway/module/domain/model"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
	"log"
	"net/http"
	"time"
)

type AuthenticationHandler struct {
	l                   *log.Logger
	service             *services.UserService
	validator           *validator.Validate
	passwordUtil        *helpers.PasswordUtil
	passwordlessService *services.PasswordLessService
}

func NewAuthenticationHandler(l *log.Logger, service *services.UserService, validator *validator.Validate,
	passwordUtil *helpers.PasswordUtil, passwordlessService *services.PasswordLessService) Handler {
	return &AuthenticationHandler{l, service, validator, passwordUtil, passwordlessService}
}

func (a AuthenticationHandler) Init(mux *runtime.ServeMux) {
	err := mux.HandlePath("POST", "/users/login/user", a.LoginUser)
	err2 := mux.HandlePath("POST", "/users/login/passwordless", a.PasswordLessLoginReq)
	if err != nil {
		panic(err)
	}
	if err2 != nil {
		panic(err)
	}
}

func (a AuthenticationHandler) LoginUser(rw http.ResponseWriter, r *http.Request, params map[string]string) {
	a.l.Println("Handling LOGIN Users")

	var loginRequest dto.LoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(rw, "Error decoding loginRequest:"+err.Error(), http.StatusBadRequest)
		return
	}

	user, err := a.service.GetByUsername(context.TODO(), loginRequest.Username)
	if err != nil {
		http.Error(rw, "User not found! "+err.Error(), http.StatusBadRequest)
		return
	}
	if !user.IsConfirmed {
		fmt.Println("account not activated")
		http.Error(rw, "User account not activated! ", http.StatusBadRequest)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	if err != nil {
		fmt.Println(err)
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}

	var claims = &interceptor.JwtClaims{}
	claims.Username = loginRequest.Username

	userRoles, err := a.service.GetUserRole(loginRequest.Username)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return

	}

	claims.Roles = append(claims.Roles, userRoles)
	var tokenCreationTime = time.Now().UTC()
	var tokenExpirationTime = tokenCreationTime.Add(time.Duration(30) * time.Minute)

	token, err := auth.GenerateToken(claims, tokenExpirationTime)

	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return

	}

	var roleString string

	if user.Role == modelGateway.Admin {
		roleString = "Admin"
	} else if user.Role == modelGateway.Agent {
		roleString = "Agent"
	} else if user.Role == modelGateway.Regular {
		roleString = "Regular"
	}

	logInResponse := dto.LogInResponseDto{
		Token:    token,
		Role:     roleString,
		Email:    user.Email,
		Username: user.Username,
	}

	logInResponseJson, _ := json.Marshal(logInResponse)
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(logInResponseJson)

}

func (a AuthenticationHandler) PasswordLessLoginReq(rw http.ResponseWriter, r *http.Request, params map[string]string) {
	a.l.Println("Handling PasswordLessLoginReq")

	var loginRequest dto.PasswordLessLoginRequest
	err := json.NewDecoder(r.Body).Decode(&loginRequest)
	if err != nil {
		http.Error(rw, "Error decoding loginRequest:"+err.Error(), http.StatusBadRequest)
		return
	}

	user, err := a.service.GetByUsername(context.TODO(), loginRequest.Username)
	if err != nil {
		http.Error(rw, "User not found! "+err.Error(), http.StatusBadRequest)
		return
	}
	if !user.IsConfirmed {
		fmt.Println("account not activated")
		http.Error(rw, "User account not activated! ", http.StatusBadRequest)
		return
	}
	a.passwordlessService.SendMagicLink(context.TODO(), "https://localhost:4200", "http://localhost:9090/", user)

	//err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(loginRequest.Password))
	//if err != nil {
	//	fmt.Println(err)
	//	http.Error(rw, err.Error(), http.StatusBadRequest)
	//	return
	//}
	//
	//var claims = &interceptor.JwtClaims{}
	//claims.Username = loginRequest.Username
	//
	//userRoles, err := a.service.GetUserRole(loginRequest.Username)
	//if err != nil {
	//	http.Error(rw, err.Error(), http.StatusBadRequest)
	//	return
	//
	//}
	//
	//claims.Roles = append(claims.Roles, userRoles)
	//var tokenCreationTime = time.Now().UTC()
	//var tokenExpirationTime = tokenCreationTime.Add(time.Duration(30) * time.Minute)
	//
	//token, err := auth.GenerateToken(claims, tokenExpirationTime)
	//
	//if err != nil {
	//	http.Error(rw, err.Error(), http.StatusBadRequest)
	//	return
	//
	//}
	//
	//var roleString string
	//
	//if user.Role == modelGateway.Admin {
	//	roleString = "Admin"
	//} else if user.Role == modelGateway.Agent {
	//	roleString = "Agent"
	//} else if user.Role == modelGateway.Regular {
	//	roleString = "Regular"
	//}
	//
	logInResponse := dto.LogInResponseDto{
		Token:    "",
		Role:     "",
		Email:    "",
		Username: "",
	}

	logInResponseJson, _ := json.Marshal(logInResponse)
	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(logInResponseJson)

}
