<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { onMount } from "svelte";
  import { Button, Card, Heading, Label, Radio, StepIndicator, P, Tooltip } from "flowbite-svelte";

  let currentStep = 0;
  let steps: string[] = [];
  let mainDecisions: any[] = [];
  let decision_points: any[] = [];
  let decisions_table: any[] = [];
  let userDecisions: any = {};
  const vectorBeginning = "SSVCv2/";
  let vector: string;

  onMount(async () => {
    if (true) {
      loadDecisionTree();
    }
  });

  function loadDecisionTree(): Promise<void> {
    return new Promise((resolve) => {
      resetDecisions();
      fetch("CISA-Coordinator.json").then((response) => {
        response.json().then((json) => {
          const addedPoints: string[] = [];
          decision_points.push(...json.decision_points);
          decisions_table.push(...json.decisions_table);
          for (let i = decision_points.length - 1; i >= 0; i--) {
            const decision = decision_points[i];
            if (!addedPoints.includes(decision.label)) {
              mainDecisions.push(decision);
              if (decision.decision_type === "complex") {
                for (const child of decision.children) {
                  addedPoints.push(child.label);
                }
              } else {
                addedPoints.push(decision.label);
              }
            }
          }
          mainDecisions = mainDecisions.reverse();
          steps = mainDecisions.map((element) => element.label);
          resolve();
        });
      });
    });
  }

  function resetDecisions() {
    resetUserDecisions();
    steps = [];
    mainDecisions = [];
    decision_points = [];
    decisions_table = [];
  }

  function resetUserDecisions() {
    userDecisions = {};
    currentStep = 0;
    vector = vectorBeginning;
  }

  function getDecision(label: string): any {
    return decision_points.find((element) => element.label === label);
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

  function calculateComplexOption(event: any) {
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

  function calculateResult(): any {
    let filteredDecisions = decisions_table;
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
    const isoString = new Date().toISOString();
    extendVector(`${mainDecisions[currentStep].key}:${option.key}/${isoString}`);
    return Object.values(finalDecision)[0];
  }

  function stepBack() {
    if (currentStep === steps.length - 1) {
      // Delete ISO string and last key pair
      vector = vector.slice(0, -28);
    } else {
      // Delete key pair
      vector = vector.slice(0, -4);
    }
    currentStep--;
    const keyOfLastDecision = Object.keys(userDecisions)[Object.keys(userDecisions).length - 1];
    delete userDecisions[keyOfLastDecision];
  }
</script>

<div id="ssvc-calc" class="pe-4">
  <h1 class="mb-3 text-lg">Calculate SSVC</h1>
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
  <Heading class="mb-6 max-w-fit">{steps[currentStep]}</Heading>
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
        <div class="flex gap-3">
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
            <div class="flex flex-row gap-x-6 gap-y-4">
              {#each mainDecisions[currentStep].children as child}
                {@const childOptions = getDecision(child.label).options}
                <div class="min-w-60">
                  <h5 class="mb-2 text-2xl font-bold tracking-tight text-gray-900 dark:text-white">
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
    {:else}
      <Label class="text-lg">Result: {calculateResult()}</Label>
      <Label class="text-lg">Vector: {vector}</Label>
      <Button>
        <i class="bx bx-save me-2 text-xl"></i>Save</Button
      >
    {/if}
  {/if}
</div>
