// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
//  Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

import type { DocModel } from "$lib/Advisories/CSAFWebview/docmodel/docmodeltypes";
import { ADMIN, AUDITOR, EDITOR, IMPORTER, REVIEWER, SOURCE_MANAGER } from "./workflow";
import { MESSAGE } from "./Messages/messagetypes";
import { UserManager, type UserProfile } from "oidc-client-ts";
import { SvelteSet } from "svelte/reactivity";
import type { SearchParameters } from "./Search/search";
import type { SEARCHTYPES } from "./Queries/query";

type ErrorMessage = {
  id: string;
  type: string;
  message: string;
};

export type ProfileWithRoles = UserProfile & {
  realm_access: {
    roles: string[];
  };
};

type AppStore = {
  app: {
    config: any;
    userProfile: {
      firstName: string;
      lastName: string;
    };
    expiryTime: string;
    isUserLoggedIn: boolean;
    sessionExpired: boolean;
    sessionExpiredMessage: string | null;
    tokenParsed: ProfileWithRoles | null;
    userManager: UserManager | null;
    errors: ErrorMessage[];
    documents: any[] | null;
    documentsToDelete: any[] | null;
    selectedDocumentIDs: SvelteSet<number>;
    isToolboxOpen: boolean;
    isDeleteModalOpen: boolean;
    diff: {
      docA_ID: string | undefined;
      docB_ID: string | undefined;
      docA: any | undefined;
      docB: any | undefined;
    };
    isDarkMode: boolean;
    search: {
      parameters: {
        advisories: SearchParameters | undefined;
        documents: SearchParameters | undefined;
        events: SearchParameters | undefined;
      };
      searchURL: string | undefined;
    };
  };
  webview: {
    doc: DocModel | null;
    providerMetadata: any;
    currentFeed: any;
    four_cves: any;
    ui: {
      docToggleExpandAll: boolean;
      feedErrorMsg: string;
      loading: boolean;
      singleErrorMsg: string;
      isGeneralSectionVisible: boolean;
      isRevisionHistoryVisible: boolean;
      isVulnerabilitiesOverviewVisible: boolean;
      isFeedSectionOpen: boolean;
      selectedCVE: string;
      selectedProduct: string;
      uploadedFile: boolean;
      history: string[];
    };
  };
};

const generateMessage = (msg: string, type: string) => {
  return {
    id: crypto.randomUUID(),
    type: type,
    message: msg
  };
};

const generateInitialState = (): AppStore => {
  return {
    app: {
      config: undefined,
      userProfile: {
        firstName: "",
        lastName: ""
      },
      diff: {
        docA_ID: undefined,
        docB_ID: undefined,
        docA: undefined,
        docB: undefined
      },
      sessionExpired: false,
      sessionExpiredMessage: null,
      expiryTime: "",
      isUserLoggedIn: false,
      tokenParsed: null,
      userManager: null,
      errors: [],
      documents: null,
      documentsToDelete: null,
      isDeleteModalOpen: false,
      selectedDocumentIDs: new SvelteSet<number>(),
      isToolboxOpen: false,
      isDarkMode: document.firstElementChild?.classList.contains("dark") ?? false,
      search: {
        parameters: {
          advisories: undefined,
          documents: undefined,
          events: undefined
        },
        searchURL: undefined
      }
    },
    webview: {
      doc: null,
      providerMetadata: null,
      currentFeed: null,
      four_cves: [],
      ui: {
        docToggleExpandAll: false,
        feedErrorMsg: "",
        loading: false,
        singleErrorMsg: "",
        isGeneralSectionVisible: true,
        isRevisionHistoryVisible: false,
        isVulnerabilitiesOverviewVisible: true,
        isFeedSectionOpen: false,
        selectedCVE: "",
        selectedProduct: "",
        uploadedFile: false,
        history: []
      }
    }
  };
};

// Create the store using Svelte 5 runes
const state = $state(generateInitialState());

// Derived values
const roles = $derived(state.app.tokenParsed?.realm_access?.roles ?? []);

// Store methods
export const appStore = {
  // Getters
  get state() {
    return state;
  },
  get roles() {
    return roles;
  },

  // Role checks
  isImporter: () => roles.includes(IMPORTER),
  isEditor: () => roles.includes(EDITOR),
  isReviewer: () => roles.includes(REVIEWER),
  isAdmin: () => roles.includes(ADMIN),
  isAuditor: () => roles.includes(AUDITOR),
  isSourceManager: () => roles.includes(SOURCE_MANAGER),

  // User management
  getUserManager: () => state.app.userManager,
  getIsUserLoggedIn: () => state.app.isUserLoggedIn,

  // Configuration
  getOption: (option: string) => state.app.config?.[option],
  getIdleTimeout: () => state.app.config?.idle_timeout,
  getKeycloakClientID: () => state.app.config?.keycloak_client_id,
  getKeycloakRealm: () => state.app.config?.keycloak_realm,
  getKeycloakURL: () => state.app.config?.keycloak_url,
  getUpdateInterval: () => state.app.config?.update_interval,

  // Setters
  setFourCVEs: (cves: any) => {
    state.webview.four_cves = cves;
  },

  clearFourCVEs: () => {
    state.webview.four_cves = [];
  },

  setSessionExpired: (expired: boolean) => {
    state.app.sessionExpired = expired;
  },

  setSessionExpiredMessage: (message: string) => {
    state.app.sessionExpiredMessage = message;
  },

  setExpiryTime: (newExpiryTime: string) => {
    state.app.expiryTime = newExpiryTime;
  },

  setIsUserLoggedIn: (isUserLoggedIn: boolean) => {
    state.app.isUserLoggedIn = isUserLoggedIn;
  },

  setTokenParsed: (tokenParsed: ProfileWithRoles) => {
    state.app.tokenParsed = tokenParsed;
  },

  toggleDocExpandAll: () => {
    state.webview.ui.docToggleExpandAll = !state.webview.ui.docToggleExpandAll;
  },

  setFeedSectionOpen: () => {
    state.webview.ui.isFeedSectionOpen = true;
  },

  setFeedSectionClosed: () => {
    state.webview.ui.isFeedSectionOpen = false;
  },

  setLoading: (option: boolean) => {
    state.webview.ui.loading = option;
  },

  setSingleErrorMsg: (msg: string) => {
    state.webview.ui.singleErrorMsg = msg;
  },

  setFeedErrorMsg: (msg: string) => {
    state.webview.ui.feedErrorMsg = msg;
  },

  setDocument: (data: any) => {
    state.webview.doc = data;
  },

  setSelectedCVE: (cve: string) => {
    state.webview.ui.selectedCVE = cve;
  },

  resetSelectedCVE: () => {
    state.webview.ui.selectedCVE = "";
  },

  setSelectedProduct: (product: string) => {
    state.webview.ui.selectedProduct = product;
  },

  resetSelectedProduct: () => {
    state.webview.ui.selectedProduct = "";
  },

  setGeneralSectionVisible: () => {
    state.webview.ui.isGeneralSectionVisible = true;
  },

  setGeneralSectionInvisible: () => {
    state.webview.ui.isGeneralSectionVisible = false;
  },

  setVulnerabilitiesOverviewVisible: () => {
    state.webview.ui.isVulnerabilitiesOverviewVisible = true;
  },

  setVulnerabilitiesOverviewInvisible: () => {
    state.webview.ui.isVulnerabilitiesOverviewVisible = false;
  },

  setUploadedFile: () => {
    state.webview.ui.uploadedFile = true;
  },

  updateDarkMode: () => {
    state.app.isDarkMode = document.documentElement.classList.contains("dark") ?? false;
  },

  clearUploadedFile: () => {
    state.webview.ui.uploadedFile = false;
  },

  clearHistory: () => {
    state.webview.ui.history = [];
  },

  unshiftHistory: (id: string) => {
    state.webview.ui.history.unshift(id);
  },

  shiftHistory: () => {
    if (state.webview.ui.history.length > 0) {
      state.webview.ui.history.shift();
    }
  },

  setUserManager: (userManager: UserManager) => {
    state.app.userManager = userManager;
  },

  setUserProfile: (userProfile: any) => {
    const { firstName, lastName } = userProfile;
    state.app.userProfile.firstName = firstName;
    state.app.userProfile.lastName = lastName;
  },

  toggleToolbox: () => {
    state.app.isToolboxOpen = !state.app.isToolboxOpen;
  },

  openToolbox: () => {
    state.app.isToolboxOpen = true;
  },

  setDiffDocA_ID: (id: string | number | undefined) => {
    state.app.diff.docA_ID = id?.toString();
  },

  setDiffDocB_ID: (id: string | number | undefined) => {
    state.app.diff.docB_ID = id?.toString();
  },

  setDiffDocA: (doc: any) => {
    state.app.diff.docA = doc;
  },

  setDiffDocB: (doc: any) => {
    state.app.diff.docB = doc;
  },

  setDocuments: (newDocuments: any[]) => {
    state.app.documents = newDocuments;
  },

  setDocumentsToDelete: (documents: any[]) => {
    state.app.documentsToDelete = documents;
  },

  setIsDeleteModalOpen: (isOpen: boolean) => {
    state.app.isDeleteModalOpen = isOpen;
  },

  addSelectedDocumentID: (id: number) => {
    state.app.selectedDocumentIDs.add(id);
  },

  removeSelectedDocumentID: (id: number) => {
    state.app.selectedDocumentIDs.delete(id);
  },

  clearSelectedDocumentIDs: () => {
    state.app.selectedDocumentIDs.clear();
  },

  displayErrorMessage: (msg: string) => {
    const errorMessage = generateMessage(msg, MESSAGE.ERROR);
    state.app.errors = [errorMessage, ...state.app.errors];
  },

  displayWarningMessage: (msg: string) => {
    const errorMessage = generateMessage(msg, MESSAGE.WARNING);
    state.app.errors = [errorMessage, ...state.app.errors];
  },

  displaySuccessMessage: (msg: string) => {
    const errorMessage = generateMessage(msg, MESSAGE.SUCCESS);
    state.app.errors = [errorMessage, ...state.app.errors];
  },

  displayInfoMessage: (msg: string) => {
    const errorMessage = generateMessage(msg, MESSAGE.INFO);
    state.app.errors = [errorMessage, ...state.app.errors];
  },

  removeError: (id: string) => {
    state.app.errors = state.app.errors.filter((msg: ErrorMessage) => {
      return msg.id !== id;
    });
  },

  setConfig: (newConfig: any) => {
    state.app.config = newConfig;
  },

  setSearchURL: (newSearchURL: string | undefined) => {
    state.app.search.searchURL = newSearchURL;
  },

  setSearchParametersForType: (type: SEARCHTYPES, parameters: SearchParameters) => {
    state.app.search.parameters[type] = parameters;
  },

  getSearchParametersForType: (type: SEARCHTYPES) => {
    return state.app.search.parameters[type];
  },

  reset: () => {
    Object.assign(state, generateInitialState());
  },

  // Add missing methods that were in the old store
  getRoles: () => {
    return state.app.tokenParsed?.realm_access?.roles ?? [];
  }
};
