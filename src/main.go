package main

import (
	"fmt"
	"github.com/Ryan-Har/site-monitor/src/handlers"
	"github.com/Ryan-Har/site-monitor/src/internal/auth"
	"github.com/Ryan-Har/site-monitor/src/internal/database"
	"net/http"
)

func main() {
	fb := auth.NewServer()
	_, err := database.NewSQLiteHandler()
	if err != nil {
		fmt.Println("unable to initialise database", err)
	}

	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.HandleFunc("POST /verifylogin", fb.VerifyLogin)
	http.HandleFunc("POST /updateauthcookie", fb.UpdateAuthCookie)
	http.HandleFunc("GET /signup", handlers.NewGetSignupHandler().ServeHTTP)
	http.HandleFunc("GET /login", handlers.NewGetLoginHandler().ServeHTTP)
	http.Handle("GET /monitors", fb.AuthMiddleware(http.HandlerFunc(handlers.NewGetMonitorOverviewHandler().ServeHTTP)))
	http.Handle("GET /monitors/new", fb.AuthMiddleware(http.HandlerFunc(handlers.NewGetMonitorFormHandler().ServeHTTP)))
	http.Handle("GET /monitors/{monitorid}", fb.AuthMiddleware(http.HandlerFunc(handlers.NewGetMonitorByID().ServeHTTP)))
	http.HandleFunc("GET /monitors/getCreateFormInfo", handlers.NewGetMonitorFormHandler().ServeFormContent)
	http.Handle("GET /maintenance", fb.AuthMiddleware(http.HandlerFunc(handlers.NewGetMaintenanceHandler().ServeHTTP)))
	http.Handle("GET /incidents", fb.AuthMiddleware(http.HandlerFunc(handlers.NewGetIncidentsHandler().ServeHTTP)))
	http.Handle("GET /settings/account", fb.AuthMiddleware(http.HandlerFunc(handlers.NewGetAccountSettingsHandler().ServeHTTP)))
	http.Handle("GET /settings/notifications", fb.AuthMiddleware(http.HandlerFunc(handlers.NewGetNotificationSettingsHandler().ServeHTTP)))
	http.Handle("GET /settings/security", fb.AuthMiddleware(http.HandlerFunc(handlers.NewGetSecuritySettingsHandler().ServeHTTP)))

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}
