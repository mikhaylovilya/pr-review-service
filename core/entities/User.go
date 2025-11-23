package entities

type User struct {
	Id       string
	Name     string
	TeamName string
	IsActive bool
}

func NewUser(id string, name string, teamName string, isActive bool) *User {
	return &User{
		Id:       id,
		Name:     name,
		TeamName: teamName,
		IsActive: isActive,
	}
}

func (u *User) SetStatus(isActive bool) {
	u.IsActive = isActive
}
