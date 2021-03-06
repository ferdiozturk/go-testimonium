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
// This file contains logic executed if the command "verify block" is typed in.
// Authors: Marten Sigwart, Philipp Frauenthaler

package cmd

import (
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"log"

	"github.com/spf13/cobra"
)

// verifyBlockCmd represents the block command
var verifyBlockCmd = &cobra.Command{
	Use:   "block [blockHash]",
	Short: "Verifies a block",
	Long: `Verifies a block from the target chain on the verifying chain

The command queries the block information belonging to the specified block hash ('blockHash') stored on the 
verifying blockchain and verifies if the information is correct by comparing it to the block information
on the target chain.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		blockHash := common.HexToHash(args[0])	// omit the first two chars "0x"

		testimoniumClient = createTestimoniumClient()
		headerExists, err := testimoniumClient.BlockHeaderExists(blockHash, verifyFlagDestChain)
		if err != nil {
			log.Fatal("Could not verify block header on verifying chain: " + err.Error())
		}
		if !headerExists {
			fmt.Printf("No header stored for block %s on verifying chain\n", ShortHexString(args[0]))
			return
		}
		_, err = testimoniumClient.OriginalBlockHeader(blockHash, verifyFlagSrcChain)
		if err != nil {
			log.Fatal("Could not get original block on target chain: " + err.Error())
		}
		fmt.Printf("Block %s is valid\n", ShortHexString(args[0]))
	},
}

func init() {
	verifyCmd.AddCommand(verifyBlockCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// verifyBlockCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// verifyBlockCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func ShortHexString(hex string) string {
	if len(hex) <= 12 {
		return hex
	}
	return fmt.Sprintf("%s...%s", hex[:6], hex[len(hex)-4:])
}
