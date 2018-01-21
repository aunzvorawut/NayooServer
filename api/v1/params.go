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

	//-------------------  keyjson  ------------------

	ID string = "id"

	TITLE                   string = "title"
	DESCRIPTION             string = "description"
	PRICE                   string = "price"
	IMAGE                   string = "image"
	IMAGES                  string = "images"
	VIP_TYPE                string = "vip_type"
	SALE_STR                string = "sale_str"
	SALE_STR_MARK           string = "sale_str_mark"
	RENT_STR_MONTH          string = "rent_str_month"
	RENT_STR_DAY            string = "rent_str_day"
	RENT_STR_MONTH_MARK     string = "rent_str_month_mark"
	RENT_STR_DAY_MARK       string = "rent_str_day_mark"
	RESIDENT_TYPE_LIST      string = "resident_type_list"
	RESIDENT_ADDRESS        string = "resident_address"
	TAG_TYPE_LIST           string = "tag_type_list"
	TAG_TYPE_STR            string = "tag_type_str"
	COUNT_BEDROOM           string = "count_bedroom"
	COUNT_BATHROOM          string = "count_bath_room"
	AREA_LAND               string = "area_land"
	AREA_USEFUL             string = "area_useful"
	IS_PROMOTON_NAYOO       string = "is_promotion_nayuu"
	IS_GURU                 string = "is_guru"
	IS_VIDEO_360            string = "is_video_360"
	IS_FAVORITE             string = "is_favorite"
	IS_VERIFY_PERSON_IDCARD string = "is_verify_person_idcard"
	IS_VERIFY_PERSON_POLICY string = "is_verify_person_policy"
	IS_VERIFY_COMPANY       string = "is_verify_company"
	IS_VERIFY_APARTMENT     string = "is_verify_apartment"
	LIST_RECOMMEND_VIEW     string = "list_recommend_view"
	LIST_VIDEO_VIEW         string = "list_video_view"
	LIST_REVIEW_VIEW        string = "list_review_view"
	LIST_POSTING_VIEW       string = "list_posting_view"
	LIST_RELATE_VIEW        string = "list_relate_view"
	LIST_BANNER_A_VIEW      string = "list_banner_a_view"
	LIST_BANNER_B_VIEW      string = "list_banner_b_view"
	LIST_BANNER_C_VIEW      string = "list_banner_c_view"
	LIST_BANNER_OWN_PROJECT string = "list_banner_own_project"
	REVIEW_STAR             string = "review_star"
	REVIEW_COUNT            string = "review_count"
	VIDEO_LINK              string = "video_link"

	PROJECT_FININSH     string = "project_finish"
	PROJECT_BRAND_IMAGE string = "project_brand_image"

	LAT = "lat"
	LNG = "lng"

	ICON = "icon"
	TEXT = "text"

	COUNT_RESULT string = "count_result"
	LIST_RESULT  string = "list_result"
	PARAMS       string = "params"
)

type DataParameter struct {
	Data ValueParam `json:"data"`
	jwt.StandardClaims
}

type ValueParam struct {
	Username        string `json:"username"`
	Password        string `json:"password"`
	PasswordConfirm string `json:"password_confirm"`
	FirstName       string `json:"first_name"`
	LastName        string `json:"last_name"`
	FullName        string `json:"full_name"`
	Birthdate       string `json:"birthdate"`
	TitleName       string `json:"title_name"` // mr , ms , mrs
	MobilePhone     string `json:"mobile_phone"`
	LineId          string `json:"line_id"`
	FacebookId      string `json:"facebook_id"`
	FacebookToken   string `json:"facebook_token"`
	Token           string `json:"token"`
	Nonce           string `json:"nonce"`
	AccessToken     string `json:"access_token"`
	ResetToken      string `json:"reset_token"`
	LANGUAGE        string `json:"language"`

	Max            int64 `json:"max"`
	Offset         int64 `json:"offset"`
	TimeStamp      int64 `json:"timestamp"`
	HouseSaleId    int64 `json:"house_sale_id"`
	HouseRentId    int64 `json:"house_rent_id"`
	HouseProjectId int64 `json:"house_project_id"`
	OwnProjectId   int64 `json:"own_project_id"`
	AgentId        int64 `json:"agent_id"`
	EntrepreneurId int64 `json:"entrepreneur_id"`
	ProvinceId     int64 `json:"province_id"`

	JSON_PARAMS JsonParams `json:"json_params"`
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

	if (err == nil || err.Error() == "Token used before issued") && token.Valid {
		return token, true
	} else {
		return nil, false
	}
}

func GetJsonData(ctx *context.Context) (bool, ValueParam) {
	claims, ok := ctx.Input.GetData(JWT_NEW_ASSIGN_VALUE).(DataParameter)
	return ok, claims.Data
}

func IntToFloat64(i int) float64 {
	return float64(i)
}
func Float64ToString(input_num float64, digit int) string {
	return strconv.FormatFloat(input_num, 'f', digit, 64)
}
