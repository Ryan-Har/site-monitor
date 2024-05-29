package auth

import (
	"context"
	//"encoding/json"
	//"io"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"fmt"
	"github.com/Ryan-Har/site-monitor/src/config"
	"github.com/Ryan-Har/site-monitor/src/models"
	"google.golang.org/api/option"
	"net/http"
	"time"
)

type Server struct {
	app *firebase.App
}

func NewServer() *Server {
	opt := option.WithCredentialsFile(config.GetConfig().FIREBASE_SERVICE_ACCOUNT_LOCATION)
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		fmt.Printf("error initializing firebase app: %v\n", err)
	}
	return &Server{app: app}
}

func (s *Server) verifyIDToken(idToken string) (*auth.Token, error) {
	client, err := s.app.Auth(context.Background())
	if err != nil {
		return nil, err
	}
	token, err := client.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		return nil, err
	}
	return token, nil
}

func (s *Server) setAuthCookie(w http.ResponseWriter, token string) {
	cookie := http.Cookie{
		Name:     "auth-token",
		Value:    token,
		HttpOnly: true,
		Secure:   false,
		Path:     "/",
		SameSite: http.SameSiteStrictMode,
		Expires:  time.Now().Add(60 * time.Minute),
	}
	http.SetCookie(w, &cookie)
}

func (s *Server) VerifyLogin(w http.ResponseWriter, r *http.Request) {
	idToken := r.Header.Get("Authorization")

	token, err := s.verifyIDToken(idToken)
	if err != nil {
		http.Error(w, "Invalid token", http.StatusUnauthorized)
		return
	}
	fmt.Printf("Verified token: %v", token)

	s.setAuthCookie(w, idToken)
	w.Write([]byte(`{"Response": "Token Valid"}`))
}

func (s *Server) UpdateAuthCookie(w http.ResponseWriter, r *http.Request) {
	idToken := r.Header.Get("Authorization")
	s.setAuthCookie(w, idToken)
	w.Write([]byte(`{"Response": "Token Updated"}`))
}

func (s *Server) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		cookie, err := r.Cookie("auth-token")
		if err != nil {
			fmt.Println("error retrieving auth-token cookie for request", err)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		//verify the token is still valid
		token, err := s.verifyIDToken(cookie.Value)
		if err != nil {
			fmt.Println("id token not currently valid, need to reauth", err)
			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		name, ok := token.Claims["name"].(string)
		if !ok {
			name = "user"
		}

		email, ok := token.Claims["email"].(string)
		if !ok {
			email = "email"
		}

		userInfo := models.UserInfo{
			UUID:  token.UID,
			Name:  name,
			Email: email,
		}

		// Add the variable to the request context
		ctx := context.WithValue(r.Context(), models.UserInfoKey, userInfo)

		// Call the next handler with the new context
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
