<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Button, Label, Input } from "flowbite-svelte";
  import {
    createIsoTimeStringForSSVC,
    getDecision,
    parseDecisionTree,
    type SSVCDecision,
    type SSVCOption
  } from "./SSVCCalculator";
  import { createEventDispatcher, onMount } from "svelte";
  import { request } from "$lib/utils";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { getErrorMessage } from "$lib/Errors/error";
  import ComplexDecision from "./ComplexDecision.svelte";

  const dispatch = createEventDispatcher();

  export let disabled = false;
  export let documentID: string;
  export let allowEditing: any;
  let startedCalculation = false;
  export let isEditing = false;
  let isComplex = false;
  let currentStep = 0;
  let steps: string[] = [];
  let mainDecisions: SSVCDecision[] = [];
  let decisionPoints: any[] = [];
  let decisionsTable: any[] = [];
  let userDecisions: any = {};
  const vectorBeginning = "SSVCv2/";
  let vector: string;
  export let vectorInput = "";
  let result: any;
  let saveSSVCError: string;
  $: resultStyle = result?.color ? `color: ${result.color}` : "";

  onMount(async () => {
    loadDecisionTree();
  });

  function loadDecisionTree() {
    resetDecisions();
    ({ decisionPoints, decisionsTable, mainDecisions, steps } = parseDecisionTree());
  }

  function resetDecisions() {
    resetUserDecisions();
    steps = [];
    mainDecisions = [];
    decisionPoints = [];
    decisionsTable = [];
  }

  function resetUserDecisions() {
    userDecisions = {};
    currentStep = 0;
    vector = vectorBeginning;
    isComplex = false;
  }

  function getOption(decision: SSVCDecision, label: string): SSVCOption | undefined {
    return decision.options.find((element: SSVCOption) => element.label === label);
  }

  function extendVector(text: string) {
    vector = vector.concat(text);
  }

  function selectOption(option: SSVCOption) {
    userDecisions[mainDecisions[currentStep].label] = option.label;
    extendVector(`${mainDecisions[currentStep].key}:${option.key}/`);
    currentStep++;
    if (currentStep === mainDecisions.length - 1) {
      calculateResult();
    }
  }

  function doesContainChildCombo(selectedOptions: any, childCombos: any[]): boolean {
    const test = childCombos.find((combos: any[]) => {
      let count = 0;
      for (const option in selectedOptions) {
        const result = combos.find((combo: any) => {
          return (
            combo.child_label === option &&
            combo.child_option_labels.includes(selectedOptions[option])
          );
        });
        if (result) count++;
      }
      return count === Object.keys(selectedOptions).length;
    });
    return !!test;
  }

  function calculateComplexOption() {
    const selectedChildOptions: any = {};
    mainDecisions[currentStep].children?.forEach((child: any) => {
      const checkedRadioButton: any = document.querySelector(
        `input[name="${child.label}"]:checked`
      );
      if (checkedRadioButton) {
        selectedChildOptions[child.label] = checkedRadioButton.value;
      }
    });
    let selectedOption: SSVCOption | null = null;
    mainDecisions[currentStep].options.forEach((option: SSVCOption) => {
      if (option.child_combinations) {
        if (doesContainChildCombo(selectedChildOptions, option.child_combinations)) {
          selectedOption = option;
        }
      }
    });
    Object.keys(selectedChildOptions).forEach((decisionLabel) => {
      const decision = getDecision(decisionPoints, decisionLabel);
      if (decision) {
        const option = getOption(decision, selectedChildOptions[decisionLabel]);
        extendVector(`${decision.key}:${option?.key}/`);
      }
    });
    if (selectedOption) {
      selectOption(selectedOption);
    }
  }

  function calculateResult() {
    let filteredDecisions = decisionsTable;
    for (const key of Object.keys(userDecisions)) {
      filteredDecisions = filteredDecisions.filter((decision) => {
        return decision[key] && userDecisions[key] && decision[key] === userDecisions[key];
      });
    }
    const finalDecision = structuredClone(filteredDecisions[0]);
    for (const key of Object.keys(userDecisions)) {
      delete finalDecision[key];
    }
    const option = getOption(mainDecisions[currentStep], finalDecision.Decision);
    extendVector(
      `${mainDecisions[currentStep].key}:${option?.key}/${createIsoTimeStringForSSVC()}/`
    );
    const resultText: any = Object.values(finalDecision)[0];
    const color = getOption(mainDecisions[currentStep], resultText)?.color;
    result = {
      text: resultText,
      color: color
    };
  }

  function resetError() {
    saveSSVCError = "";
  }

  function stepBack() {
    if (currentStep === steps.length - 1) {
      // Delete ISO string and "/"
      vector = vector.slice(0, -1 - createIsoTimeStringForSSVC().length);
      // Delete final decision
      vector = vector.slice(0, -4);
    }
    // Find out if user did select child options of a decision. If yes we need to cut-off
    // more than the key pair for the parent decision.
    if (mainDecisions[currentStep - 1].children) {
      const children = mainDecisions[currentStep - 1].children;
      let tmpVector = vector;
      // Cut-off parent
      tmpVector = tmpVector.slice(0, -4);
      const keyPairs: string[] = [];
      if (children) {
        children.forEach(() => {
          const splittedVector = tmpVector.split("/");
          keyPairs.push(splittedVector[splittedVector.length - 2]);
          tmpVector = tmpVector.slice(0, -4);
        });
        let didUserChooseChildren = true;
        keyPairs.forEach((pair) => {
          const splittedPair = pair.split(":");
          let isChild = false;
          children.forEach((child: any) => {
            const childDecision = getDecision(decisionPoints, child.label);
            if (childDecision && childDecision.key !== splittedPair[0]) return;
            const optionsKeys = childDecision?.options.map((option) => option.key);
            if (optionsKeys?.includes(splittedPair[1])) isChild = true;
          });
          if (!isChild) didUserChooseChildren = false;
        });
        if (didUserChooseChildren) {
          vector = vector.slice(0, -(4 * children.length));
        }
      }
    }
    // Delete (parent) key pair
    vector = vector.slice(0, -4);
    currentStep--;
    const keyOfLastDecision = Object.keys(userDecisions)[Object.keys(userDecisions).length - 1];
    delete userDecisions[keyOfLastDecision];
  }

  async function saveSSVC(vector: string) {
    await allowEditing();
    resetError();
    const encodedUrl = encodeURI(`/api/ssvc/${documentID}?vector=${vector}`);
    const response = await request(encodedUrl, "PUT");
    if (response.ok) {
      isEditing = false;
      startedCalculation = false;
      resetUserDecisions();
      dispatch("updateSSVC");
    } else if (response.error) {
      if (response.error === "400") {
        saveSSVCError = `An error occured: ${response.content}.`;
      } else {
        saveSSVCError = getErrorMessage(response.error);
      }
    }
  }

  const toggleEditing = () => {
    isEditing = !isEditing;
  };

  $: onChangeDisabled(disabled);

  const onChangeDisabled = (disabled: boolean) => {
    if (disabled) {
      startedCalculation = false;
      resetUserDecisions();
    }
  };
</script>

<div id="ssvc-calc" class="flex flex-row items-center gap-x-3">
  {#if !startedCalculation}
    <div class="flex flex-col gap-x-3">
      <div>
        {#if isEditing}
          <Label class="mb-4">
            <span>Enter a SSVC directly</span>
          </Label>
          <Input
            autofocus
            class="h-6 w-96"
            disabled={disabled || !isEditing}
            on:keyup={(e) => {
              if (e.key === "Enter") saveSSVC(vectorInput);
              if (e.key === "Escape") toggleEditing();
            }}
            on:input={resetError}
            type="text"
            bind:value={vectorInput}
          />
        {:else if !vectorInput}
          <span class="h-6 text-lg">Please enter a SSVC</span>
        {/if}
      </div>
      {#if !isEditing}
        <button class="h-6" {disabled} on:click={toggleEditing}
          ><i class="bx bx-edit-alt ml-1"></i></button
        >
      {/if}
      {#if isEditing}
        <div class="ml-auto mt-4 flex flex-row gap-x-3">
          <Button
            outline
            size="xs"
            class="h-8"
            title="Undo"
            on:click={() => {
              isEditing = false;
            }}
          >
            Undo
          </Button>
          <Button
            outline
            size="xs"
            class="h-8"
            title="Evaluate"
            {disabled}
            on:click={() => (startedCalculation = true)}>Evaluate</Button
          >
          <Button size="xs" class="h-8" title="Save" on:click={() => saveSSVC(vectorInput)}
            >Save</Button
          >
        </div>
      {/if}
    </div>
  {:else}
    <span class="pt-[0.3rem] font-mono text-gray-400" color="gray"
      >{currentStep + 1}/{steps.length}</span
    >
    {#if steps[currentStep]}
      <span class="text-nowrap">{steps[currentStep]}</span>
    {/if}
    {#if mainDecisions[currentStep]}
      {#if currentStep < mainDecisions.length - 1}
        {#if mainDecisions[currentStep].decision_type === "simple"}
          <div class="flex flex-row items-baseline gap-3">
            {#each mainDecisions[currentStep].options as option}
              <Button
                outline
                size="xs"
                title={option.description}
                on:click={() => selectOption(option)}
                class="h-6"
                >{option.label}
              </Button>
            {/each}
          </div>
        {:else if mainDecisions[currentStep].decision_type === "complex"}
          {#if !isComplex}
            <div class="flex flex-row gap-x-3">
              {#each mainDecisions[currentStep].options as option}
                <Button
                  class="h-6"
                  outline
                  title={option.description}
                  size="xs"
                  on:click={() => selectOption(option)}>{option.label}</Button
                >
              {/each}
            </div>
            or
            <Button
              class="h-6"
              outline
              title="Custom"
              size="xs"
              on:click={() => {
                isComplex = true;
              }}>Custom</Button
            >
          {:else}
            <ComplexDecision
              on:calculateComplexOption={calculateComplexOption}
              children={mainDecisions[currentStep].children}
              {decisionPoints}
            ></ComplexDecision>
          {/if}
        {/if}
      {:else if result}
        <Label
          >Result:
          <span style={resultStyle}>{result.text}</span>
        </Label>
        <Label class="text-gray-400">Vector: {vector}</Label>
        <button title="Save" on:click={() => saveSSVC(vector)}>
          <i class="bx bx-save me-2 text-xl"></i></button
        >
      {/if}
    {/if}
    <div class="flex flex-row items-baseline gap-x-1">
      {#if currentStep > 0}
        <button title="Undo" class="h-6" color="light" on:click={stepBack}>
          <i class="bx bx-undo me-2 text-xl"></i>
        </button>
      {/if}
      <button
        title="Start over"
        class="h-6 text-nowrap"
        color="light"
        on:click={resetUserDecisions}
      >
        <i class="bx bx-reset me-2 text-xl"></i>
      </button>
      <button
        class="h-6"
        color="light"
        title="Back"
        on:click={() => {
          resetUserDecisions();
          startedCalculation = false;
        }}><i class="bx bx-arrow-back me-2 text-xl"></i></button
      >
    </div>
  {/if}
  <ErrorMessage message={saveSSVCError}></ErrorMessage>
</div>
