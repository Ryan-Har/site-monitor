// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.707
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"github.com/Ryan-Har/site-monitor/src/templates/partials"
)

func Signup() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var1 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var1 == nil {
			templ_7745c5c3_Var1 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"container d-flex flex-column align-items-center justify-content-center h-50 w-50 mt-5\"><div class=\"mw-md mx-auto text-center\"><a class=\"d-inline-block mb-32\" href=\"#\"><img src=\"flaro-assets/logos/flaro-logo-black-xl.svg\" alt=\"\"></a><h3 class=\"mb-4\">Join for free</h3><p class=\"fs-8 text-secondary mb-12\">Create an account to start monitoring your sites.</p><form onsubmit=\"signUp(); return false;\"><div class=\"mb-2\"><input class=\"form-control\" id=\"signup-name\" type=\"text\" placeholder=\"Name\"></div><div class=\"mb-2\"><input class=\"form-control\" id=\"signup-email\" type=\"email\" placeholder=\"Email address\"></div><div class=\"mb-2\"><input class=\"form-control\" id=\"signup-password\" type=\"password\" placeholder=\"Create Password\"></div><button class=\"btn w-100 mb-8 btn-primary shadow\" type=\"submit\">Create Account</button><p class=\"d-flex flex-wrap align-items-center justify-content-center\"><span class=\"me-1\">Already have an account?</span> <a class=\"btn px-0 btn-link fw-bold\" href=\"/login\">Login</a></p></form></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = partials.ImportFirebaseScripts().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = partials.FirebaseConfig().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = partials.InitializeFirebaseApp().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = partials.SignupFunction().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = partials.RefreshAndVerifyTokenFunction().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = partials.VerifyTokenFunction().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
