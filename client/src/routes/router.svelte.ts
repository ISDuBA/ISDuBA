// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2026 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2026 Intevation GmbH <https://intevation.de>

import { push as routerPush } from "svelte-spa-router";

const routerState = $state({
  didPush: false
});

const push = (location: string) => {
  routerState.didPush = true;
  routerPush(location);
};

export { push, routerState };
