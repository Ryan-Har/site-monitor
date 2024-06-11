// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.707
package templates

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import ()

func ResetPassword() templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<div class=\"container d-flex flex-column align-items-center justify-content-center h-50 mt-5\"><div class=\"col-7 mx-auto text-center\"><h3 class=\"mb-4\">Forgotten Password?</h3><p class=\"fs-8 text-secondary mb-12\">Please enter your email to receive a forgotten password link</p><form id=\"passwordResetForm\"><div class=\"mb-2\"><label for=\"email\" class=\"form-label\"></label> <input class=\"form-control\" id=\"email\" name=\"email\" type=\"email\" placeholder=\"Email address\"></div><p id=\"passwordResetMessage\"></p><button class=\"btn w-100 mb-8 btn-primary shadow\">Send Reset link</button></form></div></div><script>\n  const resetForm = document.getElementById('passwordResetForm');\n  resetForm.addEventListener('submit', (e) => {\n    e.preventDefault();\n    const email = document.getElementById('email').value;\n    const messageElement = document.getElementById('passwordResetMessage');\n\n    firebase.auth().sendPasswordResetEmail(email).then(() => {\n      messageElement.textContent = 'An email has been sent if an account with the provided email exists, please check your inbox';\n    }).catch((error) => {\n      messageElement.textContent = 'An error occured, please try again later';\n      console.log('Error: ' + error.message);\n    });\n  });\n</script>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
