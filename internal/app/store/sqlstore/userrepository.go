package sqlstore

import (
	"database/sql"
	"errors"
	"fmt"
	"http-rest-api/internal/app/model"
	"http-rest-api/internal/app/store"
	"strings"

	"github.com/lib/pq"
)

type UserRepository struct {
	store *Store
}

func (r *UserRepository) GetAll(limit, offset int) ([]*model.User, error) {
	rows, err := r.store.db.Query(
		"SELECT id, name, email, age, created_at, updated_at FROM users ORDER BY id LIMIT $1 OFFSET $2",
		limit,
		offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*model.User, 0, limit)
	for rows.Next() {
		u := &model.User{}
		if err := rows.Scan(
			&u.ID,
			&u.Name,
			&u.Email,
			&u.Age,
			&u.CreatedAt,
			&u.UpdatedAt,
		); err != nil {
			return nil, err
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *UserRepository) Create(u *model.User) error {
	if err := u.Validate(); err != nil {
		return store.ErrInvalidData
	}

	err := r.store.db.QueryRow(
		"INSERT INTO users (name, email, age) VALUES ($1, $2, $3) RETURNING id, name, email, age, created_at, updated_at",
		u.Name,
		u.Email,
		u.Age,
	).Scan(&u.ID,
		&u.Name,
		&u.Email,
		&u.Age,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		var pqErr *pq.Error
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return store.ErrDataConflict
		}
		return err
	}

	return nil
}

func (r *UserRepository) FindByID(id int) (*model.User, error) {
	u := &model.User{}
	if err := r.store.db.QueryRow(
		"SELECT id, name, email, age, created_at, updated_at FROM users WHERE id = $1",
		id,
	).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Age,
		&u.CreatedAt,
		&u.UpdatedAt,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}

		return nil, err
	}

	return u, nil
}

func (r *UserRepository) Update(id int, upd model.UserUpdate) (*model.User, error) {
	set := []string{}
	args := []any{}
	i := 1

	if upd.Name != nil {
		set = append(set, fmt.Sprintf("name = $%d", i))
		args = append(args, *upd.Name)
		i++
	}
	if upd.Email != nil {
		set = append(set, fmt.Sprintf("email = $%d", i))
		args = append(args, *upd.Email)
		i++
	}
	if upd.Age != nil {
		set = append(set, fmt.Sprintf("age = $%d", i))
		args = append(args, *upd.Age)
		i++
	}

	if len(set) == 0 {
		return r.FindByID(id)
	}

	//set = append(set, "updated_at = NOW()")
	args = append(args, id)

	query := fmt.Sprintf(
		"UPDATE users SET %s WHERE id = $%d RETURNING id, name, email, age, created_at, updated_at",
		strings.Join(set, ", "),
		i,
	)

	u := &model.User{}
	err := r.store.db.QueryRow(query, args...).Scan(
		&u.ID,
		&u.Name,
		&u.Email,
		&u.Age,
		&u.CreatedAt,
		&u.UpdatedAt,
	)

	if err != nil {
		var pqErr *pq.Error
		if err == sql.ErrNoRows {
			return nil, store.ErrRecordNotFound
		}
		if errors.As(err, &pqErr) && pqErr.Code == "23505" {
			return nil, store.ErrDataConflict
		}
		return nil, err
	}

	if err := u.Validate(); err != nil {
		return nil, store.ErrInvalidData
	}

	return u, nil
}

func (r *UserRepository) Delete(id int) error {
	res, err := r.store.db.Exec(
		"DELETE FROM users WHERE id = $1",
		id,
	)
	if err != nil {
		return err
	}

	affected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if affected == 0 {
		return store.ErrRecordNotFound
	}

	return nil
}
