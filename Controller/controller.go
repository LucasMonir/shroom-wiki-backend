package controller

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	auth "shroom-wiki-backend/Auth"
	middleware "shroom-wiki-backend/Middleware"
	"time"

	"github.com/google/uuid"
)

type Handler struct {
	UserData map[uuid.UUID]auth.User
}

func (h *Handler) Init() {
	h.UserData = make(map[uuid.UUID]auth.User)
}

func (h *Handler) CreateUser(writter http.ResponseWriter, request *http.Request) {
	requestBody, _ := io.ReadAll(request.Body)
	var tempUser auth.User

	json.Unmarshal(requestBody, &tempUser)

	log.Printf("%v", tempUser)

	tempUser.Password, _ = auth.HashPassword(tempUser.Password)
	userId := uuid.New()
	h.UserData[userId] = tempUser

	generatedToken := auth.GenerateToken(userId, tempUser)

	http.SetCookie(writter, &http.Cookie{
		Name:     "jwt",
		Value:    generatedToken,
		Expires:  time.Now().Add(time.Hour * 3),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
	})

	middleware.JSON(writter, http.StatusCreated, "OK")
}

func (h *Handler) Auth(writter *http.ResponseWriter, request *http.Request) {
	middleware.JSON(*writter, http.StatusOK, "OK")
}

func (h *Handler) Logout(writter http.ResponseWriter, request *http.Request) {
	http.SetCookie(writter, &http.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
		SameSite: http.SameSiteStrictMode,
		Secure:   true,
	})

	middleware.JSON(writter, http.StatusOK, "OK")
}

func (h *Handler) GetUser(writter http.ResponseWriter, request *http.Request) {
	token, _ := request.Cookie("jwt")
	uid, err := auth.GetTokenIdFromCookie(token)

	if err != nil {
		print(err.Error())
	}

	userData := h.UserData[uid]
	middleware.JSON(writter, http.StatusOK, map[string]string{
		"username": userData.Username,
		"email":    userData.Email,
	})
}

func (h *Handler) Login(writter http.ResponseWriter, request *http.Request) {
	requestBody, _ := io.ReadAll(request.Body)

	var tempUser auth.User

	json.Unmarshal(requestBody, &tempUser)
	log.Printf("%v", tempUser)

	for uid, value := range h.UserData {
		if value.Email == tempUser.Email || value.Username == tempUser.Username {
			generatedTokend := auth.GenerateToken(uid, value)

			http.SetCookie(writter, &http.Cookie{
				Name:     "jwt",
				Value:    generatedTokend,
				Expires:  time.Now().Add(time.Hour * 3),
				HttpOnly: true,
				SameSite: http.SameSiteStrictMode,
				Secure:   true,
			})

			middleware.JSON(writter, http.StatusCreated, "OK")
			return
		}
	}

	middleware.JSON(writter, http.StatusUnauthorized, "")
}
