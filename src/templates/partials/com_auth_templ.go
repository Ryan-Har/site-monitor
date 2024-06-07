// Code generated by templ - DO NOT EDIT.

// templ: version: v0.2.707
package partials

//lint:file-ignore SA4006 This context is only used if a nested component is present.

import "github.com/a-h/templ"
import "context"
import "io"
import "bytes"

import (
	"github.com/Ryan-Har/site-monitor/src/config"
	// "encoding/json"
	"fmt"
	"html/template"
)

func ImportFirebaseScripts() templ.Component {
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
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<!-- Firebase App (the core Firebase SDK) --><script src=\"https://www.gstatic.com/firebasejs/8.10.0/firebase-app.js\"></script><!-- Firebase Authentication --><script src=\"https://www.gstatic.com/firebasejs/8.10.0/firebase-auth.js\"></script>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

// templ makes it difficult to embed information within the javascript variable so we make use of html templating
func firebaseConfig() *template.Template {
	tmpl := template.New("firebaseConfig")

	tmpl.Funcs(template.FuncMap{
		"rawJS": func(s string) template.JS {
			return template.JS(s)
		},
	})

	_, err := tmpl.Parse("<script> var firebaseConfig = {{. | rawJS }} </script>")
	if err != nil {
		fmt.Println("error parsing firebaseConfig script", err)
	}

	return tmpl
}

func FirebaseConfig() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var2 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var2 == nil {
			templ_7745c5c3_Var2 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = ImportFirebaseScripts().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		templ_7745c5c3_Err = templ.FromGoHTML(firebaseConfig(), config.GetConfig().FirebaseConfigAsJsonString()).Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func InitializeFirebaseApp() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var3 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var3 == nil {
			templ_7745c5c3_Var3 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = FirebaseConfig().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<script>\n    firebase.initializeApp(firebaseConfig);\n    </script>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func MonitorAuthState() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var4 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var4 == nil {
			templ_7745c5c3_Var4 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		templ_7745c5c3_Err = InitializeFirebaseApp().Render(ctx, templ_7745c5c3_Buffer)
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<script>\n  async function updateAuthCookie(idToken) {\n     try {\n      const response = await fetch('/updateauthcookie', {\n        method: 'POST',\n        headers: {\n          'Content-Type': 'application/json',\n          'Authorization': idToken\n        },\n        body: JSON.stringify({ idToken: idToken })\n      });\n      \n      const data = await response.json();\n      console.log('Success:', data);\n      return data;\n    } catch (error) {\n      console.error('Error:', error);\n      throw error;  // re-throw the error so that the calling function can handle it\n    }\n  };\n\n    async function expireAuthCookie() {\n     try {\n      const response = await fetch('/expireauthcookie', {\n        method: 'POST',\n        headers: {\n          'Content-Type': 'application/json',\n        },\n      });\n      \n      const data = await response.json();\n      console.log('Success:', data);\n      return data;\n    } catch (error) {\n      console.error('Error:', error);\n      throw error;  // re-throw the error so that the calling function can handle it\n    }\n  };\n\n  firebase.auth().onAuthStateChanged((user) => {\n  if (user) {\n    user.getIdToken(true).then((idToken) => {\n      updateAuthCookie(idToken).then((response) => {\n        console.log(response)\n        const path = window.location.pathname;\n        if (path.includes('login') || path.includes('signup')) {\n          window.location.href = '/monitors';\n        }\n      }) \n    }).catch((error) => {\n      console.error('Error refreshing token:', error);\n    });\n  } else {\n    console.log('No user is signed in');\n  }\n  });\n\n  function signUp() {\n      var name = document.getElementById('name').value;\n      var email = document.getElementById('email').value;\n      var password = document.getElementById('password').value;\n\n      firebase.auth().createUserWithEmailAndPassword(email, password)\n        .then((userCredential) => {\n          // Signed up successfully\n          var user = userCredential.user;\n          user.updateProfile({\n              displayName: name\n          }).then(() => {\n            console.log('User signed up: ', user);\n            login(email, password)\n            })\n          })\n        .catch((error) => {\n          var errorCode = error.code;\n          var errorMessage = error.message;\n          console.error('Error: ', errorCode, errorMessage);\n        });\n  };\n\n  function login(eml, psswd) {\n      var email = eml || document.getElementById('email').value;\n      var password = psswd || document.getElementById('password').value;\n\n    firebase.auth().signInWithEmailAndPassword(email, password)\n      .then((userCredential) => {\n        const user = userCredential.user;\n        console.log('logged in');\n      })\n      .catch((error) => {\n        const errorCode = error.code;\n        const errorMessage = error.message;\n        console.error('Error: ', errorCode, errorMessage);\n      });\n    };\n\n  function logout() {\n    firebase.auth().signOut().then(() => {\n    expireAuthCookie()\n    }, (error) => {\n    console.log(\"error occured logging out\")\n    });\n  }\n  </script>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func ChangeDisplayNameJS() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var5 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var5 == nil {
			templ_7745c5c3_Var5 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<script>\n  const changeDisplayNameForm = document.getElementById(\"changeDisplayNameForm\")\n  const fullNameInput = document.getElementById(\"fullName\")\n  const displayNameNotificationElement = document.getElementById(\"changeDisplayNameNotification\")\n\n  changeDisplayNameForm.addEventListener('submit', async (event) => {\n  event.preventDefault();\n\n  const newDisplayName = fullNameInput.value;\n\n  // Clear any previous notifications\n  displayNameNotificationElement.textContent = '';\n\n  try {\n    // Get the currently signed-in user\n    const user = firebase.auth().currentUser;\n\n    // Update the user's display name (without re-authentication)\n    const profileUpdates = { displayName: newDisplayName };\n    await user.updateProfile(profileUpdates);\n\n    displayNameNotificationElement.textContent = 'Display name changed successfully!';\n    fullNameInput.value = ''; // Clear display name field after success\n    user.getIdToken(true).then((idToken) => {\n      updateAuthCookie(idToken)});\n  } catch (error) {\n    console.error('Error changing display name:', error);\n    displayNameNotificationElement.textContent = 'An error occurred. Please try again.';\n  }\n});\n  </script>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}

func ChangePasswordJS() templ.Component {
	return templ.ComponentFunc(func(ctx context.Context, templ_7745c5c3_W io.Writer) (templ_7745c5c3_Err error) {
		templ_7745c5c3_Buffer, templ_7745c5c3_IsBuffer := templ_7745c5c3_W.(*bytes.Buffer)
		if !templ_7745c5c3_IsBuffer {
			templ_7745c5c3_Buffer = templ.GetBuffer()
			defer templ.ReleaseBuffer(templ_7745c5c3_Buffer)
		}
		ctx = templ.InitializeContext(ctx)
		templ_7745c5c3_Var6 := templ.GetChildren(ctx)
		if templ_7745c5c3_Var6 == nil {
			templ_7745c5c3_Var6 = templ.NopComponent
		}
		ctx = templ.ClearChildren(ctx)
		_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteString("<script>\n  const changePasswordForm = document.getElementById(\"changePasswordForm\")\n  const currentPasswordInput = document.getElementById(\"currentPassword\")\n  const newPasswordInput = document.getElementById(\"newPassword\")\n  const repeatNewPasswordInput = document.getElementById(\"repeatNewPassword\")\n  const passwordNotificationElement = document.getElementById(\"passwordNotification\")\n\n  changePasswordForm.addEventListener('submit', async (event) => {\n    event.preventDefault();\n\n    const currentPassword = currentPasswordInput.value;\n    const newPassword = newPasswordInput.value;\n    const repeatPassword = repeatNewPasswordInput.value;\n\n    // Clear any previous notifications\n    passwordNotificationElement.textContent = '';\n\n    // Validate passwords\n    if (newPassword !== repeatPassword) {\n      passwordNotificationElement.textContent = 'New passwords do not match.';\n      return;\n    } else if (newPassword.length < 8) {\n      passwordNotificationElement.textContent = 'New password must be at least 8 characters long.';\n      return;\n    }\n\n    try {\n      // Get the currently signed-in user\n      const user = firebase.auth().currentUser;\n\n      // Re-authenticate the user with their current password before updating\n      const credential = firebase.auth.EmailAuthProvider.credential(user.email, currentPassword);\n      await user.reauthenticateWithCredential(credential);\n\n      // Update the user's password with the new one\n      await user.updatePassword(newPassword);\n\n      passwordNotificationElement.textContent = 'Password changed successfully!';\n      currentPasswordInput.value = ''; // Clear password fields after successful change\n      newPasswordInput.value = '';\n      repeatNewPasswordInput.value = '';\n    } catch (error) {\n        if (error.code === 'auth/internal-error') {\n        passwordNotificationElement.textContent = 'Current password incorrect.';\n      } else {\n        console.error('Error changing password:', error);\n        passwordNotificationElement.textContent = 'An error occurred. Please try again.';\n      }\n    }\n  });\n  </script>")
		if templ_7745c5c3_Err != nil {
			return templ_7745c5c3_Err
		}
		if !templ_7745c5c3_IsBuffer {
			_, templ_7745c5c3_Err = templ_7745c5c3_Buffer.WriteTo(templ_7745c5c3_W)
		}
		return templ_7745c5c3_Err
	})
}
