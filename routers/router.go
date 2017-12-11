package routers

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"gitlab.com/wisdomvast/NayooServer/api/v1"
	"gitlab.com/wisdomvast/NayooServer/controllers"
	"gitlab.com/wisdomvast/NayooServer/models"
	"strconv"
	"time"
)

func init() {
	beego.Router("/", &controllers.MainController{})

	nsAPI := beego.NewNamespace("/v1",
		beego.NSBefore(FilterDebug),

		beego.NSNamespace("/api",
			//beego.NSBefore(FilterNonce),
			beego.NSNamespace("/user",
				beego.NSRouter("/register", &v1.UserController{}, "get,post:Register"),
				beego.NSRouter("/loginbyemail", &v1.UserController{}, "get,post:Authenticate"),
				beego.NSRouter("/loginbyfacebook", &v1.UserController{}, "get,post:LoginByFacebook"),
				beego.NSRouter("/verify", &v1.UserController{}, "get:VerifyEmailUser"),
				beego.NSRouter("/forgotpassword", &v1.UserController{}, "post:ForgotPassword"),
				beego.NSRouter("/resetpassword", &v1.UserController{}, "post:ResetPassword"),
				beego.NSRouter("/getUserProfile", &v1.UserController{}, "get,post:GetUserProfile"),
				beego.NSRouter("/updateUserProfile", &v1.UserController{}, "get,post:UpdateUserProfile"),
			),
		),
	)

	beego.AddNamespace(nsAPI)

	nsDebug := beego.NewNamespace("/debug",
		beego.NSRouter("/", &v1.DebugController{}, "get,post:Main"),
	)

	beego.AddNamespace(nsDebug)

}

var FilterNonce = func(ctx *context.Context) {

	nonce := v1.ToString(ctx.Request.FormValue("nonce"))
	timestamp, _ := strconv.ParseInt(ctx.Request.FormValue("timestamp"), 10, 64)

	w := ctx.ResponseWriter
	if nonce == "" || timestamp == 0 {
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(400)

		results := map[string]interface{}{
			"statusCode":  400,
			"status":      "error",
			"description": "bad request",
			"data":        make(map[string]interface{}, 0),
		}

		response, _ := json.Marshal(results)

		w.Write([]byte(response))
		return
	}

	isUsed := models.IsUsedNonce(nonce)

	if isUsed {

		beego.Error("Bad FilterNonce used")
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(429)
		results := map[string]interface{}{
			"statusCode":  429,
			"status":      "error",
			"description": "duplicate request",
			"data":        make(map[string]interface{}, 0),
		}
		response, _ := json.Marshal(results)
		w.Write([]byte(response))

	} else {

		var expiryTime int64 = 60 * 60 * 12
		now := time.Now().Unix()
		expires := timestamp + expiryTime
		if expires-now < 0 { // check expire
			//expired
			beego.Error("Bad FilterNonce timed out")
			w.Header().Set("Content-Type", "application/json;charset=UTF-8")
			w.WriteHeader(408)

			results := map[string]interface{}{
				"statusCode":  408,
				"status":      "error",
				"description": "request timeout",
				"data":        make(map[string]interface{}, 0),
			}
			response, _ := json.Marshal(results)

			w.Write([]byte(response))
		}
	}

}

var FilterDebug = func(ctx *context.Context) {

	beego.Debug("body:", string(ctx.Input.RequestBody))
	beego.Debug("params:", ctx.Input.Params())
	beego.Debug("form:", ctx.Request.Form)
	beego.Debug("postform:", ctx.Request.PostForm)
	beego.Debug("RequestURI", ctx.Request.RequestURI)
	for name, value := range ctx.Request.Header {
		beego.Debug(name, ":", value)
	}
	return
}
