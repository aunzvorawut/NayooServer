package v1

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"gitlab.com/wisdomvast/NayooServer/models"
	"math/rand"
	"strconv"
)

type HouseProjectController struct {
	GlobalApi
}

func (this *HouseProjectController) Main() {

	params := this.GlobalParamsJWT()

	allRelateHouseProject, countRelateHouseProject := models.GetAllHouseProjectOnClientByEnabledAndStartAndExpired(-1, 0)
	listRelateHouseProjectResult := make([]map[string]interface{}, len(allRelateHouseProject))
	for i, v := range allRelateHouseProject {
		re := CreateOneHouseProjectContentRelateView(v, params)
		listRelateHouseProjectResult[i] = re
	}

	result := map[string]interface{}{

		LIST_RECOMMEND_VIEW: map[string]interface{}{
			TITLE:        "โครงการแนะนำ",
			COUNT_RESULT: Int64ToString(countRelateHouseProject),
			LIST_RESULT:  listRelateHouseProjectResult,
		},

		LIST_BANNER_OWN_PROJECT: map[string]interface{}{
			TITLE:        "เจ้าของโครงการ",
			COUNT_RESULT: 5,
			LIST_RESULT:  CreateMockyBanner(5),
		},

		LIST_VIDEO_VIEW: map[string]interface{}{
			TITLE:        "วีดิโอแนะนำ",
			COUNT_RESULT: 5,
			LIST_RESULT:  CreateMockyVideoList(5),
		},
		LIST_REVIEW_VIEW: map[string]interface{}{
			TITLE:        "รีวิวโครงการ",
			COUNT_RESULT: countRelateHouseProject,
			LIST_RESULT:  listRelateHouseProjectResult,
		},

		LIST_BANNER_A_VIEW: CreateMockyBanner(1),
		LIST_BANNER_B_VIEW: CreateMockyBanner(2),
		LIST_BANNER_C_VIEW: CreateMockyBanner(2),
	}

	this.ResponseJSON(result, 200, "success")
	return

}

func (this *HouseProjectController) List() {

	params := this.GlobalParamsJWT()

	allHouseProject, countAllHouseProject := models.GetAllHouseProjectOnClientByEnabledAndStartAndExpired(-1, 0)
	listHouseProjectResult := make([]map[string]interface{}, len(allHouseProject))
	for i, v := range allHouseProject {
		re := CreateOneHouseProjectContentMainView(v, params)
		listHouseProjectResult[i] = re
	}

	//allBannerSale := models.GetAllBannerByBannerTypeAndEnabledAndStartAndExpired(TYPE_SALE)
	//for i, v := range allBannerSale {
	//	re := this.CreateOneHouseBanner(v)
	//	list_house_result[i] = re
	//}

	allRelateHouseProject, countRelateHouseProject := models.GetAllHouseProjectOnClientByEnabledAndStartAndExpired(-1, 0)
	listRelateHouseProjectResult := make([]map[string]interface{}, len(allRelateHouseProject))
	for i, v := range allRelateHouseProject {
		re := CreateOneHouseProjectContentRelateView(v, params)
		listRelateHouseProjectResult[i] = re
	}

	result := map[string]interface{}{
		LIST_POSTING_VIEW: map[string]interface{}{
			COUNT_RESULT: Int64ToString(countAllHouseProject),
			LIST_RESULT:  listHouseProjectResult,
		},
		LIST_BANNER_A_VIEW: CreateMockyBanner(1),
		LIST_BANNER_B_VIEW: CreateMockyBanner(2),
		LIST_BANNER_C_VIEW: CreateMockyBanner(2),
		LIST_RELATE_VIEW: map[string]interface{}{
			COUNT_RESULT: Int64ToString(countRelateHouseProject),
			LIST_RESULT:  listRelateHouseProjectResult,
		},
		PARAMS: params,
	}

	this.ResponseJSON(result, 200, "success")

}

func (this *HouseProjectController) TinyDetail() {

	params := this.GlobalParamsJWT()
	houseProjectId := params.HouseProjectId
	houseProjectObj , _ := models.GetHouseProjectById(houseProjectId)
	if houseProjectObj == nil{
		this.ResponseJSON(nil,400,GetStringByLanguage(BAD_REQUEST_TH,BAD_REQUEST_TH,BAD_REQUEST_ENG,params))
		return
	}

	result := CreateOneHouseProjectContentRelateView(houseProjectObj,params)

	this.ResponseJSON(result,200,GetStringByLanguage(SUCCESS_TH,SUCCESS_TH,SUCCESS_ENG,params))
	return

}

func (this *HouseProjectController) MockListMap() {
	params := this.GlobalParamsJWT()
	vipTypes := [8]string{"silver", "gold", "standard", "bronze", "silver", "gold", "standard", "bronze"}
	for i:= 0; i<=100; i++ {
		vipType := vipTypes[i%4]
		lat := rand.Float64() + IntToFloat64(16)
		lng := rand.Float64() + IntToFloat64(102)
		house := models.HouseProject{
			Enabled: true,
			VipType: vipType,
			Image:  "static/img/default_home.jpg",
			TitleTh:  "สินธาราบ้านโจด",
			TitleEng:  "",
			Lat: lat,
			Lng: lng,
			Location: "POINT(" + strconv.FormatFloat(lat, 'f', 14, 64) + " " + strconv.FormatFloat(lng, 'f', 14, 64) + ")",
		}
		_, err := models.AddHouseProject(&house)
		if err != nil {
			beego.Error(err)
		}
	}
	this.ResponseJSON(nil,200,GetStringByLanguage(SUCCESS_TH,SUCCESS_TH,SUCCESS_ENG,params))
	return
}

func (this *HouseProjectController) TonListMap() {
	params := this.GlobalParamsJWT()
	origLat := 16.60466028797962
	origlng := 102.94050908804502
	radius := 20

	results := models.GetHouseProjectNearByLocation(origLat, origlng, radius)
	this.ResponseJSON(map[string]interface{}{
		"s": results,
	}, 200,GetStringByLanguage(SUCCESS_TH,SUCCESS_TH,SUCCESS_ENG,params))
}

func (this *HouseProjectController) ListMap() {

	params := this.GlobalParamsJWT()
	origLat := params.LAT
	origlng := params.LNG
	radius := params.Radius

	beego.Debug(origLat , origlng , radius)

	result := make([]map[string]interface{},3)

	re := map[string]interface{}{
		ID : 1,
		LAT : 13.732,
		LNG : 100.569,
	}
	result[0] = re

	re = map[string]interface{}{
		ID : 2,
		LAT : 13.742,
		LNG : 100.579,
	}
	result[1] = re

	re = map[string]interface{}{
		ID : 3,
		LAT : 13.752,
		LNG : 100.589,
	}
	result[2] = re

	this.ResponseJSON(result,200,GetStringByLanguage(SUCCESS_TH,SUCCESS_TH,SUCCESS_ENG,params))
	return

}

func (this *HouseProjectController) ToggleFavorite() {

	params := this.GlobalParamsJWT()

	nonce := params.Nonce
	timeStamp := params.TimeStamp

	defer addUsedNonce(nonce, timeStamp)
	accessToken := params.AccessToken
	userObj := GetUserByToken(accessToken)

	if userObj != nil {
		houseProjectId := params.HouseProjectId

		houseProjectObj, err := models.GetHouseProjectById(houseProjectId)
		if err != nil || houseProjectObj == nil {
			beego.Error("err != nil || houseProjectObj == nil")
			this.ResponseJSON(nil, 401, GetStringByLanguage(BAD_REQUEST_TH, BAD_REQUEST_TH, BAD_REQUEST_ENG, params))
			return
		}

		isFavorite, err := ToggleFavoriteHouseProject(userObj, houseProjectObj)
		if err != nil {
			beego.Error("isFavorite, err := ToggleFavoriteHouseProject(userObj, houseProjectObj)")
			this.ResponseJSON(nil, 500, GetStringByLanguage(SERVER_ERROR_TH, SERVER_ERROR_TH, SERVER_ERROR_ENG, params))
			return

		} else {
			this.ResponseJSON(map[string]interface{}{
				IS_FAVORITE: isFavorite,
			}, 200, GetStringByLanguage(SUCCESS_TH, SUCCESS_TH, SUCCESS_ENG, params))
			return
		}

	} else {
		beego.Error("userObj != nil")
		this.ResponseJSON(nil, 401, GetStringByLanguage(LOGIN_FAIL_TH, LOGIN_FAIL_TH, LOGIN_FAIL_ENG, params))
		return
	}

}

func CreateOneHouseProjectContentMainView(HouseProjectObj *models.HouseProject, params ValueParam) map[string]interface{} {
	result := map[string]interface{}{
		TITLE:               GetStringByLanguage(HouseProjectObj.TitleTh, HouseProjectObj.TitleTh, HouseProjectObj.TitleEng, params),
		IMAGE:               GetHostNayooName() + HouseProjectObj.Image,
		SALE_STR:            "1.62 ล้านบาท",
		SALE_STR_MARK:       "",
		RENT_STR_MONTH:      "",
		RENT_STR_MONTH_MARK: "",
		RENT_STR_DAY:        "",
		RENT_STR_DAY_MARK:   "",
		RESIDENT_TYPE_LIST:  CreateMockyResidentType(HouseProjectObj.Id),
		RESIDENT_ADDRESS:    "ต.บ้านเป็ด อ.เมืองขแนแก่น จ.ขอนแก่น",
		PROJECT_FININSH:     "สร้างเสร็จปี 2017",
		PROJECT_BRAND_IMAGE: GetHostNayooName() + HouseProjectObj.Image,
		IS_PROMOTON_NAYOO:   IsPromotionNaYooOnHouseProject(HouseProjectObj),
		IS_GURU:             IsGuruOnHouseProject(HouseProjectObj),
		VIP_TYPE:            HouseProjectObj.VipType,
		LAT:                 rand.Float64() + IntToFloat64(16),
		LNG:                 rand.Float64() + IntToFloat64(102),
	}
	return result
}

func CreateOneHouseProjectContentRelateView(HouseProjectObj *models.HouseProject, params ValueParam) map[string]interface{} {
	result := map[string]interface{}{
		TITLE:               GetStringByLanguage(HouseProjectObj.TitleTh, HouseProjectObj.TitleTh, HouseProjectObj.TitleEng, params),
		DESCRIPTION:         "Lorem ipsum dolor sit amet, consectetur adipiscing elit, sed do eiusmod tempor incididunt ut labore et dolore magna aliqua.",
		SALE_STR:            "เริ่มที่ 2.6 ล้านบาท",
		IMAGE:               GetHostNayooName(),
		IS_PROMOTON_NAYOO:   IsPromotionNaYooOnHouseProject(HouseProjectObj),
		IS_GURU:             IsGuruOnHouseProject(HouseProjectObj),
		RESIDENT_ADDRESS:    "ต.บ้านเป็ด อ.เมืองขแนแก่น จ.ขอนแก่น",
		PROJECT_BRAND_IMAGE: GetHostNayooName() + HouseProjectObj.Image,
		LAT:                 rand.Float64() + IntToFloat64(16),
		LNG:                 rand.Float64() + IntToFloat64(102),
	}
	return result
}

func IsPromotionNaYooOnHouseProject(houseProjectObj *models.HouseProject) bool {
	switch os := houseProjectObj.Id % 2; os {
	case 1:
		return true
	default:
		return false
	}
}

func IsGuruOnHouseProject(houseProjectObj *models.HouseProject) bool {
	switch os := houseProjectObj.Id % 3; os {
	case 0:
		return false
	case 1:
		return true
	default:
		return false
	}
}

func ToggleFavoriteHouseProject(userObj *models.User, houseProjectObj *models.HouseProject) (bool, error) {
	isFavorite := false
	var err error
	if userObj != nil && houseProjectObj != nil {
		isFavorite = IsFavoriteHouseProject(userObj, houseProjectObj)
		o := orm.NewOrm()
		sqlStr := ""
		if isFavorite {
			sqlStr = "delete from user_house_projects where user_id =" + Int64ToString(userObj.Id) + " and house_project_id = " + Int64ToString(houseProjectObj.Id)
			isFavorite = false
		} else {
			sqlStr = "insert ignore into user_house_projects (user_id, house_project_id) values(" + Int64ToString(userObj.Id) + ", " + Int64ToString(houseProjectObj.Id) + ")"
			isFavorite = true
		}
		_, err = o.Raw(sqlStr).Exec()
		if err != nil {
			beego.Error(err)
		}
	}

	return isFavorite, err
}

func IsFavoriteHouseProject(userObj *models.User, houseProjectObj *models.HouseProject) bool {
	isFavorite := false
	o := orm.NewOrm()
	sqlStr := "select count(*) from user_house_projects where user_id =" + Int64ToString(userObj.Id) + " and house_project_id=" + Int64ToString(houseProjectObj.Id)
	count := 0
	o.Raw(sqlStr).QueryRow(&count)
	if count > 0 {
		isFavorite = true
	}
	return isFavorite
}
