package contexts

import (
	"fmt"
	"github.com/Global-Optima/zeep-web/backend/pkg/utils"
	"github.com/gin-gonic/gin"
	"reflect"
)

const CUSTOMER_CONTEXT = "Customer"

func GetCustomerClaimsFromCtx(c *gin.Context) (*utils.CustomerClaims, error) {
	var claims *utils.CustomerClaims

	ctx, ok := c.Get(CUSTOMER_CONTEXT)
	if !ok {
		return nil, fmt.Errorf("no employee context found")
	}

	claims, ok = ctx.(*utils.CustomerClaims)
	if !ok {
		wrappedErr := fmt.Errorf("error getting customer context: type assertion failed, from <%v> to <%v>", reflect.TypeOf(ctx), reflect.TypeOf(claims))
		return nil, wrappedErr
	}

	return claims, nil
}

func SetCustomerCtx(c *gin.Context, claims *utils.CustomerClaims) {
	c.Set(CUSTOMER_CONTEXT, claims)
}
