<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import KeyValue from "$lib/Advisories/CSAFWebview/KeyValue.svelte";
  import type { xGenericURI } from "./producttypes";

  interface Props {
    x_generic_uris: xGenericURI[];
  }
  let { x_generic_uris }: Props = $props();
</script>

<p>X Generic URIs</p>
<div>
  {#each x_generic_uris as uri (uri.uri)}
    <KeyValue keys={["Namespace", "URI"]} values={[uri.namespace, uri.uri]} />
  {/each}
</div>
