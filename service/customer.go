package service

import (
	"fmt"
	"github.com/ellis90/assessment-bg/datastore"
	"github.com/ellis90/assessment-bg/datastore/model"
	"github.com/ellis90/assessment-bg/entity"
	"github.com/ellis90/assessment-bg/utils"
	"github.com/labstack/echo/v4"
	"github.com/sirupsen/logrus"
	"net/http"
)

const Successful = "successful"

type CustomerConfiguration func(us *CustomerService) error

type CustomerService struct {
	userRepo datastore.UserRepository
}

func NewCustomerServices(cfgs ...CustomerConfiguration) (*CustomerService, error) {
	cs := &CustomerService{}
	for _, cfg := range cfgs {
		if err := cfg(cs); err != nil {
			return nil, err
		}
	}
	return cs, nil
}

func WithCustomerRepository(ur datastore.UserRepository, err error) CustomerConfiguration {
	return func(us *CustomerService) error {
		if err != nil {
			return err
		}
		us.userRepo = ur
		return err
	}
}

func WithPGXConfiguration(logger *logrus.Logger, src string) CustomerConfiguration {
	db, err := datastore.NewStore(logger, src)
	return WithCustomerRepository(db, err)
}

// handlers

func (cs *CustomerService) Create(ctx echo.Context) error {
	user := new(entity.User)
	logrus.Info("entry Binding")
	if err := ctx.Bind(user); err != nil {
		return utils.JSON(ctx, "bind", http.StatusBadRequest, err)
	}
	logrus.Info(user, "user gotten")
	cus, err := model.NewCustomer(user)
	if err != nil {
		return utils.JSON(ctx, "validation", http.StatusBadRequest, err)
	}
	logrus.Info(cus, "new customer gotten")
	out, err := cs.userRepo.Create(cus)
	if err != nil {
		return utils.JSON(ctx, "save", http.StatusBadRequest, err)
	}
	return utils.JSON(ctx, Successful, http.StatusOK, out.GetExportedCustomer())
}

func (cs *CustomerService) Update(ctx echo.Context) error {
	user := new(entity.User)
	if err := ctx.Bind(user); err != nil {
		return utils.JSON(ctx, "user", http.StatusBadRequest, err)
	}
	cus, err := model.NewCustomer(user)
	if err != nil {
		return utils.JSON(ctx, "validation", http.StatusBadRequest, err)
	}
	out, err := cs.userRepo.Update(cus)
	if err != nil {
		return utils.JSON(ctx, "update", http.StatusBadRequest, err)
	}
	return utils.JSON(ctx, Successful, http.StatusOK, out.GetExportedCustomer())
}

func (cs *CustomerService) FetchAll(ctx echo.Context) error {
	allCus, err := cs.userRepo.Get()
	logrus.Info(allCus, "get all customer gotten")
	if err != nil {
		return utils.JSON(ctx, "fetch all", http.StatusBadRequest, err)
	}
	return utils.JSON(ctx, Successful, http.StatusOK, allCus.GetExportedCustomers())
}

func (cs *CustomerService) DeleteById(ctx echo.Context) error {
	id := ctx.Param("id")
	err := cs.userRepo.Delete(id)
	if err != nil {
		return utils.JSON(ctx, fmt.Sprintf("delete %s", id), http.StatusBadRequest, err)
	}
	return utils.JSON(ctx, Successful, http.StatusOK, nil)
}
