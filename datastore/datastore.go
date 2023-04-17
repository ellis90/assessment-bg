package datastore

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/Masterminds/squirrel"
	"github.com/ellis90/assessment-bg/datastore/model"
	"github.com/ellis90/assessment-bg/entity"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/log/logrusadapter"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/sirupsen/logrus"
)

type Store struct {
	Logger     *logrus.Logger
	DB         *sql.DB
	SQLBuilder squirrel.StatementBuilderType
}

// NewStore is a factory function that open a connection to db
func NewStore(logger *logrus.Logger, src string) (*Store, error) {
	var (
		err  error
		conn *pgx.ConnConfig
	)
	conn, err = pgx.ParseConfig(src)
	if err != nil {
		return nil, err
	}
	conn.Logger = logrusadapter.NewLogger(logger)
	db := stdlib.OpenDB(*conn)
	err = validateSchema(db)
	if err != nil {
		return nil, err
	}

	return &Store{
		Logger:     logger,
		DB:         db,
		SQLBuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(db),
	}, err
}

//Close connection
func (s *Store) Close() error {
	return s.DB.Close()
}

// Queries to communicate with the DB

// Create add new entity to the db
func (s *Store) Create(cus model.Customer) (model.Customer, error) {
	row := s.SQLBuilder.Insert(usersSchema).SetMap(map[string]any{
		"user_name":   cus.GetUserName(),
		"first_name":  cus.GetFirstName(),
		"last_name":   cus.GetLastName(),
		"email":       cus.GetEmail(),
		"department":  cus.GetDepartment(),
		"user_status": cus.GetUserStatus(),
	}).Suffix(`RETURNING "id"`).QueryRow()

	var Id string
	if err := row.Scan(&Id); err != nil {
		return model.Customer{}, fmt.Errorf(errorMsg, ErrFailedToCreateCustomer, err)
	}
	cus.SetID(Id)
	s.Logger.Info("customer created successfully")
	return cus, nil
}

func (s *Store) Update(cus model.Customer) (model.Customer, error) {
	_, err := s.SQLBuilder.Update(
		usersSchema,
	).SetMap(
		map[string]interface{}{
			"user_name":   cus.GetUserName(),
			"first_name":  cus.GetFirstName(),
			"last_name":   cus.GetLastName(),
			"email":       cus.GetEmail(),
			"department":  cus.GetDepartment(),
			"user_status": cus.GetUserStatus(),
		},
	).Where(
		squirrel.Eq{"id": cus.GetID()},
	).Exec()
	if err != nil {
		return model.Customer{}, fmt.Errorf(errorMsg, ErrUpdateCustomer, err)
	}
	return cus, nil
}

func (s *Store) Get() (model.Customers, error) {
	var customers model.Customers
	rows, err := s.SQLBuilder.Select("id, user_name, first_name, last_name, email, department, user_status").From(usersSchema).Query()
	if err != nil {
		return nil, fmt.Errorf(errorMsg, ErrFetchCustomer, err)
	}

	defer func() {
		if err := rows.Close(); err != nil {
			s.Logger.Error(err.Error())
		}
	}()
	s.Logger.Info("done fetching users")
	for rows.Next() {
		s.Logger.Info("done next row")
		sr, err := scanUserRows(rows)
		if err != nil {
			return nil, err
		}
		customers = append(customers, sr)
	}
	s.Logger.Info(customers)
	return customers, nil
}

func (s *Store) Delete(id string) error {
	_, err := s.SQLBuilder.Delete(
		usersSchema,
	).Where(squirrel.Eq{"id": id}).Exec()
	if err != nil {
		fmt.Errorf(errorMsg, ErrDeleteCustomer, err)
	}
	return nil
}

func scanUserRows(row squirrel.RowScanner) (model.Customer, error) {
	as := new(entity.User)
	err := row.Scan(
		&as.ID,
		&as.UserName,
		&as.FirstName,
		&as.LastName,
		&as.Email,
		&as.Department,
		(*statusWrapper)(&as.UserStatus),
	)
	logrus.Info(as, "userpage")
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return model.Customer{}, err
		}
		return model.Customer{}, err
	}
	return model.AddCustomer(as), nil
}
