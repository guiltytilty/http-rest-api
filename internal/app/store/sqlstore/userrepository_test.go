package sqlstore_test

import (
	"fmt"
	"http-rest-api/internal/app/model"
	"http-rest-api/internal/app/store"
	"http-rest-api/internal/app/store/sqlstore"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUserRepository_Create(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)
	u := model.TestUser(t)
	assert.NoError(t, s.User().Create(u))
	assert.EqualError(t, s.User().Create(u), store.ErrDataConflict.Error())
	assert.NotNil(t, u)
}

func TestUserRepository_FindByID(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)
	id := 1
	_, err := s.User().FindByID(id)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	s.User().Create(&model.User{
		Name:  "user",
		Email: "user@example.org",
	})
	u, err := s.User().FindByID(id)
	assert.NoError(t, err)
	assert.NotNil(t, u)
}

func TestUserRepository_GetAll(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)

	users, err := s.User().GetAll(10, 0)
	assert.NoError(t, err)
	assert.Empty(t, users)

	for i := 0; i < 5; i++ {
		err := s.User().Create(&model.User{
			Name:  fmt.Sprintf("user%d", i),
			Email: fmt.Sprintf("user%d@example.org", i),
		})
		assert.NoError(t, err)
	}

	users, err = s.User().GetAll(3, 0)
	assert.NoError(t, err)
	assert.Len(t, users, 3)

	users, err = s.User().GetAll(3, 3)
	assert.NoError(t, err)
	assert.Len(t, users, 2)

	users, err = s.User().GetAll(10, 100)
	assert.NoError(t, err)
	assert.Empty(t, users)

	users, err = s.User().GetAll(1, 0)
	assert.NoError(t, err)
	assert.Len(t, users, 1)
	assert.NotZero(t, users[0].ID)
	assert.NotEmpty(t, users[0].Name)
	assert.NotEmpty(t, users[0].Email)
	assert.False(t, users[0].CreatedAt.IsZero())
	assert.False(t, users[0].UpdatedAt.IsZero())
}

func TestUserRepository_Update(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)

	name := "new"
	_, err := s.User().Update(1, model.UserUpdate{Name: &name})
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	age := 20
	err = s.User().Create(&model.User{
		Name:  "user",
		Email: "user@example.org",
		Age:   &age,
	})
	assert.NoError(t, err)

	before, err := s.User().FindByID(1)
	assert.NoError(t, err)
	require.NotNil(t, before.Age)
	assert.Equal(t, 20, *before.Age)

	newName := "updated"
	updated, err := s.User().Update(1, model.UserUpdate{Name: &newName})
	assert.NoError(t, err)
	assert.Equal(t, "updated", updated.Name)
	assert.Equal(t, before.Email, updated.Email)
	require.NotNil(t, updated.Age)
	assert.Equal(t, 20, *updated.Age)

	assert.False(t, updated.UpdatedAt.Equal(updated.CreatedAt),
		"updated_at should differ from created_at after update")

	newEmail := "new@example.org"
	newAge := 35
	updated, err = s.User().Update(1, model.UserUpdate{
		Email: &newEmail,
		Age:   &newAge,
	})
	assert.NoError(t, err)
	assert.Equal(t, "updated", updated.Name)
	assert.Equal(t, "new@example.org", updated.Email)
	require.NotNil(t, updated.Age)
	assert.Equal(t, 35, *updated.Age)

	assert.False(t, updated.UpdatedAt.Equal(updated.CreatedAt),
		"updated_at should differ from created_at after update")

	zero := 0
	updated, err = s.User().Update(1, model.UserUpdate{Age: &zero})
	assert.NoError(t, err)
	require.NotNil(t, updated.Age)
	assert.Equal(t, 0, *updated.Age)
	assert.Equal(t, "new@example.org", updated.Email)

	got, err := s.User().Update(1, model.UserUpdate{})
	assert.NoError(t, err)
	assert.Equal(t, updated.Name, got.Name)
	assert.Equal(t, updated.Email, got.Email)
	require.NotNil(t, got.Age)
	assert.Equal(t, 0, *got.Age)
	assert.False(t, got.UpdatedAt.Equal(got.CreatedAt),
		"updated_at should differ from created_at")

	fromDB, err := s.User().FindByID(1)
	assert.NoError(t, err)
	assert.Equal(t, "updated", fromDB.Name)
	assert.Equal(t, "new@example.org", fromDB.Email)
	require.NotNil(t, fromDB.Age)
	assert.Equal(t, 0, *fromDB.Age)
	assert.False(t, fromDB.UpdatedAt.Equal(fromDB.CreatedAt),
		"updated_at should differ from created_at in DB")
}

func TestUserRepository_Delete(t *testing.T) {
	db, teardown := sqlstore.TestDB(t, databaseURL)
	defer teardown("users")

	s := sqlstore.New(db)

	err := s.User().Delete(1)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	err = s.User().Create(&model.User{
		Name:  "user",
		Email: "user@example.org",
	})
	assert.NoError(t, err)

	_, err = s.User().FindByID(1)
	assert.NoError(t, err)

	err = s.User().Delete(1)
	assert.NoError(t, err)

	_, err = s.User().FindByID(1)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())

	err = s.User().Delete(1)
	assert.EqualError(t, err, store.ErrRecordNotFound.Error())
}
