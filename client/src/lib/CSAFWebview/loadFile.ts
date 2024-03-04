// This file is Free Software under the MIT License
// without warranty, see README.md and LICENSES/MIT.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2023 Intevation GmbH <https://intevation.de>

import { convertToDocModel } from "$lib/CSAFWebview/docmodel/docmodel";
import { appStore } from "$lib/store";
/**
 * loadFile loads files via FileReader.
 * @param csafFile
 */
const loadFile = (csafFile: File) => {
  const fileReader: FileReader = new FileReader();
  let jsonDocument = {};
  fileReader.onload = (e: ProgressEvent<FileReader>) => {
    if (e.target) {
      try {
        jsonDocument = JSON.parse(e.target.result as string);
      } catch (_) {
        /*
						Treat unparsable documents as empty documents
					   	The according errors will be reflected in the converted
					   	DocModel.
					*/
      }
      const docModel = convertToDocModel(jsonDocument);
      appStore.setDocument(docModel);
    }
  };
  fileReader.readAsText(csafFile);
};

export { loadFile };
