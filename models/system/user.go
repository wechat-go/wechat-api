package system

import (
	"strconv"
	"time"
	db "zyjsxy-api/database"
	"zyjsxy-api/util"
	"zyjsxy-api/util/aes"

	"github.com/jinzhu/gorm"
)

type User struct {
	uid         uint `gorm:"primary_key"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
	DeletedAt   *time.Time `sql:"index"`
	Username    string
	Password    string
	Avatar      string    `gorm:"default:'/img/user.jpg'"`
	Roles       []Role    `gorm:"many2many:user_role;save_associations:false"`
	Realname    string    `form:"Realname" json:"Realname" binding:"required"`
	Sex         int       `form:"Sex" json:"Sex,string"`
	Birth       time.Time `form:"Birth" json:"Birth" time_format:"2006-01-02"`
	Email       string    `form:"Email" json:"Email"`
	Webchat     string    `form:"Webchat" json:"Webchat"`
	Qq          string    `form:"Qq" json:"Qq"`
	Phone       string    `form:"Phone" json:"Phone"`
	Tel         string    `form:"Tel" json:"Tel"`
	Address     string    `form:"Address" json:"Address"`
	Emercontact string    `form:"Emercontact" json:"Emercontact"`
	Emerphone   string    `form:"Emerphone" json:"Emerphone"`
	Departid    int64
	Positionid  int64
	Lognum      int
	Ip          string
	Lasted      int64
}

type UserProfile struct {
	ID          int
	Username    string
	Avatar      string
	Realname    string
	Sex         int
	Birth       time.Time
	Email       string
	Webchat     string
	Qq          string
	Phone       string
	Tel         string
	Address     string
	Emercontact string
	Emerphone   string
	Departid    int64
	Positionid  int64
}

type UserInfo struct {
	ID        uint
	Avatar    string
	Profile   ProfileTemp
	ProfileID uint
	Roles     []uint
}

type ProfileTemp struct {
	ID          uint
	Realname    string    `form:"Realname" json:"Realname" binding:"required"`
	Sex         int       `form:"Sex" json:"Sex,string"`
	Birth       time.Time `form:"Birth" json:"Birth" time_format:"2006-01-02"`
	Email       string    `form:"Email" json:"Email"`
	Webchat     string    `form:"Webchat" json:"Webchat"`
	Qq          string    `form:"Qq" json:"Qq"`
	Phone       string    `form:"Phone" json:"Phone"`
	Tel         string    `form:"Tel" json:"Tel"`
	Address     string    `form:"Address" json:"Address"`
	Emercontact string    `form:"Emercontact" json:"Emercontact"`
	Emerphone   string    `form:"Emerphone" json:"Emerphone"`
	Departid    int64
	Positionid  int64
}

type PostUserInfo struct {
	Username string
	Password string
	UserInfo
}

func init() {
	if !db.Orm.HasTable(&User{}) {
		db.Orm.CreateTable(&User{})
	}
}

func AddUser(user *User) (err error) {
	err = db.Orm.Create(user).Error
	return err
}

func AddUserRole(user PostUserInfo) (err error) {
	var u User
	u.Username = user.Username
	u.Password = aes.Hashsalt(user.Password)
	u.Avatar = user.Avatar

	u.Profile.Realname = user.Profile.Realname
	u.Profile.Sex = user.Profile.Sex
	u.Profile.Birth = user.Profile.Birth
	u.Profile.Email = user.Profile.Email
	u.Profile.Webchat = user.Profile.Webchat
	u.Profile.Qq = user.Profile.Qq
	u.Profile.Phone = user.Profile.Phone
	u.Profile.Tel = user.Profile.Tel
	u.Profile.Address = user.Profile.Address
	u.Profile.Emercontact = user.Profile.Emercontact
	u.Profile.Emerphone = user.Profile.Emerphone
	u.Profile.Departid = user.Profile.Departid
	u.Profile.Positionid = user.Profile.Positionid

	err = db.Orm.Create(&u).Error
	if err == nil {
		str := ""
		for index, obj := range user.Roles {
			if index == 0 {
				str += "(" + strconv.Itoa(int(u.ID)) + "," + strconv.Itoa(int(obj)) + ")"
			} else {
				str += ",(" + strconv.Itoa(int(u.ID)) + "," + strconv.Itoa(int(obj)) + ")"
			}
		}
		err = db.Orm.Exec("INSERT INTO user_role (user_id,role_id)VALUES " + str).Error
	}
	return err
}

func GetSelfInfo(id int) (u User, err error) {
	err = db.Orm.Preload("Profile").First(&u, id).Error
	ids, err := getRoleIds(id)
	roles, err := getRoles(ids)
	u.Roles = roles
	return u, err
}

func GetUserInfo(id int) (u UserInfo, err error) {
	var user User
	err = db.Orm.Preload("Profile").Preload("Roles").First(&user, id).Error
	intarr := make([]uint, 0)

	for _, obj := range user.Roles {
		intarr = append(intarr, obj.ID)
	}
	u.ID = user.ID
	u.Avatar = user.Avatar
	u.ProfileID = user.ProfileID
	u.Roles = intarr

	u.Profile.ID = user.Profile.ID
	u.Profile.Realname = user.Profile.Realname
	u.Profile.Sex = user.Profile.Sex
	u.Profile.Birth = user.Profile.Birth
	u.Profile.Email = user.Profile.Email
	u.Profile.Webchat = user.Profile.Webchat
	u.Profile.Qq = user.Profile.Qq
	u.Profile.Phone = user.Profile.Phone
	u.Profile.Tel = user.Profile.Tel
	u.Profile.Address = user.Profile.Address
	u.Profile.Emercontact = user.Profile.Emercontact
	u.Profile.Emerphone = user.Profile.Emerphone
	u.Profile.Departid = user.Profile.Departid
	u.Profile.Positionid = user.Profile.Positionid
	return u, err
}

func (u *User) UpdateUser() (err error) {
	err = db.Orm.Model(&u).Updates(map[string]interface{}{
		"name": u.Username,
		"age":  u.Password,
	}).Error
	return err
}

func DeleteUser(sa []string) (err error) {
	tx := db.Orm.Begin()
	if err = db.Orm.Where("id in (?)", sa).Delete(&User{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err = db.Orm.Where("id in (?)", sa).Delete(&Profile{}).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Exec("DELETE FROM user_role WHERE user_id in (?)", sa).Error; err != nil {
		tx.Rollback()
		return err
	}

	tx.Commit()
	return nil
}

func UpdataAvatar(u User) (err error) {
	err = db.Orm.Model(&u).Updates(map[string]interface{}{
		"avatar": u.Avatar,
	}).Error
	return err
}

func GetPassword(id int) string {
	var u User
	_ = db.Orm.First(&u, id).Error
	return u.Password
}

func UpdataPassword(u User) (err error) {
	err = db.Orm.Model(&u).Updates(map[string]interface{}{
		"password": u.Password,
	}).Error
	return err
}

func FetchUserList(u User, o util.Query) (results map[string]interface{}, err error) {
	var sql string
	if u.Username != "" {
		sql += "users.username='" + u.Username + "' AND "
	}
	if u.Profile.Realname != "" {
		sql += "profiles.realname='" + u.Profile.Realname + "' AND "
	}
	if u.Profile.Phone != "" {
		sql += "profiles.phone='" + u.Profile.Phone + "' AND "
	}
	if u.Profile.Sex != 0 {
		sql += "profiles.sex=" + strconv.Itoa(u.Profile.Sex) + " AND "
	}

	sql = sql + "users.deleted_at is NULL"
	results = make(map[string]interface{})
	var userprofile []UserProfile
	userprofile = make([]UserProfile, 0)
	var total int
	err = db.Orm.Table("users").
		Select("users.id,users.username,users.avatar,profiles.realname,profiles.sex,profiles.email,profiles.webchat,profiles.qq,profiles.phone,profiles.tel,profiles.address,profiles.emercontact,profiles.emerphone,profiles.departid,profiles.positionid").
		Joins("left join profiles on users.profile_id = profiles.id").
		Where(sql).
		Count(&total).
		Limit(o.Limit).Offset(o.Offset).
		Scan(&userprofile).Error
	results["items"] = userprofile
	results["total"] = total

	return results, err
}

func UpdataUserInfo(u User) (err error) {
	err = db.Orm.Save(&u).Error
	return err
}

func UserNameUinque(s string) (result bool, err error) {
	var count int
	err = db.Orm.Model(&User{}).Where("username = ?", s).Count(&count).Error
	if count == 0 {
		return true, err
	} else {
		return false, err
	}
}

func UpdataUser(user PostUserInfo) (err error) {
	var p Profile

	p.ID = user.Profile.ID
	p.Realname = user.Profile.Realname
	p.Sex = user.Profile.Sex
	p.Birth = user.Profile.Birth
	p.Email = user.Profile.Email
	p.Webchat = user.Profile.Webchat
	p.Qq = user.Profile.Qq
	p.Phone = user.Profile.Phone
	p.Tel = user.Profile.Tel
	p.Address = user.Profile.Address
	p.Emercontact = user.Profile.Emercontact
	p.Emerphone = user.Profile.Emerphone
	p.Departid = user.Profile.Departid
	p.Positionid = user.Profile.Positionid
	str := ""
	for index, obj := range user.Roles {
		if index == 0 {
			str += "(" + strconv.Itoa(int(user.ID)) + "," + strconv.Itoa(int(obj)) + ")"
		} else {
			str += ",(" + strconv.Itoa(int(user.ID)) + "," + strconv.Itoa(int(obj)) + ")"
		}
	}
	tx := db.Orm.Begin()
	if err = tx.Model(&p).
		Updates(Profile{Realname: p.Realname,
			Sex:         p.Sex,
			Birth:       p.Birth,
			Email:       p.Email,
			Webchat:     p.Webchat,
			Qq:          p.Qq,
			Phone:       p.Phone,
			Tel:         p.Tel,
			Address:     p.Address,
			Emercontact: p.Emercontact,
			Departid:    p.Departid,
			Positionid:  p.Positionid,
			Emerphone:   p.Emerphone}).
		Error; err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Exec("DELETE FROM user_role WHERE user_id = " + strconv.Itoa(int(user.ID))).Error; err != nil {
		tx.Rollback()
		return err
	}

	if err = tx.Exec("INSERT INTO user_role (user_id,role_id)VALUES " + str).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func GetUserFrist(user User) (u User, err error) {
	err = db.Orm.Where(&user).First(&u).Error
	return u, err
}
