package domain

type User struct {
	id   int    `json:"id"`
	name string `json:"name"`
	age  uint8  `json:"age"`
}

func NewUser(id int, name string, age uint8) *User {
	return &User{id: id, name: name, age: age}
}

func (user *User) SetId(id int) {
	user.id = id
}

func (user *User) SetName(name string) {
	user.name = name
}

func (user *User) SetAge(age uint8) {
	user.age = age
}

func (user *User) ID() int {
	return user.id
}

func (user *User) Name() string {
	return user.name
}

func (user *User) Age() uint8 {
	return user.age
}
