package main

import "time"

type (
	User struct {
		ID int64
	}
	Content struct {
		Caption string
		Type    string
		Data    []byte
		Creator int64 //user-id
		Created time.Time
		Updated time.Time
	}
	History struct {
		ID      int64
		Actor   User
		Action  string //CRUD...
		Subject string
		OldVal  any
		NewVal  any
		Time    time.Time
	}
)

/*
权限拟考虑RBAC，基于role来控制。以下层级高的role拥有低级role的所有权限，
且用户的角色不重叠（例如，一个用户如果已经是moderator了，就不能是该顶级票
范围内的developer了）：

- admin
  - 可创建、编辑、删除用户
  - 可设定用户的角色
  - 可开顶级票（即创建项目，本系统拟用顶级票来表示项目），且自动成为任何顶级票的moderator
- moderator
  - 权限限定在某个顶级票范围内（一个用户可以有多个顶级票的moderator角色）
  - 可在一个顶级票的范围内移动票
  - 可将票在本顶级票与分诊区之间移动
  - 可以设定票的状态
  - 可以分配票给自己或他人
  - 可以设定票的priority
  - 可编辑其他人创建的票/回复的内容、删除他人的内容
  - 可设定票的due date
- developer：
  - 权限限定在某个顶级票范围内（一个用户可以有多个顶级票的developer角色）
  - 可在顶级票下开子票
  - 可编辑/删除自己创建的内容
  - 可为票增加/删除tag
  - 可编辑分配给自己的票的status（具有一定限制，比如不能关闭）
- guest: 只读权限

其他设定：

- 系统默认有一个顶级票称为“分诊区”，是所有新票的默认入口（如果没有设置顶级票的话），
- 所有用户在分诊区具有他本人所拥有的最大权限，但至少为developer。
- 所有用户默认拥有所有顶级票的guest权限（即：只要本系统定义的用户，无需任何权限分配，
  至少拥有guest权限）
- status的值默认为0，表示“新建”，-1表示“关闭”，其他负值表示“删除”，其他正值可以
  自由定义。负值是“终结值”(即不可再设置为其他非负值了，票也不能再编辑了)。
- metrics是预留用作各种属性的，比如：
  - priority：票的优先级
  - estimate：预估工作时长
  - progress：当前百分比进度
  这些属性应该都是浮点数，且应该是定义后才使用。
*/
