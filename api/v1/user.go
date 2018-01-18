package v1

import (
	"github.com/astaxie/beego"
	"github.com/dgrijalva/jwt-go"
	fb "github.com/huandu/facebook"
	"gitlab.com/wisdomvast/NayooServer/models"
	"golang.org/x/crypto/bcrypt"
	"time"
)

type UserController struct {
	GlobalApi
}

const (
	FACEBOOK_APP_ID     = "136599500258818"
	FACEBOOK_APP_SECRET = "a3c8aed413483fbed4bb7a58f6f9df34"
)

var globalApp = fb.New(FACEBOOK_APP_ID, FACEBOOK_APP_SECRET)

func (this *UserController) Register() {

	params := this.GlobalParamsJWT()

	fullName := params.FullName
	username := params.Username // email
	TitleName := params.TitleName
	MobilePhone := params.MobilePhone
	LineID := params.LineId
	birthDate := params.Birthdate
	password := params.Password
	password2 := params.PasswordConfirm
	nonce := params.Nonce
	timestamp := params.TimeStamp

	if fullName == "" || username == "" || password == "" || password2 == "" || TitleName == "" || MobilePhone == "" || LineID == "" || birthDate == "" {
		beego.Error("error")
		this.ResponseJSON(nil, 400, GetStringByLanguage(BAD_REQUEST_TH, BAD_REQUEST_TH, BAD_REQUEST_ENG , params))
		return
	}
	if password != password2 {
		beego.Error("error")
		this.ResponseJSON(nil, 400, GetStringByLanguage(BAD_REQUEST_TH, BAD_REQUEST_TH, BAD_REQUEST_ENG , params))
		return
	}

	if IsUsernameAvailable(username) == false {
		beego.Error("error")
		this.ResponseJSON(nil, 400, GetStringByLanguage(DUPLICATE_USER_TH, DUPLICATE_USER_TH, DUPLICATE_USER_ENG , params))
		return
	}
	if saveUser(username, password, fullName) == false {
		beego.Error("error")
		this.ResponseJSON(nil, 400, GetStringByLanguage(SERVER_ERROR_TH, SERVER_ERROR_TH, SERVER_ERROR_ENG , params))
		return
	}

	registerCode := newRegistrationCode(username, "register")

	if registerCode == nil {
		beego.Error("error")
		this.ResponseJSON(nil, 400, GetStringByLanguage(SERVER_ERROR_TH, SERVER_ERROR_TH, SERVER_ERROR_ENG,params))
		return
	}

	go SendRegisterMail(username, registerCode.Token, "register")
	addUsedNonce(nonce, timestamp)
	this.ResponseJSON(nil, 200, GetStringByLanguage(SUCCESS_TH, SUCCESS_TH, SUCCESS_ENG,params))
	return
}

func (this *UserController) LoginByFacebook() {

	params := this.GlobalParamsJWT()

	facebookId := params.FacebookId
	facebookToken := params.FacebookToken
	nonce := params.Nonce
	timestamp := params.TimeStamp
	if facebookId == "" {
		this.ResponseJSON(nil, 400, "Bad Request")
		return
	}
	user := loginByFacebook(facebookId, facebookToken)

	if user != nil {
		authToken := NewAuthenticationToken(user)
		if authToken != nil {
			addUsedNonce(nonce, timestamp)
			this.ResponseJSON([]map[string]interface{}{
				map[string]interface{}{
					"accessToken": authToken.AccessToken,
				},
			}, 200, GetStringByLanguage(SUCCESS_TH, SUCCESS_TH, SUCCESS_ENG,params))
			return
		} else {
			this.ResponseJSON(nil, 400, GetStringByLanguage(BAD_REQUEST_TH, BAD_REQUEST_TH, BAD_REQUEST_ENG , params))
			return
		}
	} else if user != nil && !user.Verified {
		this.ResponseJSON(nil, 400, GetStringByLanguage(BAD_REQUEST_TH, BAD_REQUEST_TH, BAD_REQUEST_ENG , params))
		return
	} else {
		this.ResponseJSON(nil, 400, GetStringByLanguage(BAD_REQUEST_TH, BAD_REQUEST_TH, BAD_REQUEST_ENG , params))
		return
	}
}

func (this *UserController) Authenticate() {

	params := this.GlobalParamsJWT()

	username := params.Username
	password := params.Password
	nonce := params.Nonce
	timestamp := params.TimeStamp

	if username == "" || password == "" {
		this.ResponseJSON(nil, 400, GetStringByLanguage(BAD_REQUEST_TH,BAD_REQUEST_TH,BAD_REQUEST_ENG , params))
		return
	}
	//beego.Debug(username, password)
	user := authenticate(username, password)
	beego.Debug("user:", user)
	if user != nil && user.Verified {
		authToken := NewAuthenticationToken(user)
		if authToken != nil {
			beego.Debug(authToken.AccessToken)
			addUsedNonce(nonce, timestamp)

			this.ResponseJSON(map[string]interface{}{
				"accessToken": authToken.AccessToken,
			}, 200, "success")
			return

		} else {
			this.ResponseJSON(nil, 500, GetStringByLanguage(SYSTEM_ERROR_TH, SYSTEM_ERROR_TH, SYSTEM_ERROR_ENG , params))
			return
		}
	} else if user != nil && !user.Verified {
		this.ResponseJSON(nil, 401, GetStringByLanguage(USER_NOT_VERIFY_TH, USER_NOT_VERIFY_TH, USER_NOT_VERIFY_ENG , params))
		return
	} else {
		this.ResponseJSON(nil, 401, GetStringByLanguage(LOGIN_FAIL_TH, LOGIN_FAIL_TH, LOGIN_FAIL_ENG , params))
		return
	}
}

func (this *UserController) VerifyEmailUser() {

	params := this.GlobalParamsJWT()
	token := params.Token
	registerCode := models.GetRegisterCodeByToken(token, "register")
	if registerCode != nil {
		user := models.GetUserByUsername(registerCode.Username)
		if user != nil {
			user.Verified = true
			if err := models.UpdateUserById(user); err == nil {
				// add free trial 15 days
				// TODO remove after 1 Oct 2016
				//v1.AddFreeSubscription(user, 15)
				if errDel := models.DeleteRegistrationCode(registerCode.Id); errDel != nil {
					beego.Error("Cannot delete ReisterCode("+user.Username+"): ", errDel.Error())
				}

			} else {
				beego.Error("try again")
				this.ResponseJSON(nil, 400, GetStringByLanguage(TRY_AGAIN_TH, TRY_AGAIN_TH, TRY_AGAIN_ENG , params))
				return
			}

		} else {
			beego.Error("user not found")
			this.ResponseJSON(nil, 400, GetStringByLanguage(USER_NOT_FOUND_TH, USER_NOT_FOUND_TH, USER_NOT_FOUND_ENG , params))
			return
		}
	} else {
		beego.Error("register code not found")
		this.ResponseJSON(nil, 400, GetStringByLanguage(TRY_AGAIN_TH, TRY_AGAIN_TH, TRY_AGAIN_ENG , params))
		return
	}
	this.ResponseJSON(nil, 200, GetStringByLanguage(WELCOME_TH, WELCOME_TH, WELCOME_ENG , params))
	return
}

func (this *UserController) ForgotPassword() {

	params := this.GlobalParamsJWT()
	nonce := params.Nonce
	timestamp := params.TimeStamp
	username := params.Username

	userObj := models.GetUserByUsername(username)

	if username == "" || userObj == nil {
		this.ResponseJSON(nil, 400, GetStringByLanguage(USER_NOT_FOUND_TH, USER_NOT_FOUND_TH, USER_NOT_FOUND_ENG , params))
		return
	}

	registerCode := newRegistrationCode(userObj.Username, "resetpassword")
	if registerCode == nil {
		this.ResponseJSON(nil, 400, GetStringByLanguage(ERROR_MESSAGE_TH, ERROR_MESSAGE_TH, ERROR_MESSAGE_ENG , params))
		return
	}
	go SendRegisterMail(username, registerCode.Token, "resetpassword")
	addUsedNonce(nonce, timestamp)

	this.ResponseJSON(nil, 200, GetStringByLanguage(RESET_PASSWORD_MESSAGE_TH, RESET_PASSWORD_MESSAGE_TH, RESET_PASSWORD_MESSAGE_ENG , params))
	return
}

func (this *UserController) ResetPassword() {

	params := this.GlobalParamsJWT()

	nonce := params.Nonce
	timestamp := params.TimeStamp
	resetToken := params.ResetToken
	password := params.Password
	confirmPassword := params.PasswordConfirm

	if resetToken == "" || password == "" || confirmPassword == "" {
		this.ResponseJSON(nil, 400, GetStringByLanguage(ERROR_MESSAGE_TH, ERROR_MESSAGE_TH, ERROR_MESSAGE_ENG , params))
		return
	}

	addUsedNonce(nonce, timestamp)

	if password != confirmPassword {
		this.ResponseJSON(nil, 400, GetStringByLanguage(PASSWORD_MISMATCH_TH, PASSWORD_MISMATCH_TH, PASSWORD_MISMATCH_ENG , params))
		return
	}

	registerCode := models.GetRegisterCodeByToken(resetToken, "resetpassword")
	if registerCode == nil {
		this.ResponseJSON(nil, 400, GetStringByLanguage(ERROR_MESSAGE_TH, ERROR_MESSAGE_TH, ERROR_MESSAGE_ENG , params))
		return
	}

	user := models.GetUserByUsername(registerCode.Username)
	if user == nil {
		this.ResponseJSON(nil, 400, GetStringByLanguage(ERROR_MESSAGE_TH, ERROR_MESSAGE_TH, ERROR_MESSAGE_ENG , params))
		return
	}

	if hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost); err == nil {
		user.Password = string(hashedPassword)
		user.Verified = true
	} else {
		beego.Error(err)
	}

	if err := models.UpdateUserById(user); err == nil {
		this.Data["json"] = "SUCCESS"
		if errDel := models.DeleteRegistrationCode(registerCode.Id); errDel != nil {
			beego.Error("Cannot delete ReisterCode("+user.Username+"): ", errDel.Error())
		}
		this.ResponseJSON(nil, 200, GetStringByLanguage(SUCCESS_TH, SUCCESS_ENG, SUCCESS_ENG , params))
		return

	} else {
		beego.Error("Cannot update user: ", err.Error())
		this.ResponseJSON(nil, 400, GetStringByLanguage(ERROR_MESSAGE_TH, ERROR_MESSAGE_TH, ERROR_MESSAGE_ENG , params))
		return
	}

}

func (this *UserController) GetUserProfile() {

	params := this.GlobalParamsJWT()

	accessToken := params.AccessToken
	user := GetUserByToken(accessToken)
	if user != nil {
		this.ResponseJSON(this.GenerateUserDetailJson(user), 200, "success")
		return
	} else {
		this.ResponseJSON(nil, 401, "UNAUTHORIZED")
		return
	}

}

func (this *UserController) UpdateUserProfile() {

	//if errDel := models.DeleteRegistrationCode(registerCode.Id); errDel != nil {
	//	beego.Error("Cannot delete ReisterCode("+user.Username+"): ", errDel.Error())
	//}

	accessToken := this.GetString("accessToken")
	user := GetUserByToken(accessToken)
	if user != nil {

		file, handler, err := this.GetFile("image")
		checkUpload := true
		reasonUpload := ""

		if file != nil {
			checkUpload, reasonUpload = this.UploadProfileImageGlobal(0, file, handler, err, "userProfile")

			if checkUpload == false {
				this.ResponseJSON(nil, 400, "BAD UPLOAD")
				return
			}
		} else {
			reasonUpload = user.Image
		}

		user.Image = reasonUpload
		fullname := ""

		if firstname := this.GetString("firstname"); firstname != "" {
			user.FirstName = firstname
			fullname += firstname + " "
		}
		if lastname := this.GetString("lastname"); lastname != "" {
			user.LastName = lastname
			fullname += "" + lastname
		}
		user.FullName = fullname
		if career := this.GetString("career"); career != "" {
			user.Career = career
		}
		if mobileNumber := this.GetString("mobileNumber"); mobileNumber != "" {
			user.MobileNumber = mobileNumber
		}
		if phoneNumber := this.GetString("phoneNumber"); phoneNumber != "" {
			user.PhoneNumber = phoneNumber
		}
		if address := this.GetString("address"); address != "" {
			user.Address = address
		}
		if province := this.GetString("province"); province != "" {
			user.Province = province
		}
		if distinct := this.GetString("distinct"); distinct != "" {
			user.District = distinct
		}
		if subdistinct := this.GetString("subdistinct"); subdistinct != "" {
			user.SubDistrict = subdistinct
		}

		models.UpdateUserById(user)

		this.ResponseJSON(this.GenerateUserDetailJson(user), 200, "success")
		return

	} else {
		this.ResponseJSON(nil, 401, "UNAUTHORIZED")
		return
	}

}

//=========================================================================================

func loginByFacebook(facebookId, token string) *models.User {

	sessionClient := globalApp.Session(token)
	sessionClient.EnableAppsecretProof(true)
	beego.Debug(sessionClient)
	beego.Debug(facebookId)

	check_token, err := sessionClient.Inspect()
	if err != nil {
		beego.Error(err)
		return nil
	}
	user_id := ""
	_ = check_token.DecodeField("user_id", &user_id)
	is_valid := false
	_ = check_token.DecodeField("is_valid", &is_valid)
	app_id := ""
	_ = check_token.DecodeField("app_id", &app_id)
	if !is_valid || user_id == "" || app_id != globalApp.AppId {
		beego.Error("Invalid token")
		return nil
	}

	if facebookId != user_id {
		beego.Error("Invalid facebookId")
		return nil
	}

	getPath := ("/" + user_id)
	res, err := sessionClient.Get(getPath, fb.Params{"fields": "first_name,last_name,name,email,gender,birthday"})
	if err != nil {
		beego.Error(err)
	}

	email := ""
	err = res.DecodeField("email", &email)
	if err != nil {
		beego.Error(err)
	}
	beego.Debug(email)
	firstName := ""
	_ = res.DecodeField("first_name", &firstName)
	lastName := ""
	_ = res.DecodeField("last_name", &lastName)
	displayName := ""
	_ = res.DecodeField("name", &displayName)
	gender := ""
	_ = res.DecodeField("gender", &gender)
	birthday := ""
	_ = res.DecodeField("birthday", &birthday)

	user := saveFacebookUser(facebookId, email, firstName, lastName, displayName, gender, birthday)
	return user
}

func saveFacebookUser(facebookId, email, firstName, lastName, displayName, gender, birthday string) (result *models.User) {

	if hashedPassword, err := bcrypt.GenerateFromPassword([]byte(facebookId), bcrypt.DefaultCost); err == nil {

		birthDate, _ := time.Parse("01/02/2006", birthday)
		if result = models.GetUserByFacebookId(facebookId); result != nil {
			result.FirstName = firstName
			result.LastName = lastName
			result.FullName = firstName + lastName
			result.DisplayName = displayName
			result.Gender = gender
			result.BirthDate = birthDate
			err = models.UpdateUserById(result)
			if err != nil {
				beego.Error(err)
			}
		} else if result = models.GetUserByUsername(email); result != nil {
			result.FacebookId = facebookId
			result.FirstName = firstName
			result.LastName = lastName
			result.FullName = firstName + lastName
			result.DisplayName = displayName
			result.Gender = gender
			result.BirthDate = birthDate
			err = models.UpdateUserById(result)
			if err != nil {
				beego.Error(err)
			}
		} else {
			username := email
			if email == "" {
				username = facebookId
			}
			result = &models.User{
				Username:     username,
				FacebookId:   facebookId,
				Password:     string(hashedPassword),
				Email:        email,
				FirstName:    firstName,
				LastName:     lastName,
				FullName:     firstName + lastName,
				DisplayName:  displayName,
				Gender:       gender,
				BirthDate:    birthDate,
				RegisterDate: time.Now(),
			}

			id, err := models.AddUser(result)

			if err != nil {
				beego.Error(err)
			} else {
				result.Id = id
				// add free trial 15 days
				// TODO remove after 1 Oct 2016
				//AddFreeSubscription(&result, 15)
			}

		}

	}

	return result
}

func NewAuthenticationToken(user *models.User) *models.AuthenticationToken {

	auth := models.GetAuthTokenByUserId(user.Id)
	if auth != nil {
		return auth
	} else {
		accessToken := NewAccessToken()
		if accessToken == nil {
			return nil
		}

		auth = &models.AuthenticationToken{
			AccessToken: *accessToken,
			User:        user,
		}

		_, err := models.AddAuthenticationToken(auth)
		if err != nil {
			beego.Error("DB Error cannot add new AuthenticationToken | ", err.Error())
			return nil
		}

	}

	return auth

}

func NewAccessToken() *string {
	//token := jwt.New(jwt.SigningMethodHS512)
	// Set some claims
	// Create the Claims
	claims := &jwt.StandardClaims{
		ExpiresAt: time.Now().Add(time.Hour * 72).Unix(),
		Issuer:    "nayoo.com",
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	//token.Claims["time"] = time.Now().Unix()
	//token.Claims["exp"] = time.Now().Add(time.Hour * 72).Unix()
	// Sign and get the complete encoded token as a string
	tokenString, err := token.SignedString([]byte(AccessKey))
	if err != nil {
		beego.Error("Error cannot gen new token | ", err.Error())
		return nil
	}
	return &tokenString
}

func IsUsernameAvailable(username string) bool {
	return models.IsUsernameAvailable(username)
}

func saveUser(username, password, fullName string) bool {
	result := true
	if hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost); err == nil {
		if facebookUser := models.GetFacebookUserByEmail(username); facebookUser != nil {
			facebookUser.Username = username
			facebookUser.Password = string(hashedPassword)
			facebookUser.RegisterDate = time.Now()
			err := models.UpdateUserById(facebookUser)
			if err != nil {
				result = false
			}
		} else {
			user := models.User{
				Username:     username,
				Password:     string(hashedPassword),
				FullName:     fullName,
				Email:        username,
				RegisterDate: time.Now(),
			}

			id, err := models.AddUser(&user)
			if err != nil {
				result = false
			} else {
				user.Id = id
			}
		}

	} else {
		result = false
	}

	return result
}

func newRegistrationCode(username, typeStr string) *models.RegistrationCode {
	if token, err := bcrypt.GenerateFromPassword([]byte(username+typeStr), bcrypt.DefaultCost); err == nil {

		registerCode := models.GetRegisterCodeByUsername(username, typeStr)
		if registerCode != nil {
			registerCode.Token = string(token)

			err := models.UpdateRegistrationCodeById(registerCode)
			if err != nil {
				return nil
			}
		} else {
			registerCode = &models.RegistrationCode{
				Username: username,
				Token:    string(token),
				Type:     typeStr,
			}
			_, err := models.AddRegistrationCode(registerCode)
			if err != nil {
				return nil
			}
		}

		return registerCode
	}

	return nil

}

func authenticate(username, password string) *models.User {

	user := models.GetUserByUsername(username)

	if user == nil {
		return nil
	}

	// Comparing the password with the hash
	success := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if user == nil || success != nil {
		return nil
	}

	return user
}
