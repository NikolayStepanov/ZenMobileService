package repository

import (
	"ZenMobileService/internal/domain"
	"context"
	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
	log "github.com/sirupsen/logrus"
)

const UsersTable = "users"

const ErrUsersTableExist = "ERROR: relation \"users\" does not exist (SQLSTATE 42P01)"

type UserPostgres struct {
	tableCreated bool
	db           *pgx.Conn
}

func NewUserPostgres(db *pgx.Conn) *UserPostgres {
	return &UserPostgres{db: db}
}

func (ud *UserPostgres) CreateUsersTable(ctx context.Context) error {
	sql := "CREATE TABLE IF NOT EXISTS users(id serial primary key, name varchar(1024) not null, age integer not null);"
	_, err := ud.db.Exec(ctx, sql)
	if err != nil {
		log.Println(err)
		return err
	}
	return err
}

func (ud *UserPostgres) GetUser(ctx context.Context, userID int) (domain.User, error) {
	var row pgx.Row
	user := domain.User{}
	id := 0
	age := 0
	name := ""
	sql, args, err := sq.Select("*").
		From(UsersTable).
		Where("id = ?", userID).
		PlaceholderFormat(sq.Dollar).
		ToSql()

	row = ud.db.QueryRow(ctx, sql, args...)

	err = row.Scan(&id, &name, &age)
	if err != nil {
		log.Println(err)
		return user, err
	}

	user.SetId(id)
	user.SetName(name)
	user.SetAge(uint8(age))
	return user, err
}

func (ud *UserPostgres) CreateUser(ctx context.Context, user domain.User) (int, error) {
	var row pgx.Row
	id := 0

	sql, args, err := sq.Insert(UsersTable).Columns("name", "age").
		Values(user.Name(), int(user.Age())).
		Suffix("RETURNING id").
		PlaceholderFormat(sq.Dollar).
		ToSql()

	row = ud.db.QueryRow(ctx, sql, args...)

	err = row.Scan(&id)
	if err != nil {
		log.Infoln(err.Error())
		if err.Error() == ErrUsersTableExist {
			err = ud.CreateUsersTable(ctx)
			if err != nil {
				log.Error(err)
				return 0, err
			}
			row = ud.db.QueryRow(ctx, sql, args...)
			err = row.Scan(&id)
			if err != nil {
				log.Error(err)
				return id, err
			}
		}
	}
	return id, err
}
