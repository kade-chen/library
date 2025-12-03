package user

import (
	// "github.com/infraboard/mcube/v2/tools/pretty"
	"github.com/kade-chen/library/tools/format"
)

func NewUserSet() *UserSet {
	return &UserSet{
		Items: []*User{},
	}
}

func (u *UserSet) Add(item *User) {
	u.Items = append(u.Items, item)
}

func (u *UserSet) Tojson() string {
	return format.ToJSON(u)
}

// func (u *UserSet) Tojson1() string {
// 	return pretty.ToJSON(u)
// }

func (u *UserSet) UserIds() (uids []string) {
	// for i := range u.Items {
	// uids = append(uids, u.Items[i].Meta.Id)
	// }
	return
}
