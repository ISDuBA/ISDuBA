// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2023 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
// Software-Engineering: 2023 Intevation GmbH <https://intevation.de>

import {
  ProductStatusSymbol,
  type Vulnerability,
  type Product,
  type Relationship,
  type ProductStatus_t,
  type ProductStatus_t_Key,
  type VulnerabilitesExtractionResult
} from "./productvulnerabilitiestypes";

/**
 * generateProductVulnerabilities generates data for product vulnerabilites overview.
 * @param jsonDocument
 * @param products all products
 * @param productLookup {id:name}
 * @returns product vulnerabilities crosstable as [[]]
 * E.g. [
 * [
 *   'Product',
 *   'Total result',
 *   'CVE-2016-0173',
 *   'CVE-2018-0172',
 *   'CVE-2019-0171',
 *   'CVE-2020-0174'
 * ],
 * [ '123', 'K', '', '', 'K', '' ],
 * [ '3456', 'K', '', '', 'K', '' ],
 * [ '8910', 'K', '', 'K', '', '' ],
 * [ '1112', 'F', '', '', '', 'F' ],
 * [ '1314', 'N', 'NR', '', '', '' ]
 * ]
 */
const generateProductVulnerabilities = (jsonDocument: any, products: any, productLookup: any) => {
  const { vulnerabilities, relevantProducts } = extractVulnerabilities(jsonDocument);
  products = products.filter((product: Product) => {
    return relevantProducts[product.product_id];
  });
  vulnerabilities.sort((vuln1: Vulnerability, vuln2: Vulnerability) => {
    if (vuln1.cve < vuln2.cve) return -1;
    if (vuln1.cve > vuln2.cve) return 1;
    return 0;
  });
  const result = generateCrossTableFrom(products, vulnerabilities, productLookup);
  return result;
};

/**
 *
 * @param products mentioned products
 * @param vulnerabilities all vulnerabilites
 * @param productLookup {id:name}
 * @returns a crosstable as [[]]
 * E.g. [
 * [
 *   {name: 'Product' content:'Product'},
 *   {name: 'Total result' content: 'Total result'},
 *   {name: 'CVE-2016-0173' content: 'CVE-2016-0173'},
 *   {name: 'CVE-2018-0172' content: 'CVE-2018-0172'},
 *   {name: 'CVE-2019-0171', content: 'CVE-2019-0171'},
 *   {name: 'CVE-2020-0174', content: 'CVE-2020-0174'}
 * ],
 * [
 *  {name: 'Product' content:'123'},
 *  {name: 'Total result' content:'K'},
 *  {name: 'CVE-2016-0173' content: ''},
 *  {name: 'CVE-2018-0172' content:''},
 *  {name: 'CVE-2019-0171', content:'K'},
 *  {name: 'CVE-2020-0174', content: ''}
 *  ],
 * ...
 * ]
 */
const generateCrossTableFrom = (
  products: Product[],
  vulnerabilities: Vulnerability[],
  productLookup: any
) => {
  let result: any = [];
  let header: any = [
    { name: "Product", content: "Product" },
    { name: "Total result", content: "Total result" }
  ];
  const getCVE = vulnerabilities.map((vulnerability: Vulnerability) => {
    return { name: vulnerability.cve, content: vulnerability.cve };
  });
  header = header.concat(getCVE);
  result.push(header);
  const productLines = products.map((product: Product) => {
    let line = [{ name: `Product`, content: `${product.product_id}` }];
    let result = generateLineWith(product, vulnerabilities);
    result = ["Total", ...vulnerabilities].map((r: any, i: number) => {
      if (i === 0)
        return {
          name: "Total",
          content: result[i]
        };
      return { name: r.cve, content: result[i] };
    });
    line = line.concat(result);
    return line;
  });
  productLines.sort((line1: any, line2: any) => {
    if (productLookup[line1[0].content] < productLookup[line2[0].content]) return -1;
    if (productLookup[line1[0].content] > productLookup[line2[0]].content) return 1;
    return 0;
  });
  result = [...result, ...productLines];
  return result;
};

/**
 * generateLineWith generates columns for a line for the crosstable with symbols.
 * @param product
 * @param vulnerabilities
 * @returns Array of columns for the line.
 */
const generateLineWith = (product: Product, vulnerabilities: Vulnerability[]) => {
  const DUMMY_TOTAL = "N.A";
  const line: any = [DUMMY_TOTAL];
  vulnerabilities.forEach((vulnerability: Vulnerability) => {
    let column = "";
    if (vulnerability.fixed?.[product.product_id]) {
      column += ProductStatusSymbol.FIXED;
    }
    if (vulnerability.under_investigation?.[product.product_id]) {
      column += ProductStatusSymbol.UNDER_INVESTIGATION;
    }
    if (vulnerability.known_affected?.[product.product_id]) {
      column += ProductStatusSymbol.KNOWN_AFFECTED;
    }
    if (vulnerability.known_not_affected?.[product.product_id]) {
      column += ProductStatusSymbol.NOT_AFFECTED;
    }
    if (vulnerability.recommended?.[product.product_id]) {
      column += ProductStatusSymbol.RECOMMENDED;
    }
    line.push(column);
  });

  /**
   * calculateLineTotal calculates what final symbol should be displayed.
   * @param line
   * @returns final result
   */
  const calculateLineTotal = (line: string[]) => {
    let result = DUMMY_TOTAL;
    switch (true) {
      case line.includes("K"):
        result = "K";
        break;
      case line.includes("U"):
        result = "U";
        break;
      case line.includes("F"):
        result = "F";
        break;
      case line.includes("NR"):
      case line.includes("N"):
        result = "N";
        break;
    }
    return result;
  };
  line[0] = calculateLineTotal(line);
  return line;
};

/**
 * extractProducts retrieves all products from the product tree and adds those defined in relationships.
 * @param jsonDocument
 * @returns An array of products [{product_id:"", name}]
 */
const extractProducts = (jsonDocument: any): Product[] => {
  if (!jsonDocument.product_tree) {
    return [];
  }
  let products: any = [];
  if (jsonDocument.product_tree.branches) {
    const productsFromBranches = jsonDocument.product_tree.branches.reduce(parseBranch, []);
    products = products.concat(productsFromBranches);
  }
  if (jsonDocument.product_tree["full_product_names"]) {
    products = products.concat(jsonDocument.product_tree["full_product_names"]);
  }
  const productsFromRelationships: Product[] = getProductsFromRelationships(jsonDocument);
  return products.concat(productsFromRelationships);
};

/**
 * getProductsFromRelationships retrieves the products from relationships.
 * @param jsonDocument
 * @returns An array of products [{product_id:"", name}]
 */
const getProductsFromRelationships = (jsonDocument: any): Product[] => {
  if (!jsonDocument.product_tree.relationships) return [];
  return jsonDocument.product_tree.relationships.map((relationship: Relationship) => {
    return {
      product_id: relationship.full_product_name.product_id,
      name: relationship.full_product_name.name
    };
  });
};

/**
 * parseBranch parses recursively branches of the product tree for products.
 * @param acc an array of products [{product_id:"", name}]
 * @param branch branch element of product tree
 * @returns acc as an array of products [{product_id:"", name}]
 */
const parseBranch = (acc: Product[], branch: any) => {
  if (branch.branches) {
    branch.branches.forEach((subbranch: any) => {
      acc.concat(parseBranch(acc, subbranch));
    });
  } else {
    if (isProduct(branch)) {
      acc.push({ product_id: branch.product.product_id, name: branch.product.name });
    }
  }
  return acc;
};

/**
 * isProduct determines when a branch is a product branch.
 * @param branch
 * @returns true | false
 */
const isProduct = (branch: any) => {
  return branch.product && branch.product.product_id && branch.product.name;
};

/**
 * generateDictFrom generates a lookup from productstatus and section.
 * @param productStatus
 * @param section
 * @returns dict
 */
const generateDictFrom = (productStatus: ProductStatus_t, section: ProductStatus_t_Key) => {
  return productStatus[section]?.reduce((o: any, n: string) => {
    o[n] = n;
    return o;
  }, {});
};

/**
 * extractVulnerabilities retrieves the vulnerabilites from a CSAF document and collects relevant products (IDs).
 * @param jsonDocument
 * @returns vulnerabilities
 */
const extractVulnerabilities = (jsonDocument: any): VulnerabilitesExtractionResult => {
  const extractionResult: VulnerabilitesExtractionResult = {
    vulnerabilities: [],
    relevantProducts: {}
  };
  if (!jsonDocument.vulnerabilities) {
    return extractionResult;
  }
  /**
   * vulnerabilities as a collection of all vulnerabilities found
   */
  const vulnerabilities = jsonDocument.vulnerabilities.reduce(
    (acc: Vulnerability[], vulnerability: any) => {
      if (!vulnerability.cve) {
        return acc;
      }
      const result: Vulnerability = {
        cve: vulnerability.cve
      };
      if (vulnerability.product_status) {
        if (vulnerability.product_status.known_affected) {
          result.known_affected = generateDictFrom(vulnerability.product_status, "known_affected");
          Object.assign(extractionResult.relevantProducts, result.known_affected);
        }
        if (vulnerability.product_status.fixed) {
          result.fixed = generateDictFrom(vulnerability.product_status, "fixed");
          Object.assign(extractionResult.relevantProducts, result.fixed);
        }
        if (vulnerability.product_status.under_investigation) {
          result.under_investigation = generateDictFrom(
            vulnerability.product_status,
            "under_investigation"
          );
          Object.assign(extractionResult.relevantProducts, result.under_investigation);
        }
        if (vulnerability.product_status.known_not_affected) {
          result.known_not_affected = generateDictFrom(
            vulnerability.product_status,
            "known_not_affected"
          );
          Object.assign(extractionResult.relevantProducts, result.known_not_affected);
        }
        if (vulnerability.product_status.recommended) {
          result.recommended = generateDictFrom(vulnerability.product_status, "recommended");
          Object.assign(extractionResult.relevantProducts, result.recommended);
        }
      }
      acc.push(result);
      return acc;
    },
    []
  );
  extractionResult.vulnerabilities = vulnerabilities;
  return extractionResult;
};

export { extractProducts, extractVulnerabilities, generateProductVulnerabilities };
