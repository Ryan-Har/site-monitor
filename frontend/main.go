package main

import (
	"fmt"
	"net/http"

	"github.com/Ryan-Har/site-monitor/frontend/templates"
	"github.com/a-h/templ"
)

func main() {

	http.Handle("/signup", templ.Handler(templates.SignUp()))
	http.Handle("/login", templ.Handler(templates.Login()))
	http.Handle("/monitors", templ.Handler(templates.Monitors()))

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}
