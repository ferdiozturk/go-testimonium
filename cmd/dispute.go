// This file contains logic executed if the command "dispute" is typed in.
// Authors: Marten Sigwart, Philipp Frauenthaler

package cmd

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/pantos-io/go-testimonium/ethereum/ethash"
	"github.com/pantos-io/go-testimonium/testimonium"
	"github.com/spf13/cobra"
	"log"
)

var disputeFlagSrcChain uint8
var disputeFlagDestChain uint8

// disputeCmd represents the dispute command
var disputeCmd = &cobra.Command{
	Use:   "dispute [txHash]",
	Short: "Disputes a submitted block header",
	Long: `Disputes the submitted block header that was submitted through the specified transaction (txHash)`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		txHash := common.HexToHash(args[0]) // omit the first two chars "0x"

		// get blockNumber, nonce and RlpHeaderHashWithoutNonce and generate dataSetLookup and witnessForLookup
		testimoniumClient = createTestimoniumClient()
		header, err := testimoniumClient.GetHeaderFromTxData(txHash, disputeFlagDestChain)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("create DAG, compute dataSetLookup and witnessForLookup")
		// get DAG and compute dataSetLookup and witnessForLookup
		rlpHeaderWithoutNonce, err := testimonium.RlpHeaderWithoutNonce(header)
		if err != nil {
			log.Fatal(err)
		}
		rlpHeaderHashWithoutNonce := crypto.Keccak256Hash(rlpHeaderWithoutNonce)
		blockMetaData := ethash.NewBlockMetaData(header.Number.Uint64(), header.Nonce.Uint64(), rlpHeaderHashWithoutNonce)
		dataSetLookup := blockMetaData.DAGElementArray()
		witnessForLookup := blockMetaData.DAGProofArray()
		parent, err := testimoniumClient.OriginalBlockHeader(header.ParentHash, disputeFlagSrcChain)
		if err != nil {
			log.Fatal(err)
		}
		testimoniumClient.DisputeBlock(header, parent.Header(), dataSetLookup, witnessForLookup, disputeFlagDestChain)
	},
}

func init() {
	rootCmd.AddCommand(disputeCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// disputeCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// disputeCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
	disputeCmd.Flags().Uint8VarP(&disputeFlagSrcChain, "source", "s", 0, "source chain")
	disputeCmd.Flags().Uint8VarP(&disputeFlagDestChain, "chain", "d", 1, "destination chain")
}
