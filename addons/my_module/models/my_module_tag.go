
package models

import (
	"sumeru/core/sdk"
)

type MyModuleTag struct {
	sdk.Model `sumeru:"model=my.module.tag"`

	Name  sdk.String             `sumeru:"required,unique,index,string=Name"`
	Color sdk.String `sumeru:"string=Color,default=gray,selection=gray:Gray,blue:Blue,green:Green,orange:Orange,red:Red"`
}
