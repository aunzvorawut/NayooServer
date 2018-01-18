package v1

import (
	"gitlab.com/wisdomvast/NayooServer/models"
)

type HousesaleController struct {
	GlobalApi
}

func (this *HousesaleController) Main() {

}

func (this *HousesaleController) List() {

	params := this.GlobalParamsJWT()

	allHouseale, countAllHousesale := models.GetAllHousesaleOnClientByEnabledAndStartAndExpired(-1, 0)
	listHousesaleResult := make([]map[string]interface{}, len(allHouseale))
	for i, v := range allHouseale {
		re := CreateOneHousesaleContentMainView(v, params)
		listHousesaleResult[i] = re
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
			COUNT_RESULT: Int64ToString(countAllHousesale),
			LIST_RESULT:  listHousesaleResult,
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

func CreateOneHousesaleContentMainView(housesaleObj *models.Housesale, params ValueParam) map[string]interface{} {
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

func IsPromotionNaYooOnHouseSale(houseSaleObj *models.Housesale) bool {
	switch os := houseSaleObj.Id % 2; os {
	case 1:
		return true
	default:
		return false
	}
}