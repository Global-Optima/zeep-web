package contexts

import (
	"fmt"
	"reflect"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/auth/types"
	"github.com/gin-gonic/gin"
)

const EMPLOYEE_CONTEXT = "EMPLOYEE_CONTEXT"

func GetEmployeeClaimsFromCtx(c *gin.Context) (*types.EmployeeClaims, error) {
	var claims *types.EmployeeClaims

	ctx, ok := c.Get(EMPLOYEE_CONTEXT)
	if !ok {
		return nil, fmt.Errorf("no employee context found")
	}

	claims, ok = ctx.(*types.EmployeeClaims)
	if !ok {
		wrappedErr := fmt.Errorf("error getting employee context: type assertion failed, from <%v> to <%v>", reflect.TypeOf(ctx), reflect.TypeOf(claims))
		return nil, wrappedErr
	}

	return claims, nil
}

func SetEmployeeCtx(c *gin.Context, claims *types.EmployeeClaims) {
	c.Set(EMPLOYEE_CONTEXT, claims)
}

func GetEmployeeIDFromCtx(c *gin.Context) (uint, error) {
	var claims *types.EmployeeClaims

	ctx, ok := c.Get(EMPLOYEE_CONTEXT)
	if !ok {
		return 0, fmt.Errorf("no employee context found")
	}

	claims, ok = ctx.(*types.EmployeeClaims)
	if !ok {
		wrappedErr := fmt.Errorf("error getting employee context: type assertion failed, from <%v> to <%v>", reflect.TypeOf(ctx), reflect.TypeOf(claims))
		return 0, wrappedErr
	}

	return claims.EmployeeClaimsData.ID, nil
}
