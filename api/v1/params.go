package v1

import (
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/iballbar/beegoAPI"
	"strconv"
)

type ParamsController struct {
	beegoAPI.API
}

const (

	// ----- parameter -----

	USERNAME         string = "username"
	PASSWORD         string = "password"
	PASSWORD_CONFIRM string = "password_confirm"
	FIRSTNAME        string = "firstname"
	LASTNAME         string = "lastname"
	FACEBOOK_ID      string = "facebook_id"
	FACEBOOK_TOKEN   string = "facebook_token"
	NONCE            string = "nonce"
	TOKEN            string = "token"
	TIMESTAMP        string = "timestamp"
	JSON_PARAMS      string = "json_params"
	LANGUAGE         string = "language"
	MAX              string = "max"
	OFFSET           string = "offset"

	//keyjson

)

type ValueParam struct {
	Username        string
	Password        string
	PasswordConfirm string
	FirstName       string
	LastName        string
	FacebookId      string
	FacebookToken   string
	Token           string
	Nonce           string
	ResetToken      string
	TimeStamp       int64
	Max             int64
	Offset          int64
	SimpleJson      SimpleJson `json:"simple_json"`
}

type SimpleJson struct {
	Image string `json:"image"`
}

func (this *ParamsController) GlobalParams() ValueParam {

	jsonStr := this.GetString(JSON_PARAMS, "")
	var simpleJson SimpleJson

	if jsonStr != "" {
		err := json.Unmarshal([]byte(jsonStr), &simpleJson)
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
		SimpleJson:      simpleJson,
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
