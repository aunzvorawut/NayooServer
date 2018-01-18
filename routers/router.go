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

	beego.NSBefore(FilterDebug)

	nsAPI := beego.NewNamespace("/v1",
		beego.NSBefore(FilterJwt),

		beego.NSNamespace("/api",
			//beego.NSBefore(FilterNonce),
			beego.NSNamespace("/user",
				beego.NSRouter("/register", &v1.UserController{}, "get,post:Register"),
				beego.NSRouter("/loginbyemail", &v1.UserController{}, "get,post:Authenticate"),
				beego.NSRouter("/loginbyfacebook", &v1.UserController{}, "get,post:LoginByFacebook"),
				beego.NSRouter("/verify", &v1.UserController{}, "get,post:VerifyEmailUser"),
				beego.NSRouter("/forgotpassword", &v1.UserController{}, "post:ForgotPassword"),
				beego.NSRouter("/resetpassword", &v1.UserController{}, "post:ResetPassword"),
				beego.NSRouter("/getUserProfile", &v1.UserController{}, "get,post:GetUserProfile"),
				beego.NSRouter("/updateUserProfile", &v1.UserController{}, "get,post:UpdateUserProfile"),
			),

			beego.NSNamespace("/housesale",
				beego.NSRouter("/list", &v1.HousesaleController{}, "get,post:List"),
			),

			beego.NSNamespace("/houserent",
				beego.NSRouter("/list", &v1.HouserentController{}, "get,post:List"),
			),

			beego.NSNamespace("/houseproject",
				beego.NSRouter("/list", &v1.HouseProjectController{}, "get,post:List"),
				beego.NSRouter("/main", &v1.HouseProjectController{}, "get,post:Main"),
			),

			beego.NSNamespace("/ownproject",
				beego.NSRouter("/list", &v1.OwnProjectController{}, "get,post:List"),
			),

			beego.NSNamespace("/agent",
				beego.NSRouter("/list", &v1.AgentController{}, "get,post:List"),
			),

			beego.NSNamespace("/entrepreneur",
				beego.NSRouter("/list", &v1.EntrepreneurController{}, "get,post:List"),
			),
		),
	)

	beego.AddNamespace(nsAPI)

	nsDebug := beego.NewNamespace("/debug",
		beego.NSRouter("/", &v1.DebugController{}, "get,post:Main"),
	)

	beego.AddNamespace(nsDebug)

}

var FilterJwt = func(ctx *context.Context) {

	beego.Debug("body:", string(ctx.Input.RequestBody))
	beego.Debug("params:", ctx.Input.Params())
	beego.Debug("form:", ctx.Request.Form)
	beego.Debug("postform:", ctx.Request.PostForm)
	beego.Debug("RequestURI", ctx.Request.RequestURI)
	for name, value := range ctx.Request.Header {
		beego.Debug(name, ":", value)
	}

	data := ctx.Request.FormValue("data")
	token, valid := v1.IsJwtTokenValid(data)

	w := ctx.ResponseWriter

	beego.Debug(token)
	beego.Debug(valid)

	if token == nil && valid == false {

		beego.Error("error")

		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(400)

		results := map[string]interface{}{
			"code":           400,
			"message":        "error",
			"responseObject": make(map[string]interface{}, 0),
		}
		response, _ := json.Marshal(results)

		w.Write([]byte(response))
		return

	}

	claims, ok := token.Claims.(*v1.DataParameter)
	if !ok {
		valid = false
	}

	res2B, _ := json.Marshal(claims)
	beego.Debug(string(res2B))

	if valid {
		ctx.Input.SetData(v1.JWT_NEW_ASSIGN_VALUE, *claims)
		return
	} else {

		beego.Error("error")

		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(400)

		results := map[string]interface{}{
			"code":           400,
			"message":        "error",
			"responseObject": make(map[string]interface{}, 0),
		}
		response, _ := json.Marshal(results)

		w.Write([]byte(response))
		return

	}
}

var FilterNonce = func(ctx *context.Context) {

	nonce := v1.ToString(ctx.Request.FormValue("nonce"))
	timestamp, _ := strconv.ParseInt(ctx.Request.FormValue("timestamp"), 10, 64)

	w := ctx.ResponseWriter
	if nonce == "" || timestamp == 0 {
		beego.Error("error")
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(400)

		results := map[string]interface{}{
			"code":           400,
			"message":        "error",
			"responseObject": make(map[string]interface{}, 0),
		}

		response, _ := json.Marshal(results)

		w.Write([]byte(response))
		return
	}

	isUsed := models.IsUsedNonce(nonce)

	if isUsed {
		beego.Error("error")
		beego.Error("Bad FilterNonce used")
		w.Header().Set("Content-Type", "application/json;charset=UTF-8")
		w.WriteHeader(400)
		results := map[string]interface{}{
			"code":           400,
			"message":        "error",
			"responseObject": make(map[string]interface{}, 0),
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
			beego.Error("error")
			w.Header().Set("Content-Type", "application/json;charset=UTF-8")
			w.WriteHeader(400)

			results := map[string]interface{}{
				"code":           400,
				"message":        "error",
				"responseObject": make(map[string]interface{}, 0),
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
