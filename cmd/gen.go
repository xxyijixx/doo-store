/*
Copyright © 2024 xxyijixx@gmail.com

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program. If not, see <http://www.gnu.org/licenses/>.
*/
package cmd

import (
	"doo-store/backend/app/cmd/gen"
	"doo-store/backend/config"

	"github.com/spf13/cobra"
)

// genCmd represents the gen command
var genCmd = &cobra.Command{
	Use:   "gen",
	Short: "A brief description of your command",
	Run: func(cmd *cobra.Command, args []string) {
		gen.Generate()
	},
}

func init() {
	if config.EnvConfig.ENV == "dev" {
		rootCmd.AddCommand(genCmd)
	}

}
