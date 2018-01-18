package v1

import (
	"gitlab.com/wisdomvast/NayooServer/models"
	"math/rand"
)

type OwnProjectController struct {
	GlobalApi
}

func (this *OwnProjectController) Main() {

}

func (this *OwnProjectController) List() {

	params := this.GlobalParamsJWT()

	allOwnProject, countAllOwnProject := models.GetAllOwnProjectOnClientByEnabledAndStartAndExpired(-1, 0)
	listOwnProjectResult := make([]map[string]interface{}, len(allOwnProject))
	for i, v := range allOwnProject {
		re := CreateOneOwnProjectContentMainView(v , params)
		listOwnProjectResult[i] = re
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
			COUNT_RESULT: Int64ToString(countAllOwnProject),
			LIST_RESULT:  listOwnProjectResult,
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

func CreateOneOwnProjectContentMainView(ownProjectObj *models.OwnProject , params ValueParam) map[string]interface{} {
	result := map[string]interface{}{
		TITLE:            GetStringByLanguage(ownProjectObj.TitleTh, ownProjectObj.TitleTh, ownProjectObj.TitleEng , params),
		IMAGE:            GetHostNayooName() + ownProjectObj.Image,
		RESIDENT_ADDRESS: "ต.บ้านเป็ด อ.เมืองขแนแก่น จ.ขอนแก่น",
		TAG_TYPE_LIST:    CreateMockyTagTypeList(3),
		IS_FAVORITE:      true,
		VIP_TYPE:         ownProjectObj.VipType,
		LAT:              rand.Float64() + IntToFloat64(13),
		LNG:              rand.Float64() + IntToFloat64(100),
	}
	return result
}