package main

import (
	"fmt"
	"net/http"

	"github.com/Ryan-Har/site-monitor/frontend/handlers"
)

func main() {
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	http.HandleFunc("/signup", handlers.NewGetSignupHandler().ServeHTTP)
	http.HandleFunc("/login", handlers.NewGetLoginHandler().ServeHTTP)
	http.HandleFunc("/monitors", handlers.NewGetMonitorOverviewHandler().ServeHTTP)
	http.HandleFunc("/monitors/new", handlers.NewGetMonitorFormHandler().ServeHTTP)
	http.HandleFunc("/monitors/getCreateFormInfo", handlers.NewGetMonitorFormHandler().ServeFormContent)
	http.HandleFunc("/maintenance", handlers.NewGetMaintenanceHandler().ServeHTTP)
	http.HandleFunc("/incidents", handlers.NewGetIncidentsHandler().ServeHTTP)
	http.HandleFunc("/settings/account", handlers.NewGetAccountSettingsHandler().ServeHTTP)
	http.HandleFunc("/settings/notifications", handlers.NewGetNotificationSettingsHandler().ServeHTTP)
	http.HandleFunc("/settings/security", handlers.NewGetSecuritySettingsHandler().ServeHTTP)

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}
