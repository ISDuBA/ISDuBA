// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

// Package web implements the endpoints of the web server.
package web

//go:generate swag init -g doc.go --parseDependency --parseInternal --outputTypes go,json

//	@description	This is the ISDuBA API.
//	@title			ISDuBA API
//	@version		1.0

//	@contact.name	ISDuBA team
//	@contact.url	https://github.com/ISDuBA/ISDuBA
//	@contact.email	info@intevation.de

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@BasePath	/api
