package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type Apply struct {
	Id         int       `orm:"column(apply_id);auto"`
	Openid     string    `orm:"column(openid);size(32)"`
	Username   string    `orm:"column(username);size(32);null"`
	Tel        string    `orm:"column(tel);size(32);null"`
	Company    string    `orm:"column(company);size(128);null"`
	Status     int       `orm:"column(status);null"`
	Operator   string    `orm:"column(operator);size(32);null"`
	CreateTime time.Time `orm:"column(create_time);type(datetime);null;auto_now_add"`
	UpdateTime time.Time `orm:"column(update_time);type(datetime);null"`
}

func (t *Apply) TableName() string {
	return "apply"
}

func init() {
	orm.RegisterModel(new(Apply))
}

// AddApply insert a new Apply into database and returns
// last inserted Id on success.
func AddApply(m *Apply) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetApplyById retrieves Apply by Id. Returns error if
// Id doesn't exist
func GetApplyById(id int) (v *Apply, err error) {
	o := orm.NewOrm()
	v = &Apply{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllApply retrieves all Apply matches certain condition. Returns empty list if
// no records exist
func GetAllApply(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Apply))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		if strings.Contains(k, "isnull") {
			qs = qs.Filter(k, (v == "true" || v == "1"))
		} else {
			qs = qs.Filter(k, v)
		}
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

	var l []Apply
	qs = qs.OrderBy(sortFields...)
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

// UpdateApply updates Apply by Id and returns error if
// the record to be updapply.go
//user.goated doesn't exist
func UpdateApplyById(m *Apply) (err error) {
	o := orm.NewOrm()
	v := Apply{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteApply deletes Apply by Id and returns error if
// the record to be deleted doesn't exist
func DeleteApply(id int) (err error) {
	o := orm.NewOrm()
	v := Apply{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Apply{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}

// 根据openid删除
func DeleteApplyByOpenid(openid string) (err error) {
	o := orm.NewOrm()
	if num, err := o.QueryTable("apply").Filter("openid", openid).Delete(); err == nil {
		fmt.Println("delete apply ,Number of records deleted in database:", num)
	}

	return
}
