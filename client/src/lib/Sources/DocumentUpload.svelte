<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { getErrorDetails, type ErrorDetails } from "$lib/Errors/error";
  import ErrorMessage from "$lib/Errors/ErrorMessage.svelte";
  import Upload from "$lib/Upload.svelte";
  import { request } from "$lib/request";

  let uploadError: ErrorDetails | null;

  const uploadDocuments = async (files: FileList) => {
    for (const file of files) {
      const formData = new FormData();
      formData.append("file", file);
      const resp = await request(`/api/documents`, "POST", formData);
      if (resp.error) {
        uploadError = getErrorDetails(`Could not upload file`, resp);
      }
    }
  };
</script>

<Upload upload={uploadDocuments} label="Upload a document"></Upload>
<ErrorMessage error={uploadError}></ErrorMessage>
