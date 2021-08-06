package main

import (
	"context"
	"github.com/jackc/pgx/v4/pgxpool"
)

type User struct {
	id int
	userName string
	firstName string
	lastName string
	email string
	phone string
}

func instantiateUser(
	id int,
	userName *string,
	firstName *string,
	lastName *string,
	email *string,
	phone *string) *User {

	user := User{id, *userName, *firstName, *lastName, *email, *phone}
	return &user
}

func (user *User) Create(pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	_ = conn.QueryRow(
		context.Background(),
		"INSERT INTO client (userName, firstName, lastName, email, phone) VALUES($1, $2, $3, $4, $5)",
		user.userName, user.firstName, user.lastName, user.email, user.phone)

	return nil
}

func RetrieveById(id int, pool *pgxpool.Pool) (*User, error) {
	var userName, firstName, lastName, email, phone string

	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return nil, err
	}
	defer conn.Release()

	row =: conn.QueryRow("SELECT * FROM client WHERE id = $1", id)

	err = row.Scan(&userName, &firstName, &lastName, &email, &phone)
	if err != nil {
		return nil, err
	}

	return instantiateUser(id, &userName, &firstName, &lastName, &email, &phone), nil
}

func (user *User) Update(id int, pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(
		context.Background(),
		"UPDATE client SET firstName = $1, lastName = $2, email = $3, phone = $4 WHERE id = $5",
		user.firstName, user.lastName, user.email, user.phone, id)
	if err != nil {
		return err
	}

	return nil
}

func DeleteById(id int, pool *pgxpool.Pool) error {
	conn, err := pool.Acquire(context.Background())
	if err != nil {
		return err
	}
	defer conn.Release()

	_, err = conn.Exec(
		context.Background(),
		"DELETE FROM client WHERE id = $1",
		id)
	if err != nil {
		return err
	}

	return nil
}


