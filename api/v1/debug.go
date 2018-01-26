package v1

import (
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego"
	"encoding/json"
	"time"
	"strconv"
)

type DebugController struct {
	GlobalApi
}

func (this *DebugController) Main() {

	//SELECT 3956 * 2 * ASIN(SQRT( POWER(SIN((13.732 - 15.732) *  pi()/180 / 2), 2) +COS(13.732 * pi()/180) * COS(15.732 * pi()/180) * POWER(SIN((100.569 -105.569) * pi()/180 / 2), 2) )) as distance
	//from agent

	beego.Debug(Distance( 13.732 , 100.569 , 15.732 , 105.569  ))

	//testGoogleMap()

	// ===== request =====

	req := httplib.Post("http://www.mocky.io/v2/5a0a6ce22e0000ab1b489ca9")
	req.Param("customerOrderID", "1234")

	req.Header("testheader", "headervalue1")


	resultRequest, err := req.Debug(true).Response()

	// ทำ debug ตัว struct request
	beego.Debug(req)
	beego.Debug(resultRequest)

	// ====other====

	longHeaderString := "{\"aa\":\"aass\"}"

	var f map[string]interface{}
	err = json.Unmarshal([]byte(longHeaderString), &f)

	beego.Debug(err)

	for i,v := range f{
		beego.Debug(i)
		beego.Debug(v)
	}

	newDistance := Haversine(112312312, 0.1, 0.1, 0.1)
	beego.Debug(newDistance)

	t := time.Now()
	t.AddDate(0,0,1)
	beego.Debug(t.Weekday())
	t.AddDate(0,0,1)
	beego.Debug(t.Weekday())

	for i:=1 ; i<=7 ; i++ {
		t.AddDate(0,0,1)
		beego.Debug(t.AddDate(0,0,i).Weekday())
	}

	beego.Debug(t.Format("15:04:05"))


	beego.Debug("debug")

	nanoTime := time.Now().UnixNano()

	transaction := strconv.FormatInt(nanoTime,10)

	transaction = transaction[10:15]

	beego.Debug(transaction)

	timeTransaction := t.Format("20060102150405")

	this.ResponseJSON(timeTransaction+transaction,200,"success")

}