package v1

import (
	"gitlab.com/wisdomvast/NayooServer/models"
	"math/rand"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
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

func (this *AgentController) ToggleFavorite(){

	params := this.GlobalParamsJWT()

	nonce := params.Nonce
	timeStamp := params.TimeStamp

	defer addUsedNonce(nonce,timeStamp)
	accessToken := params.AccessToken
	userObj := GetUserByToken(accessToken)

	if userObj != nil {
		agentId := params.AgentId

		agentObj, err := models.GetAgentById(agentId)
		if err != nil || agentObj == nil {
			beego.Error("err != nil || agentObj == nil")
			this.ResponseJSON(nil , 401 , GetStringByLanguage(BAD_REQUEST_TH,BAD_REQUEST_TH,BAD_REQUEST_ENG,params))
			return
		}

		isFavorite, err := ToggleFavoriteAgent(userObj, agentObj)
		if err != nil {
			beego.Error("isFavorite, err := ToggleFavoriteAgent(userObj, agentObj)")
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

func ToggleFavoriteAgent(userObj *models.User, agentObj *models.Agent) (bool, error) {
	isFavorite := false
	var err error
	if userObj != nil && agentObj != nil {
		isFavorite = IsFavoriteAgent(userObj, agentObj)
		o := orm.NewOrm()
		sqlStr := ""
		if isFavorite {
			sqlStr = "delete from user_agents where user_id =" + Int64ToString(userObj.Id) + " and agent_id = " + Int64ToString(agentObj.Id)
			isFavorite = false
		} else {
			sqlStr = "insert ignore into user_agents (user_id, agent_id) values(" + Int64ToString(userObj.Id) + ", " + Int64ToString(agentObj.Id) + ")"
			isFavorite = true
		}
		_, err = o.Raw(sqlStr).Exec()
		if err != nil {
			beego.Error(err)
		}
	}

	return isFavorite, err
}

func IsFavoriteAgent(userObj *models.User, agentObj *models.Agent) bool {
	isFavorite := false
	o := orm.NewOrm()
	sqlStr := "select count(*) from user_agents where user_id =" + Int64ToString(userObj.Id) + " and agent_id=" + Int64ToString(agentObj.Id)
	count := 0
	o.Raw(sqlStr).QueryRow(&count)
	if count > 0 {
		isFavorite = true
	}
	return isFavorite
}