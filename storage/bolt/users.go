package bolt

import (
	"reflect"

	"github.com/asdine/storm"
	"github.com/filebrowser/filebrowser/v2/errors"
	"github.com/filebrowser/filebrowser/v2/users"
)

type usersBackend struct {
	db *storm.DB
}

func (st usersBackend) GetByID(id uint) (*users.User, error) {
	user := &users.User{}
	err := st.db.One("ID", id, user)
	if err == storm.ErrNotFound {
		return nil, errors.ErrNotExist
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (st usersBackend) GetByUsername(username string) (*users.User, error) {
	user := &users.User{}
	err := st.db.One("Username", username, user)
	if err == storm.ErrNotFound {
		return nil, errors.ErrNotExist
	}

	if err != nil {
		return nil, err
	}

	return user, nil
}

func (st usersBackend) Gets() ([]*users.User, error) {
	users := []*users.User{}
	err := st.db.All(&users)
	if err == storm.ErrNotFound {
		return nil, errors.ErrNotExist
	}

	if err != nil {
		return users, err
	}

	return users, err
}

func (st usersBackend) Update(user *users.User, fields ...string) error {
	if len(fields) == 0 {
		return st.Save(user)
	}

	for _, field := range fields {
		val := reflect.ValueOf(user).Elem().FieldByName(field).Interface()
		if err := st.db.UpdateField(user, field, val); err != nil {
			return err
		}
	}

	return nil
}

func (st usersBackend) Save(user *users.User) error {
	err := st.db.Save(user)
	if err == storm.ErrAlreadyExists {
		return errors.ErrExist
	}
	return err
}

func (st usersBackend) DeleteByID(id uint) error {
	return st.db.DeleteStruct(&users.User{ID: id})
}

func (st usersBackend) DeleteByUsername(username string) error {
	user, err := st.GetByUsername(username)
	if err != nil {
		return err
	}

	return st.db.DeleteStruct(user)
}
