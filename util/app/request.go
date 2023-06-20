package app

import (
	"net/http"

	"github.com/astaxie/beego/validation"
	"github.com/gin-gonic/gin"
	errorspkg "github.com/pkg/errors"
	"github.com/yulecd/pp-common/perrors"
	"github.com/yulecd/pp-common/plog"
)

func BindAndValid(c *gin.Context, form interface{}) (int, int) {
	plog.Warn(c, "deprecated validate func:BindAndValid, Please use BindReqAndValid replace it")
	err := BindReqAndValid(c, form)
	if err != nil {
		if errorspkg.Is(err, perrors.NewValidateError(err)) {
			return http.StatusBadRequest, 400
		} else {
			return http.StatusInternalServerError, 500
		}
	}

	return http.StatusOK, 200
}

func BindReqAndValid(c *gin.Context, form interface{}) error {
	if err := c.ShouldBind(form); err != nil {
		plog.Errorf(c, "bind req err: %v", err)
		return err
	}

	if err := c.ShouldBindHeader(form); err != nil {
		plog.Errorf(c, "bind header err: %v", err)
		return err
	}

	xvalid := Validation{}
	ok, err := xvalid.Check(form, xvalidFun...)
	if !ok {
		plog.Errorf(c, "valid req err: %v", err)
		return perrors.NewValidateError(err)
	}

	valid := validation.Validation{}
	check, err := valid.Valid(form)
	if err != nil {
		plog.Errorf(c, "valid req err: %v", err)
		return perrors.NewValidateError(err)
	}
	if !check {
		err = ValidateMarkErrors(c, valid.Errors)
		return perrors.NewValidateError(err)
	}

	return nil
}

func ValidateMarkErrors(c *gin.Context, errors []*validation.Error) error {
	errwrap := errorspkg.New("invalid parameter")
	for _, err := range errors {
		plog.Error(c, err.Key, err.Message)
		errwrap = errorspkg.Wrap(err, errwrap.Error())
	}

	return errwrap
}
