<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Button, Label, StepIndicator } from "flowbite-svelte";
  import {
    createIsoTimeStringForSSVC,
    getDecision,
    parseDecisionTree,
    vectorStart,
    type SSVCAction,
    type SSVCDecision,
    type SSVCDecisionChild,
    type SSVCDecisionChildCombinationItem,
    type SSVCDecisionCombination,
    type SSVCObject,
    type SSVCOption
  } from "./SSVCCalculator";
  import { createEventDispatcher, onMount } from "svelte";
  import { request } from "$lib/request";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import { getErrorDetails, type ErrorDetails } from "$lib/Errors/error";
  import ComplexDecision from "./ComplexDecision.svelte";
  import SsvcInput from "./SSVCInput.svelte";

  const dispatch = createEventDispatcher();

  export let disabled = false;
  export let documentID: string;
  export let allowEditing: any;
  export let isEditing = false;
  export let vectorInput = "";

  let isVectorInputValid = false;
  let startedCalculation = false;
  let isComplex = false;
  let currentStep = 0;
  let steps: string[] = [];
  let mainDecisions: SSVCDecision[] = [];
  let decisionPoints: SSVCDecision[] = [];
  let decisionsTable: SSVCDecisionCombination[] = [];
  let userDecisions: SSVCDecisionCombination;
  let vector: string;
  let result: SSVCObject | null = null;
  let saveSSVCError: ErrorDetails | null;
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
    vector = vectorStart;
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

  /** Checks if `childCombos` contains `selectedOptions` */
  function doesContainChildCombo(
    selectedOptions: SSVCDecisionCombination,
    childCombos: SSVCDecisionChildCombinationItem[][]
  ): boolean {
    const doesContainCombo = childCombos.find((combos: SSVCDecisionChildCombinationItem[]) => {
      let count = 0;
      for (const option in selectedOptions) {
        const result = combos.find((combo: SSVCDecisionChildCombinationItem) => {
          return (
            combo.child_label === option &&
            combo.child_option_labels.includes(selectedOptions[option])
          );
        });
        if (result) count++;
      }
      return count === Object.keys(selectedOptions).length;
    });
    return !!doesContainCombo;
  }

  /** Finds out the selected option of a complex decision and then calls selectOption for further processing. */
  function calculateComplexOption() {
    const selectedChildOptions: SSVCDecisionCombination = {};
    mainDecisions[currentStep].children?.forEach((child: SSVCDecisionChild) => {
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
    const resultText: SSVCAction = Object.values(finalDecision)[0];
    const color = getOption(mainDecisions[currentStep], resultText)?.color;
    result = {
      vector: vector,
      label: resultText,
      color: color || "#000"
    };
  }

  function resetError() {
    saveSSVCError = null;
  }

  /** Removes part of the vector representing the last decision and removes the decision from the array
      containing all made decisions. */
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

  /** Save SSVC vector to the backend. */
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
      saveSSVCError = getErrorDetails(`Could not save SSVC.`, response);
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

<div id="ssvc-calc" class="flex w-full flex-col">
  {#if !isEditing}
    {#if !disabled}
      <div class="flex flex-row">
        {#if !vectorInput}
          <span class="h-6 text-lg">Please enter a SSVC</span>
        {/if}
        <button class="mr-auto h-6" {disabled} on:click={toggleEditing} aria-label="Edit SSVC"
          ><i class="bx bx-edit-alt ml-1"></i></button
        >
      </div>
    {/if}
  {:else if !startedCalculation}
    <div class="flex flex-col">
      <Label class="mb-4">
        <span>Enter a SSVC directly</span>
      </Label>
      <SsvcInput
        autofocus
        disabled={disabled || !isEditing}
        on:keyup={(e) => {
          if (e.detail.key === "Enter") saveSSVC(vectorInput);
          if (e.detail.key === "Escape") toggleEditing();
        }}
        on:input={resetError}
        bind:value={vectorInput}
        bind:isValid={isVectorInputValid}
      ></SsvcInput>
      <div class="mb-2 ml-auto flex flex-row gap-x-3">
        <Button
          color="light"
          outline
          size="xs"
          class="h-8"
          on:click={() => {
            isEditing = false;
          }}
        >
          Cancel
        </Button>
        <Button
          color="light"
          outline
          size="xs"
          class="h-8"
          {disabled}
          on:click={() => (startedCalculation = true)}
        >
          Evaluate
        </Button>
        <Button
          color="green"
          disabled={!isVectorInputValid}
          size="xs"
          class="h-8"
          on:click={() => saveSSVC(vectorInput)}>Save</Button
        >
      </div>
    </div>
  {:else}
    <div class="flex flex-row">
      <Label class="mb-4">Evaluate a SSVC</Label>
    </div>
    <div class="w-full">
      <span class="pt-[0.3rem] font-mono text-gray-400" color="gray"
        >Step {currentStep + 1}/{steps.length}</span
      >
      <StepIndicator
        glow={false}
        completedCustom=""
        currentCustom=""
        color="primary"
        size="h-1"
        currentStep={currentStep + 1}
        {steps}
        hideLabel
      />
      {#if steps[currentStep]}
        <div class="mt-2 mb-2">
          <span class="text-xl text-nowrap">{steps[currentStep]}</span>
        </div>
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
                or
                <Button
                  class="h-6"
                  outline
                  size="xs"
                  on:click={() => {
                    isComplex = true;
                  }}>Custom</Button
                >
              </div>
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
            <span style={resultStyle}>{result.label}</span>
          </Label>
          <Label class="text-xs text-gray-400">{vector}</Label>
        {/if}
      {/if}
      <div class="mt-4 flex flex-col">
        <div class="ml-auto flex flex-row gap-x-3">
          {#if currentStep > 0}
            <Button color="light" size="xs" class="h-6 p-3" on:click={stepBack}>Back</Button>
          {/if}
          <Button size="xs" color="light" class="h-6 p-3 text-nowrap" on:click={resetUserDecisions}>
            Restart
          </Button>
          {#if isComplex || result}
            <Button
              color={currentStep === steps.length - 1 ? "green" : "primary"}
              on:click={() => {
                if (currentStep === steps.length - 1) {
                  saveSSVC(vector);
                } else {
                  calculateComplexOption();
                }
              }}
              class="h-6 p-3"
              title="Calculate">Save</Button
            >
          {/if}
          <Button
            color="light"
            size="xs"
            class="h-6 p-3"
            title="Cancel SSVC input"
            on:click={() => {
              resetUserDecisions();
              startedCalculation = false;
              isEditing = false;
            }}>Cancel</Button
          >
        </div>
      </div>
    </div>
  {/if}
  <ErrorMessage error={saveSSVCError}></ErrorMessage>
</div>
