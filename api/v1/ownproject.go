package v1

import (
	"gitlab.com/wisdomvast/NayooServer/models"
	"math/rand"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
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

func (this *OwnProjectController) ToggleFavorite(){

	params := this.GlobalParamsJWT()

	nonce := params.Nonce
	timeStamp := params.TimeStamp

	defer addUsedNonce(nonce,timeStamp)
	accessToken := params.AccessToken
	userObj := GetUserByToken(accessToken)

	if userObj != nil {
		ownProjectId := params.OwnProjectId

		ownProjectObj, err := models.GetOwnProjectById(ownProjectId)
		if err != nil || ownProjectObj == nil {
			beego.Error("err != nil || ownProjectObj == nil")
			this.ResponseJSON(nil , 401 , GetStringByLanguage(BAD_REQUEST_TH,BAD_REQUEST_TH,BAD_REQUEST_ENG,params))
			return
		}

		isFavorite, err := ToggleFavoriteOwnProject(userObj, ownProjectObj)
		if err != nil {
			beego.Error("isFavorite, err := ToggleFavoriteOwnProject(userObj, ownProjectObj)")
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

func ToggleFavoriteOwnProject(userObj *models.User, ownProjectObj *models.OwnProject) (bool, error) {
	isFavorite := false
	var err error
	if userObj != nil && ownProjectObj != nil {
		isFavorite = IsFavoriteOwnProject(userObj, ownProjectObj)
		o := orm.NewOrm()
		sqlStr := ""
		if isFavorite {
			sqlStr = "delete from user_own_projects where user_id =" + Int64ToString(userObj.Id) + " and own_project_id = " + Int64ToString(ownProjectObj.Id)
			isFavorite = false
		} else {
			sqlStr = "insert ignore into user_own_projects (user_id, own_project_id) values(" + Int64ToString(userObj.Id) + ", " + Int64ToString(ownProjectObj.Id) + ")"
			isFavorite = true
		}
		_, err = o.Raw(sqlStr).Exec()
		if err != nil {
			beego.Error(err)
		}
	}

	return isFavorite, err
}

func IsFavoriteOwnProject(userObj *models.User, ownProjectObj *models.OwnProject) bool {
	isFavorite := false
	o := orm.NewOrm()
	sqlStr := "select count(*) from user_own_projects where user_id =" + Int64ToString(userObj.Id) + " and own_project_id=" + Int64ToString(ownProjectObj.Id)
	count := 0
	o.Raw(sqlStr).QueryRow(&count)
	if count > 0 {
		isFavorite = true
	}
	return isFavorite
}