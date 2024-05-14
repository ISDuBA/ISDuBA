# Errorhandling

## Requests

Every request is done centralized in `client/src/utils.ts`.

### Requests per component


## Error categories

The messages are defined in `error.ts`.

### Not Found

When a (sub-) page could not be found we display `NotFound.svelte`.

### Network errors

Regarding Network errors we inform the user in the following categories:

  `500`: `An error occured on the server. Please contact an administrator.`,
  `400`: `The request sent from the client could not be understood. Please contact an administrator.`,
  `401`: `You are unauthorized. Please re-login.`,
  `402`: `You are not allowed to do this. Please contact an administrator.`,

### Special case
  When client is unable to parse the JSON returned from the server: `The response from the server is not parsable. Please contact an administrator.`

### General Error
Every other error case is answered with `A general error occured. Please contact an administrator.`

