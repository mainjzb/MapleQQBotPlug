package service

import (
	_ "database/sql"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/mainjzb/MapleQQBotPlug/config"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
	"strconv"
)

var xCasbin *casbin.Enforcer = nil

const (
	QA     = "QA"
	Admin2 = "admin2"
	Admin5 = "admin5"
)

func NewCasbin(gdb *gorm.DB) *casbin.Enforcer {
	if xCasbin != nil {
		return xCasbin
	}

	a, err := gormadapter.NewAdapterByDB(gdb)
	if err != nil {
		logrus.Panic(err)
	}

	xCasbin, err = casbin.NewEnforcer("config/rbac_model.conf", a)
	if err != nil {
		logrus.Panic(err)
	}

	err = xCasbin.LoadPolicy()
	if err != nil {
		logrus.Panic(err)
	}

	xCasbin.RemoveFilteredPolicy(0, "")
	xCasbin.SavePolicy()

	for _, group := range config.Instance.QAEditQQGroup {
		groupString := strconv.Itoa(int(group))
		xCasbin.AddPolicy("admin2", groupString, QA, "write")
		xCasbin.AddPolicy("admin2", groupString, QA, "read")

		xCasbin.AddPolicy("admin5", groupString, QA, "write")
		xCasbin.AddPolicy("admin5", groupString, QA, "read")
		xCasbin.AddPolicy("admin5", groupString, Admin2, "write")
		xCasbin.AddPolicy("admin5", groupString, Admin2, "read")

		xCasbin.AddPolicy(config.Instance.MasterQQ, groupString, QA, "write")
		xCasbin.AddPolicy(config.Instance.MasterQQ, groupString, QA, "read")
		xCasbin.AddPolicy(config.Instance.MasterQQ, groupString, Admin2, "write")
		xCasbin.AddPolicy(config.Instance.MasterQQ, groupString, Admin2, "read")
		xCasbin.AddPolicy(config.Instance.MasterQQ, groupString, Admin5, "write")
		xCasbin.AddPolicy(config.Instance.MasterQQ, groupString, Admin5, "read")
	}
	return xCasbin
}
func AddLevel(fromGroup, fromQQ, authQQ, level int) bool {
	if IsEditAdmin(fromGroup, fromQQ, level) && !IsEditAdmin(fromGroup, authQQ, level) {
		xCasbin.AddGroupingPolicy(strconv.Itoa(authQQ), "admin"+strconv.Itoa(level), strconv.Itoa(fromGroup))
		return true
	}
	return false
}

func DeleteLevel(fromGroup, fromQQ, authQQ int) bool {
	authQQInfo := xCasbin.GetFilteredGroupingPolicy(0, strconv.Itoa(authQQ), "", strconv.Itoa(fromGroup))
	if len(authQQInfo) != 1 {
		return true
	}
	authQQLevel := authQQInfo[0][1]

	if IsWrite(fromGroup, fromQQ, authQQLevel) && !IsWrite(fromGroup, authQQ, authQQLevel) {
		xCasbin.RemoveFilteredGroupingPolicy(0, strconv.Itoa(authQQ), authQQLevel, strconv.Itoa(fromGroup))
		return true
	}
	return false
}

func IsEditAdmin(fromGroup, fromQQ, level int) bool {
	ok, _ := xCasbin.Enforce(strconv.Itoa(fromQQ), strconv.Itoa(fromGroup), "admin"+strconv.Itoa(level), "write")
	return ok
}

func IsEditQA(fromGroup, fromQQ int) bool {
	ok, _ := xCasbin.Enforce(strconv.Itoa(fromQQ), strconv.Itoa(fromGroup), QA, "write")
	return ok
}
func IsQueryQA(fromGroup, fromQQ int) bool {
	ok, _ := xCasbin.Enforce(strconv.Itoa(fromQQ), strconv.Itoa(fromGroup), QA, "write")
	return ok
}

func IsWrite(fromGroup, fromQQ int, level string) bool {
	ok, _ := xCasbin.Enforce(strconv.Itoa(fromQQ), strconv.Itoa(fromGroup), level, "write")
	return ok
}
