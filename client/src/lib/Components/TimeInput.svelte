<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Input } from "flowbite-svelte";

  interface Time {
    hours: number;
    minutes: number;
  }

  interface Props {
    roundEnd?: boolean;
    hours?: number | string | undefined;
    minutes?: number | string | undefined;
    onTimeChanged: (time: Time) => void;
  }

  let {
    roundEnd = true,
    hours = $bindable(undefined),
    minutes = $bindable(undefined),
    onTimeChanged
  }: Props = $props();

  const minutesInputClass = `w-20 rounded-s-none ${roundEnd ? "" : "rounded-e-none"}`;
  const numberRegex = /\d/i;

  const timeChanged = () => {
    onTimeChanged({
      hours: hours === "" ? 0 : Number(hours),
      minutes: minutes === "" ? 0 : Number(minutes)
    });
  };

  const isNumber = (text: string) => {
    return text.match(numberRegex);
  };

  const hasValidLength = (text: string) => {
    return text.length <= 2;
  };

  const hoursChanged = () => {
    const hoursString = `${hours}`;
    if (!isNumber(hoursString) || !hasValidLength(hoursString)) {
      hours = hoursString.slice(0, -1);
    }
    if (Number(hours) > 23) {
      hours = 23;
    }
    if (minutes === undefined || minutes === "") {
      minutes = 0;
    }
    timeChanged();
  };

  const minutesChanged = () => {
    const minutesString = `${minutes}`;
    if (!isNumber(minutesString) || !hasValidLength(minutesString)) {
      minutes = minutesString.slice(0, -1);
    }
    if (Number(minutes) > 59) {
      minutes = 59;
    }
    if (hours === undefined || hours === "") {
      hours = 0;
    }
    timeChanged();
  };

  export const clearInput = () => {
    hours = undefined;
    minutes = undefined;
  };
</script>

<Input
  oninput={hoursChanged}
  bind:value={hours}
  class="w-20 rounded-s-none rounded-e-none"
  placeholder="hh"
  title="Hours"
  type="number"
/>
<Input
  oninput={minutesChanged}
  bind:value={minutes}
  class={minutesInputClass}
  placeholder="mm"
  title="Minutes"
  type="number"
/>
