package system

import (
	db "zyjsxy-api/database"
	"zyjsxy-api/util"

	"github.com/jinzhu/gorm"
)

type Depart struct {
	gorm.Model
	Name     string
	Desc     string
	ParentId int
}

type UDepart struct {
	ID       uint
	Name     string
	ParentId int
}

func init() {
	if !db.Orm.HasTable(&Depart{}) {
		db.Orm.CreateTable(&Depart{})
	}
}

func GetDepartList(d Depart, q util.Query) (results map[string]interface{}, err error) {
	results = make(map[string]interface{})
	var departs []Depart
	departs = make([]Depart, 0)
	var total int
	err = db.Orm.Where("name LIKE ?", "%"+d.Name+"%").Limit(q.Limit).Offset(q.Offset).Find(&departs).Error
	db.Orm.Table("departs").Where("name LIKE ? AND deleted_at IS NULL", "%"+d.Name+"%").Count(&total)

	results["items"] = departs
	results["total"] = total

	return results, err
}

//获取所有部门，选父级所用
func GetAllDepart() (d []UDepart, err error) {
	err = db.Orm.Table("departs").Find(&d).Error
	return d, err
}

//根据id获取部门信息
func GetDepartinfo(id int) (p Depart, err error) {
	err = db.Orm.First(&p, id).Error
	return p, err
}

//唯一性验证
func DepartNameUinque(s string) (result bool, err error) {
	var count int
	err = db.Orm.Model(&Depart{}).Where("name = ?", s).Count(&count).Error
	if count == 0 {
		return true, err
	} else {
		return false, err
	}
}

//添加部门
func AddDepart(d Depart) (err error) {
	return db.Orm.Create(&d).Error
}

//更新部门
func UpdataDepart(o Depart) (err error) {
	return db.Orm.Model(&o).Updates(Depart{Name: o.Name, Desc: o.Desc, ParentId: o.ParentId}).Error
}

//删除部门
func DeleteDepart(sa []string) (err error) {
	err = db.Orm.Where("id in (?)", sa).Delete(&Depart{}).Error
	return err
}
