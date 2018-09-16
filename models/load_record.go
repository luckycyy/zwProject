package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/astaxie/beego/orm"
)

type LoadRecord struct {
	Id          int       `orm:"column(load_record_id);auto"`
	ProductName string    `orm:"column(product_name);size(128);null"`
	Num         float64   `orm:"column(num);null;digits(10);decimals(2)" json:"Num,string"`
	Station     string    `orm:"column(station);size(255);null"`
	Location    string    `orm:"column(location);size(255);null"`
	CreateTime  time.Time `orm:"column(create_time);type(datetime);null;auto_now_add"`
	Creator     string    `orm:"column(creator);size(128);null"`
	IsUnload    int       `orm:"column(is_unload);null"`
	Location2   string    `orm:"column(location2);size(255);null"`
}

func (t *LoadRecord) TableName() string {
	return "load_record"
}

func init() {
	orm.RegisterModel(new(LoadRecord))
}

// AddLoadRecord insert a new LoadRecord into database and returns
// last inserted Id on success.
func AddLoadRecord(m *LoadRecord) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetLoadRecordById retrieves LoadRecord by Id. Returns error if
// Id doesn't exist
func GetLoadRecordById(id int) (v *LoadRecord, err error) {
	o := orm.NewOrm()
	v = &LoadRecord{Id: id}
	if err = o.Read(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllLoadRecord retrieves all LoadRecord matches certain condition. Returns empty list if
// no records exist
func GetAllLoadRecord(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(LoadRecord))
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

	var l []LoadRecord
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

// UpdateLoadRecord updates LoadRecord by Id and returns error if
// the record to be updated doesn't exist
func UpdateLoadRecordById(m *LoadRecord) (err error) {
	o := orm.NewOrm()
	v := LoadRecord{Id: m.Id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		//注意！！！！！！！此处指更新指定的字段了！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！！
		if num, err = o.Update(m, "product_name", "num", "station"); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteLoadRecord deletes LoadRecord by Id and returns error if
// the record to be deleted doesn't exist
func DeleteLoadRecord(id int) (err error) {
	o := orm.NewOrm()
	v := LoadRecord{Id: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&LoadRecord{Id: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
