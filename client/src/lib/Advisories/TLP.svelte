<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import Link from "$lib/Components/Link.svelte";

  interface Props {
    tlp?: string;
  }
  let { tlp = "" }: Props = $props();
</script>

<Link href="https://www.first.org/tlp/" class={tlp ? `tlp${tlp.toLocaleLowerCase()}` : ""}>
  {tlp}
</Link>

<style>
  .tlpwhite {
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
