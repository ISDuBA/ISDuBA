<!--
 This file is Free Software under the MIT License
 without warranty, see README.md and LICENSES/MIT.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { Diff2HtmlUI, type Diff2HtmlUIConfig } from "diff2html/lib/ui/js/diff2html-ui-slim";

  export let diff: string;

  $: {
    const diffElement = document.getElementById("diff");
    if (diff?.length > 0 && diffElement) {
      const config: Diff2HtmlUIConfig = {
        colorScheme: "light",
        drawFileList: false
      };
      const diff2htmlUi = new Diff2HtmlUI(diffElement, diff, config);
      diff2htmlUi.draw();
    }
  }
</script>

<div id="diff"></div>
