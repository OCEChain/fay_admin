package admin_model

import (
	"errors"
	"github.com/henrylee2cn/faygo"
	"github.com/henrylee2cn/faygo/ext/db/xorm"
	"time"
)

type Role struct {
	Id          int    `xorm:"not null int(11) pk autoincr"`
	Name        string `xorm:"not null default('') varchar(30) comment('角色名')"`
	Status      int    `xorm:"not null default(0) tinyint comment('状态：1正常 0禁用')"`
	Power       string `xorm:"not null default('') varchar(255) comment('角色权限')"`
	Create_time int64  `xorm:"not null default(0) int(11) comment('创建时间')" json:"-"`
	Update_time int64  `xorm:"not null default(0) int(11) comment('修改时间')" json:"-"`
}

type RoleInfo struct {
	*Role
	Create_time string
	Update_time string
}

const (
	RoleTable = "role"
)

var DefaultRole *Role = new(Role)

func init() {
	err := xorm.MustDB().Table(RoleTable).Sync2(DefaultRole)
	if err != nil {
		faygo.Error(err.Error())
	}
}

var (
	NotExistRole  = errors.New("不存在的角色")
	RoleNameExits = errors.New("该角色名已存在")
)

//修改角色的权限
func (r *Role) ModifyPower(id int, cids string) (err error) {
	engine := xorm.MustDB()
	has, err := engine.Table(RoleTable).Where("id=?", id).Exist()
	if err != nil {
		err = SystemFail
		return
	}
	//角色不存在
	if !has {
		err = NotExistRole
		return
	}
	role := new(Role)
	role.Power = cids
	role.Update_time = time.Now().Unix()
	//修改角色权限
	n, err := engine.Where("id=?", id).Cols("power", "update_time").Update(role)
	if err != nil {
		err = SystemFail
		return
	}
	if n == 0 {
		err = errors.New("修改权限失败")
	}
	return
}

//添加一个角色
func (r *Role) Add(rolename string) (err error) {
	if rolename == "" {
		err = ParamEmpty
		return
	}
	engine := xorm.MustDB()
	//检查是否存在
	has, err := engine.Table(RoleTable).Where("name=?", rolename).Exist()
	if err != nil {
		err = SystemFail
		return
	}
	if has {
		err = RoleNameExits
		return
	}
	role := new(Role)
	role.Create_time = time.Now().Unix()
	role.Status = 1
	role.Name = rolename
	n, err := engine.Insert(role)
	if err != nil {
		err = SystemFail
		return
	}
	if n == 0 {
		err = errors.New("添加角色失败")
		return
	}
	return
}

//修改角色名
func (r *Role) EditRoleName(id int, rolename string) (err error) {
	engine := xorm.MustDB()
	has, err := engine.Table(RoleTable).Where("id=?", id).Exist()
	if err != nil {
		err = SystemFail
		return
	}
	//角色不存在
	if !has {
		err = NotExistRole
		return
	}
	role := new(Role)
	role.Name = rolename
	role.Update_time = time.Now().Unix()
	//修改角色名
	n, err := engine.Where("id=?", id).Cols("name", "update_time").Update(role)
	if err != nil {
		err = SystemFail
		return
	}
	if n == 0 {
		err = errors.New("修改权限失败")
	}
	return

}

//修改角色状态
func (r *Role) EditRoleStatus(id int, status int) (err error) {
	engine := xorm.MustDB()
	has, err := engine.Table(RoleTable).Where("id=?", id).Exist()
	if err != nil {
		err = SystemFail
		return
	}
	//角色不存在
	if !has {
		err = NotExistRole
		return
	}
	role := new(Role)
	role.Status = status
	role.Update_time = time.Now().Unix()
	//修改角色名
	n, err := engine.Where("id=?", id).Cols("status", "update_time").Update(role)
	if err != nil {
		err = SystemFail
		return
	}
	if n == 0 {
		err = errors.New("修改权限失败")
	}
	return

}

//删除角色
func (r *Role) DelRole(id int) (err error) {
	if id <= 1 {
		err = errors.New("不能删除超级管理员")
		return
	}
	engine := xorm.MustDB()
	role := new(Role)
	role.Id = id
	n, err := engine.Delete(role)
	if err != nil {
		err = SystemFail
		return
	}
	if n == 0 {
		err = errors.New("删除失败")
	}
	return
}

//根绝角色id获取所有的权限id
func (r *Role) GetPowerByRid(rid int) (role *Role, err error) {
	if rid == 0 {
		err = NotExistRole
		return
	}
	engine := xorm.MustDB()
	role = new(Role)
	has, err := engine.Where("id=?", rid).Get(role)
	if err != nil {
		err = SystemFail
		return
	}
	if !has {
		err = NotExistRole
		return
	}
	return
}

//获取所有的角色
func (r *Role) GetMapRole() (role_list map[int]*Role, err error) {
	engine := xorm.MustDB()
	rows, err := engine.Rows(r)
	if err != nil {
		err = SystemFail
		return
	}
	defer rows.Close()
	role_list = make(map[int]*Role)
	for rows.Next() {
		role := new(Role)
		err = rows.Scan(role)
		if err != nil {
			err = SystemFail
			return
		}
		role_list[role.Id] = role
	}
	return
}

func (r *Role) GetRole() (role_list []*RoleInfo, err error) {
	engine := xorm.MustDB()
	rows, err := engine.Rows(r)
	if err != nil {
		err = SystemFail
		return
	}
	defer rows.Close()
	for rows.Next() {
		role := new(Role)
		roleinfo := new(RoleInfo)
		err = rows.Scan(role)
		if err != nil {
			err = SystemFail
			return
		}
		roleinfo.Role = role
		create_time := time.Unix(roleinfo.Role.Create_time, 0).Format("2006-01-02 15:04:05")
		update_time := create_time
		if roleinfo.Role.Update_time > 0 {
			update_time = time.Unix(roleinfo.Role.Update_time, 0).Format("2006-01-02 15:04:05")
		}

		roleinfo.Create_time = create_time
		roleinfo.Update_time = update_time
		role_list = append(role_list, roleinfo)
	}
	return
}
