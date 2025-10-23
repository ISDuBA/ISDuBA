<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  interface Props {
    filled?: boolean;
  }

  let { filled = false }: Props = $props();

  const background = "bg-gray-600 dark:bg-gray-500";
  const borderColor = "gray";
  let border = $derived(`border border-${borderColor}-800 dark:border-${borderColor}-500`);
</script>

<span class={`${filled ? background : ""} ${border} inline-block h-3 w-3 rounded-[3pt]`}></span>
