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

// @Title Doo Store API Documentation
// @Version 1.0
// @Description Description of Doo Store API documentation
// @SecurityDefinitions.apiKey BearerAuth
// @in header
// @name token
// @TermsOfService http://swagger.io/terms/

// @Contact.name xxyijixx
// @Contact.email xxyijixx@gmail.com

// @License.name AGPL-3.0
// @License.url https://opensource.org/license/agpl-v3
package main

import "doo-store/cmd"

func main() {
	cmd.Execute()
}
