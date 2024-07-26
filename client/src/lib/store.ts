// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
//  Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

import { writable } from "svelte/store";
import type { DocModel } from "$lib/Advisories/CSAFWebview/docmodel/docmodeltypes";
import { ADMIN, AUDITOR, EDITOR, IMPORTER, REVIEWER, SOURCE_MANAGER } from "./workflow";
import { MESSAGE } from "./Messages/messagetypes";
import { UserManager, type UserProfile } from "oidc-client-ts";

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
    diff: {
      isDiffBoxOpen: boolean;
      docA_ID: string | undefined;
      docB_ID: string | undefined;
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
      isVulnerabilitiesSectionVisible: boolean;
      isProductTreeOpen: boolean;
      isProductTreeVisible: boolean;
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
      userProfile: {
        firstName: "",
        lastName: ""
      },
      diff: {
        isDiffBoxOpen: false,
        docA_ID: undefined,
        docB_ID: undefined
      },
      sessionExpired: false,
      sessionExpiredMessage: null,
      expiryTime: "",
      isUserLoggedIn: false,
      tokenParsed: null,
      userManager: null,
      errors: []
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
        isVulnerabilitiesSectionVisible: false,
        isProductTreeOpen: false,
        isProductTreeVisible: false,
        isFeedSectionOpen: false,
        selectedCVE: "",
        selectedProduct: "",
        uploadedFile: false,
        history: []
      }
    }
  };
};

function createStore() {
  const { subscribe, set, update } = writable(generateInitialState());
  let state: any;
  subscribe((v) => (state = v));
  return {
    subscribe,
    setFourCVEs: (cves: any) => {
      update((settings) => {
        settings.webview.four_cves = cves;
        return settings;
      });
    },
    clearFourCVEs: () => {
      update((settings) => {
        settings.webview.four_cves = [];
        return settings;
      });
    },
    setSessionExpired: (expired: boolean) => {
      update((settings) => {
        settings.app.sessionExpired = expired;
        return settings;
      });
    },
    setSessionExpiredMessage: (message: string) => {
      update((settings) => {
        settings.app.sessionExpiredMessage = message;
        return settings;
      });
    },
    setExpiryTime: (newExpiryTime: string) => {
      update((settings) => {
        settings.app.expiryTime = newExpiryTime;
        return settings;
      });
    },
    setIsUserLoggedIn: (isUserLoggedIn: boolean) => {
      update((settings) => {
        settings.app.isUserLoggedIn = isUserLoggedIn;
        return settings;
      });
    },
    setTokenParsed: (tokenParsed: ProfileWithRoles) => {
      update((settings) => {
        settings.app.tokenParsed = tokenParsed;
        return settings;
      });
    },
    toggleDocExpandAll: () => {
      update((settings) => {
        settings.webview.ui.docToggleExpandAll = !settings.webview.ui.docToggleExpandAll;
        return settings;
      });
    },
    setFeedSectionOpen: () => {
      update((settings) => {
        settings.webview.ui.isFeedSectionOpen = true;
        return settings;
      });
    },
    setFeedSectionClosed: () => {
      update((settings) => {
        settings.webview.ui.isFeedSectionOpen = false;
        return settings;
      });
    },
    setLoading: (option: boolean) => {
      update((settings) => {
        settings.webview.ui.loading = option;
        return settings;
      });
    },
    setSingleErrorMsg: (msg: string) => {
      update((settings) => {
        settings.webview.ui.singleErrorMsg = msg;
        return settings;
      });
    },
    setFeedErrorMsg: (msg: string) => {
      update((settings) => {
        settings.webview.ui.feedErrorMsg = msg;
        return settings;
      });
    },
    setDocument: (data: any) =>
      update((settings) => {
        settings.webview.doc = data;
        return settings;
      }),
    setSelectedCVE: (cve: string) => {
      update((settings) => {
        settings.webview.ui.selectedCVE = cve;
        return settings;
      });
    },
    resetSelectedCVE: () => {
      update((settings) => {
        settings.webview.ui.selectedCVE = "";
        return settings;
      });
    },
    setSelectedProduct: (product: string) => {
      update((settings) => {
        settings.webview.ui.selectedProduct = product;
        return settings;
      });
    },
    resetSelectedProduct: () => {
      update((settings) => {
        settings.webview.ui.selectedProduct = "";
        return settings;
      });
    },
    setGeneralSectionVisible: () => {
      update((settings) => {
        settings.webview.ui.isGeneralSectionVisible = true;
        return settings;
      });
    },
    setGeneralSectionInvisible: () => {
      update((settings) => {
        settings.webview.ui.isGeneralSectionVisible = false;
        return settings;
      });
    },
    setVulnerabilitiesSectionVisible: () => {
      update((settings) => {
        settings.webview.ui.isVulnerabilitiesSectionVisible = true;
        return settings;
      });
    },
    setVulnerabilitiesSectionInvisible: () => {
      update((settings) => {
        settings.webview.ui.isVulnerabilitiesSectionVisible = false;
        return settings;
      });
    },
    setVulnerabilitiesOverviewVisible: () => {
      update((settings) => {
        settings.webview.ui.isVulnerabilitiesOverviewVisible = true;
        return settings;
      });
    },
    setVulnerabilitiesOverviewInvisible: () => {
      update((settings) => {
        settings.webview.ui.isVulnerabilitiesOverviewVisible = false;
        return settings;
      });
    },
    setProductTreeOpen: () => {
      update((settings) => {
        settings.webview.ui.isProductTreeOpen = true;
        return settings;
      });
    },
    setProductTreeClosed: () => {
      update((settings) => {
        settings.webview.ui.isProductTreeOpen = false;
        return settings;
      });
    },
    setProductTreeSectionVisible: () => {
      update((settings) => {
        settings.webview.ui.isProductTreeVisible = true;
        return settings;
      });
    },
    setProductTreeSectionInVisible: () => {
      update((settings) => {
        settings.webview.ui.isProductTreeVisible = false;
        return settings;
      });
    },
    setUploadedFile: () => {
      update((settings) => {
        settings.webview.ui.uploadedFile = true;
        return settings;
      });
    },
    clearUploadedFile: () => {
      update((settings) => {
        settings.webview.ui.uploadedFile = false;
        return settings;
      });
    },
    clearHistory: () => {
      update((settings) => {
        settings.webview.ui.history = [];
        return settings;
      });
    },
    unshiftHistory: (id: string) => {
      update((settings) => {
        settings.webview.ui.history.unshift(id);
        return settings;
      });
    },
    shiftHistory: () => {
      update((settings) => {
        if (settings.webview.ui.history.length > 0) {
          settings.webview.ui.history.shift();
        }
        return settings;
      });
    },
    setUserManager: (userManager: UserManager) => {
      update((settings) => {
        settings.app.userManager = userManager;
        return settings;
      });
    },
    setUserProfile: (userProfile: any) => {
      update((settings) => {
        const { firstName, lastName } = userProfile;
        settings.app.userProfile.firstName = firstName;
        settings.app.userProfile.lastName = lastName;
        return settings;
      });
    },
    toggleDiffBox: () => {
      update((settings) => {
        settings.app.diff.isDiffBoxOpen = !settings.app.diff.isDiffBoxOpen;
        return settings;
      });
    },
    openDiffBox: () => {
      update((settings) => {
        settings.app.diff.isDiffBoxOpen = true;
        return settings;
      });
    },
    setDiffDocA_ID: (id: string | number | undefined) => {
      update((settings) => {
        settings.app.diff.docA_ID = id?.toString();
        return settings;
      });
    },
    setDiffDocB_ID: (id: string | number | undefined) => {
      update((settings) => {
        settings.app.diff.docB_ID = id?.toString();
        return settings;
      });
    },
    displayErrorMessage: (msg: string) => {
      update((settings) => {
        const errorMessage = generateMessage(msg, MESSAGE.ERROR);
        settings.app.errors = [errorMessage, ...settings.app.errors];
        return settings;
      });
    },
    displayWarningMessage: (msg: string) => {
      update((settings) => {
        const errorMessage = generateMessage(msg, MESSAGE.WARNING);
        settings.app.errors = [errorMessage, ...settings.app.errors];
        return settings;
      });
    },
    displaySuccessMessage: (msg: string) => {
      update((settings) => {
        const errorMessage = generateMessage(msg, MESSAGE.SUCCESS);
        settings.app.errors = [errorMessage, ...settings.app.errors];
        return settings;
      });
    },
    displayInfoMessage: (msg: string) => {
      update((settings) => {
        const errorMessage = generateMessage(msg, MESSAGE.INFO);
        settings.app.errors = [errorMessage, ...settings.app.errors];
        return settings;
      });
    },
    removeError: (id: string) => {
      update((settings) => {
        settings.app.errors = settings.app.errors.filter((msg) => {
          return msg.id !== id;
        });
        return settings;
      });
    },
    reset: () => {
      set(generateInitialState());
    },
    getRoles: () => state.app.tokenParsed.realm_access.roles,
    isImporter: () => appStore.getRoles().includes(IMPORTER),
    isEditor: () => appStore.getRoles().includes(EDITOR),
    isReviewer: () => appStore.getRoles().includes(REVIEWER),
    isAdmin: () => appStore.getRoles().includes(ADMIN),
    isAuditor: () => appStore.getRoles().includes(AUDITOR),
    isSourceManager: () => appStore.getRoles().includes(SOURCE_MANAGER),
    getUserManager: () => state.app.userManager,
    getIsUserLoggedIn: () => state.app.isUserLoggedIn
  };
}

export const appStore = createStore();
