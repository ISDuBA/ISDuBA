/**
 * This file is Free Software under the Apache-2.0 License
 * without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
 *
 * SPDX-License-Identifier: Apache-2.0
 *
 * SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 * Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
 */

import type { HttpResponse } from "$lib/types";

const contactAdmin = `Please contact an administrator.`;

const ERRORMESSAGES: any = {
  "500": `An error occured on the server. ${contactAdmin}`,
  "600": `A network error occured. Try again later. If the error persists: ${contactAdmin}`,
  "400": `The request sent could not be understood.`,
  "401": `You are unauthorized. Please re-login.`,
  "403": `You are not allowed to do this. ${contactAdmin}`,
  "404": `Content not found. Maybe it was deleted in the meantime.`,
  "783": `The response from the server is not parsable. ${contactAdmin}`
};

const getErrorMessage = (code: string) => {
  const standardmessage = `An error occured. ${contactAdmin}`;
  if (code) {
    if (ERRORMESSAGES[code]) return ERRORMESSAGES[code];
  }
  return standardmessage;
};

export type ErrorDetails = {
  message: string;
  details?: string;
};

export const getErrorDetails = (message: string, response?: HttpResponse): ErrorDetails => {
  const errorDetails: ErrorDetails = { message: message, details: undefined };
  if (response?.error) {
    errorDetails.message += ` ${response.error}: ${getErrorMessage(response.error)}${response.content ? "\n" + response.content : ""}`;
  }
  return errorDetails;
};
