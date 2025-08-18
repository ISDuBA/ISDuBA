<!--
 This file is Free Software under the Apache-2.0 License
 without warranty, see README.md and LICENSES/Apache-2.0.txt for details.

 SPDX-License-Identifier: Apache-2.0

 SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
 Software-Engineering: 2024 Intevation GmbH <https://intevation.de>
-->

<script lang="ts">
  import { getErrorDetails } from "$lib/Errors/error";
  import { type UploadInfo } from "$lib/Sources/source";
  import Upload from "$lib/Upload.svelte";
  import { request } from "$lib/request";

  const uploadDocuments = async (files: FileList): Promise<UploadInfo[]> => {
    let uploadInfo = [];
    for (const file of files) {
      let info: UploadInfo = { success: true };
      const formData = new FormData();
      formData.append("file", file);
      const resp = await request(`/api/documents`, "POST", formData);
      if (resp.error) {
        info.success = false;
        let details = getErrorDetails(`Could not upload file`, resp);
        info.message = `${details.message} ${details.details}`;
      }
      uploadInfo.push(info);
    }
    return uploadInfo;
  };
</script>

<Upload upload={uploadDocuments} label="Upload a document"></Upload>
