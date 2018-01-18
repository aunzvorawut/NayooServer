package v1

import (
	"gitlab.com/wisdomvast/NayooServer/models"
	"math/rand"
)

type EntrepreneurController struct {
	GlobalApi
}

func (this *EntrepreneurController) Main() {

}

func (this *EntrepreneurController) List() {

	params := this.GlobalParamsJWT()

	allEntrepreneur, countAllEntrepreneur := models.GetAllEntrepreneurOnClientByEnabledAndStartAndExpired(-1, 0)
	listEntrepreneurResult := make([]map[string]interface{}, len(allEntrepreneur))
	for i, v := range allEntrepreneur {
		re := this.CreateOneEntrepreneurContentMainView(v, params)
		listEntrepreneurResult[i] = re
	}

	//allBannerSale := models.GetAllBannerByBannerTypeAndEnabledAndStartAndExpired(TYPE_SALE)
	//for i, v := range allBannerSale {
	//	re := this.CreateOneHouseBanner(v)
	//	list_house_result[i] = re
	//}

	allRelateEntrepreneur, countRelateHouseProject := models.GetAllEntrepreneurOnClientByEnabledAndStartAndExpired(-1, 0)
	listRelateprojectResult := make([]map[string]interface{}, len(allRelateEntrepreneur))
	for i, v := range allRelateEntrepreneur {
		re := this.CreateOneEntrepreneurContentRelateView(v, params)
		listRelateprojectResult[i] = re
	}

	result := map[string]interface{}{
		LIST_POSTING_VIEW: map[string]interface{}{
			COUNT_RESULT: Int64ToString(countAllEntrepreneur),
			LIST_RESULT:  listEntrepreneurResult,
		},
		LIST_BANNER_A_VIEW: CreateMockyBanner(1),
		LIST_BANNER_B_VIEW: CreateMockyBanner(2),
		LIST_BANNER_C_VIEW: CreateMockyBanner(2),
		LIST_RELATE_VIEW: map[string]interface{}{
			COUNT_RESULT: Int64ToString(countRelateHouseProject),
			LIST_RESULT:  listRelateprojectResult,
		},
		PARAMS: params,
	}

	this.ResponseJSON(result, 200, "success")

}

func (this *EntrepreneurController) CreateOneEntrepreneurContentMainView(entrepreneurObj *models.Entrepreneur, params ValueParam) map[string]interface{} {
	result := map[string]interface{}{
		TITLE:                   GetStringByLanguage(entrepreneurObj.TitleTh, entrepreneurObj.TitleTh, entrepreneurObj.TitleEng, params),
		IMAGE:                   GetHostNayooName() + entrepreneurObj.Image,
		RESIDENT_ADDRESS:        "ต.บ้านเป็ด อ.เมืองขแนแก่น จ.ขอนแก่น",
		TAG_TYPE_LIST:           CreateMockyTagTypeList(3),
		IS_FAVORITE:             true,
		REVIEW_STAR:             4.5,
		REVIEW_COUNT:            5,
		VIP_TYPE:                entrepreneurObj.VipType,
		IS_VERIFY_PERSON_IDCARD: true,
		IS_VERIFY_PERSON_POLICY: true,
		IS_VERIFY_COMPANY:       true,
		IS_VERIFY_APARTMENT:     true,
		LAT:                     rand.Float64() + IntToFloat64(13),
		LNG:                     rand.Float64() + IntToFloat64(100),
		IS_PROMOTON_NAYOO:       IsPromotionNaYooOnEntrepreneur(entrepreneurObj),
	}
	return result
}

func (this *EntrepreneurController) CreateOneEntrepreneurContentRelateView(entrepreneurObj *models.Entrepreneur, params ValueParam) map[string]interface{} {
	result := map[string]interface{}{
		TITLE:            GetStringByLanguage(entrepreneurObj.TitleTh, entrepreneurObj.TitleTh, entrepreneurObj.TitleEng, params),
		REVIEW_STAR:      4.5,
		REVIEW_COUNT:     5,
		IMAGE:            GetHostNayooName() + entrepreneurObj.Image,
		RESIDENT_ADDRESS: "ต.บ้านเป็ด อ.เมืองขแนแก่น จ.ขอนแก่น",
		TAG_TYPE_LIST:    CreateMockyTagTypeList(3),
		LAT:              rand.Float64() + IntToFloat64(13),
		LNG:              rand.Float64() + IntToFloat64(100),
	}
	return result
}

func IsPromotionNaYooOnEntrepreneur(entrepreneurObj *models.Entrepreneur) bool {
	switch os := entrepreneurObj.Id % 2; os {
	case 1:
		return true
	default:
		return false
	}
}
