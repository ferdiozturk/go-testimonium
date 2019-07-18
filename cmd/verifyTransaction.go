// Copyright © 2019 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/spf13/cobra"
	"log"
)

var noOfConfirmations uint8

// verifyTransactionCmd represents the transaction command
var verifyTransactionCmd = &cobra.Command{
	Use:   "transaction [txHash]",
	Short: "Verify a transaction",
	Long: `Verify a transaction from the source chain on the destination chain

Behind the scene, the command queries the transaction with the specified hash ('txHash') from the source chain.
It then generates a Merkle Proof contesting the existence of the transaction within a specific block.
This information gets sent to the destination chain, where not only the existence of the block but also the Merkle Proof are verified`,
	Aliases: []string{"tx"},
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		txHash := common.HexToHash(args[0])
		txHash, blockHash, err := testimoniumClient.GenerateMerkleProof(txHash, verifyFlagSrcChain)
		if err != nil {
			log.Fatal("Failed to generate Merkle Proof: " + err.Error())
		}
		isValid := testimoniumClient.VerifyTransaction(txHash, blockHash, noOfConfirmations, verifyFlagDestChain)
		fmt.Println(isValid)
	},
}

func init() {
	verifyCmd.AddCommand(verifyTransactionCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// verifyTransactionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	verifyTransactionCmd.Flags().Uint8VarP(&noOfConfirmations, "confirmations", "c", 4, "Number of block confirmations")
}
