package dao

import (
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
)

type User struct {
	id   int
	name string
}

type Dao struct {
	db *sql.DB
}

func NewDao() *Dao {
	return &Dao{db: &sql.DB{}}
}

func (d *Dao) FindUserById(id int) (User, error) {
	//data, err := d.db.Query("select * From User WHERE id = ? ", id)
	if id == 1 {
		return User{
			id:   1,
			name: "ode",
		}, nil
	} else {
		msg := fmt.Sprintf("can not find user by id: %d", id)
		return User{}, errors.Wrap(sql.ErrNoRows, msg)
	}
}
