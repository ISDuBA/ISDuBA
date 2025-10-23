<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { type Source, fetchSourceDefaultConfig } from "$lib/Sources/source";
  import { Accordion, AccordionItem, Button, Input, Label, Select } from "flowbite-svelte";
  import CCheckbox from "$lib/Components/CCheckbox.svelte";
  import { onMount, untrack } from "svelte";
  import CFileinput from "$lib/Components/CFileinput.svelte";

  // Define enum first
  enum AgeUnit {
    hours = "h",
    days = "d",
    weeks = "w",
    months = "m",
    years = "y"
  }

  interface Props {
    formClass?: string;
    source: Source;
    oldSource?: Source | undefined;
    enableActive?: boolean;
    parseSource?: boolean;
    inputChange?: () => void;
  }

  let {
    formClass = "",
    source = $bindable(),
    oldSource = undefined,
    enableActive = false,
    parseSource = $bindable(true),
    inputChange = () => {}
  }: Props = $props();
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
    }
    ageNumber = baseNumber;
    ageUnit = baseUnit;
  };

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

  let headers: [string, string][] = $state([["", ""]]);
  let privateCert: FileList | undefined = $state();
  let privateCertReset: boolean = $state(false);
  let publicCert: FileList | undefined = $state();
  let publicCertReset: boolean = $state(false);

  let ageUnit: AgeUnit = $state(AgeUnit.years);
  let ageNumber: number | undefined = $state(undefined);
  let previousAgeNumber: number | undefined;

  let displayActiveHighlight: boolean = $state(true);

  let ratePlaceholder = $state(0);
  let slotPlaceholder = $state(2);

  const loadSourceDefaults = async () => {
    const resp = await fetchSourceDefaultConfig();
    if (resp.ok) {
      ratePlaceholder = resp.value.rate;
      slotPlaceholder = resp.value.slots;
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
    let lastElement = headers[headers.length - 1];
    if (lastElement === undefined || (lastElement[0] !== "" && lastElement[1] !== "")) {
      headers.push(["", ""]);
    }
    formatHeaders();

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
    if (source.ignore_patterns.length === 1) {
      source.ignore_patterns = [""];
    } else {
      source.ignore_patterns = source.ignore_patterns.toSpliced(index, 1);
    }
    inputChange();
  };

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
      if (header[0] !== "") source.headers.push(`${header[0]}:${header[1] ?? ""}`);
    }
  };

  const loadCerts = async () => {
    if (privateCert) {
      source.client_cert_private = await privateCert.item(0)?.text();
    } else {
      if (oldSource?.client_cert_private === "***" && privateCertReset) {
        source.client_cert_private = null;
      } else {
        source.client_cert_private = "***";
      }
    }
    if (publicCert) {
      source.client_cert_public = await publicCert.item(0)?.text();
    } else {
      if (oldSource?.client_cert_public === "***" && publicCertReset) {
        source.client_cert_public = null;
      } else {
        source.client_cert_public = "***";
      }
    }
  };
  $effect(() => {
    untrack(() => headers);
    if (source.headers && parseSource) {
      parseHeaders();
      parseSource = false;
    }
  });
</script>

<form class={formClass}>
  <Label>Name</Label>
  <Input class="mb-3" oninput={inputChange} bind:value={source.name}></Input>
  <div class={!source.active && displayActiveHighlight ? "blink" : ""}>
    {#if enableActive}
      <CCheckbox
        class="mb-3"
        onChanged={() => {
          displayActiveHighlight = false;
          inputChange();
        }}
        bind:checked={source.active}>Active</CCheckbox
      >
    {/if}
  </div>
  <Accordion>
    <AccordionItem>
      {#snippet header()}
        <span>Credentials</span>
      {/snippet}
      <Label>Private cert</Label>
      <CFileinput
        bind:files={privateCert}
        bind:isFileReset={privateCertReset}
        oldFile={oldSource?.client_cert_private}
        onChanged={inputChange}
        id="private-cert"
        titleClearButton="Remove private cert"
      ></CFileinput>
      <Label>Public cert</Label>
      <CFileinput
        bind:files={publicCert}
        bind:isFileReset={publicCertReset}
        oldFile={oldSource?.client_cert_public}
        onChanged={inputChange}
        id="public-cert"
        titleClearButton="Remove public cert"
      ></CFileinput>
      <Label>Client cert passphrase</Label>
      <div class="mb-3 inline-flex w-full">
        <Input
          class="rounded-none rounded-l-lg"
          oninput={inputChange}
          bind:value={source.client_cert_passphrase}
        />
        <Button
          onclick={() => {
            source.client_cert_passphrase = "";
          }}
          title="Remove passphrase"
          class="w-fit rounded-none rounded-r-lg border-l-0 p-1 dark:border-gray-500 dark:bg-gray-600"
          color="light"
        >
          <i class="bx bx-x"></i>
        </Button>
      </div>
    </AccordionItem>
    <AccordionItem>
      {#snippet header()}
        <span>Advanced options</span>
      {/snippet}
      <div class="mb-3 grid w-full gap-x-2 gap-y-4 md:grid-cols-[minmax(190px,1fr)_1fr_1fr]">
        <div>
          <Label>Maximum document age</Label>
          <div class="inline-flex w-full">
            <Input
              class="rounded-none rounded-l-lg"
              type="number"
              min="0"
              placeholder="0"
              oninput={onChangedAge}
              bind:value={ageNumber}
            ></Input>
            <Select
              class="rounded-none rounded-r-lg border-l-0"
              items={ageUnits}
              bind:value={ageUnit}
              onchange={onChangedAge}
            />
          </div>
        </div>
        <div>
          <Label>Downloads/s</Label>
          <Input
            type="number"
            step="0.01"
            placeholder={ratePlaceholder.toString()}
            oninput={inputChange}
            min="0"
            bind:value={source.rate}
          />
        </div>
        <div>
          <Label>Concurrent downloads</Label>
          <Input
            type="number"
            step="1"
            placeholder={slotPlaceholder.toString()}
            min="1"
            oninput={inputChange}
            bind:value={source.slots}
          />
        </div>
      </div>

      <Label>Options</Label>
      <div class="mb-3 flex w-full gap-4">
        <CCheckbox onChanged={inputChange} bind:checked={source.strict_mode}>Strict mode</CCheckbox>
        <CCheckbox onChanged={inputChange} bind:checked={source.secure}
          >Check TLS certificates</CCheckbox
        >
        <CCheckbox onChanged={inputChange} bind:checked={source.signature_check}
          >Check document OpenPGP signature</CCheckbox
        >
      </div>

      <Label>Ignore patterns</Label>
      {#each source.ignore_patterns as _pattern, index (index)}
        <div class="mb-3 inline-flex w-full">
          <Label class="grow">
            <Input
              class="rounded-none rounded-l-lg"
              oninput={onChangedIgnorePatterns}
              bind:value={source.ignore_patterns[index]}
            />
          </Label>
          <Button
            onclick={() => removePattern(index)}
            title="Remove pattern"
            class="w-fit rounded-none rounded-r-lg border-l-0 p-1 dark:border-gray-500 dark:bg-gray-600"
            color="light"
            disabled={source.ignore_patterns.length === 0 ||
              (source.ignore_patterns.length === 1 && source.ignore_patterns[0] === "")}
          >
            <i class="bx bx-x"></i>
          </Button>
        </div>
      {/each}

      <Label>HTTP headers</Label>
      {#each headers as header, index (index)}
        <div class="mb-3 grid w-full grid-cols-[1fr_auto] sm:grid-cols-[1fr_1fr_auto]">
          <Label class="col-span-2 col-start-1 row-start-1 sm:col-span-1">
            <span class="text-gray-500">Field name</span>
            <span class="text-gray-500 sm:collapse">- Field value</span>
          </Label>
          <Input
            class="col-span-2 row-start-2 rounded-none rounded-t-lg sm:col-span-1 sm:rounded-l-lg sm:rounded-tr-none"
            oninput={onChangedHeaders}
            bind:value={header[0]}
          />
          <Label class="collapse col-span-2 col-start-1 row-start-1 sm:visible sm:col-start-2">
            <span class="text-gray-500">Field value</span>
          </Label>
          <Input
            class="row-start-3 rounded-none rounded-bl-lg border-t-0 sm:row-start-2 sm:rounded-bl-none sm:border-t sm:border-l-0"
            oninput={onChangedHeaders}
            bind:value={header[1]}
          />
          {#if headers.length > 1}
            <Button
              onclick={() => removeHeader(index)}
              title="Remove field-name-field-value-pair"
              class="row-start-3 h-full w-fit rounded-none rounded-br-lg border-t-0 border-l-0 p-1 sm:row-start-2 sm:rounded-tr-lg sm:border-t dark:border-gray-500 dark:bg-gray-600"
              color="light"
            >
              <i class="bx bx-x"></i>
            </Button>
          {:else}
            <Button
              title="Remove key-value-pair"
              class=" row-start-3 h-full w-fit rounded-none rounded-br-lg border-t-0 border-l-0 p-1 sm:row-start-2 sm:rounded-tr-lg sm:border-t dark:border-gray-500 dark:bg-gray-600"
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
