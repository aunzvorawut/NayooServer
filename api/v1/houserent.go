package v1

import (
	"gitlab.com/wisdomvast/NayooServer/models"
)

type HouserentController struct {
	GlobalApi
}

func (this *HouserentController) Main() {

}

func (this *HouserentController) List() {

	params := this.GlobalParamsJWT()

	allHouserent, countAllHouserent := models.GetAllHouserentOnClientByEnabledAndStartAndExpired(-1, 0)
	listHouserentResult := make([]map[string]interface{}, len(allHouserent))
	for i, v := range allHouserent {
		re := CreateOneHouserentContentMainView(v, params)
		listHouserentResult[i] = re
	}

	//allBannerSale := models.GetAllBannerByBannerTypeAndEnabledAndStartAndExpired(TYPE_SALE)
	//for i, v := range allBannerSale {
	//	re := this.CreateOneHouseBanner(v)
	//	list_house_result[i] = re
	//}

	allRelateHouserent, countRelateHoursrent := models.GetAllHouserentOnClientByEnabledAndStartAndExpired(-1, 0)
	listRelaterentResult := make([]map[string]interface{}, len(allRelateHouserent))
	for i, v := range allRelateHouserent {
		re := CreateOneHouserentContentRelateView(v, params)
		listRelaterentResult[i] = re
	}

	result := map[string]interface{}{
		LIST_POSTING_VIEW: map[string]interface{}{
			COUNT_RESULT: Int64ToString(countAllHouserent),
			LIST_RESULT:  listHouserentResult,
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

func CreateOneHouserentContentMainView(houserentObj *models.Houserent, params ValueParam) map[string]interface{} {
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

func CreateOneHouserentContentRelateView(houserentObj *models.Houserent, params ValueParam) map[string]interface{} {
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

func IsPromotionNaYooOnHouseRent(houseRentObj *models.Houserent) bool {
	switch os := houseRentObj.Id % 2; os {
	case 1:
		return true
	default:
		return false
	}
}

func IsVideo360OnHouseRent(houseRentObj *models.Houserent) bool {
	switch os := houseRentObj.Id % 3; os {
	case 0:
		return false
	case 1:
		return true
	default:
		return false
	}
}
