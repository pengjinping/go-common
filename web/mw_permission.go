package web

import (
	"net/http"
	"oa-common/logger"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type Abilities struct {
	Name string `json:"name"`
}

func PermissionMiddleware() gin.HandlerFunc {
	// 顶一个带过期时间的用户权限结构体
	type EmployeesPermission struct {
		Expired       int
		permissionMap map[string]bool
	}
	// 此处单个用户 不存在并发问题 所以不用考虑变量同步问题
	globalEmployeesPermissionMap := make(map[string]EmployeesPermission)
	sql := "SELECT `name`,`entity_type`,`entity_id` FROM `abilities` " +
		"WHERE (EXISTS (SELECT*FROM `roles` INNER JOIN `permissions` ON `roles`.`id`=" +
		"`permissions`.`entity_id` WHERE permissions.ability_id=abilities.id AND " +
		"`permissions`.`forbidden`=0 AND `permissions`.`entity_type`='roles' AND " +
		"(EXISTS (SELECT*FROM `employees` INNER JOIN `assigned_roles` ON `employees`." +
		"`id`=`assigned_roles`.`entity_id` WHERE assigned_roles.role_id=roles.id AND `assigned_roles`." +
		"`entity_type`='App\\\\Models\\\\Employee' AND `employees`.`id`=?) OR `level`< (SELECT max(LEVEL) " +
		"FROM `roles` WHERE EXISTS (SELECT*FROM `employees` INNER JOIN `assigned_roles` ON `employees`." +
		"`id`=`assigned_roles`.`entity_id` WHERE assigned_roles.role_id=roles.id AND `assigned_roles`." +
		"`entity_type`='App\\\\Models\\\\Employee' AND `employees`.`id`=?)))) OR EXISTS (SELECT*FROM " +
		"`employees` INNER JOIN `permissions` ON `employees`.`id`=`permissions`.`entity_id` WHERE " +
		"permissions.ability_id=abilities.id AND `permissions`.`forbidden`=0 AND `permissions`." +
		"`entity_type`='App\\\\Models\\\\Employee' AND `employees`.`id`=?) OR EXISTS (SELECT*FROM " +
		"`permissions` WHERE permissions.ability_id=abilities.id AND `permissions`.`forbidden`=0 AND " +
		"`entity_id` IS NULL))"
	return func(c *gin.Context) {
		Db, exists := c.Get("DB")
		if !exists {
			logger.Error(c, "租户[%s]数据库连接不存在", zap.Any("host", c.Request.Host))
			FailWithMessage(http.StatusUnauthorized, "无效的租户数据库连接", c)
			c.Abort()
			return
		}
		//获取员工的email
		emailMix, e1 := c.Get("email")
		userIdMix, e2 := c.Get("userId")
		if !e1 || !e2 {
			logger.Error(c, "授权用户权限:用户不存在")
			FailWithMessage(http.StatusUnauthorized, "登录已失效", c)
			c.Abort()
			return
		}
		email := emailMix.(string)
		userId := userIdMix.(string)
		employeePermission, ok := globalEmployeesPermissionMap[email]
		//如果没有员工权限map或者权限map已过期
		if !ok || employeePermission.Expired <= time.Now().Second() {
			employeePermission = EmployeesPermission{
				Expired:       time.Now().Second() + 600, //有效期10分钟
				permissionMap: make(map[string]bool),
			}
			var results []Abilities
			Db.(*gorm.DB).Raw(sql, userId, userId, userId).Scan(&results)
			for _, item := range results {
				employeePermission.permissionMap[item.Name] = true
			}
			globalEmployeesPermissionMap[email] = employeePermission
		}
		url := c.Request.URL.String()
		if _, ok := employeePermission.permissionMap[url]; ok {
			c.Next()
		} else {
			FailWithMessage(http.StatusForbidden, "禁止访问", c)
			c.Abort()
			return
		}
	}
}
