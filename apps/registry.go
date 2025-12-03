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
	_ "github.com/kade-chen/google-billing-console/apps/invoice/api"
	_ "github.com/kade-chen/google-billing-console/apps/token/api"
	_ "github.com/kade-chen/google-billing-console/apps/usagedate/api"

	//impl config
	_ "github.com/kade-chen/google-billing-console/apps/token/impl"
	_ "github.com/kade-chen/google-billing-console/apps/usagedate/impl/labelkey/impl"

	_ "github.com/kade-chen/google-billing-console/apps/usagedate/impl/project/impl"
	_ "github.com/kade-chen/google-billing-console/apps/usagedate/impl/services/impl"
	_ "github.com/kade-chen/google-billing-console/apps/usagedate/impl/sku/impl"

	_ "github.com/kade-chen/google-billing-console/apps/invoice/impl/labelkey/impl"
	_ "github.com/kade-chen/google-billing-console/apps/invoice/impl/project/impl"
	_ "github.com/kade-chen/google-billing-console/apps/invoice/impl/services/impl"
	_ "github.com/kade-chen/google-billing-console/apps/invoice/impl/sku/impl"
)
