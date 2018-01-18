package v1

import (
	"gitlab.com/wisdomvast/NayooServer/models"
	"math/rand"
)

type AgentController struct {
	GlobalApi
}

func (this *AgentController) Main() {

}

func (this *AgentController) List() {

	params := this.GlobalParamsJWT()

	allAgent, countAllAgent := models.GetAllAgentOnClientByEnabledAndStartAndExpired(-1, 0)
	listAgentResult := make([]map[string]interface{}, len(allAgent))
	for i, v := range allAgent {
		re := CreateOneAgentContentMainView(v, params)
		listAgentResult[i] = re
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
			COUNT_RESULT: Int64ToString(countAllAgent),
			LIST_RESULT:  listAgentResult,
		},
		LIST_BANNER_A_VIEW: CreateMockyBanner(1),
		LIST_BANNER_B_VIEW: CreateMockyBanner(2),
		LIST_BANNER_C_VIEW: CreateMockyBanner(2),
		LIST_RELATE_VIEW: map[string]interface{}{
			COUNT_RESULT: Int64ToString(countRelateHouseProject),
			LIST_RESULT:  allRelateHouseProject,
		},
		PARAMS: params,
	}

	this.ResponseJSON(result, 200, "success")

}

func CreateOneAgentContentMainView(agentObj *models.Agent, params ValueParam) map[string]interface{} {
	result := map[string]interface{}{
		TITLE:                   GetStringByLanguage(agentObj.TitleTh, agentObj.TitleTh, agentObj.TitleEng, params),
		IMAGE:                   GetHostNayooName() + agentObj.Image,
		RESIDENT_ADDRESS:        "ต.บ้านเป็ด อ.เมืองขแนแก่น จ.ขอนแก่น",
		TAG_TYPE_LIST:           CreateMockyTagTypeList(3),
		IS_FAVORITE:             true,
		REVIEW_STAR:             4.5,
		REVIEW_COUNT:            5,
		VIP_TYPE:                agentObj.VipType,
		IS_VERIFY_PERSON_IDCARD: true,
		IS_VERIFY_PERSON_POLICY: true,
		IS_VERIFY_COMPANY:       true,
		IS_VERIFY_APARTMENT:     true,
		LAT:                     rand.Float64() + IntToFloat64(13),
		LNG:                     rand.Float64() + IntToFloat64(100),
	}
	return result
}
