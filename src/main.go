package main

import (
	"fmt"
	"net/http"

	"github.com/Ryan-Har/site-monitor/src/handlers"
)

func main() {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.HandleFunc("GET /signup", handlers.NewGetSignupHandler().ServeHTTP)
	http.HandleFunc("GET /login", handlers.NewGetLoginHandler().ServeHTTP)
	http.HandleFunc("GET /monitors", handlers.NewGetMonitorOverviewHandler().ServeHTTP)
	http.HandleFunc("GET /monitors/new", handlers.NewGetMonitorFormHandler().ServeHTTP)
	http.HandleFunc("/monitors/{monitorid}", handlers.NewGetMonitorByID().ServeHTTP)
	http.HandleFunc("GET /monitors/getCreateFormInfo", handlers.NewGetMonitorFormHandler().ServeFormContent)
	http.HandleFunc("GET /maintenance", handlers.NewGetMaintenanceHandler().ServeHTTP)
	http.HandleFunc("GET /incidents", handlers.NewGetIncidentsHandler().ServeHTTP)
	http.HandleFunc("GET /settings/account", handlers.NewGetAccountSettingsHandler().ServeHTTP)
	http.HandleFunc("GET /settings/notifications", handlers.NewGetNotificationSettingsHandler().ServeHTTP)
	http.HandleFunc("GET /settings/security", handlers.NewGetSecuritySettingsHandler().ServeHTTP)

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}
