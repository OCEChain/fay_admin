package admin_model

import (
	"github.com/go-errors/errors"
	x "github.com/go-xorm/xorm"
	"github.com/henrylee2cn/faygo"
	"github.com/henrylee2cn/faygo/ext/db/xorm"
	"sort"
	"strconv"
	"strings"
)

type Power struct {
	Id          int    `xorm:"not null int(11) pk autoincr"`
	Name        string `xorm:"not null default('') varchar(30) comment('栏目名')"`
	Action      string `xorm:"not null default('') varchar(30) comment('栏目对应地址')"`
	Pid         int    `xorm:"not null default(0) int(11) comment('父级id')"`
	Sort        int    `xorm:"not null default(0) int(11) comment('排序')"`
	Create_time int64  `xorm:"not null default(0) int(11) comment('创建时间')"`
	Update_time int64  `xorm:"not null default(0) int(11) comment('修改时间')"`
}

type Column struct {
	Power
	SonPower  []Column
	HasColumn bool
}

const (
	PowerTable = "power"
)

var DefaultPower *Power = new(Power)

var (
	NotExistPower = errors.New("用户没有任何权限")
)

func init() {
	err := xorm.MustDB().Table(PowerTable).Sync2(DefaultPower)
	if err != nil {
		faygo.Error(err.Error())
	}
}

//获取所有的栏目，如果当前传进来得权限id有该id，则给栏目打标记
func (p *Power) GetAllColumns(powerIds string) (columns []Column, err error) {
	engine := xorm.MustDB()
	var rows *x.Rows
	rows, err = engine.Rows(p)
	if err != nil {
		err = SystemFail
		return
	}
	defer rows.Close()
	power := strings.Split(powerIds, ",")
	power_map := make(map[int]int)
	for _, v := range power {
		i, err := strconv.Atoi(v)
		if err != nil {
			continue
		}
		power_map[i] = i
	}
	data_map := make(map[int]Column)
	for rows.Next() {
		power := Power{}
		column := Column{}
		err = rows.Scan(&power)
		if err != nil {
			err = SystemFail
			return
		}
		column.Power = power
		//如果有该栏目的权限，则标记
		if _, ok := power_map[column.Id]; ok {
			column.HasColumn = true
		}
		data_map[column.Id] = column
	}
	columns = p.convertData(data_map, 0)
	return
}

//通过栏目id获取权限
func (p *Power) GetColumnsByPowerIds(powerIds string) (columns []Column, actions map[string]Column, err error) {
	actions = make(map[string]Column)
	engine := xorm.MustDB()
	var rows *x.Rows
	switch powerIds {
	case "all":
		//如果是全部
		rows, err = engine.Rows(p)
	case "":
		return
	default:
		rows, err = engine.Where("id in (?)", powerIds).Rows(p)
	}

	if err != nil {
		err = SystemFail
		return
	}
	defer rows.Close()
	data_map := make(map[int]Column)

	for rows.Next() {
		power := Power{}
		column := Column{}
		err = rows.Scan(&power)
		if err != nil {
			err = SystemFail
			return
		}
		column.Power = power
		data_map[column.Id] = column
		actions[column.Action] = column
	}
	columns = p.convertData(data_map, 0)
	return
}

//将查询出的结果转换成树状结构返回
func (p *Power) convertData(data map[int]Column, pid int) (res []Column) {
	sortid := []int{}
	sort_data := make(map[int]Column)
	for _, v := range data {
		if v.Pid == pid {
			v.SonPower = p.convertData(data, v.Id)
			sortid = append(sortid, v.Sort)
			sort_data[v.Sort] = v
		}
	}
	//进行排序
	sort.IntSlice(sortid).Sort()
	for _, v := range sortid {
		col := sort_data[v]
		res = append(res, col)
	}
	return
}

//获取所有的栏目id
func (p *Power) GetAllCids() (cids_map map[int]int, cids []int, err error) {
	engine := xorm.MustDB()
	var rows *x.Rows
	//如果是全部
	rows, err = engine.Rows(p)
	defer rows.Close()
	cids_map = make(map[int]int)
	for rows.Next() {
		power := Power{}
		err = rows.Scan(&power)
		if err != nil {
			err = SystemFail
			return
		}
		cids_map[power.Id] = power.Id
		cids = append(cids, power.Id)
	}
	return
}
