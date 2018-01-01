package v1

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/dgrijalva/jwt-go"
	"github.com/iballbar/beegoAPI"
	"strconv"
)

type ParamsController struct {
	beegoAPI.API
}

const (

	// ----- jwt -----
	SecretKey_JWT        = "a2Mmh1LuclBvp"
	JWT_NEW_ASSIGN_VALUE = "token"

	// ----- parameter -----

	ACCESS_TOKEN     string = "access_token"
	USERNAME         string = "username"
	PASSWORD         string = "password"
	PASSWORD_CONFIRM string = "password_confirm"
	FIRSTNAME        string = "first_name"
	LASTNAME         string = "last_name"
	FACEBOOK_ID      string = "facebook_id"
	FACEBOOK_TOKEN   string = "facebook_token"
	TOKEN            string = "token"
	NONCE            string = "nonce"
	TIMESTAMP        string = "timestamp"
	RESET_TOKEN      string = "reset_token"
	LANGUAGE         string = "language"
	MAX              string = "max"
	OFFSET           string = "offset"
	JSON_PARAMS      string = "json_params"

	//keyjson

)

type DataParameter struct {
	Data ValueParam `json:"data"`
	jwt.StandardClaims
}

type ValueParam struct {
	Username        string     `json:"username"`
	Password        string     `json:"password"`
	PasswordConfirm string     `json:"password_confirm"`
	FirstName       string     `json:"first_name"`
	LastName        string     `json:"last_name"`
	FullName        string     `json:"full_name"`
	Birthdate       string     `json:"birthdate"`
	TitleName       string     `json:"title_name"` // mr , ms mrs
	MobilePhone     string     `json:"mobile_phone"`
	LineId          string     `json:"line_id"`
	FacebookId      string     `json:"facebook_id"`
	FacebookToken   string     `json:"facebook_token"`
	Token           string     `json:"token"`
	Nonce           string     `json:"nonce"`
	TimeStamp       int64      `json:"timestamp"`
	AccessToken     string     `json:"access_token"`
	ResetToken      string     `json:"reset_token"`
	LANGUAGE        string     `json:"language"`
	Max             int64      `json:"max"`
	Offset          int64      `json:"offset"`
	JSON_PARAMS     JsonParams `json:"json_params"`
}

type JsonParams struct {
	Image string `json:"image"`
}

//jwt
func (this *ParamsController) GlobalParamsJWT() ValueParam {

	claims, ok := this.Ctx.Input.GetData(JWT_NEW_ASSIGN_VALUE).(DataParameter)

	if !ok {
		beego.Debug(this.Ctx.Input.GetData(JWT_NEW_ASSIGN_VALUE))
		beego.Error(claims)
		beego.Error("GlobalParamsJWT")
	}

	return claims.Data

}

func (this *ParamsController) GlobalParamsNormal() ValueParam {

	jsonStr := this.GetString(JSON_PARAMS, "")
	var jsonParams JsonParams

	if jsonStr != "" {
		err := json.Unmarshal([]byte(jsonStr), &jsonParams)
		if err != nil {
			beego.Error(err)
		}
	}

	result := ValueParam{
		Username:        this.GetString(USERNAME),
		Password:        this.GetString(PASSWORD),
		PasswordConfirm: this.GetString(PASSWORD_CONFIRM),
		FirstName:       this.GetString(FIRSTNAME),
		LastName:        this.GetString(LASTNAME),
		FacebookId:      this.GetString(FACEBOOK_ID),
		FacebookToken:   this.GetString(FACEBOOK_TOKEN),
		Nonce:           this.GetString(NONCE),
		Token:           this.GetString(TOKEN), //verifymail
		TimeStamp:       this.GetInt64Req(TIMESTAMP, 0),
		Max:             this.GetInt64Req(MAX, -1),
		Offset:          this.GetInt64Req(OFFSET, 0),
		JSON_PARAMS:     jsonParams,
	}

	return result
}

func (this *ParamsController) GetInt64Req(key string, defaultValue int64) int64 {
	value, _ := this.GetInt64(key, defaultValue)
	return value

}

func (this *ParamsController) GetLanguage() string {

	return this.GetString(LANGUAGE, "th")

}

func Float64toString(amountFloat float64) string {
	return strconv.FormatFloat(amountFloat, 'f', 2, 64)
}

func StringToFloat64(amountString string) float64 {
	f, _ := strconv.ParseFloat(amountString, 64)
	return f
}

func ToString(value interface{}) string {
	var resultStr = ""

	if value != nil {
		resultStr, _ = value.(string)
	}

	return resultStr
}

func FloatToInt64(value interface{}) int64 {
	var resultInt int64 = 0

	if value != nil {
		floatValue, ok := value.(float64)
		if ok {
			resultInt = int64(floatValue)
		}
	}

	return resultInt
}

func Int64ToString(a int64) string {
	return strconv.FormatInt(a, 10)
}

func IntToString(a int) string {

	t := strconv.Itoa(a)
	return t
}

func StringToInt(s string) int {

	i, err := strconv.Atoi(s)
	if err != nil {
		return 0
	}

	return i

}

func StringToInt64(s string) int64 {

	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0
	}

	return i

}

func IsJwtTokenValid(data string) (*jwt.Token, bool) {
	beego.Debug("data = " + data)
	token, err := jwt.ParseWithClaims(data, &DataParameter{}, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			beego.Debug("Unexpected signing method: ", token.Header["alg"])
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		//beego.Debug("base64.StdEncoding.EncodeToString([]byte(SecretKey)): " ,base64.StdEncoding.EncodeToString([]byte(SecretKey)))
		//beego.Debug("jwt.EncodeSegment([]byte(SecretKey)): " ,base64.StdEncoding.EncodeToString([]byte(SecretKey)))
		//return []byte(jwt.EncodeSegment([]byte(SecretKey))), nil
		return []byte(SecretKey_JWT), nil
	})

	beego.Debug(err)

	if err == nil && token.Valid {
		return token, true
	} else {
		return nil, false
	}
}

func GetJsonData(ctx *context.Context) (bool, ValueParam) {
	claims, ok := ctx.Input.GetData(JWT_NEW_ASSIGN_VALUE).(DataParameter)
	return ok, claims.Data
}
