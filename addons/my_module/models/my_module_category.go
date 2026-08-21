
package models

import (
	"sumeru/core/sdk"
)

type MyModuleCategory struct {
	sdk.Model `sumeru:"model=my.module.category"`

	Name sdk.String `sumeru:"required,string=Name"`
}
