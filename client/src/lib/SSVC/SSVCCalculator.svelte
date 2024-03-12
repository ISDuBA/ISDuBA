<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

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
    convertVectorToLabel,
    createIsoTimeStringForSSVC,
    loadDecisionTreeFromFile
  } from "./SSVCCalculator";
  import { onMount } from "svelte";
  import { appStore } from "$lib/store";

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
  let vectorInput = "SSVCv2/E:N/A:N/T:P/M:M/D:T/2024-03-12T13:26:47Z/";
  let labelOfConvertedVector = "";
  let colorOfConvertedVector: string;
  let result: any;
  $: resultStyle = result?.color ? `color: ${result.color}` : "";
  $: convertedVectorStyle = colorOfConvertedVector ? `color: ${colorOfConvertedVector}` : "";

  onMount(async () => {
    loadDecisionTree();
  });

  function loadDecisionTree(): Promise<void> {
    return new Promise((resolve) => {
      resetDecisions();
      loadDecisionTreeFromFile().then((result: any) => {
        ({ decisionPoints, decisionsTable, mainDecisions, steps } = result);
        resolve();
      });
    });
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

  function getDecision(label: string): any {
    return decisionPoints.find((element) => element.label === label);
  }

  function getOption(decision: any, label: string): any {
    return decision.options.find((element: any) => element.label === label);
  }

  function extendVector(text: string) {
    vector = vector.concat(text);
  }

  function selectOption(option: any) {
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
    let result: string = "";
    let selectedOption: any;
    mainDecisions[currentStep].options.forEach((option: any) => {
      if (doesContainChildCombo(selectedChildOptions, option.child_combinations)) {
        result = option.label;
        selectedOption = option;
      }
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
    extendVector(`${mainDecisions[currentStep].key}:${option.key}/${createIsoTimeStringForSSVC()}`);
    const resultText = Object.values(finalDecision)[0];
    const color = getOption(mainDecisions[currentStep], resultText).color;
    result = {
      text: resultText,
      color: color
    };
  }

  function stepBack() {
    if (currentStep === steps.length - 1) {
      // Delete ISO string and last key pair
      vector = vector.slice(0, -24);
    } else {
      // Delete key pair
      vector = vector.slice(0, -4);
    }
    currentStep--;
    const keyOfLastDecision = Object.keys(userDecisions)[Object.keys(userDecisions).length - 1];
    delete userDecisions[keyOfLastDecision];
  }

  function convertVector() {
    const { color, label } = convertVectorToLabel(vectorInput, mainDecisions);
    labelOfConvertedVector = label;
    colorOfConvertedVector = color;
  }

  function saveSSVC(vector: string) {
    const encodedUrl = encodeURI(`/api/ssvc/${documentID}?vector=${vector}`);
    fetch(encodedUrl, {
      headers: {
        Authorization: `Bearer ${$appStore.app.keycloak.token}`
      },
      method: "PUT"
    }).then((response) => {
      if (response.ok) {
        console.log("response", response);
      } else {
        // Do errorhandling
      }
    });
  }
</script>

<div id="ssvc-calc" class="pe-4">
  {#if !startedCalculation}
    <Label class="mb-4 text-lg">Convert and save existing SSVC vector</Label>
    <Label class="mb-2">
      Vector:
      <Input type="text" bind:value={vectorInput} />
    </Label>
    <div class="flex gap-4">
      <Button on:click={convertVector}>Convert</Button>
      {#if labelOfConvertedVector?.length > 0}
        <P style={convertedVectorStyle}>{labelOfConvertedVector}</P>
        <Button on:click={() => saveSSVC(vectorInput)}>
          <i class="bx bx-save me-2 text-xl"></i>Save</Button
        >
      {/if}
    </div>
    <Label class="mb-4 mt-6 text-lg">Calculate new SSVC vector</Label>
    <Button on:click={() => (startedCalculation = true)}>Calculate</Button>
  {:else}
    <div class="mb-4 flex gap-4">
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
        <Button>
          <i class="bx bx-save me-2 text-xl"></i>Save</Button
        >
      {/if}
    {/if}
  {/if}
</div>
