package contexts

import (
	"fmt"
	"reflect"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/auth/types"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

const EMPLOYEE_CONTEXT = "EMPLOYEE_CONTEXT"

func GetEmployeeClaimsFromCtx(c *gin.Context) (*types.EmployeeSession, error) {
	var employeeSessionData *types.EmployeeSession

	ctx, ok := c.Get(EMPLOYEE_CONTEXT)
	if !ok {
		logrus.Info("no employee context found")
		return nil, fmt.Errorf("no employee context found")
	}

	employeeSessionData, ok = ctx.(*types.EmployeeSession)
	if !ok {
		wrappedErr := fmt.Errorf("error getting employee context: type assertion failed, from <%v> to <%v>", reflect.TypeOf(ctx), reflect.TypeOf(employeeSessionData))
		return nil, wrappedErr
	}

	return employeeSessionData, nil
}

func SetEmployeeCtx(c *gin.Context, employeeSessionData *types.EmployeeSession) {
	c.Set(EMPLOYEE_CONTEXT, employeeSessionData)
}

func GetEmployeeIDFromCtx(c *gin.Context) (uint, error) {
	var employeeSessionData *types.EmployeeSession

	ctx, ok := c.Get(EMPLOYEE_CONTEXT)
	if !ok {
		return 0, fmt.Errorf("no employee context found")
	}

	employeeSessionData, ok = ctx.(*types.EmployeeSession)
	if !ok {
		wrappedErr := fmt.Errorf("error getting employee context: type assertion failed, from <%v> to <%v>", reflect.TypeOf(ctx), reflect.TypeOf(employeeSessionData))
		return 0, wrappedErr
	}

	return employeeSessionData.EmployeeID, nil
}
