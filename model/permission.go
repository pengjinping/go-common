package model

import (
	"context"
	"fmt"
)

type Permission struct {
	BaseModel
}
type Abilities struct {
	Name string `json:"name"`
}

// Model 的用法
// 业务方
// 注意：
// 每个 Model 都【必须】有一个对应的 New 方法
// 1、方便这样的链式调用： model.NewPermission(ctx).UserPermission("uuid")
// 2、若是“默认强制使用平台库”，这里会预设 forcePlatform 为 true

func NewPermission(ctx context.Context) *Permission {
	var p Permission
	p.InitModel(ctx)
	p.InitMemCache()
	return &p
}

func (p *Permission) UserPermission(email, userId string, t int) map[string]interface{} {
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

	pMap := p.Cache().Remember(email, t, func(args ...interface{}) (interface{}, error) {
		var results []Abilities
		var pMap = make(map[string]interface{})

		p.DB().Raw(sql, userId, userId, userId).Scan(&results)
		for _, item := range results {
			pMap[item.Name] = true
		}

		if len(pMap) <= 0 {
			return nil, fmt.Errorf("用户权限不存在")
		}

		return pMap, nil
	})

	if pMap != nil {
		return pMap.(map[string]interface{})
	}

	return nil
}
