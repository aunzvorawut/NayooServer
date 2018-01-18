package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
	"time"
	"github.com/astaxie/beego"
)

type Banner struct {
	Id            int64     `orm:"pk;auto"`
	Created       time.Time `orm:"auto_now_add;type(datetime)"`
	Updated       time.Time `orm:"auto_now;type(datetime)"`
	Start         time.Time `orm:"null;type(datetime)"`
	Expired       time.Time `orm:"null;type(datetime)"`
	Enabled       bool      `orm:"null;default(true)"`
	OrderPosition int
	Grade         string `orm:"null;size(255)"`
	Image         string `orm:"null;size(255)"`
	BannerType      string `orm:"null;size(255)"`
}

func init() {
	orm.RegisterModel(new(Banner))
}

// AddBanner insert a new Banner into database and returns
// last inserted Id on success.
func AddBanner(m *Banner) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetBannerById retrieves Banner by Id. Returns error if
// Id doesn't exist
func GetBannerById(id int64) (v *Banner, err error) {
	o := orm.NewOrm()
	v = &Banner{Id: id}
	if err = o.QueryTable(new(Banner)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllBanner retrieves all Banner matches certain condition. Returns empty list if
// no records exist
func GetAllBanner(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Banner))
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

	var l []Banner
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

// UpdateBanner updates Banner by Id and returns error if
// the record to be updated doesn't exist
func UpdateBannerById(m *Banner) (err error) {
	o := orm.NewOrm()
	v := Banner{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteBanner deletes Banner by Id and returns error if
// the record to be deleted doesn't exist
func DeleteBanner(id int64) (err error) {
	o := orm.NewOrm()
	v := Banner{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Banner{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetAllBannerByBannerTypeAndEnabledAndStartAndExpired(bannerType string) (result []*Banner) {

	stringQuery := " select id as id , image as image from banner where banner_type = ? "

	beego.Debug(stringQuery)
	o := orm.NewOrm()
	_, err := o.Raw(stringQuery, bannerType).QueryRows(&result)
	if err != nil {
		beego.Error(err.Error())
	}
	return result

}