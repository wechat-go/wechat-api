package system

import (
	//	"fmt"
	db "zyjsxy-api/database"

	"zyjsxy-api/util"

	"github.com/jinzhu/gorm"
)

type Role struct {
	gorm.Model
	Name        string
	Chinesename string
	CreateUser  int
	UpdateUser  int
}

func init() {
	if !db.Orm.HasTable(&Role{}) {
		db.Orm.CreateTable(&Role{})
	}
}

func FetchRoleList() (roles []Role, err error) {
	roles = make([]Role, 0)
	err = db.Orm.Find(&roles).Error
	return roles, err
}

func GetRoleList(r Role, q util.Query) (results map[string]interface{}, err error) {
	results = make(map[string]interface{})
	var roles []Role
	roles = make([]Role, 0)
	var total int
	err = db.Orm.Where(&Role{Name: r.Name, Chinesename: r.Chinesename}).Limit(q.Limit).Offset(q.Offset).Find(&roles).Error
	db.Orm.Table("roles").Where(&Role{Name: r.Name, Chinesename: r.Chinesename}).Where("deleted_at IS NULL").Count(&total)

	results["items"] = roles
	results["total"] = total

	return results, err

}

func DeleteRole(sa []string) (err error) {
	tx := db.Orm.Begin()
	if err = db.Orm.Where("id in (?)", sa).Delete(&Role{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Exec("DELETE FROM privilege_role WHERE role_id in (?)", sa).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func GetRoleinfo(id int) (r Role, err error) {
	err = db.Orm.First(&r, id).Error
	return r, err
}

func UpdataRole(r Role) (err error) {
	return db.Orm.Model(&r).Updates(Role{Name: r.Name, Chinesename: r.Chinesename, UpdateUser: r.UpdateUser}).Error
}

func AddRole(r Role) (err error) {
	return db.Orm.Create(&r).Error
}

func UpdateRolePrivilege(id string, ids []string) (err error) {
	str := ""
	for index, obj := range ids {
		if index == 0 {
			str += "(" + id + "," + obj + ")"
		} else {
			str += ",(" + id + "," + obj + ")"
		}
	}

	tx := db.Orm.Begin()

	if err = tx.Exec("DELETE FROM privilege_role WHERE role_id = ?", id).Error; err != nil {
		tx.Rollback()
		return err
	}
	if str != "" {
		if err = tx.Exec("INSERT INTO privilege_role (role_id,privilege_id)VALUES " + str).Error; err != nil {
			tx.Rollback()
			return err
		}
	}
	tx.Commit()
	return nil
}
