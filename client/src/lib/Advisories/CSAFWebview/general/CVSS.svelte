<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2023 Intevation GmbH <https://intevation.de>
-->
<script lang="ts">
  export let baseScore = "";
  export let baseSeverity = "";

  const getSeverityClass = (severity: string) => {
    return severity.toLowerCase();
  };
</script>

<div class={"score " + getSeverityClass(baseSeverity)}>
  <span class="baseScore">{baseScore}</span>
  <span class="baseSeverity">({baseSeverity})</span>
</div>

<style>
  .score.none,
  .score.low,
  .score.medium,
  .score.high,
  .score.critical {
    color: #ffffff;
  }

  .score.none {
    background: #53aa33;
    border: 2px solid #53aa33;
  }

  .score.low {
    background: #ffcb0d;
    border: 2px solid #ffcb0d;
  }

  .score.medium {
    background: #f9a009;
    border: 2px solid #f9a009;
  }

  .score.high {
    background: #df3d03;
    border: 2px solid #df3d03;
  }

  .score.critical {
    background: #cc0500;
    border: 2px solid #cc0500;
  }

  .score span {
    font-weight: bold;
    width: 100%;
  }

  .score {
    top: -36px;
    right: 0;
    padding: 0 0.4em;
    margin: 0 15px;
    border: 2px solid #666666;
    background: #dddddd;
    font-size: 11px;
    border-radius: 10px;
    width: 100px;
    height: auto;
    line-height: 150%;
    text-align: center;
  }

  .baseScore {
    display: block;
    font-size: 32px;
    line-height: 32px;
    font-weight: normal;
    margin-top: 4px;
  }
  .baseSeverity {
    font-size: 16px;
    font-weight: normal;
    margin-bottom: 5px;
    display: block;
  }
</style>
