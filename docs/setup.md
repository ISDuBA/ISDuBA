<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

 SPDX-License-Identifier: MIT

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->


# Installing prerequisite Software
In order to setup Keycloak and isdubad,
recent versions of Java and go are needed.

The following details on how to set these up.

## Install Java
A recent version of java is required.
The following will install Java 17:
```
sudo apt install openjdk-17-jre-headless
```

## Setup Go

The following will download Go 1.22.0:
```
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
```
Extract it and place the new go version into the /usr/local directory:
```
rm -rf /usr/local/go && tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
```
### Make the profile always use this version of go:
Open the profile with a text manager.
```
vim /etc/profile
```
In there, add the line:
```
export PATH=$PATH:/usr/local/go/bin
```
The system will now use go1.22.0 when go is called.

You can check whether it was successful via:
```
go version
```

## Install unzip
The following will install unzip which can be used
to unpack ''.zip' archives
```
sudo apt install unzip
```
