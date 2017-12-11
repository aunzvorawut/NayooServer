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

type AuthenticationToken struct {
	Id          int64     `orm:"pk;auto"`
	AccessToken string    `orm:"size(128)"`
	User        *User     `orm:"rel(fk)"`
	ExpiredDate time.Time `orm:"null;type(datetime)"`
	Created     time.Time `orm:"auto_now_add;type(datetime)"`
	Updated     time.Time `orm:"auto_now;type(datetime)"`
}

func init() {
	orm.RegisterModel(new(AuthenticationToken))
}

// AddAuthenticationToken insert a new AuthenticationToken into database and returns
// last inserted Id on success.
func AddAuthenticationToken(m *AuthenticationToken) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetAuthenticationTokenById retrieves AuthenticationToken by Id. Returns error if
// Id doesn't exist
func GetAuthenticationTokenById(id int64) (v *AuthenticationToken, err error) {
	o := orm.NewOrm()
	v = &AuthenticationToken{Id: id}
	if err = o.QueryTable(new(AuthenticationToken)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllAuthenticationToken retrieves all AuthenticationToken matches certain condition. Returns empty list if
// no records exist
func GetAllAuthenticationToken(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(AuthenticationToken))
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

	var l []AuthenticationToken
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

// UpdateAuthenticationToken updates AuthenticationToken by Id and returns error if
// the record to be updated doesn't exist
func UpdateAuthenticationTokenById(m *AuthenticationToken) (err error) {
	o := orm.NewOrm()
	v := AuthenticationToken{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteAuthenticationToken deletes AuthenticationToken by Id and returns error if
// the record to be deleted doesn't exist
func DeleteAuthenticationToken(id int64) (err error) {
	o := orm.NewOrm()
	v := AuthenticationToken{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&AuthenticationToken{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetAuthToken(accessToken string) *AuthenticationToken {
	o := orm.NewOrm()
	auth := new(AuthenticationToken)
	err := o.QueryTable(new(AuthenticationToken)).Filter("AccessToken", accessToken).RelatedSel().Limit(1).One(auth)
	if err != nil {
		return nil
	}
	return auth
}

func GetAuthTokenByUserId(userId int64) *AuthenticationToken {
	o := orm.NewOrm()
	auth := new(AuthenticationToken)
	err := o.QueryTable(new(AuthenticationToken)).Filter("User__Id", userId).Limit(1).One(auth)
	if err != nil {
		return nil
	}
	return auth
}

func GetUserByAccessToken(accessToken string) *User {
	o := orm.NewOrm()
	auth := new(AuthenticationToken)
	err := o.QueryTable(new(AuthenticationToken)).Filter("AccessToken", accessToken).RelatedSel("User").Limit(1).One(auth)
	if err != nil {
		return nil
	}

	user := auth.User

	return user
}

func GetAuthTokensByUser(user User) []*AuthenticationToken {
	o := orm.NewOrm()
	var auths []*AuthenticationToken
	_, err := o.QueryTable(new(AuthenticationToken)).Filter("User__Id", user.Id).RelatedSel("Device").OrderBy("-Created").Limit(10).All(&auths)
	if err != nil {
		beego.Error(err)
		return auths
	}
	return auths
}