<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { type Source, fetchSourceDefaultConfig } from "$lib/Sources/source";
  import {
    Accordion,
    AccordionItem,
    Button,
    Fileupload,
    Input,
    Label,
    Select
  } from "flowbite-svelte";
  import CCheckbox from "$lib/Components/CCheckbox.svelte";
  import { onMount } from "svelte";
  export let formClass: string = "";
  export let source: Source;
  export let enableActive: boolean = false;
  export const updateSource = async () => {
    formatHeaders();
    await loadCerts();
  };

  const parseAge = (age?: string): [number | undefined, AgeUnit] => {
    if (age === "0s") {
      age = "0h";
    }
    let baseNumber: number | undefined = undefined;
    let baseUnit: AgeUnit = ageUnit;
    let [numStr, ...r]: string[] = (age ?? "").split("h");
    let num: number = +numStr;
    if (!(Number.isInteger(num) && !/[1-9]/.test(r.join(``)))) {
      throw Error("Expected age to be given exclusively in hours, actual value was '" + age + "'.");
    }
    if (num) {
      for (let i = ageUnits.length - 1; i >= 0; i--) {
        let unit = ageUnits[i].value;
        let len = ageUnitLengths[unit];
        if (num % len === 0) {
          baseNumber = num / len;
          baseUnit = unit;
          break;
        }
      }
    }
    return [baseNumber, baseUnit];
  };

  export const fillAgeDataFromSource = (useSource: Source) => {
    ageUnit = AgeUnit.years;
    let baseNumber: number | undefined = undefined;
    let baseUnit: AgeUnit = ageUnit;
    if (useSource.age && !["0s", "0h"].includes(useSource.age)) {
      [baseNumber, baseUnit] = parseAge(source.age);
    } else if (useSource.age && ["0s", "0h"].includes(useSource.age)) {
      baseNumber = 0;
    } else {
      let placeholder: number | undefined;
      [placeholder, baseUnit] = parseAge(ageDefaultDuration);
      agePlaceholder = placeholder ?? 0;
    }
    ageNumber = baseNumber;
    ageUnit = baseUnit;
  };

  export let inputChange = () => {};

  enum AgeUnit {
    hours = "h",
    days = "d",
    weeks = "w",
    months = "m",
    years = "y"
  }

  const ageUnits: { value: AgeUnit; name: string }[] = [
    { value: AgeUnit.hours, name: "hours" },
    { value: AgeUnit.days, name: "days" },
    { value: AgeUnit.weeks, name: "weeks" },
    { value: AgeUnit.months, name: "months" },
    { value: AgeUnit.years, name: "years" }
  ];

  const ageUnitLengths: { [unit in AgeUnit]: number } = {
    h: 1,
    d: 24,
    w: 24 * 7,
    m: 24 * 30,
    y: 24 * 365
  };

  let headers: [string, string][] = [["", ""]];
  let privateCert: FileList | undefined;
  let publicCert: FileList | undefined;

  let ageUnit: AgeUnit;
  let ageNumber: number | undefined;
  let previousAgeNumber: number | undefined;

  let displayActiveHighlight: boolean = true;

  let ratePlaceholder = 0;
  let slotPlaceholder = 2;

  let ageDefaultDuration = "1h";
  let agePlaceholder = 2;

  const loadSourceDefaults = async () => {
    const resp = await fetchSourceDefaultConfig();
    if (resp.ok) {
      ratePlaceholder = resp.value.rate;
      slotPlaceholder = resp.value.slots;
      ageDefaultDuration = resp.value.age;
    }
  };

  onMount(async () => {
    await loadSourceDefaults();
    fillAgeDataFromSource(source);
  });

  const onChangedAge = () => {
    if (!ageNumber && ageNumber !== 0) {
      source.age = "";
    } else {
      let num = ageNumber;
      num *= ageUnitLengths[ageUnit];
      source.age = num.toString() + "h";
    }
    if (ageNumber || previousAgeNumber !== ageNumber) {
      inputChange();
      previousAgeNumber = ageNumber;
    }
  };

  const onChangedHeaders = (e: Event | undefined) => {
    const lastIndex = headers.length - 1;
    if (
      (headers[lastIndex][0].length > 0 && headers[lastIndex][1].length > 0) ||
      (lastIndex - 1 >= 0 &&
        headers[lastIndex - 1][0].length > 0 &&
        headers[lastIndex - 1][1].length > 0)
    ) {
      headers.push(["", ""]);
      headers = headers;
    }
    if (e) {
      inputChange();
    }
  };

  const onChangedIgnorePatterns = () => {
    if (source.ignore_patterns.at(-1) !== "") {
      source.ignore_patterns.push("");
    }
    inputChange();
  };

  const removeHeader = (index: number) => {
    if (headers.length === 1)
      headers = [
        ["", ""],
        ["", ""]
      ];
    headers = headers.toSpliced(index, 1);
    inputChange();
  };

  const removePattern = (index: number) => {
    if (source.ignore_patterns.length === 1) source.ignore_patterns = [""];
    source.ignore_patterns = source.ignore_patterns.toSpliced(index, 1);
    inputChange();
  };

  $: if (source.headers) {
    parseHeaders();
  }

  const parseHeaders = () => {
    headers = [];
    for (const header of source.headers) {
      let h = header.split(":");
      headers.push([h[0], h[1]]);
    }
    if (headers.length === 0) {
      headers.push(["", ""]);
    }
    onChangedHeaders(undefined);
  };

  const formatHeaders = () => {
    source.headers = [];
    for (const header of headers) {
      if (header[0] !== "" && header[1] !== "") source.headers.push(`${header[0]}:${header[1]}`);
    }
  };

  const loadCerts = async () => {
    if (privateCert) {
      source.client_cert_private = await privateCert.item(0)?.text();
    }
    if (publicCert) {
      source.client_cert_public = await publicCert.item(0)?.text();
    }
  };
</script>

<form class={formClass}>
  <Label>Name</Label>
  <Input class="mb-3" on:input={inputChange} bind:value={source.name}></Input>
  <div class={!source.active && displayActiveHighlight ? "backgroundHighlight" : ""}>
    {#if enableActive}
      <CCheckbox
        class="mb-3"
        on:change={() => {
          displayActiveHighlight = false;
          inputChange();
        }}
        bind:checked={source.active}>Active</CCheckbox
      >
    {/if}
  </div>
  <Accordion>
    <AccordionItem
      ><span slot="header">Credentials</span>
      <Label>Private cert</Label>
      <div class="mb-3 inline-flex w-full">
        <Fileupload
          class="rounded-none rounded-l-lg"
          on:change={inputChange}
          bind:files={privateCert}
        ></Fileupload>
        <Button
          on:click={() => {
            source.client_cert_private = null;
            privateCert = undefined;
            inputChange();
          }}
          title="Remove private cert"
          class="w-fit rounded-none rounded-r-lg border-l-0 p-1"
          color="light"
        >
          <i class="bx bx-x"></i>
        </Button>
      </div>
      <Label>Public cert</Label>
      <div class="mb-3 inline-flex w-full">
        <Fileupload
          class="rounded-none rounded-l-lg"
          on:change={inputChange}
          bind:files={publicCert}
        ></Fileupload>
        <Button
          on:click={() => {
            source.client_cert_public = null;
            publicCert = undefined;
            inputChange();
          }}
          title="Remove public cert"
          class="w-fit rounded-none rounded-r-lg border-l-0 p-1"
          color="light"
        >
          <i class="bx bx-x"></i>
        </Button>
      </div>
      <Label>Client cert passphrase</Label>
      <div class="mb-3 inline-flex w-full">
        <Input
          class="rounded-none rounded-l-lg"
          on:input={inputChange}
          bind:value={source.client_cert_passphrase}
        />
        <Button
          on:click={() => {
            source.client_cert_passphrase = null;
          }}
          title="Remove passphrase"
          class="w-fit rounded-none rounded-r-lg border-l-0 p-1"
          color="light"
        >
          <i class="bx bx-x"></i>
        </Button>
      </div>
    </AccordionItem>
    <AccordionItem
      ><span slot="header">Advanced options</span>
      <div class="mb-3 grid w-full gap-x-2 gap-y-4 md:grid-cols-[minmax(190px,1fr)_1fr_1fr]">
        <div>
          <Label>Age</Label>
          <div class="inline-flex w-full">
            <Input
              class="rounded-none rounded-l-lg"
              type="number"
              min="0"
              placeholder={agePlaceholder.toString()}
              on:input={onChangedAge}
              bind:value={ageNumber}
            ></Input>
            <Select
              class="rounded-none rounded-r-lg border-l-0"
              items={ageUnits}
              bind:value={ageUnit}
              on:change={onChangedAge}
            />
          </div>
        </div>
        <div>
          <Label>Rate</Label>
          <Input
            type="number"
            step="0.01"
            placeholder={ratePlaceholder.toString()}
            on:input={inputChange}
            min="0"
            bind:value={source.rate}
          />
        </div>
        <div>
          <Label>Slots</Label>
          <Input
            type="number"
            step="1"
            placeholder={slotPlaceholder.toString()}
            min="1"
            on:input={inputChange}
            bind:value={source.slots}
          />
        </div>
      </div>

      <Label>Options</Label>
      <div class="mb-3 flex w-full gap-4">
        <CCheckbox on:change={inputChange} bind:checked={source.strict_mode}>Strict mode</CCheckbox>
        <CCheckbox on:change={inputChange} bind:checked={source.insecure}>Insecure</CCheckbox>
        <CCheckbox on:change={inputChange} bind:checked={source.signature_check}
          >Signature check</CCheckbox
        >
      </div>

      <Label>Ignore patterns</Label>
      {#each source.ignore_patterns as pattern, index (index)}
        <div class="mb-3 inline-flex w-full">
          <Label class="grow">
            <Input
              class="rounded-none rounded-l-lg"
              on:input={onChangedIgnorePatterns}
              bind:value={pattern}
            />
          </Label>
          <Button
            on:click={() => removePattern(index)}
            title="Remove pattern"
            class="w-fit rounded-none rounded-r-lg border-l-0 p-1"
            color="light"
            disabled={source.ignore_patterns.length <= 1}
          >
            <i class="bx bx-x"></i>
          </Button>
        </div>
      {/each}

      <Label>HTTP headers</Label>
      {#each headers as header, index (index)}
        <div class="mb-3 grid w-full grid-cols-[1fr_auto] sm:grid-cols-[1fr_1fr_auto]">
          <Label class="col-span-2 col-start-1 row-start-1 sm:col-span-1">
            <span class="text-gray-500">Key</span>
            <span class="text-gray-500 sm:collapse">- Value</span>
          </Label>
          <Input
            class="col-span-2 row-start-2 rounded-none rounded-t-lg sm:col-span-1 sm:rounded-l-lg sm:rounded-tr-none"
            on:input={onChangedHeaders}
            bind:value={header[0]}
          />
          <Label class="collapse col-span-2 col-start-1 row-start-1 sm:visible sm:col-start-2">
            <span class="text-gray-500">Value</span>
          </Label>
          <Input
            class="row-start-3 rounded-none rounded-bl-lg border-t-0 sm:row-start-2 sm:rounded-bl-none sm:border-l-0 sm:border-t"
            on:input={onChangedHeaders}
            bind:value={header[1]}
          />
          {#if headers.length > 1}
            <Button
              on:click={() => removeHeader(index)}
              title="Remove key-value-pair"
              class="row-start-3 h-full w-fit rounded-none rounded-br-lg border-l-0 border-t-0 p-1 sm:row-start-2 sm:rounded-tr-lg sm:border-t"
              color="light"
            >
              <i class="bx bx-x"></i>
            </Button>
          {:else}
            <Button
              title="Remove key-value-pair"
              class=" row-start-3 h-full w-fit rounded-none rounded-br-lg border-l-0 border-t-0 p-1 sm:row-start-2 sm:rounded-tr-lg sm:border-t"
              color="light"
              disabled={true}
            >
              <i class="bx bx-x"></i>
            </Button>
          {/if}
        </div>
      {/each}
    </AccordionItem>
  </Accordion>
  <br />
</form>

<style>
  @keyframes fadeIt {
    0% {
      background-color: #ffffff;
    }
    50% {
      background-color: #048dd9;
    }
    100% {
      background-color: #ffffff;
    }
  }

  .backgroundHighlight {
    background-image: none !important;
    animation: fadeIt 0.75s 2 ease-in-out 0.5s;
  }
</style>
