// This file is Free Software under the MIT License
// without warranty, see README.md and LICENSES/MIT.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2023 Intevation GmbH <https://intevation.de>

import { describe, it, expect } from "vitest";
import { convertToDocModel } from "$lib/Advisories/CSAFWebview/docmodel/docmodel";
import {
  EMPTY,
  Status,
  CSAFDocProps,
  type DocModel,
  type DocModelKey,
  TLP
} from "$lib/Advisories/CSAFWebview/docmodel/docmodeltypes";

const allEmpty = (docModel: DocModel, properties: DocModelKey[]) => {
  properties.forEach((p) => {
    expect(docModel[p]).toBe(EMPTY);
  });
};

const allDisabled = (docModel: DocModel, properties: DocModelKey[]) => {
  return (
    properties
      .map((p: DocModelKey) => {
        return docModel[p] as boolean;
      })
      .some((p: boolean) => p) === false
  );
};

describe("docmodel test", () => {
  it("converts an empty object", () => {
    const doc = {};
    const docModel: DocModel = convertToDocModel(doc);
    allEmpty(docModel, ["id", "lang", "lastUpdate", "published", "status", "title", "csafVersion"]);
    expect(
      allDisabled(docModel, [
        "isDocPresent",
        "isTrackingPresent",
        "isDistributionPresent",
        "isTLPPresent"
      ])
    ).true;
    const publisher = docModel.publisher;
    expect(publisher.name).toBe("");
    expect(publisher.category).toBe("");
    expect(publisher.namespace).toBe("");
    expect(docModel.tlp.label).toBe(EMPTY);
  });
});

describe("docmodel test", () => {
  it("converts an object with document property", () => {
    const doc = { [CSAFDocProps.DOCUMENT]: {} };
    const docModel: DocModel = convertToDocModel(doc);
    allEmpty(docModel, ["id", "lang", "lastUpdate", "published", "status", "title", "csafVersion"]);
    expect(allDisabled(docModel, ["isTrackingPresent", "isDistributionPresent", "isTLPPresent"]))
      .true;
    const publisher = docModel.publisher;
    expect(publisher.name).toBe("");
    expect(publisher.category).toBe("");
    expect(publisher.namespace).toBe("");
    expect(docModel.trackingVersion).toBe("");
    expect(docModel.tlp.label).toBe(EMPTY);
  });
});

describe("docmodel test", () => {
  it("converts an object with document property and title", () => {
    const TITLE = "test";
    const doc = { [CSAFDocProps.DOCUMENT]: { [CSAFDocProps.TITLE]: TITLE } };
    const docModel: DocModel = convertToDocModel(doc);
    allEmpty(docModel, ["id", "lang", "lastUpdate", "published", "status", "csafVersion"]);
    expect(docModel.title).toBe(TITLE);
    expect(allDisabled(docModel, ["isTrackingPresent", "isDistributionPresent", "isTLPPresent"]))
      .true;
    const publisher = docModel.publisher;
    expect(publisher.name).toBe("");
    expect(publisher.category).toBe("");
    expect(publisher.namespace).toBe("");
    expect(docModel.trackingVersion).toBe("");
    expect(docModel.tlp.label).toBe(EMPTY);
  });
});

describe("docmodel test", () => {
  it("converts an object with document property and csafVersion", () => {
    const version = "2.0";
    const doc = { [CSAFDocProps.DOCUMENT]: { [CSAFDocProps.CSAFVERSION]: version } };
    const docModel: DocModel = convertToDocModel(doc);
    allEmpty(docModel, ["id", "lang", "lastUpdate", "published", "status", "title"]);
    expect(docModel.csafVersion).toBe(version);
    expect(allDisabled(docModel, ["isTrackingPresent", "isDistributionPresent", "isTLPPresent"]))
      .true;
    const publisher = docModel.publisher;
    expect(publisher.name).toBe("");
    expect(publisher.category).toBe("");
    expect(publisher.namespace).toBe("");
    expect(docModel.trackingVersion).toBe("");
    expect(docModel.tlp.label).toBe(EMPTY);
  });
});

describe("docmodel test", () => {
  it("converts an object with document property and publisher information", () => {
    const version = "2.0";
    const publisherName = "ABC";
    const publisherCategory = "coordinator";
    const publisherNameSpace = "https://www.example.com";
    const publisherIssuingAuthority =
      "This service is provided as it is. It is free for everybody.";
    const doc = {
      [CSAFDocProps.DOCUMENT]: {
        [CSAFDocProps.CSAFVERSION]: version,
        [CSAFDocProps.PUBLISHER]: {
          [CSAFDocProps.PUBLISHER_NAME]: publisherName,
          [CSAFDocProps.PUBLISHER_CATEGORY]: publisherCategory,
          [CSAFDocProps.PUBLISHER_NAMESPACE]: publisherNameSpace,
          [CSAFDocProps.ISSUING_AUTHORITY]: publisherIssuingAuthority
        }
      }
    };
    const docModel: DocModel = convertToDocModel(doc);
    allEmpty(docModel, ["id", "lang", "lastUpdate", "published", "status", "title"]);
    expect(docModel.csafVersion).toBe(version);
    expect(allDisabled(docModel, ["isTrackingPresent", "isDistributionPresent", "isTLPPresent"]))
      .true;
    const publisher = docModel.publisher;
    expect(publisher.name).toBe(publisherName);
    expect(publisher.category).toBe(publisherCategory);
    expect(publisher.namespace).toBe(publisherNameSpace);
    expect(publisher.issuing_authority).toBe(publisherIssuingAuthority);
    expect(docModel.trackingVersion).toBe("");
    expect(docModel.tlp.label).toBe(EMPTY);
  });
});

describe("docmodel test", () => {
  it("converts an object with document property and language", () => {
    const lang = "de_DE";
    const doc = { [CSAFDocProps.DOCUMENT]: { [CSAFDocProps.LANG]: lang } };
    const docModel: DocModel = convertToDocModel(doc);
    allEmpty(docModel, ["id", "lastUpdate", "published", "status", "title", "csafVersion"]);
    expect(docModel.lang).toBe(lang);
    expect(allDisabled(docModel, ["isTrackingPresent", "isDistributionPresent", "isTLPPresent"]))
      .true;
    const publisher = docModel.publisher;
    expect(publisher.name).toBe("");
    expect(publisher.category).toBe("");
    expect(publisher.namespace).toBe("");
    expect(docModel.trackingVersion).toBe("");
    expect(docModel.tlp.label).toBe(EMPTY);
  });
});

describe("docmodel test", () => {
  it("converts an object with document property and category", () => {
    const category = "csaf_security_advisory";
    const doc = { [CSAFDocProps.DOCUMENT]: { [CSAFDocProps.CATEGORY]: category } };
    const docModel: DocModel = convertToDocModel(doc);
    allEmpty(docModel, ["id", "lastUpdate", "published", "status", "title", "csafVersion", "lang"]);
    expect(docModel.category).toBe(category);
    expect(allDisabled(docModel, ["isTrackingPresent", "isDistributionPresent", "isTLPPresent"]))
      .true;
    const publisher = docModel.publisher;
    expect(publisher.name).toBe("");
    expect(publisher.category).toBe("");
    expect(publisher.namespace).toBe("");
    expect(docModel.trackingVersion).toBe("");
    expect(docModel.tlp.label).toBe(EMPTY);
  });
});

describe("docmodel test", () => {
  it("converts an object with tracking property", () => {
    const doc = { [CSAFDocProps.DOCUMENT]: { [CSAFDocProps.TRACKING]: {} } };
    const docModel: DocModel = convertToDocModel(doc);
    allEmpty(docModel, ["id", "lang", "lastUpdate", "published", "title", "csafVersion"]);
    expect(allDisabled(docModel, ["isDistributionPresent", "isTLPPresent"])).true;
    expect(docModel.status).toBe(Status.ERROR);
    const publisher = docModel.publisher;
    expect(publisher.name).toBe("");
    expect(publisher.category).toBe("");
    expect(publisher.namespace).toBe("");
    expect(docModel.trackingVersion).toBe("");
    expect(docModel.tlp.label).toBe(EMPTY);
  });
});

describe("docmodel test", () => {
  it("converts an object with tracking property and id", () => {
    const testId = "123";
    const doc = {
      [CSAFDocProps.DOCUMENT]: { [CSAFDocProps.TRACKING]: { [CSAFDocProps.ID]: testId } }
    };
    const docModel: DocModel = convertToDocModel(doc);
    allEmpty(docModel, ["lang", "lastUpdate", "published", "title", "csafVersion"]);
    expect(allDisabled(docModel, ["isDistributionPresent", "isTLPPresent"])).true;
    expect(docModel.status).toBe(Status.ERROR);
    expect(docModel.id).toBe(testId);
    const publisher = docModel.publisher;
    expect(publisher.name).toBe("");
    expect(publisher.category).toBe("");
    expect(publisher.namespace).toBe("");
    expect(docModel.trackingVersion).toBe("");
    expect(docModel.tlp.label).toBe(EMPTY);
  });
});

describe("docmodel test", () => {
  it("converts an object with tracking property and status", () => {
    const status: string = Status.final;
    const doc = {
      [CSAFDocProps.DOCUMENT]: { [CSAFDocProps.TRACKING]: { [CSAFDocProps.STATUS]: status } }
    };
    const docModel: DocModel = convertToDocModel(doc);
    allEmpty(docModel, ["id", "lang", "lastUpdate", "published", "title", "csafVersion"]);
    expect(allDisabled(docModel, ["isDistributionPresent", "isTLPPresent"])).true;
    expect(docModel.status).toBe(Status.final);
    const publisher = docModel.publisher;
    expect(publisher.name).toBe("");
    expect(publisher.category).toBe("");
    expect(publisher.namespace).toBe("");
    expect(docModel.trackingVersion).toBe("");
    expect(docModel.tlp.label).toBe(EMPTY);
  });
});

describe("docmodel test", () => {
  it("converts an object with tracking property and trackingVersion", () => {
    const version = 1;
    const doc = {
      [CSAFDocProps.DOCUMENT]: {
        [CSAFDocProps.TRACKING]: { [CSAFDocProps.TRACKINGVERSION]: version }
      }
    };
    const docModel: DocModel = convertToDocModel(doc);
    allEmpty(docModel, ["id", "lang", "lastUpdate", "published", "title", "csafVersion"]);
    expect(allDisabled(docModel, ["isDistributionPresent", "isTLPPresent"])).true;
    expect(docModel.status).toBe(Status.ERROR);
    const publisher = docModel.publisher;
    expect(publisher.name).toBe("");
    expect(publisher.category).toBe("");
    expect(publisher.namespace).toBe("");
    expect(docModel.trackingVersion).toBe(version);
    expect(docModel.tlp.label).toBe(EMPTY);
  });
});

describe("docmodel test", () => {
  it("converts an object with tracking property and published", () => {
    const published: string = new Date().toISOString();
    const doc = {
      [CSAFDocProps.DOCUMENT]: {
        [CSAFDocProps.TRACKING]: { [CSAFDocProps.INITIALRELEASEDATE]: published }
      }
    };
    const docModel: DocModel = convertToDocModel(doc);
    allEmpty(docModel, ["id", "lang", "lastUpdate", "title", "csafVersion"]);
    expect(allDisabled(docModel, ["isDistributionPresent", "isTLPPresent"])).true;
    expect(docModel.status).toBe(Status.ERROR);
    const publisher = docModel.publisher;
    expect(publisher.name).toBe("");
    expect(publisher.category).toBe("");
    expect(publisher.namespace).toBe("");
    expect(docModel.trackingVersion).toBe("");
    expect(docModel.tlp.label).toBe(EMPTY);
  });
});

describe("docmodel test", () => {
  it("converts an object with tracking property and lastUpdate", () => {
    const lastUpdate: string = new Date().toISOString();
    const doc = {
      [CSAFDocProps.DOCUMENT]: {
        [CSAFDocProps.TRACKING]: { [CSAFDocProps.CURRENTRELEASEDATE]: lastUpdate }
      }
    };
    const docModel: DocModel = convertToDocModel(doc);
    allEmpty(docModel, ["id", "lang", "published", "title", "csafVersion"]);
    expect(allDisabled(docModel, ["isDistributionPresent", "isTLPPresent"])).true;
    expect(docModel.status).toBe(Status.ERROR);
    const publisher = docModel.publisher;
    expect(publisher.name).toBe("");
    expect(publisher.category).toBe("");
    expect(publisher.namespace).toBe("");
    expect(docModel.trackingVersion).toBe("");
    expect(docModel.tlp.label).toBe(EMPTY);
  });
});

describe("docmodel test", () => {
  it("converts an object with document property and distribution", () => {
    const doc = { [CSAFDocProps.DOCUMENT]: { [CSAFDocProps.DISTRIBUTION]: {} } };
    const docModel: DocModel = convertToDocModel(doc);
    allEmpty(docModel, ["id", "lang", "lastUpdate", "published", "status", "title", "csafVersion"]);
    expect(allDisabled(docModel, ["isTrackingPresent", "isTLPPresent"])).true;
    const publisher = docModel.publisher;
    expect(publisher.name).toBe("");
    expect(publisher.category).toBe("");
    expect(publisher.namespace).toBe("");
    expect(docModel.trackingVersion).toBe("");
    expect(docModel.tlp.label).toBe(EMPTY);
  });
});

describe("docmodel test", () => {
  it("converts an object with document property and distribution and TLP", () => {
    const doc = {
      [CSAFDocProps.DOCUMENT]: { [CSAFDocProps.DISTRIBUTION]: { [CSAFDocProps.TLP]: {} } }
    };
    const docModel: DocModel = convertToDocModel(doc);
    allEmpty(docModel, ["id", "lang", "lastUpdate", "published", "status", "title", "csafVersion"]);
    expect(allDisabled(docModel, ["isTrackingPresent"])).true;
    const publisher = docModel.publisher;
    expect(publisher.name).toBe("");
    expect(publisher.category).toBe("");
    expect(publisher.namespace).toBe("");
    expect(docModel.trackingVersion).toBe("");
    expect(docModel.tlp.label).toBe(EMPTY);
  });
});

describe("docmodel test", () => {
  it("converts an object with document property and distribution and TLP and label", () => {
    const tlpLabel: string = TLP.RED;
    const doc = {
      [CSAFDocProps.DOCUMENT]: {
        [CSAFDocProps.DISTRIBUTION]: { [CSAFDocProps.TLP]: { [CSAFDocProps.LABEL]: tlpLabel } }
      }
    };
    const docModel: DocModel = convertToDocModel(doc);
    allEmpty(docModel, ["id", "lang", "lastUpdate", "published", "status", "title", "csafVersion"]);
    expect(allDisabled(docModel, ["isTrackingPresent"])).true;
    expect(docModel.tlp.label).toBe(TLP.RED);
    const publisher = docModel.publisher;
    expect(publisher.name).toBe("");
    expect(publisher.category).toBe("");
    expect(publisher.namespace).toBe("");
    expect(docModel.trackingVersion).toBe("");
  });
});

describe("docmodel test", () => {
  it("converts an object with document property and distribution and TLP and wrong label", () => {
    const tlpLabel = "bananas";
    const doc = {
      [CSAFDocProps.DOCUMENT]: {
        [CSAFDocProps.DISTRIBUTION]: { [CSAFDocProps.TLP]: { [CSAFDocProps.LABEL]: tlpLabel } }
      }
    };
    const docModel: DocModel = convertToDocModel(doc);
    allEmpty(docModel, ["id", "lang", "lastUpdate", "published", "status", "title", "csafVersion"]);
    expect(allDisabled(docModel, ["isTrackingPresent"])).true;
    expect(docModel.tlp.label).toBe(TLP.ERROR);
    const publisher = docModel.publisher;
    expect(publisher.name).toBe("");
    expect(publisher.category).toBe("");
    expect(publisher.namespace).toBe("");
    expect(docModel.trackingVersion).toBe("");
  });
});
