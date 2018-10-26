package tokens

import (
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"

	"github.com/BiJie/BinanceChain/app/router"
	"github.com/BiJie/BinanceChain/plugins/tokens/burn"
	"github.com/BiJie/BinanceChain/plugins/tokens/freeze"
	"github.com/BiJie/BinanceChain/plugins/tokens/issue"
	"github.com/BiJie/BinanceChain/plugins/tokens/store"
)

func Routes(tokenMapper store.Mapper, accKeeper auth.AccountKeeper, keeper bank.Keeper) map[string]router.Handler {
	routes := make(map[string]router.Handler)
	routes[issue.Route] = issue.NewHandler(tokenMapper, keeper)
	routes[burn.BurnRoute] = burn.NewHandler(tokenMapper, keeper)
	routes[freeze.FreezeRoute] = freeze.NewHandler(tokenMapper, accKeeper, keeper)
	return routes
}
