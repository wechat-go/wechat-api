package system

import (
	"strconv"
	//	"fmt"
	db "zyjsxy-api/database"
	"zyjsxy-api/util"

	"github.com/jinzhu/gorm"
)

type Privilege struct {
	gorm.Model
	Name       string
	Path       string
	ParentId   int
	CreateUser int
	Roles      []Role `gorm:"many2many:privilege_role"`
}

type UserRole struct {
	UserId int
	RoleId int
}

type PrivilegeRole struct {
	PrivilegeId int
	RoleId      int
}

type Privilegesingle struct {
	ID       uint
	Name     string
	Path     string
	ParentId int
}

func init() {
	if !db.Orm.HasTable(&Privilege{}) {
		db.Orm.CreateTable(&Privilege{})
	}
}
func GetUserPrivilege(id int) (privileges []Privilege, err error) {
	ids, err := getRoleIds(id)
	err = db.Orm.Table("privilege_role").
		Select("DISTINCT `privileges`.id,`privileges`.`name`,`privileges`.parent_id,`privileges`.path,`privileges`.create_user").
		Joins("LEFT JOIN `privileges` ON privilege_role.privilege_id = `privileges`.id").
		Where("privilege_role.role_id IN (?)", ids).
		Order("`privileges`.id").Scan(&privileges).Error
	if err != nil {
		return nil, err
	}
	return privileges, err
}

func getRoleIds(id int) (ids string, err error) {
	userroles := make([]PrivilegeRole, 0)
	err = db.Orm.Table("user_role").Select("user_id,role_id").Where("user_id=?", id).Scan(&userroles).Error
	if err != nil {
		return "", err
	}
	for index, value := range userroles {
		if index == 0 {
			ids = strconv.Itoa(value.RoleId)
		} else {
			ids = ids + "," + strconv.Itoa(value.RoleId)
		}
	}
	return ids, err
}

func getRoles(ids string) (roles []Role, err error) {
	err = db.Orm.Table("roles").Where("id IN (?)", ids).Scan(&roles).Error
	return roles, err
}

func GetAllPrivilege() (p []Privilegesingle, err error) {
	err = db.Orm.Table("privileges").Find(&p).Error
	return p, err
}

func GetRolePrivilege(id int) (p []int, err error) {
	pr := make([]PrivilegeRole, 0)
	err = db.Orm.Table("privilege_role").Select("privilege_id").Where("role_id=?", id).Scan(&pr).Error
	for _, obj := range pr {
		p = append(p, obj.PrivilegeId)
	}
	return p, err
}

func GetPrivilegeList(p Privilege, q util.Query) (results map[string]interface{}, err error) {
	results = make(map[string]interface{})
	var privileges []Privilege
	privileges = make([]Privilege, 0)
	var total int
	err = db.Orm.Where("name LIKE ?", "%"+p.Name+"%").Limit(q.Limit).Offset(q.Offset).Find(&privileges).Error
	db.Orm.Table("privileges").Where("name LIKE ? AND deleted_at IS NULL", "%"+p.Name+"%").Count(&total)

	results["items"] = privileges
	results["total"] = total

	return results, err

}

func UpdataPrivilege(p Privilege) (err error) {
	return db.Orm.Model(&p).Updates(Privilege{Name: p.Name, Path: p.Path, ParentId: p.ParentId}).Error
}

func AddPrivilege(p Privilege) (err error) {
	return db.Orm.Create(&p).Error
}

func GetPrivilegeinfo(id int) (p Privilege, err error) {
	err = db.Orm.First(&p, id).Error
	return p, err
}

func DeletePrivilege(sa []string) (err error) {
	err = db.Orm.Where("id in (?)", sa).Delete(&Privilege{}).Error
	return err
}
