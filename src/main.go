package main

import (
	"fmt"
	"github.com/Ryan-Har/site-monitor/src/handlers"
	"github.com/Ryan-Har/site-monitor/src/internal/auth"
	"github.com/Ryan-Har/site-monitor/src/internal/database"
	"github.com/Ryan-Har/site-monitor/src/internal/scheduler"
	"net/http"
)

func main() {
	fb := auth.NewServer()
	dbh, err := database.NewSQLiteHandler()
	if err != nil {
		fmt.Println("unable to initialise database", err)
	}
	go scheduler.StartSchedulers(*dbh)
	//statics
	http.Handle("/assets/", http.StripPrefix("/assets/", http.FileServer(http.Dir("assets"))))

	//serve html
	http.HandleFunc("GET /signup", handlers.NewGetSignupHandler().ServeHTTP)
	http.HandleFunc("GET /login", handlers.NewGetLoginHandler().ServeHTTP)
	http.HandleFunc("GET /forgottenpassword", handlers.NewGetResetPasswordHandler().ServeHTTP)
	http.Handle("GET /monitors", fb.AuthMiddleware(http.HandlerFunc(handlers.NewGetMonitorOverviewHandler(*dbh).ServeHTTP)))
	http.Handle("GET /monitors/new", fb.AuthMiddleware(http.HandlerFunc(handlers.NewGetMonitorFormHandler().ServeHTTP)))
	http.Handle("GET /monitors/{monitorid}", fb.AuthMiddleware(http.HandlerFunc(handlers.NewGetMonitorByID(*dbh).ServeHTTP)))

	// http.Handle("GET /monitors/{monitorid}/edit", fb.AuthMiddleware(http.HandlerFunc(handlers.NewGetMonitorByID().ServeHTTP)))
	http.HandleFunc("GET /monitors/getCreateFormInfo", handlers.NewGetMonitorFormHandler().ServeFormContent)
	http.Handle("GET /maintenance", fb.AuthMiddleware(http.HandlerFunc(handlers.NewGetMaintenanceHandler().ServeHTTP)))
	http.Handle("GET /incidents", fb.AuthMiddleware(http.HandlerFunc(handlers.NewGetIncidentsHandler().ServeHTTP)))
	http.Handle("GET /settings/account", fb.AuthMiddleware(http.HandlerFunc(handlers.NewGetAccountSettingsHandler().ServeHTTP)))
	http.Handle("GET /settings/notifications", fb.AuthMiddleware(http.HandlerFunc(handlers.NewGetNotificationSettingsHandler(*dbh).ServeHTTP)))
	http.Handle("GET /settings/notifications/{notificationid}", fb.AuthMiddleware(http.HandlerFunc(handlers.NewGetEditNotificationByID(*dbh).ServeHTTP)))

	http.Handle("GET /settings/security", fb.AuthMiddleware(http.HandlerFunc(handlers.NewGetSecuritySettingsHandler().ServeHTTP)))
	http.HandleFunc("GET /settings/getNotificationFormInfo", handlers.NewGetNotificationSettingsHandler(*dbh).ServeFormContent)

	//serve json
	http.Handle("GET /monitors/{monitorid}/responsetime", fb.AuthMiddleware(http.HandlerFunc(handlers.NewGetMonitorByID(*dbh).ServeResponseTimes)))

	//perform actions
	http.HandleFunc("POST /verifylogin", fb.VerifyLogin)
	http.HandleFunc("POST /updateauthcookie", fb.UpdateAuthCookie)
	http.HandleFunc("POST /notifications/sendtest", handlers.NewGetNotificationSettingsHandler(*dbh).SendTestNotification)
	http.Handle("DELETE /monitors/{monitorid}", fb.AuthMiddleware(http.HandlerFunc(handlers.NewDeleteMonitorByID(*dbh).ServeHTTP)))
	http.Handle("PUT /settings/notifications/{notificationid}", fb.AuthMiddleware(http.HandlerFunc(handlers.NewPostNotificationSettingsHandler(*dbh).ByID)))
	http.Handle("DELETE /settings/notifications/{notificationid}", fb.AuthMiddleware(http.HandlerFunc(handlers.NewDeleteNotificationSettingsHandler(*dbh).ByID)))

	//forms
	http.Handle("POST /monitors/new", fb.AuthMiddleware(http.HandlerFunc(handlers.NewPostFormHandler(*dbh).NewMonitorForm)))
	http.Handle("POST /notifications/new", fb.AuthMiddleware(http.HandlerFunc(handlers.NewPostFormHandler(*dbh).NewNotificationForm)))

	//form validations
	http.HandleFunc("POST /validation/monitorlocationhttp", handlers.NewValidationFormHandler().ValidateMonitorLocationHttp)
	http.HandleFunc("POST /validation/monitorlocationiporhost", handlers.NewValidationFormHandler().ValidateMonitorLocationIpOrHost)
	http.HandleFunc("POST /validation/monitorportnumber", handlers.NewValidationFormHandler().ValidateMonitorPortNumber)

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}
