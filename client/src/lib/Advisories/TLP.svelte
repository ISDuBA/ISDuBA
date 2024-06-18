<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  export let tlp = "";
</script>

<span
  class:tlpred={tlp === "RED"}
  class:tlpamber={tlp === "AMBER"}
  class:tlpclear={tlp === "WHITE"}
  class:tlpgreen={tlp === "GREEN"}>{tlp}</span
>

<style>
  .tlpclear {
    font-size: small;
    padding: 0.25em;
    background: #000;
    color: #fff;
  }
  .tlpred {
    font-size: small;
    padding: 0.25em;
    background: #000;
    color: #ff2b2b;
  }
  .tlpamber {
    font-size: small;
    padding: 0.25em;
    background: #000;
    color: #ffc000;
  }
  .tlpgreen {
    font-size: small;
    padding: 0.25em;
    background: #000;
    color: #33ff00;
  }
</style>
