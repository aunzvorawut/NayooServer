package v1

import (
	"gitlab.com/wisdomvast/NayooServer/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

type HouseRentController struct {
	GlobalApi
}

func (this *HouseRentController) Main() {

}

func (this *HouseRentController) List() {

	params := this.GlobalParamsJWT()

	allHouseRent, countAllHouseRent := models.GetAllHouseRentOnClientByEnabledAndStartAndExpired(-1, 0)
	listHouseRentResult := make([]map[string]interface{}, len(allHouseRent))
	for i, v := range allHouseRent {
		re := CreateOneHouseRentContentMainView(v, params)
		listHouseRentResult[i] = re
	}

	//allBannerSale := models.GetAllBannerByBannerTypeAndEnabledAndStartAndExpired(TYPE_SALE)
	//for i, v := range allBannerSale {
	//	re := this.CreateOneHouseBanner(v)
	//	list_house_result[i] = re
	//}

	allRelateHouseRent, countRelateHoursrent := models.GetAllHouseRentOnClientByEnabledAndStartAndExpired(-1, 0)
	listRelaterentResult := make([]map[string]interface{}, len(allRelateHouseRent))
	for i, v := range allRelateHouseRent {
		re := CreateOneHouseRentContentRelateView(v, params)
		listRelaterentResult[i] = re
	}

	result := map[string]interface{}{
		LIST_POSTING_VIEW: map[string]interface{}{
			COUNT_RESULT: Int64ToString(countAllHouseRent),
			LIST_RESULT:  listHouseRentResult,
		},
		LIST_BANNER_A_VIEW: CreateMockyBanner(1),
		LIST_BANNER_B_VIEW: CreateMockyBanner(2),
		LIST_BANNER_C_VIEW: CreateMockyBanner(2),
		LIST_RELATE_VIEW: map[string]interface{}{
			COUNT_RESULT: Int64ToString(countRelateHoursrent),
			LIST_RESULT:  listRelaterentResult,
		},
		PARAMS: params,
	}

	this.ResponseJSON(result, 200, "success")

}

func (this *HouseRentController) ToggleFavorite(){

	params := this.GlobalParamsJWT()

	nonce := params.Nonce
	timeStamp := params.TimeStamp

	defer addUsedNonce(nonce,timeStamp)
	accessToken := params.AccessToken
	userObj := GetUserByToken(accessToken)

	if userObj != nil {
		houseRateId := params.HouseRentId

		houseRentObj, err := models.GetHouseRentById(houseRateId)
		if err != nil || houseRentObj == nil {
			beego.Error("err != nil || houseRentObj == nil")
			this.ResponseJSON(nil , 401 , GetStringByLanguage(BAD_REQUEST_TH,BAD_REQUEST_TH,BAD_REQUEST_ENG,params))
			return
		}

		isFavorite, err := ToggleFavoriteHouseRent(userObj, houseRentObj)
		if err != nil {
			beego.Error("isFavorite, err := ToggleFavoriteHouseRent(userObj, houseRentObj)")
			this.ResponseJSON(nil , 500 , GetStringByLanguage(SERVER_ERROR_TH,SERVER_ERROR_TH,SERVER_ERROR_ENG , params))
			return

		} else {
			this.ResponseJSON(map[string]interface{}{
				IS_FAVORITE:isFavorite,
			} , 200 , GetStringByLanguage(SUCCESS_TH,SUCCESS_TH,SUCCESS_ENG , params))
			return
		}

	} else {
		beego.Error("userObj != nil")
		this.ResponseJSON(nil , 401 , GetStringByLanguage(LOGIN_FAIL_TH,LOGIN_FAIL_TH,LOGIN_FAIL_ENG , params))
		return
	}

}

func CreateOneHouseRentContentMainView(houserentObj *models.HouseRent, params ValueParam) map[string]interface{} {
	result := map[string]interface{}{
		TITLE:                   GetStringByLanguage(houserentObj.TitleTh, houserentObj.TitleTh, houserentObj.TitleEng, params),
		IMAGE:                   GetHostNayooName() + houserentObj.Image,
		SALE_STR:                "",
		SALE_STR_MARK:           "",
		RENT_STR_MONTH:          "4,000 บาท/เดือน",
		RENT_STR_MONTH_MARK:     "(เช่า)",
		RENT_STR_DAY:            "600 บาท/วัน",
		RENT_STR_DAY_MARK:       "(เช่า)",
		RESIDENT_TYPE_LIST:      CreateMockyResidentType(houserentObj.Id),
		RESIDENT_ADDRESS:        "ต.บ้านเป็ด",
		COUNT_BEDROOM:           "2",
		COUNT_BATHROOM:          "2",
		IS_PROMOTON_NAYOO:       IsPromotionNaYooOnHouseRent(houserentObj),
		IS_VIDEO_360:            IsVideo360OnHouseRent(houserentObj),
		IS_FAVORITE:             true,
		IS_VERIFY_PERSON_IDCARD: true,
		IS_VERIFY_PERSON_POLICY: true,
		IS_VERIFY_COMPANY:       true,
		IS_VERIFY_APARTMENT:     true,
		VIP_TYPE:                houserentObj.VipType,
	}
	return result
}

func CreateOneHouseRentContentRelateView(houserentObj *models.HouseRent, params ValueParam) map[string]interface{} {
	result := map[string]interface{}{
		TITLE:               GetStringByLanguage(houserentObj.TitleTh, houserentObj.TitleTh, houserentObj.TitleEng, params),
		IMAGE:               GetHostNayooName(),
		IS_PROMOTON_NAYOO:   IsPromotionNaYooOnHouseRent(houserentObj),
		IS_VIDEO_360:        IsVideo360OnHouseRent(houserentObj),
		RESIDENT_ADDRESS:    "ต.บ้านเป็ด",
		RENT_STR_MONTH:      "2,000 - 3,500",
		RENT_STR_MONTH_MARK: "บาท/เดือน",
		RENT_STR_DAY:        "200 - 400",
		RENT_STR_DAY_MARK:   "บาท/วัน",
		RESIDENT_TYPE_LIST:  CreateMockyResidentType(houserentObj.Id),
	}
	return result
}

func IsPromotionNaYooOnHouseRent(houseRentObj *models.HouseRent) bool {
	switch os := houseRentObj.Id % 2; os {
	case 1:
		return true
	default:
		return false
	}
}

func IsVideo360OnHouseRent(houseRentObj *models.HouseRent) bool {
	switch os := houseRentObj.Id % 3; os {
	case 0:
		return false
	case 1:
		return true
	default:
		return false
	}
}

func ToggleFavoriteHouseRent(userObj *models.User, houseRentObj *models.HouseRent) (bool, error) {
	isFavorite := false
	var err error
	if userObj != nil && houseRentObj != nil {
		isFavorite = IsFavoriteHouseRent(userObj, houseRentObj)
		o := orm.NewOrm()
		sqlStr := ""
		if isFavorite {
			sqlStr = "delete from user_house_rents where user_id =" + Int64ToString(userObj.Id) + " and house_rent_id = " + Int64ToString(houseRentObj.Id)
			isFavorite = false
		} else {
			sqlStr = "insert ignore into user_house_rents (user_id, house_rent_id) values(" + Int64ToString(userObj.Id) + ", " + Int64ToString(houseRentObj.Id) + ")"
			isFavorite = true
		}
		_, err = o.Raw(sqlStr).Exec()
		if err != nil {
			beego.Error(err)
		}
	}

	return isFavorite, err
}

func IsFavoriteHouseRent(userObj *models.User, houseRentObj *models.HouseRent) bool {
	isFavorite := false
	o := orm.NewOrm()
	sqlStr := "select count(*) from user_house_rents where user_id =" + Int64ToString(userObj.Id) + " and house_rent_id=" + Int64ToString(houseRentObj.Id)
	count := 0
	o.Raw(sqlStr).QueryRow(&count)
	if count > 0 {
		isFavorite = true
	}
	return isFavorite
}