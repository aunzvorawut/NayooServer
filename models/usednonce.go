package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego"
)

type UsedNonce struct {
	Id        int64  `orm:"auto;pk"`
	//Nonce     string `orm:"unique;size(255)"`
	Nonce     string `orm:"size(255)"`
	Timestamp int64
}

func init() {
	orm.RegisterModel(new(UsedNonce))
}

// AddUsedNonce insert a new UsedNonce into database and returns
// last inserted Id on success.
func AddUsedNonce(m *UsedNonce) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetUsedNonceById retrieves UsedNonce by Id. Returns error if
// Id doesn't exist
func GetUsedNonceById(id int64) (v *UsedNonce, err error) {
	o := orm.NewOrm()
	v = &UsedNonce{Id: id}
	if err = o.QueryTable(new(UsedNonce)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllUsedNonce retrieves all UsedNonce matches certain condition. Returns empty list if
// no records exist
func GetAllUsedNonce(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(UsedNonce))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []UsedNonce
	qs = qs.OrderBy(sortFields...).RelatedSel()
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateUsedNonce updates UsedNonce by Id and returns error if
// the record to be updated doesn't exist
func UpdateUsedNonceById(m *UsedNonce) (err error) {
	o := orm.NewOrm()
	v := UsedNonce{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteUsedNonce deletes UsedNonce by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUsedNonce(id int64) (err error) {
	o := orm.NewOrm()
	v := UsedNonce{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&UsedNonce{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func IsUsedNonce(nonce string) bool {
	o := orm.NewOrm()
	isUsed := true
	count, err := o.QueryTable(new(UsedNonce)).Filter("Nonce", nonce).Count()
	if err != nil {
		beego.Error("Check Nonce error | ", err.Error())
		return isUsed
	}
	if count == 0 {
		isUsed = false
	}
	return isUsed
}
