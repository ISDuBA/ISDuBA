<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { appStore } from "$lib/store";
  import SectionHeader from "$lib/SectionHeader.svelte";
  import NewAdvisories from "./NewAdvisories.svelte";
  import NewEvents from "./NewEvents.svelte";
</script>

{#if $appStore.app.keycloak.authenticated}
  <SectionHeader title="New Events"></SectionHeader>
  <NewEvents></NewEvents>
  <SectionHeader title="New Advisories"></SectionHeader>
  <NewAdvisories></NewAdvisories>
{/if}
