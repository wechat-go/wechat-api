package system

import (
	db "zyjsxy-api/database"
	"zyjsxy-api/util"

	"github.com/jinzhu/gorm"
)

type Position struct {
	gorm.Model
	Name string
	Desc string
}

type UPosition struct {
	ID   uint
	Name string
}

func init() {
	if !db.Orm.HasTable(&Position{}) {
		db.Orm.CreateTable(&Position{})
	}
}

//获取所有
func GetAllPosition() (d []UPosition, err error) {
	err = db.Orm.Table("positions").Find(&d).Error
	return d, err
}

//列表
func GetPositionList(d Position, q util.Query) (results map[string]interface{}, err error) {
	results = make(map[string]interface{})
	var positions []Position
	positions = make([]Position, 0)
	var total int
	err = db.Orm.Where("name LIKE ?", "%"+d.Name+"%").Limit(q.Limit).Offset(q.Offset).Find(&positions).Error
	db.Orm.Table("positions").Where("name LIKE ? AND deleted_at IS NULL", "%"+d.Name+"%").Count(&total)

	results["items"] = positions
	results["total"] = total

	return results, err
}

//根据id获取信息
func GetPositioninfo(id int) (p Position, err error) {
	err = db.Orm.First(&p, id).Error
	return p, err
}

//唯一性验证
func PositionNameUinque(s string) (result bool, err error) {
	var count int
	err = db.Orm.Model(&Position{}).Where("name = ?", s).Count(&count).Error
	if count == 0 {
		return true, err
	} else {
		return false, err
	}
}

//添加
func AddPosition(d Position) (err error) {
	return db.Orm.Create(&d).Error
}

//更新
func UpdataPosition(o Position) (err error) {
	return db.Orm.Model(&o).Updates(Position{Name: o.Name, Desc: o.Desc}).Error
}

//删除
func DeletePosition(sa []string) (err error) {
	err = db.Orm.Where("id in (?)", sa).Delete(&Position{}).Error
	return err
}
