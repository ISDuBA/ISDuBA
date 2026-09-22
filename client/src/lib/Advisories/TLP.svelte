<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  import Link from "$lib/Components/Link.svelte";
  import type { LabelOfTLP as LabelOfTLP2_0 } from "./types/csaf-2.0";
  import type { LabelOfTLP as LabelOfTLP2_1 } from "./types/csaf-2.1";

  interface Props {
    tlp: LabelOfTLP2_0 | LabelOfTLP2_1;
  }
  let { tlp }: Props = $props();

  const getTLPClass = (tlpLabel: LabelOfTLP2_0 | LabelOfTLP2_1): string => {
    return `tlp${tlpLabel.toLocaleLowerCase().replace("+", "-")}`;
  };
</script>

<Link href="https://www.first.org/tlp/" class={getTLPClass(tlp)}>
  {tlp}
</Link>
