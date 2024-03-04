<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

 SPDX-License-Identifier: MIT

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import KeyValue from "$lib/CSAFWebview/KeyValue.svelte";
  import type { xGenericURI } from "./producttypes";

  export let x_generic_uris: xGenericURI[];
</script>

<p>X Generic URIs</p>
<div>
  {#each x_generic_uris as uri}
    <KeyValue keys={["Namespace", "URI"]} values={[uri.namespace, uri.uri]} />
  {/each}
</div>
