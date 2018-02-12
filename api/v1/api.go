package v1

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/utils"
	"gitlab.com/wisdomvast/NayooServer/models"
	"io"
	"math"
	"mime/multipart"
	"net/textproto"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type GlobalApi struct {
	ParamsController
}

var fileSeparator = string(filepath.Separator)
var loc, _ = time.LoadLocation("Asia/Bangkok")

const (
	USERNAME_MAIL = "allabouthome.am@gmail.com"
	PASSWORD_MAIL = "aaham123456"
	HOST_MAIL     = "smtp.gmail.com"
	PORT_MAIL     = "587"
	FORM_MAIL     = "allabouthome.am@gmail.com"
	SecretKey     = "uB2Ph3YVzTjXjKhZ58tv"
	AccessKey     = "Ub2Ph3YVzTjXjKhZ58TV"

	TYPE_SALE         = "type_sale"
	TYPE_RENT         = "type_rent"
	TYPE_PROJECT      = "type_project"
	TYPE_ENTERPRENEUR = "type_enterpreneur"
	TYPE_AGENT        = "type_agent"


	DEFAULT_HOME = "default_home.jpg"
	DEFAULT_LOGO = "default_banner.png"
)

const emailTemplate = `Verify your email to start using Nayoo / กรุณายืนยัน  Email ของคุณ,<br/>
<br/>
<hr/>
<br/>
<br/>
Welcome to Nayoo Just verify your email to get started. We do this as a security precaution to verify your credentials.
<br/>
ยินดีต้อนรับสู่บริการ Nayoo กรุณายืนยันตัวตนของคุณด้วยการ Click  Verify Email ด้านล่าง<br/>
&nbsp;<a href="$url">กดที่นี่เพื่อยืนยัน</a><br/>`

const forgetPasswordTpl = `
Hi $user.username,<br/>
<br/>
You (or someone pretending to be you) requested that your password be reset.<br/>
<br/>
If you didn't make this request then ignore the email; no changes have been made.<br/>
<br/>
If you did make the request, then click <a href="$url">here</a> to reset your password.`

const subscribeTpl = `
Hi $user.username,<br/>
<br/>
ผลการจ่ายเงินผ่านบัตรทรูมันนี่<br/>
<br/>
$emailMsg<br/>
<br/>`

func Haversine(lonFrom float64, latFrom float64, lonTo float64, latTo float64) (distance float64) {
	var deltaLat = (latTo - latFrom) * (math.Pi / 180)
	var deltaLon = (lonTo - lonFrom) * (math.Pi / 180)

	var a = math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(latFrom*(math.Pi/180))*math.Cos(latTo*(math.Pi/180))*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	var c = 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return c
}

func addUsedNonce(nonce string, timestamp int64) {
	AddUsedNonce(nonce, timestamp)
}

func AddUsedNonce(nonce string, timestamp int64) {
	usedNonce := &models.UsedNonce{
		Nonce:     nonce,
		Timestamp: timestamp,
	}

	models.AddUsedNonce(usedNonce)

}

func GetUserByToken(accessToken string) *models.User {
	user := models.GetUserByAccessToken(accessToken)
	if user == nil {
		return nil
	}

	return user
}

func SendRegisterMail(username, token, typeStr string) {
	clientUrl := beego.AppConfig.String("client_register_url")
	urlToken := clientUrl + "?token=" + token
	emailTpl := strings.Replace(emailTemplate, "$url", urlToken, 1)
	subj := "Nayoo Registration Email"
	if typeStr == "resetpassword" {
		clientUrl = beego.AppConfig.String("client_reset_password_url")
		urlToken = clientUrl + "?token=" + token
		emailTpl = strings.Replace(forgetPasswordTpl, "$user.username", username, 1)
		emailTpl = strings.Replace(emailTpl, "$url", urlToken, 1)
		subj = "Nayoo Password Reset Email"
	}

	SendEmail(username, subj, emailTpl)
}

func SendSubscriptionMail(username, status string) {

	subj := "Nayoo ผลการทำรายการ"
	emailTpl := strings.Replace(subscribeTpl, "$user.username", username, 1)
	emailMsg := ""
	if status == "SUCCESS" {
		emailMsg = "คุณชำระเงินสำเร็จ คุณสามารถรับชมได้ที่ https://nayoo.com"
	} else if status == "USED" {
		emailMsg = "คุณชำระเงินไม่สำเร็จ บัตรถูกใช้งานไปแล้ว"
	} else if status == "INVALID_PIN" {
		emailMsg = "คุณชำระเงินไม่สำเร็จ เลขบัตรไม่ถูกต้อง"
	} else if status == "WRONG_TYPE" {
		emailMsg = "คุณชำระเงินไม่สำเร็จ บัตรของคุณไม่ใช่บัตรเงินสด"
	}

	emailTpl = strings.Replace(emailTpl, "$emailMsg", emailMsg, 1)
	SendEmail(username, subj, emailTpl)
}

func SendEmail(to, subj, message string) {

	email := utils.Email{
		Headers:  textproto.MIMEHeader{},
		Username: "allabouthome.am@gmail.com",
		Password: "aaham123456",
		Host:     "smtp.gmail.com",
		Port:     587,
		From:     "allabouthome.am@gmail.com",
		//Username: "support@wisdomcloud.net",
		//Password: "u794jkQnCy8mgDkm",
		//Host:     "mail.wisdomcloud.net",
		//Port:     587,
		//From:     "support@wisdomcloud.net",
		To:      []string{to},
		Subject: subj,
		HTML:    message,
	}

	err := email.Send()
	if err != nil {
		beego.Error("Send email error: ", err.Error())
	}
}

func (this *GlobalApi) UploadProfileImageGlobal(id int64, file multipart.File, handler *multipart.FileHeader, err error, typeImage string) (bool, string) {

	imgPath := "." + beego.AppConfig.String("proxyPath") + fileSeparator + "static" + fileSeparator + "img" + fileSeparator + typeImage // ทีี่เก็บรูป
	imgPathv2 := fileSeparator + "static" + fileSeparator + "img" + fileSeparator + typeImage + fileSeparator                           // ที่ลง db

	//if typeImage == "product" {
	//	productObj := models.GetProductByIdRelate(id)
	//	if productObj != nil {
	//		imgPath := "." + beego.AppConfig.String("proxyPath") + fileSeparator +
	//			"static" + fileSeparator + "img" + fileSeparator + typeImage
	//		DeleteOldImageFromUrl(imgPath, GetFileNameFromUrl(productObj.CoverImage.UrlPath))
	//	}
	//}

	err2 := os.MkdirAll(imgPath, 0755)
	if err2 != nil {
		RecordError("Can't create directory", "UploadProfileImage", err2)
		beego.Error("UploadProfileImage: ", err.Error())
		return false, "Can't create directory"
	}
	var title, movieIdStr string
	var movieId int64
	this.Ctx.Input.Bind(&title, "title")
	this.Ctx.Input.Bind(&movieIdStr, "movieIdStr")
	this.Ctx.Input.Bind(&movieId, "movieId")
	title = GetSafeFileName(movieIdStr + "_" + title)
	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	imgName := Int64ToString(id) + title + "_" + timestamp

	if err != nil {
		RecordError("Can't Get File", "UploadProfileImage", err2)
		beego.Error("UploadProfileImage: ", err.Error())
		return false, "Can't Get File"
	}
	defer file.Close()
	imgName = imgName + GetFileExtension(handler.Filename)

	//fmt.Fprintf(w, "%v", handler.Header)
	f, err := os.OpenFile(imgPath+fileSeparator+imgName, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		RecordError("Can't Open File", "UploadProfileImage", err)
		beego.Error("UploadProfileImage: ", err.Error())
		return false, "Can't Open File"
	}
	defer f.Close()
	_, err = io.Copy(f, file)
	if err != nil {
		RecordError("Can't write file", "UploadProfileImage", err)
		beego.Error("UploadProfileImage: ", err.Error())
		return false, "Can't write file"
	}

	//hostName := GetHostName(ctx)
	dir := GetCurrentPath()
	realPath := dir + beego.AppConfig.String("proxyPath") + "/static/img/" + typeImage + "/" + imgName

	beego.Debug("realPathGlobal = ", realPath)

	file1, err1 := os.Open(realPath)
	if err1 != nil {

		RecordError("Server Error", "UploadProfileImage", err1)
		beego.Error("UploadProfileImage: ", err1.Error())
		return false, "Server Error"

	}
	fi, err := file1.Stat()

	beego.Debug("fi = ", fi)

	if err != nil {
		RecordError("Internal Server Error", "UploadProfileImage", err)
		beego.Error("UploadProfileImage: ", err.Error())
		return false, "Internal Server Error"
	}

	return true, imgPathv2 + imgName

}

func RecordError(descripiton string, functioName string, err error) {
	beego.Debug(err)
	newErrorLogObj := models.ErrorLog{
		TypeError:   functioName,
		Description: descripiton,
	}
	_, err1 := models.AddErrorLog(&newErrorLogObj)
	if err1 != nil {
		beego.Debug("Add Error = ", err1)
	}
}

func GetStringByLanguage(stringDefault, stringTh, stringEng string , paramsJwt ValueParam) string {

	language := paramsJwt.LANGUAGE
	result := stringDefault

	if language == "th" || language == "TH" || language == "Th" {

		if stringTh == "" {
			result = stringDefault
		} else {
			result = stringTh
		}

	} else if language == "eng" || language == "ENG" || language == "Eng" {

		if stringEng == "" {
			result = stringDefault
		} else {
			result = stringTh
		}
	}

	if result == "" {
		result = stringTh
	}
	if result == "" {
		result = stringEng
	}

	return result

}

func GetCurrentPath() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0])) // get current path
	if err != nil {
		return ""
	}
	return dir
}

func GetFileExtension(fileName string) string {
	var extension = filepath.Ext(fileName)
	return extension
}

func GetSafeFileName(input string) string {
	reg, err := regexp.Compile("[^A-Za-z0-9_]+")
	if err != nil {
		beego.Error("")
	}

	safe := reg.ReplaceAllString(input, "-")
	safe = strings.ToLower(strings.Trim(safe, "-"))
	return safe
}

func GetSafeName(input string) string {
	reg, err := regexp.Compile("[^A-Za-z0-9]+")
	if err != nil {
		beego.Error("")
	}

	safe := reg.ReplaceAllString(input, "")
	safe = strings.ToLower(strings.Trim(safe, "-"))
	return safe
}

//=========== LIB CONTENT JSON DOWN HERE =================
//=========== LIB CONTENT JSON DOWN HERE =================
//=========== LIB CONTENT JSON DOWN HERE =================
//=========== LIB CONTENT JSON DOWN HERE =================
//=========== LIB CONTENT JSON DOWN HERE =================
//=========== LIB CONTENT JSON DOWN HERE =================


func CreateMockyBanner(size int) []map[string]interface{} {

	result := make([]map[string]interface{}, size)

	for i := 1; i <= size; i++ {
		re := map[string]interface{}{
			ID:    i,
			IMAGE: GetHostNayooName() + "static/img/"+DEFAULT_HOME,
		}

		result[i-1] = re
	}
	return result

}

func CreateMockyTagTypeList(size int) []map[string]interface{}{

	result := make([]map[string]interface{} , size)
	for i := 0 ; i< size ; i++{
		result[i] = map[string]interface{}{
			ID:i+1,
			TAG_TYPE_STR:"ผ้าม่าน/วอลเปเปอร์/ฉากกั้นห้อง",
		}
	}
	return result
}

func CreateMockyVideoList(size int) []map[string]interface{}{
	result := make([]map[string]interface{} , size)
	for i := 0 ; i< size ; i++{
		result[i] = map[string]interface{}{
			ID:i+1,
			VIDEO_LINK:"http://techslides.com/demos/sample-videos/small.mp4",
		}
	}
	return result

}

func CreateMockyArticle(size int) []map[string]interface{}{
	result := make([]map[string]interface{} , size)
	for i := 0 ; i< size ; i++{
		result[i] = map[string]interface{}{
			ID:i+1,
			IMAGE: GetHostNayooName() + "static/img/"+DEFAULT_HOME,
			TITLE: "lorem title",
			DESCRIPTION : "lorem description",
			DATE : "12 กันยายน 2560",
		}
	}
	return result

}

func CreateMockyResidentType(id int64) []map[string]interface{} {

	switch os := id % 5; os {
	case 1:
		return []map[string]interface{}{
			map[string]interface{}{
				ICON:"",
				TEXT:"บ้านเดี่ยว",
			},
			map[string]interface{}{
				ICON:"",
				TEXT:"คอนโด",
			},
		}
	case 2:
		return []map[string]interface{}{
			map[string]interface{}{
				ICON:"",
				TEXT:"บ้านเดี่ยว",
			},
			map[string]interface{}{
				ICON:"",
				TEXT:"ทาวน์โฮม",
			},
		}
	case 3:
		return []map[string]interface{}{
			map[string]interface{}{
				ICON:"",
				TEXT:"บ้านเดี่ยว",
			},
		}
	case 4:
		return []map[string]interface{}{
			map[string]interface{}{
				ICON:"",
				TEXT:"ทาวน์โฮม",
			},
		}
	default:
		return []map[string]interface{}{
			map[string]interface{}{
				ICON:"",
				TEXT:"คอนโด",
			},
		}
	}

}

func GetHostNayooName() string {
	return beego.AppConfig.String("nayooServerName")
}

func hsin(theta float64) float64 {
	return math.Pow(math.Sin(theta/2), 2)
}

func Distance(lat1, lon1, lat2, lon2 float64) float64 {
	// convert to radians
	// must cast radius as float to multiply later
	var la1, lo1, la2, lo2, r float64
	la1 = lat1 * math.Pi / 180
	lo1 = lon1 * math.Pi / 180
	la2 = lat2 * math.Pi / 180
	lo2 = lon2 * math.Pi / 180

	r = 6378100 // Earth radius in METERS

	// calculate
	h := hsin(la2-la1) + math.Cos(la1)*math.Cos(la2)*hsin(lo2-lo1)

	return 2 * r * math.Asin(math.Sqrt(h))
}