package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego/orm"
	"time"
)

type HouseSale struct {
	Id      int64     `orm:"pk;auto"`
	Created time.Time `orm:"auto_now_add;type(datetime)"`
	Updated time.Time `orm:"auto_now;type(datetime)"`
	Start   time.Time `orm:"null;type(datetime)"`
	Expired time.Time `orm:"null;type(datetime)"`
	Enabled bool      `orm:"null;default(true)"`
	VipType string    `orm:"null;size(255)"` // gold , silver , bronze
	Image   string    `orm:"null;size(255)"`

	TitleTh         string `orm:"null;size(255)"`
	TitleEng        string `orm:"null;size(255)"`
	ResidentTypeTh  string `orm:"null;size(255)"` // ขาย,เช่า,ขาย/เช่า
	ResidentTypeEng string `orm:"null;size(255)"` // sale,rent,sale/rent

	UserFavorites []*User   `orm:"reverse(many)"`
}

func init() {
	orm.RegisterModel(new(HouseSale))
}

// AddHouseSale insert a new HouseSale into database and returns
// last inserted Id on success.
func AddHouseSale(m *HouseSale) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetHouseSaleById retrieves HouseSale by Id. Returns error if
// Id doesn't exist
func GetHouseSaleById(id int64) (v *HouseSale, err error) {
	o := orm.NewOrm()
	v = &HouseSale{Id: id}
	if err = o.QueryTable(new(HouseSale)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllHouseSale retrieves all HouseSale matches certain condition. Returns empty list if
// no records exist
func GetAllHouseSale(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(HouseSale))
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

	var l []HouseSale
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

// UpdateHouseSale updates HouseSale by Id and returns error if
// the record to be updated doesn't exist
func UpdateHouseSaleById(m *HouseSale) (err error) {
	o := orm.NewOrm()
	v := HouseSale{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteHouseSale deletes HouseSale by Id and returns error if
// the record to be deleted doesn't exist
func DeleteHouseSale(id int64) (err error) {
	o := orm.NewOrm()
	v := HouseSale{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&HouseSale{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetAllHouseSaleOnClientByEnabledAndStartAndExpired(max, offset int) (ml []*HouseSale, count int64) {

	currentTime := time.Now()

	o := orm.NewOrm()
	qs := o.QueryTable(new(HouseSale))
	qs.Filter("start", currentTime).Filter("expired__gt", currentTime).Filter("enabled", true).RelatedSel()
	count, _ = qs.Count()

	if _, err := qs.Limit(max, offset).All(&ml); err == nil {
		return ml, count
	}

	return nil, count

}