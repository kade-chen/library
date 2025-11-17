package apps

import (

	//注册所有provider
	_ "github.com/kade-chen/google-billing-console/apps/domain/impl"
	_ "github.com/kade-chen/google-billing-console/apps/token/impl"
	_ "github.com/kade-chen/google-billing-console/apps/token/provider/all"
	_ "github.com/kade-chen/google-billing-console/apps/user/impl"

	//impl config
	_ "github.com/kade-chen/google-billing-console/apps/configs/impl"

	//api config
	_ "github.com/kade-chen/google-billing-console/apps/token/api"
	_ "github.com/kade-chen/google-billing-console/apps/project/api"
	_ "github.com/kade-chen/google-billing-console/apps/services/api"
	_ "github.com/kade-chen/google-billing-console/apps/sku/api"

	//impl config
	_ "github.com/kade-chen/google-billing-console/apps/token/impl"
	_ "github.com/kade-chen/google-billing-console/apps/project/impl"
	_ "github.com/kade-chen/google-billing-console/apps/services/impl"
	_ "github.com/kade-chen/google-billing-console/apps/sku/impl"
)
