package admin_model

import (
	. "fay_admin/admin_common"
	"github.com/go-errors/errors"
	"github.com/henrylee2cn/faygo"
	"github.com/henrylee2cn/faygo/ext/db/xorm"
	"time"
)

type Admin struct {
	Id          int    `xorm:"not null INT(11) pk autoincr"`
	User        string `xorm:"not null unique default('') VARCHAR(50) comment('用户名')"`
	Pass        string `xorm:"not null default('') char(32) comment('密码')" json:"-"`
	Nickname    string `xorm:"not null default('') VARCHAR(30) comment('昵称')"`
	Face        string `xorm:"not null default('') varchar(255) comment('用户头像')"`
	Role_id     int    `xorm:"not null default(0) int comment('角色id')"`
	Status      int    `xorm:"not null default(0) tinyint comment('状态：1可用 0禁用')"`
	Create_time int64  `xorm:"not null default(0) int comment('注册时间')" json:"-"`
	Login_ip    string `xorm:"not null default('') char(15) comment('最后一次登陆的ip')"`
	Login_time  int64  `xorm:"not null default(0) int comment('最后一次登陆的时间')" json:"-"`
}

type User struct {
	*Admin
	Role_name   string
	Create_time string
	Login_time  string
}

const (
	Admin_TABLE = "admin"
)

var DefaultAdmin *Admin = new(Admin)

func init() {
	err := xorm.MustDB().Table(Admin_TABLE).Sync2(DefaultAdmin)
	if err != nil {
		faygo.Error(err.Error())
	}
}

var (
	NotExistAdmin = errors.New("不存在的用户")
	PassError     = errors.New("密码错误")
	AdminDisable  = errors.New("该账号已被禁用")
)

//添加用户
func (a *Admin) Add(user, pass string, role, status int) (err error) {
	engine := xorm.MustDB()
	admin := new(Admin)
	admin.User = user
	admin.Pass = pass
	admin.Role_id = role
	admin.Status = status
	admin.Create_time = time.Now().Unix()
	n, err := engine.Insert(admin)
	if err != nil {
		err = SystemFail
		return
	}
	if n == 0 {
		err = errors.New("用户添加失败")
	}
	return
}

func (a *Admin) Del(id int) (err error) {
	engine := xorm.MustDB()
	admin := new(Admin)
	admin.Id = id
	n, err := engine.Delete(admin)
	if err != nil {
		err = SystemFail
		return
	}
	if n == 0 {
		err = errors.New("删除失败")
	}
	return
}

func (a *Admin) List(offset, limit int) (list []*User, err error) {
	engine := xorm.MustDB()
	rows, err := engine.Limit(limit, offset).Rows(a)
	if err != nil {
		err = SystemFail
		return
	}
	defer rows.Close()
	role_list, err := DefaultRole.GetMapRole()
	if err != nil {
		return
	}
	for rows.Next() {
		admin := new(Admin)
		user := new(User)
		err = rows.Scan(admin)
		if err != nil {
			err = SystemFail
			return
		}
		user.Admin = admin
		if role, ok := role_list[admin.Role_id]; ok {
			user.Role_name = role.Name
		} else {
			user.Role_name = "暂无身份"
		}
		create_time := time.Unix(user.Admin.Create_time, 0).Format("2006-01-02 15:04:05")
		login_time := create_time
		if user.Admin.Login_time > 0 {
			login_time = time.Unix(user.Admin.Login_time, 0).Format("2006-01-02 15:04:05")
		}

		user.Create_time = create_time
		user.Login_time = login_time

		list = append(list, user)
	}
	return
}

//获取管理员的总数
func (a *Admin) Count() (count int64, err error) {
	engine := xorm.MustDB()
	count, err = engine.Count(a)
	if err != nil {
		err = SystemFail
		return
	}
	return
}

//检查用户名密码是否ok，ok则返回用户信息
func (a *Admin) Check(user, pass string) (admin Admin, err error) {
	if user == "" || pass == "" {
		err = ParamEmpty
		return
	}
	engine := xorm.MustDB()
	admin = Admin{}
	has, err := engine.Where("user=?", user).Get(&admin)
	if err != nil {
		err = SystemFail
		return
	}
	if !has {
		err = NotExistAdmin
		return
	}
	if admin.Pass != MakeMd5([]byte(pass)) {
		faygo.Debug(admin.Pass)
		faygo.Debug(MakeMd5([]byte(pass)))
		err = PassError
		return
	}
	if admin.Status != 1 {
		err = AdminDisable
		return
	}
	return
}

//更新用户最后一次登陆的信息
func (a *Admin) UpdateLastLogin(user string, ip string, now_time int64) (err error) {
	if user == "" || now_time == 0 {
		err = ParamEmpty
		return
	}
	engine := xorm.MustDB()
	admin := new(Admin)
	admin.Login_ip = ip
	admin.Login_time = now_time
	n, err := engine.Where("user=?", user).Cols("login_ip", "login_time").Update(admin)
	if err != nil {
		err = SystemFail
		return
	}
	if n == 0 {
		err = errors.New("更新用户信息失败")
		return
	}
	return
}

//通过id获取管理员
func (a *Admin) GetAdminById(id int) (admin *Admin, err error) {
	admin = new(Admin)
	engine := xorm.MustDB()
	has, err := engine.Where("id=?", id).Get(admin)
	if err != nil {
		err = SystemFail
		return
	}
	if !has {
		err = NotExistAdmin
	}
	return
}

func (a *Admin) SetStatus(id, status int) (err error) {
	admin, err := a.GetAdminById(id)
	if err != nil {
		return
	}
	//如果是超级管理员
	if admin.Role_id == 1 {
		err = errors.New("不能修改超级管理员的状态")
		return
	}
	admin.Status = status
	engine := xorm.MustDB()
	n, err := engine.Where("id=?", id).Cols("status").Update(admin)
	if err != nil {
		err = SystemFail
		return
	}
	if n == 0 {
		err = errors.New("修改失败")
	}
	return
}

func (a *Admin) Edit(id int, nickname string, face string) (err error) {
	engine := xorm.MustDB()
	admin := new(Admin)
	admin.Face = face
	admin.Nickname = nickname
	n, err := engine.Where("id=?", id).Cols("nickname", "face").Update(admin)
	if err != nil {
		err = SystemFail
		return
	}
	if n == 0 {
		err = errors.New("修改信息失败")
	}
	return
}

func (a *Admin) EditPass(id int, pass string) (err error) {
	engine := xorm.MustDB()
	admin := new(Admin)
	admin.Pass = MakeMd5([]byte(pass))
	n, err := engine.Where("id=?", id).Cols("pass").Update(admin)
	if err != nil {
		err = SystemFail
		return
	}
	if n == 0 {
		err = errors.New("修改信息失败")
	}
	return
}
