// This file is Free Software under the MIT License
// without warranty, see README.md and LICENSES/MIT.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
//  Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

import type { PageLoad } from "./$types";

export const load: PageLoad = ({ params }) => {
  // TODO: Use the slug to receive the advisory the user wants to see.
  return {
    id: params.slug
  };
};
