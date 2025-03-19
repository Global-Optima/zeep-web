package contexts

import (
	"fmt"
	"reflect"

	"github.com/Global-Optima/zeep-web/backend/internal/modules/auth/types"
	"github.com/gin-gonic/gin"
)

const CUSTOMER_CONTEXT = "CUSTOMER_CONTEXT"

func GetCustomerClaimsFromCtx(c *gin.Context) (*types.CustomerClaims, error) {
	var claims *types.CustomerClaims

	ctx, ok := c.Get(CUSTOMER_CONTEXT)
	if !ok {
		return nil, fmt.Errorf("no employee context found")
	}

	claims, ok = ctx.(*types.CustomerClaims)
	if !ok {
		wrappedErr := fmt.Errorf("error getting customer context: type assertion failed, from <%v> to <%v>", reflect.TypeOf(ctx), reflect.TypeOf(claims))
		return nil, wrappedErr
	}

	return claims, nil
}

func SetCustomerCtx(c *gin.Context, claims *types.CustomerClaims) {
	c.Set(CUSTOMER_CONTEXT, claims)
}
