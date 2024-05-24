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
  "600": `A network error occured. Try again later. If the error persists: ${contactAdmin}`,
  "400": `The request sent could not be understood.`,
  "401": `You are unauthorized. Please re-login.`,
  "402": `You are not allowed to do this. ${contactAdmin}`,
  "404": `Content not found. ${contactAdmin}`,
  "783": `The response from the server is not parsable. ${contactAdmin}`
};

const getErrorMessage = (code: string) => {
  const standardmessage = `An error occured. ${contactAdmin}`;
  if (ERRORMESSAGES[code]) return ERRORMESSAGES[code];
  return standardmessage;
};

export { getErrorMessage };
