package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"time"
)

type User struct {
	Id                int64     `orm:"pk;auto"`
	Created           time.Time `orm:"auto_now_add;type(datetime)"`
	Start             time.Time `orm:"auto_now_add;type(datetime)"`
	Expired           time.Time `orm:"auto_now_add;type(datetime)"`
	Updated           time.Time `orm:"auto_now;type(datetime)"`
	Enabled           bool      `orm:"null;default(true)"`
	Username          string    `orm:"null;size(128)"`
	PasswordNoEncrypt string    `orm:"null;size(128)"`
	Password          string    `orm:"null;size(128)"`
	Email             string    `orm:"null;size(128)"`
	Bronze            int       `orm:"null"`
	Silver            int       `orm:"null"`
	Gold              int       `orm:"null"`
	Image             string    `orm:"null;size(128)"`
	IsAgent           bool
	AgentName         string `orm:"null;size(128)"`
	AgentEmail        string `orm:"null;size(128)"`
	AgentImage        string `orm:"null;size(128)"`
	AgentLineId       string `orm:"null;size(128)"`
	AgentUrlFacebook  string `orm:"null;size(128)"`
	AgentWebsite      string `orm:"null;size(128)"`
	AgentResume       string `orm:"null;size(128)"`
	Verified          bool   `orm:"default(false)"`

	FacebookId           string                 `orm:"null;size(50)"`
	FirstName            string                 `orm:"null;size(255)"`
	LastName             string                 `orm:"null;size(255)"`
	DisplayName          string                 `orm:"null;size(255)"`
	Career               string                 `orm:"null;size(255)"`
	MobileNumber         string                 `orm:"null;size(255)"`
	PhoneNumber          string                 `orm:"null;size(255)"`
	Address              string                 `orm:"null;size(255)"`
	Province             string                 `orm:"null;size(255)"`
	District             string                 `orm:"null;size(255)"`
	SubDistrict          string                 `orm:"null;size(255)"`
	Gender               string                 `orm:"null;size(50)"`
	BirthDate            time.Time              `orm:"null;type(datetime)"`
	RegisterDate         time.Time              `orm:"type(datetime)"`
	AuthenticationTokens []*AuthenticationToken `orm:"reverse(many)"`
}

func init() {
	orm.RegisterModel(new(User))
}

// AddUser insert a new User into database and returns
// last inserted Id on success.
func AddUser(m *User) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetUserById retrieves User by Id. Returns error if
// Id doesn't exist
func GetUserById(id int64) (v *User, err error) {
	o := orm.NewOrm()
	v = &User{Id: id}
	if err = o.QueryTable(new(User)).Filter("Id", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllUser retrieves all User matches certain condition. Returns empty list if
// no records exist
func GetAllUser(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(User))
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

	var l []User
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

// UpdateUser updates User by Id and returns error if
// the record to be updated doesn't exist
func UpdateUserById(m *User) (err error) {
	o := orm.NewOrm()
	v := User{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteUser deletes User by Id and returns error if
// the record to be deleted doesn't exist
func DeleteUser(id int64) (err error) {
	o := orm.NewOrm()
	v := User{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&User{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

func GetUserByReferrerToken(referrerToken string) *User {
	o := orm.NewOrm()
	v := new(User)
	beego.Debug(referrerToken)
	err := o.QueryTable(v).Filter("ReferrerToken", referrerToken).Limit(1).One(v)
	if err != nil {
		return nil
	}

	return v
}

func GetFacebookUserByEmail(email string) *User {
	o := orm.NewOrm()
	user := new(User)
	err := o.QueryTable(new(User)).Filter("Email", email).Filter("FacebookId__isnull", false).Limit(1).One(user)
	if err != nil {
		return nil
	}
	return user
}

func GetUserByFacebookId(facebookId string) *User {
	o := orm.NewOrm()
	user := new(User)
	err := o.QueryTable(new(User)).Filter("FacebookId", facebookId).Limit(1).One(user)
	if err != nil {
		return nil
	}
	return user
}

func GetUserByUsername(username string) *User {
	o := orm.NewOrm()
	user := new(User)
	err := o.QueryTable(new(User)).Filter("Username", username).Limit(1).One(user)
	if err != nil {
		return nil
	}
	return user
}

func IsUsernameAvailable(username string) bool {
	o := orm.NewOrm()
	var isAvailable = false
	user := &User{}
	err := o.QueryTable(new(User)).Filter("Username", username).Limit(1).One(user)
	if err != nil {
		//panic("errrrrrr: " + err.Error())
		fmt.Println("IsUsernameAvailable error: ", err.Error())
		isAvailable = true

	}
	if user == nil {
		isAvailable = true
	} else {
		// check registerDate < 7 days
		var expiryTime int64 = 60 * 60 * 24 * 7
		now := time.Now().Unix()
		timestamp := user.RegisterDate.Unix()
		expires := timestamp + expiryTime
		if user.Verified == false && expires-now < 0 {
			isAvailable = true
		} else if user.FacebookId != "" && user.Verified == false {
			isAvailable = true
		}
	}
	return isAvailable
}
