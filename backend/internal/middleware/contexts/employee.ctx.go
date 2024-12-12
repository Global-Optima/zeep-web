package contexts

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"reflect"
)

const EmployeeCtx = "Employee"

func GetEmployeeClaimsFromCtx(c *gin.Context) (*utils.EmployeeClaims, error) {
	var claims *utils.EmployeeClaims

	ctx, ok := c.Get(EmployeeCtx)
	if !ok {
		return nil, fmt.Errorf("no employee context found")
	}

	claims, ok = ctx.(*utils.EmployeeClaims)
	if !ok {
		wrappedErr := fmt.Errorf("error getting employee context: type assertion failed, from <%v> to <%v>", reflect.TypeOf(ctx), reflect.TypeOf(claims))
		return nil, wrappedErr
	}

	return claims, nil
}

func SetEmployeeCtx(c *gin.Context, claims *utils.EmployeeClaims) {
	c.Set(EmployeeCtx, claims)
}
