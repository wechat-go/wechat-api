package system

import (
	db "zyjsxy-api/database"

	"github.com/jinzhu/gorm"
)

type Fs struct {
	gorm.Model
	Src       string
	Name      string
	Uid       int
	UpdataUid int
}

func init() {
	if !db.Orm.HasTable(&Fs{}) {
		db.Orm.CreateTable(&Fs{})
	}
}

//根据id获取附件列表
func GetFsList(q []uint) (results []Fs, err error) {
	//根据主键获取附件列表
	err = db.Orm.Where(q).Find(&results).Error

	return results, err
}

//添加附件
func AddFs(d *Fs) (err error) {
	return db.Orm.Create(d).Error
}

//更新附件
func UpdataFs(o Fs) (err error) {
	return db.Orm.Model(&o).Updates(Fs{Name: o.Name, Src: o.Src, UpdataUid: o.UpdataUid}).Error
}

//删除附件
func DeleteFs(sa []string) (err error) {
	err = db.Orm.Where("id in (?)", sa).Delete(&Fs{}).Error
	return err
}
