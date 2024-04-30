<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import "diff2html/bundles/css/diff2html.min.css";
  import { Diff2HtmlUI, type Diff2HtmlUIConfig } from "diff2html/lib/ui/js/diff2html-ui-slim";
  import { onMount } from "svelte";
  import { appStore } from "$lib/store";
  import ErrorMessage from "$lib/Messages/ErrorMessage.svelte";
  import { request } from "$lib/utils";
  import { ColorSchemeType } from "diff2html/lib/types";

  let diff: string;
  let error: string;

  onMount(async () => {
    if ($appStore.app.isUserLoggedIn) {
      error = "";
      const response = await request("advisory.diff", "GET");
      if (response.ok) {
        diff = response.content;
      } else if (response.error) {
        error = response.error;
      }

      const diffElement = document.getElementById("diff");
      if (diff?.length > 0 && diffElement) {
        const config: Diff2HtmlUIConfig = {
          colorScheme: ColorSchemeType.LIGHT,
          drawFileList: false
        };
        const diff2htmlUi = new Diff2HtmlUI(diffElement, diff, config);
        diff2htmlUi.draw();
      }
    }
  });
</script>

<div>
  <ErrorMessage message={error}></ErrorMessage>
  <div id="diff"></div>
</div>
