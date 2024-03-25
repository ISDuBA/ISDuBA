<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import {
    Button,
    Card,
    Heading,
    Label,
    Radio,
    StepIndicator,
    P,
    Tooltip,
    Input
  } from "flowbite-svelte";
  import {
    createIsoTimeStringForSSVC,
    parseDecisionTree,
    type SSVCDecision,
    type SSVCOption
  } from "./SSVCCalculator";
  import { createEventDispatcher, onMount } from "svelte";
  import { appStore } from "$lib/store";
  import { request } from "$lib/utils";

  const dispatch = createEventDispatcher();

  export let documentID: string;
  let startedCalculation = false;
  let currentStep = 0;
  let steps: string[] = [];
  let mainDecisions: any[] = [];
  let decisionPoints: any[] = [];
  let decisionsTable: any[] = [];
  let userDecisions: any = {};
  const vectorBeginning = "SSVCv2/";
  let vector: string;
  export let vectorInput = "";
  let result: any;
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
  }

  function getDecision(label: string): SSVCDecision {
    return decisionPoints.find((element) => element.label === label);
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

  function doesContainChildCombo(selectedOptions: any, childCombos: any[]): Boolean {
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
    if (test) return true;
    return false;
  }

  function calculateComplexOption() {
    const selectedChildOptions: any = {};
    mainDecisions[currentStep].children.forEach((child: any) => {
      const checkedRadioButton: any = document.querySelector(
        `input[name="${child.label}"]:checked`
      );
      if (checkedRadioButton) {
        selectedChildOptions[child.label] = checkedRadioButton.value;
      }
    });
    let selectedOption: SSVCOption;
    mainDecisions[currentStep].options.forEach((option: SSVCOption) => {
      if (doesContainChildCombo(selectedChildOptions, option.child_combinations)) {
        selectedOption = option;
      }
    });
    Object.keys(selectedChildOptions).forEach((decisionLabel) => {
      const decision = getDecision(decisionLabel);
      const option = getOption(decision, selectedChildOptions[decisionLabel]);
      extendVector(`${decision.key}:${option?.key}/`);
    });
    selectOption(selectedOption);
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
    const resultText = Object.values(finalDecision)[0];
    const color = getOption(mainDecisions[currentStep], resultText).color;
    result = {
      text: resultText,
      color: color
    };
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
          const childDecision = getDecision(child.label);
          if (childDecision.key !== splittedPair[0]) return;
          const optionsKeys = childDecision.options.map((option) => option.key);
          if (optionsKeys.includes(splittedPair[1])) isChild = true;
        });
        if (!isChild) didUserChooseChildren = false;
      });
      if (didUserChooseChildren) {
        vector = vector.slice(0, -(4 * children.length));
      }
    }
    // Delete (parent) key pair
    vector = vector.slice(0, -4);
    currentStep--;
    const keyOfLastDecision = Object.keys(userDecisions)[Object.keys(userDecisions).length - 1];
    delete userDecisions[keyOfLastDecision];
  }

  async function saveSSVC(vector: string) {
    const encodedUrl = encodeURI(`/api/ssvc/${documentID}?vector=${vector}`);
    const response = await request(encodedUrl, "PUT");
    if (response) {
      dispatch("updateSSVC");
      appStore.displaySuccessMessage("SSVC updated");
    }
  }
</script>

<div id="ssvc-calc" class="pe-4">
  {#if !startedCalculation}
    <Label class="mb-4">Enter SSVC vector manually</Label>
    <Label class="mb-2">
      <Input
        on:keyup={(e) => {
          if (e.key === "Enter") saveSSVC(vectorInput);
        }}
        type="text"
        bind:value={vectorInput}
      />
      <small class="text-slate-400">Example: SSVCv2/E:N/A:N/T:P/M:M/D:T/2024-03-12T13:26:47Z/</small
      >
    </Label>
    <Label class="mb-4 mt-6">Calculate new SSVC vector</Label>
    <Button on:click={() => (startedCalculation = true)}>Calculate</Button>
  {:else}
    <div class="mb-4 flex gap-4">
      <Button
        color="light"
        on:click={() => {
          resetUserDecisions();
          startedCalculation = false;
        }}><i class="bx bx-arrow-back me-2 text-xl"></i>Back</Button
      >
      <Button color="light" on:click={resetUserDecisions}>
        <i class="bx bx-reset me-2 text-xl"></i>
        Restart</Button
      >
      {#if currentStep > 0}
        <Button class="h-10" color="light" on:click={stepBack}>
          <i class="bx bx-undo me-2 text-xl"></i>
          Undo
        </Button>
      {/if}
    </div>
    <StepIndicator
      class="mb-4 w-3/6"
      color="gray"
      hideLabel={true}
      currentStep={currentStep + 1}
      {steps}
    ></StepIndicator>
    {#if steps[currentStep]}
      <Heading class="mb-6 max-w-fit text-xl">{steps[currentStep]}</Heading>
    {/if}
    {#if mainDecisions[currentStep]}
      {#if currentStep < mainDecisions.length - 1}
        {#if mainDecisions[currentStep].decision_type === "simple"}
          <div class="flex flex-wrap gap-3">
            {#each mainDecisions[currentStep].options as option}
              <Card class="flex w-80 justify-between">
                <div class="flex flex-col">
                  <Button on:click={() => selectOption(option)} class="mb-2">{option.label}</Button>
                  <p class="mb-4">{option.description}</p>
                </div>
              </Card>
            {/each}
          </div>
        {:else if mainDecisions[currentStep].decision_type === "complex"}
          <div class="flex flex-wrap gap-3">
            {#each mainDecisions[currentStep].options as option}
              <Card class="flex justify-between">
                <div class="flex flex-col">
                  <Button on:click={() => selectOption(option)} class="mb-2">{option.label}</Button>
                  <p class="mb-4">{option.description}</p>
                </div>
              </Card>
            {/each}
          </div>
          <div class="my-2 flex">
            <div class="flex flex-col justify-between">
              <i class="bx bx-up-arrow-alt text-xl"></i>
              <i class="bx bx-down-arrow-alt text-xl"></i>
            </div>
            <div class="flex flex-col">
              <P>Choose</P>
              <P class="ml-2">- or -</P>
              <P>Calculate</P>
            </div>
          </div>
          <Card class="flex min-w-fit justify-between">
            <form on:submit={calculateComplexOption}>
              <div class="flex flex-row flex-wrap gap-x-6 gap-y-4">
                {#each mainDecisions[currentStep].children as child}
                  {@const childOptions = getDecision(child.label).options}
                  <div class="min-w-60">
                    <h5
                      class="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white"
                    >
                      {child.label}
                    </h5>
                    <div class="mb-4">
                      {#each childOptions as option}
                        <div
                          class="mb-2 cursor-pointer rounded border border-gray-200 dark:border-gray-700"
                        >
                          <Radio name={child.label} value={option.label} class="w-full p-4"
                            >{option.label}</Radio
                          >
                        </div>
                        <Tooltip class="max-w-96">{option.description}</Tooltip>
                      {/each}
                    </div>
                  </div>
                {/each}
              </div>
              <Button type="submit">Calculate</Button>
            </form>
          </Card>
        {/if}
      {:else if result}
        <Label class="me-1 text-lg"
          >Result:
          <span style={resultStyle}>{result.text}</span>
        </Label>
        <Label class="text-lg">Vector: {vector}</Label>
        <Button on:click={() => saveSSVC(vector)}>
          <i class="bx bx-save me-2 text-xl"></i>Save</Button
        >
      {/if}
    {/if}
  {/if}
</div>
