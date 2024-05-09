/**
 * This file is Free Software under the Apache-2.0 License
 * without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
 *
 * SPDX-License-Identifier: Apache-2.0
 *
 * SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 * Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
 */

const contactAdmin = `Please contact an administrator.`;

const ERRORMESSAGES: any = {
  "500": `An error occured on the server. ${contactAdmin}`,
  "400": `The request sent from the client could not be understood. ${contactAdmin}`,
  "401": `You are not allowed to do this. ${contactAdmin}`,
  "783": `The response from the server is not parsable. ${contactAdmin}`
};

const getErrorMessage = (code: string) => {
  const standardmessage = `A general error occured. ${contactAdmin}`;
  if (ERRORMESSAGES[code]) return ERRORMESSAGES[code];
  return standardmessage;
};

export { getErrorMessage };