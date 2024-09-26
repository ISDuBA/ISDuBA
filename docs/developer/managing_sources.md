# Design considerations and notes about getting CSAF documents into ISDuBA

**work in progress**


## User needs

Let us consider two typical users (or //personas//).

### Peter

Is a source manager and will get the request from others in the
organisation to add, remove and change the sources.


### Magreth

Is a system administrator with access to the command line of
the ISDuBA system and also has an account with source management role.
Magreth stands in for Peter if Peter is out of office.


## technical

Downloading is a seperate task ideally a different application (component).

We want to schedule over all sources.

Using HTTP 1.1 Etags allows us to query so fast that
we do not need to specify intervals for regular and interim documents.
The idea is to query every 10-15 minutes, which would be 144 times a day.
