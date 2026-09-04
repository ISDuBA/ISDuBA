// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2023 Intevation GmbH <https://intevation.de>

import { describe, it, expect } from "vitest";
import {
  extractProducts,
  extractVulnerabilities,
  generateProductVulnerabilities
} from "./productvulnerabilities";
import { ProductStatusSymbol } from "./productvulnerabilitiestypes";
import type {
  CSAFDocumentv2_0,
  ListOfBranches,
  ProductTree,
  Vulnerabilities
} from "$lib/Advisories/types/csaf-2.0";
import type { DocumentLevelMetaData } from "$lib/Advisories/types/csaf-2.0";

const oneProductNotNested: ProductTree = {
  branches: [
    {
      name: "Product A",
      category: "product_version",
      product: {
        product_id: "123",
        name: "Product A"
      }
    }
  ]
};

const simpleNested: ProductTree = {
  branches: [
    {
      name: "Product C",
      category: "product_version",
      branches: [
        {
          name: "Product B",
          category: "product_version",
          product: {
            product_id: "123",
            name: "Product A"
          }
        }
      ]
    }
  ]
};

const complexNested: ProductTree = {
  branches: [
    {
      name: "Product ABC",
      category: "product_version",
      branches: [
        {
          name: "Product AB",
          category: "product_version",
          branches: [
            {
              name: "Product C",
              category: "product_version",
              product: {
                product_id: "8910",
                name: "Product C"
              }
            }
          ]
        },
        {
          name: "Product A",
          category: "product_version",
          product: {
            product_id: "123",
            name: "Product A"
          }
        },
        {
          name: "Product B",
          category: "product_version",
          product: {
            product_id: "3456",
            name: "Product B"
          }
        }
      ]
    }
  ]
};

const noVulnerabilities = {
  vulnerabilities: []
};

const vulnerability_wo_CVE = {
  vulnerabilities: [{}]
};

const vulnerability_empty_product_status = {
  vulnerabilities: [
    {
      cve: "CVE-2018-0171",
      product_status: {}
    }
  ]
};

const vulnerability_known_affected_empty = {
  vulnerabilities: [
    {
      cve: "CVE-2018-0171",
      product_status: {
        known_affected: []
      }
    }
  ]
};

const vulnerability_known_affected_filled = {
  vulnerabilities: [
    {
      cve: "CVE-2018-0171",
      product_status: {
        known_affected: ["123", "456"]
      }
    }
  ]
};

const documentDocument: DocumentLevelMetaData = {
  acknowledgments: [
    {
      names: ["qXhxe"],
      organization: "uTxaM",
      summary: "hXaqvYynsvb",
      urls: ["https://example.org"]
    }
  ],
  category: "ඊ",
  csaf_version: "2.0",
  lang: "X-2QpA4-P-IsDkdUY-Uy-W",
  notes: [{ category: "general", text: "QUOclBKPQXh" }],
  publisher: {
    category: "other",
    contact_details: "0z8eAxevc0G",
    name: "ozdkC6xS",
    namespace: "http://example.com"
  },
  source_lang: "X-KVzPzau",
  title: "xMJO",
  tracking: {
    aliases: ["1"],
    current_release_date: "2020-02-09T07:37:51.351498222Z",
    id: "test-csaf-2.0",
    initial_release_date: "2021-05-02T20:55:21.370337232Z",
    revision_history: [
      {
        date: "2022-11-15T21:06:17.574099808Z",
        legacy_version: "b",
        number:
          "6182776.3973455.6327-0.0.5940485.3761237y.9291154745.514072.0.3l9LNF2CL-c5.563701721.690DTGGso4y3a.0+BRQJC.PpY4gqc-e.EEv.cDEYK.si",
        summary: "1rL"
      }
    ],
    status: "final",
    version: "0.1707764392.0-9.0244721290y"
  }
};

const branches: ListOfBranches = [
  {
    name: "ABCDE",
    category: "product_version",
    branches: [
      {
        name: "Name CD",
        category: "product_version",
        branches: [
          {
            name: "Name C",
            category: "product_version",
            product: {
              product_id: "8910",
              name: "Product C"
            }
          },
          {
            name: "Name D",
            category: "product_version",
            product: {
              product_id: "1112",
              name: "Product D"
            }
          }
        ]
      },
      {
        name: "Name A",
        category: "product_version",
        product: {
          product_id: "123",
          name: "Product A"
        }
      },
      {
        name: "Name B",
        category: "product_version",
        product: {
          product_id: "3456",
          name: "Product B"
        }
      },
      {
        name: "Name E",
        category: "product_version",
        product: {
          product_id: "1314",
          name: "Product E"
        }
      }
    ]
  }
];

const jsonDocument: CSAFDocumentv2_0 = {
  document: { ...documentDocument },
  product_tree: {
    branches
  },
  vulnerabilities: [
    {
      cve: "CVE-2020-0174",
      product_status: {
        fixed: ["1112"]
      }
    },
    {
      cve: "CVE-2019-0171",
      product_status: {
        known_affected: ["123", "3456"]
      }
    },
    {
      cve: "CVE-2018-0172",
      product_status: {
        known_affected: ["8910"]
      }
    },
    {
      cve: "CVE-2016-0173",
      product_status: {
        known_not_affected: ["1314"],
        recommended: ["1314"]
      }
    }
  ]
};

describe("Productvulnerabilities test", () => {
  it("Product: parses empty object", () => {
    const result = extractProducts(null);
    expect(result.length).toBe(0);
  });
});

describe("Productvulnerabilities test", () => {
  it("Product: parses non nested list of products", () => {
    const clone = structuredClone(jsonDocument);
    if (clone.product_tree) {
      clone.product_tree.branches = oneProductNotNested.branches;
    }
    const result = extractProducts(clone);
    expect(result.length).toBe(1);
    expect(result[0].product_id).toBe("123");
    expect(result[0].name).toBe("Product A");
  });
});

describe("Productvulnerabilities test", () => {
  it("Product: parses simple nested list of products", () => {
    const clone = structuredClone(jsonDocument);
    if (clone.product_tree) {
      clone.product_tree.branches = simpleNested.branches;
    }
    const result = extractProducts(clone);
    expect(result.length).toBe(1);
    expect(result[0].product_id).toBe("123");
    expect(result[0].name).toBe("Product A");
  });
});

describe("Productvulnerabilities test", () => {
  it("Product: parses complex nested list of products", () => {
    const clone = structuredClone(jsonDocument);
    if (clone.product_tree) {
      clone.product_tree.branches = complexNested.branches;
    }
    const result = extractProducts(clone);
    expect(result.length).toBe(3);
    expect(result[0].product_id).toBe("8910");
    expect(result[0].name).toBe("Product C");
    expect(result[1].product_id).toBe("123");
    expect(result[1].name).toBe("Product A");
    expect(result[2].product_id).toBe("3456");
    expect(result[2].name).toBe("Product B");
  });
});

describe("Productvulnerabilities test", () => {
  it("Vulnerability: parses empty object", () => {
    const { vulnerabilities } = extractVulnerabilities({});
    expect(vulnerabilities.length).toBe(0);
  });
});

describe("Productvulnerabilities test", () => {
  it("Vulnerability: parses no vulnerabilities", () => {
    const { vulnerabilities } = extractVulnerabilities(noVulnerabilities);
    expect(vulnerabilities.length).toBe(0);
  });
});

describe("Productvulnerabilities test", () => {
  it("Vulnerability: parses vulnerability without cve", () => {
    const { vulnerabilities } = extractVulnerabilities(vulnerability_wo_CVE);
    expect(vulnerabilities.length).toBe(1);
  });
});

describe("Productvulnerabilities test", () => {
  it("Vulnerability: parses vulnerability with empty product_status", () => {
    const { vulnerabilities } = extractVulnerabilities(vulnerability_empty_product_status);
    expect(vulnerabilities.length).toBe(1);
  });
});

describe("Productvulnerabilities test", () => {
  it("Vulnerability: parses vulnerability with empty known_affected", () => {
    const { vulnerabilities } = extractVulnerabilities(vulnerability_known_affected_empty);
    expect(vulnerabilities.length).toBe(1);

    expect(Object.keys(vulnerabilities[0].known_affected!).length).toBe(0);
  });
});

describe("Productvulnerabilities test", () => {
  it("Vulnerability: parses vulnerability with filled known_affected", () => {
    const { vulnerabilities } = extractVulnerabilities(vulnerability_known_affected_filled);
    const value = vulnerabilities[0];
    expect(vulnerabilities.length).toBe(1);

    expect(Object.keys(value.known_affected!).length).toBe(2);

    expect(value.known_affected!["123"]).toBe("123");

    expect(value.known_affected!["456"]).toBe("456");
  });
});

describe("Productvulnerabilities test", () => {
  it("Crosstable: generate headers", () => {
    const clone = structuredClone(jsonDocument);
    const products = extractProducts(clone);
    const productLookup = products.reduce((o: any, n: any) => {
      o[n.product_id] = n.name;
      return o;
    }, {});
    const result = generateProductVulnerabilities(clone, products, productLookup);
    const header = result[0].map((c: any) => c.content);
    const expectedHeader = [
      "Product",
      "Total result",
      "CVE-2016-0173",
      "CVE-2018-0172",
      "CVE-2019-0171",
      "CVE-2020-0174"
    ];
    expect(result.length).toBeGreaterThan(0);
    expect(header).toStrictEqual(expectedHeader);
    expect(header.length).toBe((jsonDocument.vulnerabilities as Vulnerabilities).length + 2);
  });
});

describe("Productvulnerabilities test", () => {
  it("Crosstable: generate body", () => {
    const products = extractProducts(jsonDocument);
    const productLookup = products.reduce((o: any, n: any) => {
      o[n.product_id] = n.name;
      return o;
    }, {});
    const result = generateProductVulnerabilities(jsonDocument, products, productLookup);
    const line1 = result[1];
    const line2 = result[2];
    const line3 = result[3];
    const line4 = result[4];
    const line5 = result[5];
    const PRODUCT_COLUMN = 0;
    const TOTAL_COLUMN = 1;
    const CVE_2016_0173_COLUMN = 2;
    const CVE_2018_0172_COLUMN = 3;
    const CVE_2019_0171_COLUMN = 4;
    const CVE_2020_0174_COLUMN = 5;
    expect(result.length).toBe(6);
    // Product A
    expect(line1[PRODUCT_COLUMN].content).toBe("123");
    expect(line1[TOTAL_COLUMN].content).toBe("K");
    expect(line1[CVE_2016_0173_COLUMN].content).toBe("");
    expect(line1[CVE_2018_0172_COLUMN].content).toBe("");
    expect(line1[CVE_2019_0171_COLUMN].content).toBe(ProductStatusSymbol.KNOWN_AFFECTED);
    expect(line1[CVE_2020_0174_COLUMN].content).toBe("");
    // Product B
    expect(line2[PRODUCT_COLUMN].content).toBe("3456");
    expect(line2[TOTAL_COLUMN].content).toBe("K");
    expect(line2[CVE_2016_0173_COLUMN].content).toBe("");
    expect(line2[CVE_2018_0172_COLUMN].content).toBe("");
    expect(line2[CVE_2019_0171_COLUMN].content).toBe(ProductStatusSymbol.KNOWN_AFFECTED);
    expect(line2[CVE_2020_0174_COLUMN].content).toBe("");
    // Product C
    expect(line3[PRODUCT_COLUMN].content).toBe("8910");
    expect(line3[TOTAL_COLUMN].content).toBe("K");
    expect(line3[CVE_2016_0173_COLUMN].content).toBe("");
    expect(line3[CVE_2018_0172_COLUMN].content).toBe(ProductStatusSymbol.KNOWN_AFFECTED);
    expect(line3[CVE_2019_0171_COLUMN].content).toBe("");
    expect(line3[CVE_2020_0174_COLUMN].content).toBe("");
    // Product D
    expect(line4[PRODUCT_COLUMN].content).toBe("1112");
    expect(line4[TOTAL_COLUMN].content).toBe("F");
    expect(line4[CVE_2016_0173_COLUMN].content).toBe("");
    expect(line4[CVE_2018_0172_COLUMN].content).toBe("");
    expect(line4[CVE_2019_0171_COLUMN].content).toBe("");
    expect(line4[CVE_2020_0174_COLUMN].content).toBe(ProductStatusSymbol.FIXED);
    //Product E
    expect(line5[PRODUCT_COLUMN].content).toBe("1314");
    expect(line5[TOTAL_COLUMN].content).toBe("N");
    expect(line5[CVE_2016_0173_COLUMN].content).toBe(
      ProductStatusSymbol.NOT_AFFECTED + ProductStatusSymbol.RECOMMENDED
    );
    expect(line5[CVE_2018_0172_COLUMN].content).toBe("");
    expect(line5[CVE_2019_0171_COLUMN].content).toBe("");
    expect(line5[CVE_2020_0174_COLUMN].content).toBe("");
  });
});
