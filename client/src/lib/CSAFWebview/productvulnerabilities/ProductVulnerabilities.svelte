<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { appStore } from "$lib/store";
  import { ProductStatusSymbol } from "./productvulnerabilitiestypes";
  let headerColumns: string[] = [];
  let productLines: string[][];
  $: if ($appStore.webview.doc) {
    const vulnerabilities = [...$appStore.webview.doc.productVulnerabilities];
    // eslint-disable-next-line  @typescript-eslint/no-non-null-assertion
    headerColumns = vulnerabilities.shift()!;
    productLines = vulnerabilities;
  }
  /**
   * openProduct opens the according product given via href.
   * @param e
   */
  const openProduct = (e: Event) => {
    // eslint-disable-next-line  @typescript-eslint/no-non-null-assertion
    let product: string = (e.target as Element).getAttribute("href")!;
    appStore.setProductTreeSectionVisible();
    appStore.setSelectedProduct(product);
    appStore.unshiftHistory((e.target as Element).id);
    e.preventDefault();
  };

  /**
   * openCVE opens the CVE given via href.
   * @param e
   */
  const openCVE = (e: Event) => {
    // eslint-disable-next-line  @typescript-eslint/no-non-null-assertion
    let CVE: string = (e.target as Element).getAttribute("href")!;
    appStore.setSelectedCVE(CVE);
    appStore.unshiftHistory((e.target as Element).id);
    appStore.setVulnerabilitiesSectionVisible();
    e.preventDefault();
  };
</script>

<div class="crosstable-overview">
  {#if productLines.length > 0}
    <div class="crosstable-container">
      <div class="crosstable">
        <table>
          <thead>
            <tr class="crosstable-header-row">
              {#each headerColumns as column, index}
                {#if index == 0}
                  <th class="productname">{column}</th>
                {:else if index == 1}
                  <th class="total">{column}</th>
                {:else}
                  <th class="cve"
                    ><a id={crypto.randomUUID()} on:click={openCVE} href={column}>{column}</a></th
                  >
                {/if}
              {/each}
            </tr>
          </thead>
          <tbody>
            {#each productLines as line}
              <tr>
                {#each line as column, index}
                  {#if index < 1}
                    <td class="productname"
                      ><a id={crypto.randomUUID()} on:click={openProduct} href={column}
                        >{$appStore.webview.doc.productsByID[column]} ({column})</a
                      ></td
                    >
                  {:else if column === "N.A"}
                    <td class="affectionstate">{column}</td>
                  {:else}
                    <td class="affectionstate">
                      {#if column === ProductStatusSymbol.NOT_AFFECTED + ProductStatusSymbol.RECOMMENDED}
                        <i class="bx bx-heart" />
                        <i class="bx b-minus" />
                      {:else}
                        <i
                          class:bx={true}
                          class:bx-x={column === ProductStatusSymbol.KNOWN_AFFECTED}
                          class:bx-check={column === ProductStatusSymbol.FIXED}
                          class:bx-error={column === ProductStatusSymbol.UNDER_INVESTIGATION}
                          class:bx-minus={column === ProductStatusSymbol.NOT_AFFECTED}
                          class:bx-heart={column === ProductStatusSymbol.RECOMMENDED}
                        />
                      {/if}
                    </td>
                  {/if}
                {/each}
              </tr>
            {/each}
          </tbody>
        </table>
      </div>
    </div>
    <div class="legend">
      <h5>Legend</h5>
      <dl>
        <dt><i class="bx bx-check" /></dt>
        <dd>Fixed</dd>
        <dt><i class="bx bx-error" /></dt>
        <dd>Under investigation</dd>
        <dt><i class="bx bx-x" /></dt>
        <dd>Known affected</dd>
        <dt><i class="bx bx-minus" /></dt>
        <dd>Not affected</dd>
        <dt><i class="bx bx-heart" /></dt>
        <dd>Recommended</dd>
      </dl>
    </div>
  {/if}
</div>
