package v1

import (
	"gitlab.com/wisdomvast/NayooServer/models"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

type HouseSaleController struct {
	GlobalApi
}

func (this *HouseSaleController) Main() {

}

func (this *HouseSaleController) List() {

	params := this.GlobalParamsJWT()

	allHouseale, countAllHouseSale := models.GetAllHouseSaleOnClientByEnabledAndStartAndExpired(-1, 0)
	listHouseSaleResult := make([]map[string]interface{}, len(allHouseale))
	for i, v := range allHouseale {
		re := CreateOneHouseSaleContentMainView(v, params)
		listHouseSaleResult[i] = re
	}

	//allBannerSale := models.GetAllBannerByBannerTypeAndEnabledAndStartAndExpired(TYPE_SALE)
	//for i, v := range allBannerSale {
	//	re := this.CreateOneHouseBanner(v)
	//	list_house_result[i] = re
	//}

	allRelateHouseProject, countRelateHouseProject := models.GetAllHouseProjectOnClientByEnabledAndStartAndExpired(-1, 0)
	listRelateHouseProjectResult := make([]map[string]interface{}, len(allRelateHouseProject))
	for i, v := range allRelateHouseProject {
		re := CreateOneHouseProjectContentRelateView(v , params)
		listRelateHouseProjectResult[i] = re
	}

	result := map[string]interface{}{
		LIST_POSTING_VIEW: map[string]interface{}{
			COUNT_RESULT: Int64ToString(countAllHouseSale),
			LIST_RESULT:  listHouseSaleResult,
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

func (this *HouseSaleController) ToggleFavorite(){

	params := this.GlobalParamsJWT()

	nonce := params.Nonce
	timeStamp := params.TimeStamp

	defer addUsedNonce(nonce,timeStamp)
	accessToken := params.AccessToken
	userObj := GetUserByToken(accessToken)

	if userObj != nil {
		housesaleId := params.HouseSaleId

		houseSaleObj, err := models.GetHouseSaleById(housesaleId)
		if err != nil || houseSaleObj == nil {
			beego.Error("err != nil || houseSaleObj == nil")
			this.ResponseJSON(nil , 401 , GetStringByLanguage(BAD_REQUEST_TH,BAD_REQUEST_TH,BAD_REQUEST_ENG,params))
			return
		}

		isFavorite, err := ToggleFavoriteHouseSale(userObj, houseSaleObj)
		if err != nil {
			beego.Error("isFavorite, err := ToggleFavoriteHouseSale(userObj, houseSaleObj)")
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

func CreateOneHouseSaleContentMainView(housesaleObj *models.HouseSale, params ValueParam) map[string]interface{} {
	result := map[string]interface{}{
		TITLE:                   GetStringByLanguage(housesaleObj.TitleTh, housesaleObj.TitleTh, housesaleObj.TitleEng, params),
		IMAGE:                   GetHostNayooName() + housesaleObj.Image,
		SALE_STR:                "2.62 ล้านบาท",
		SALE_STR_MARK:           "(ขาย/เช่า)",
		RENT_STR_MONTH:          "4,000 บาท/เดือน",
		RENT_STR_MONTH_MARK:     "(เช่า)",
		RENT_STR_DAY:            "600 บาท/วัน",
		RENT_STR_DAY_MARK:       "(เช่า)",
		RESIDENT_TYPE_LIST:      CreateMockyResidentType(housesaleObj.Id),
		RESIDENT_ADDRESS:        "ต.บ้านเป็ด",
		COUNT_BEDROOM:           "2",
		COUNT_BATHROOM:          "2",
		AREA_LAND:               "10ตร.ม.",
		AREA_USEFUL:             "22 ตร.วา",
		IS_PROMOTON_NAYOO:       IsPromotionNaYooOnHouseSale(housesaleObj),
		IS_FAVORITE:             true,
		IS_VERIFY_PERSON_IDCARD: true,
		IS_VERIFY_PERSON_POLICY: true,
		IS_VERIFY_COMPANY:       true,
		IS_VERIFY_APARTMENT:     true,
		VIP_TYPE:                housesaleObj.VipType,
	}
	return result
}

func IsPromotionNaYooOnHouseSale(houseSaleObj *models.HouseSale) bool {
	switch os := houseSaleObj.Id % 2; os {
	case 1:
		return true
	default:
		return false
	}
}

func ToggleFavoriteHouseSale(userObj *models.User, houseSaleObj *models.HouseSale) (bool, error) {
	isFavorite := false
	var err error
	if userObj != nil && houseSaleObj != nil {
		isFavorite = IsFavoriteHouseSale(userObj, houseSaleObj)
		o := orm.NewOrm()
		sqlStr := ""
		if isFavorite {
			sqlStr = "delete from user_house_sales where user_id =" + Int64ToString(userObj.Id) + " and house_sale_id = " + Int64ToString(houseSaleObj.Id)
			isFavorite = false
		} else {
			sqlStr = "insert ignore into user_house_sales (user_id, house_sale_id) values(" + Int64ToString(userObj.Id) + ", " + Int64ToString(houseSaleObj.Id) + ")"
			isFavorite = true
		}
		_, err = o.Raw(sqlStr).Exec()
		if err != nil {
			beego.Error(err)
		}
	}

	return isFavorite, err
}

func IsFavoriteHouseSale(userObj *models.User, houseSaleObj *models.HouseSale) bool {
	isFavorite := false
	o := orm.NewOrm()
	sqlStr := "select count(*) from user_house_sales where user_id =" + Int64ToString(userObj.Id) + " and house_sale_id=" + Int64ToString(houseSaleObj.Id)
	count := 0
	o.Raw(sqlStr).QueryRow(&count)
	if count > 0 {
		isFavorite = true
	}
	return isFavorite
}