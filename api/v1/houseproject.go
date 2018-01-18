package v1

import (
	"gitlab.com/wisdomvast/NayooServer/models"
	"math/rand"
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
		LAT:                 rand.Float64() + IntToFloat64(13),
		LNG:                 rand.Float64() + IntToFloat64(100),
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
		LAT:                 rand.Float64() + IntToFloat64(13),
		LNG:                 rand.Float64() + IntToFloat64(100),
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
