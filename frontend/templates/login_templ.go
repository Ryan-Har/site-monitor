// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.696
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

func Login() templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"container d-flex flex-column align-items-center justify-content-center h-50 w-50 mt-5\"><div class=\"mw-md mx-auto text-center\"><a class=\"d-inline-block mb-32\" href=\"#\"><img src=\"flaro-assets/logos/flaro-logo-black-xl.svg\" alt=\"\"></a><h3 class=\"mb-4\">Sign In</h3><p class=\"fs-8 text-secondary mb-12\">Welcome back!</p><form action=\"\" data-bitwarden-watching=\"1\"><div class=\"mb-2\"><input class=\"form-control\" type=\"email\" placeholder=\"Email address\"></div><div class=\"mb-2 position-relative\"><input class=\"form-control\" type=\"password\" placeholder=\"Password\"> <a class=\"position-absolute top-50 end-0 me-4 translate-middle-y btn p-0 btn-link fs-9\" href=\"#\">Forgot Password?</a></div><a class=\"btn w-100 mb-8 btn-primary shadow\" href=\"#\">Sign In</a><p class=\"d-flex flex-wrap align-items-center justify-content-center\"><span class=\"me-1\">Don’t have an account?</span> <a class=\"btn px-0 btn-link fw-bold\" href=\"#\">Create free account</a></p></form></div></div>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
