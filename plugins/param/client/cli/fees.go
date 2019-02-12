package cli

import (
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/cosmos/cosmos-sdk/client/context"
	"github.com/cosmos/cosmos-sdk/client/utils"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authcmd "github.com/cosmos/cosmos-sdk/x/auth/client/cli"
	authtxb "github.com/cosmos/cosmos-sdk/x/auth/client/txbuilder"
	"github.com/cosmos/cosmos-sdk/x/gov"

	"github.com/binance-chain/node/app"
	"github.com/binance-chain/node/plugins/param/types"
)

const (
	flagTitle       = "title"
	flagDescription = "description"
	flagDeposit     = "deposit"

	//Fee flag
	flagFeeParamFile = "fee-param-file"
)

func GetCmdSubmitFeeChangeProposal(cdc *codec.Codec) *cobra.Command {
	feeParam := types.FeeChangeParams{[]types.FeeParam{}, ""}
	cmd := &cobra.Command{
		Use:   "submit-fee-change-proposal",
		Short: "Submit a fee or fee rate change proposal",
		RunE: func(cmd *cobra.Command, args []string) error {
			txBldr := authtxb.NewTxBuilderFromCLI().WithCodec(cdc)
			cliCtx := context.NewCLIContext().
				WithCodec(cdc).
				WithAccountDecoder(authcmd.GetAccountDecoder(cdc))
			title := viper.GetString(flagTitle)
			initialDeposit := viper.GetString(flagDeposit)
			feeParamFile := viper.GetString(flagFeeParamFile)
			feeParam.Description = viper.GetString(flagDescription)
			if feeParamFile != "" {
				bz, err := ioutil.ReadFile(feeParamFile)
				if err != nil {
					return err
				}
				err = cdc.UnmarshalJSON(bz, &(feeParam.FeeParams))
				if err != nil {
					return err
				}
			}
			err := feeParam.Check()
			if err != nil {
				return err
			}
			fromAddr, err := cliCtx.GetFromAddress()
			if err != nil {
				return err
			}
			amount, err := sdk.ParseCoins(initialDeposit)
			if err != nil {
				return err
			}
			// feeParam get interface field, use amino
			feeParamsBz, err := app.Codec.MarshalJSON(feeParam)
			if err != nil {
				return err
			}
			msg := gov.NewMsgSubmitProposal(title, string(feeParamsBz), gov.ProposalTypeFeeChange, fromAddr, amount)
			err = msg.ValidateBasic()
			if err != nil {
				return err
			}
			cliCtx.PrintResponse = true
			return utils.CompleteAndBroadcastTxCli(txBldr, cliCtx, []sdk.Msg{msg})
		},
	}
	cmd.Flags().String(flagFeeParamFile, "", "the file of fee params (json format)")
	cmd.Flags().String(flagTitle, "", "title of proposal")
	cmd.Flags().String(flagDescription, "", "description of proposal")
	cmd.Flags().String(flagDeposit, "", "deposit of proposal")
	cmd.Flags().Var(&feeParam, "fee-param", "Set the operate fee, '{param type}/{param map}'. e.g: 'operate/{\"send\": \"\",\"fee\":100000,\"fee_for\":1}'")
	return cmd
}