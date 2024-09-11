// This file is Free Software under the Apache-2.0 License
// without warranty, see README.md and LICENSES/Apache-2.0.txt for details.
//
// SPDX-License-Identifier: Apache-2.0
//
// SPDX-FileCopyrightText: 2024 German Federal Office for Information Security (BSI) <https://www.bsi.bund.de>
//  Software-Engineering: 2024 Intevation GmbH <https://intevation.de>

import { getErrorDetails, type ErrorDetails } from "$lib/Errors/error";
import { request } from "$lib/request";
import type { Result } from "$lib/types";
import type { CSAFProviderMetadata } from "$lib/provider";

type Source = {
  id?: number;
  name: string;
  url: string;
  active?: boolean;
  rate?: number | null;
  slots?: number | null;
  headers: string[];
  strict_mode?: boolean;
  insecure?: boolean;
  signature_check?: boolean;
  age?: string;
  ignore_patterns: string[];
  client_cert_public?: string | null;
  client_cert_private?: string | null;
  client_cert_passphrase?: string | null;
  stats?: {
    downloading: number;
    waiting: number;
  };
};

enum LogLevel {
  debug = "debug",
  info = "info",
  warn = "warn",
  error = "error"
}

type Feed = {
  id?: number;
  enable?: boolean;
  url: string;
  label: string;
  rolie: boolean;
  log_level: LogLevel;
  edit?: boolean;
  stats?: {
    downloading: number;
    waiting: number;
  };
};

const logLevels = [
  { value: LogLevel.error, name: "Error" },
  { value: LogLevel.info, name: "Info" },
  { value: LogLevel.warn, name: "Warn" },
  { value: LogLevel.debug, name: "Debug" }
];

const saveSource = async (source: Source): Promise<Result<Source, ErrorDetails>> => {
  let method = "POST";
  let path = `/api/sources`;
  const formData = new FormData();
  if (source.id) {
    method = "PUT";
    path += `/${source.id}`;
    formData.append("id", source.id.toString());
  } else {
    formData.append("url", source.url);
  }
  formData.append("name", source.name);
  if (source.active !== undefined) {
    formData.append("active", source.active.toString());
  }
  if (source.rate && source.rate !== 0) {
    formData.append("rate", source.rate.toString());
  }
  if (source.slots && source.slots !== 0) {
    formData.append("slots", source.slots.toString());
  }
  if (source.strict_mode !== undefined) {
    formData.append("strict_mode", source.strict_mode.toString());
  }
  if (source.insecure !== undefined) {
    formData.append("insecure", source.insecure.toString());
  }
  if (source.signature_check !== undefined) {
    formData.append("signature_check", source.signature_check.toString());
  }
  if (source.age != undefined) {
    formData.append("age", source.age.toString());
  }
  if (source.client_cert_public) {
    formData.append("client_cert_public", source.client_cert_public);
  }
  if (source.client_cert_private) {
    formData.append("client_cert_private", source.client_cert_private);
  }
  if (source.client_cert_passphrase && source.client_cert_passphrase !== "***") {
    formData.append("client_cert_passphrase", source.client_cert_passphrase);
  }
  for (const header of source.headers) {
    if (header != "") {
      formData.append("headers", header);
    }
  }
  for (const pattern of source.ignore_patterns) {
    if (pattern != "") {
      formData.append("ignore_patterns", pattern);
    }
  }
  const resp = await request(path, method, formData);
  if (resp.ok) {
    if (resp.content.id) {
      source.id = resp.content.id;
    }
    return {
      ok: true,
      value: source
    };
  } else {
    return {
      ok: false,
      error: getErrorDetails(`Could not save source`, resp)
    };
  }
};

const parseFeeds = (pmd: CSAFProviderMetadata, currentFeeds: Feed[]): Feed[] => {
  const feeds: Feed[] = [];

  const dist = pmd.distributions ?? [];

  for (const entry of dist) {
    if (entry.rolie) {
      for (const feed of entry.rolie.feeds) {
        let label = "";
        if (feed.summary) {
          label = feed.summary;
        } else {
          label = feed.tlp_label + " " + feed.url.split("/").pop();
        }
        feeds.push({
          url: feed.url,
          label: label,
          log_level: LogLevel.error,
          rolie: true,
          enable: true
        });
      }
    }
    if (entry.directory_url) {
      const splitUrl = entry.directory_url.split("/");
      let label = splitUrl.pop() ?? "Default label";
      // If the url ends with '/'
      if (label.length === 0) {
        label = splitUrl.pop() ?? "Default label";
      }
      feeds.push({
        url: entry.directory_url,
        label: label,
        log_level: LogLevel.error,
        rolie: false,
        enable: true
      });
    }
    const existingLabels = new Set(currentFeeds.map((f) => f.label));
    for (const feed of feeds) {
      while (existingLabels.has(feed.label)) {
        feed.label += "#";
      }
      existingLabels.add(feed.label);
    }
  }

  return feeds;
};

const saveFeeds = async (
  source: Source,
  feeds: Feed[]
): Promise<Result<number[], ErrorDetails>> => {
  const ids: number[] = [];
  for (const feed of feeds) {
    if (!feed.enable) {
      if (feed.id) {
        deleteFeed(feed.id);
      }
      continue;
    }
    const formData = new FormData();

    let path = `/api/sources`;
    let method = "POST";

    if (feed.id) {
      method = "PUT";
      path += `/feeds/${feed.id}`;
    } else {
      path += `/${source.id}/feeds`;
      formData.append("url", feed.url);
    }
    formData.append("label", feed.label);
    formData.append("log_level", feed.log_level);

    const resp = await request(path, method, formData);
    if (resp.error) {
      return {
        ok: false,
        error: getErrorDetails(`Could not save feed.`, resp)
      };
    } else {
      ids.push(resp.content.id);
    }
  }
  return { ok: true, value: ids };
};

const getSourceName = async (pmd: CSAFProviderMetadata): Promise<string> => {
  const result = await fetchSources();
  if (!result.ok) {
    return `New Source #${Math.floor(Math.random() * 1000)}`;
  }
  const sources = result.value;
  const name = pmd.publisher.name;
  if (!sources.find((s) => s.name === name)) {
    return name;
  }
  for (let i = 1; ; i++) {
    const customName = `${name} #${i}`;
    if (!sources.find((s) => s.name === customName)) {
      return customName;
    }
  }
};

const fetchPMD = async (url: string): Promise<Result<CSAFProviderMetadata, ErrorDetails>> => {
  const resp = await request(`/api/pmd?url=${encodeURIComponent(url)}`, "GET");
  if (resp.ok) {
    return {
      ok: true,
      value: resp.content
    };
  } else {
    return {
      ok: false,
      error: getErrorDetails(`Could not load PMD.`, resp)
    };
  }
};

const fetchSource = async (
  id: number,
  showStats: boolean = false
): Promise<Result<Source, ErrorDetails>> => {
  const resp = await request(`/api/sources/${id}?stats=${showStats}`, "GET");
  if (resp.ok) {
    if (resp.content) {
      const source = resp.content;
      if (!source.ignore_patterns) {
        source.ignore_patterns = [""];
      }
      if (source.rate === undefined) {
        source.rate = null;
      }
      if (source.slots === undefined) {
        source.slots = null;
      }
      return {
        ok: true,
        value: source
      };
    }
  }
  return {
    ok: false,
    error: getErrorDetails(`Could not load source`, resp)
  };
};

const fetchSources = async (
  showStats: boolean = false
): Promise<Result<Source[], ErrorDetails>> => {
  const resp = await request(`/api/sources?stats=${showStats}`, "GET");
  if (resp.ok) {
    if (resp.content.sources) {
      return {
        ok: true,
        value: resp.content.sources
      };
    }
  }
  return {
    ok: false,
    error: getErrorDetails(`Could not load source`, resp)
  };
};

const fetchFeed = async (
  id: number,
  showStats: boolean = false
): Promise<Result<Feed, ErrorDetails>> => {
  const resp = await request(`/api/sources/feeds/${id}?stats=${showStats}`, "GET");
  if (resp.ok) {
    return {
      ok: true,
      value: resp.content
    };
  }
  return {
    ok: false,
    error: getErrorDetails(`Could not load feed`, resp)
  };
};

const fetchFeeds = async (
  id: number,
  showStats: boolean = false
): Promise<Result<Feed[], ErrorDetails>> => {
  const resp = await request(`/api/sources/${id}/feeds?stats=${showStats}`, "GET");
  if (resp.ok) {
    if (resp.content.feeds) {
      return {
        ok: true,
        value: resp.content.feeds
      };
    }
  }
  return {
    ok: false,
    error: getErrorDetails(`Could not load feed`, resp)
  };
};

const deleteFeed = async (id: number): Promise<Result<null, ErrorDetails>> => {
  const resp = await request(`/api/sources/feeds/${id}`, "DELETE");
  if (resp.error) {
    return {
      ok: false,
      error: getErrorDetails(`Could not load feed`, resp)
    };
  } else {
    return {
      ok: true,
      value: null
    };
  }
};

const fetchFeedLogs = async (
  id: number,
  offset: number,
  limit: number,
  count: boolean = false
): Promise<Result<[any[], number], ErrorDetails>> => {
  const resp = await request(
    `/api/sources/feeds/${id}/log?limit=${limit}&offset=${offset}&count=${count}`,
    "GET"
  );
  if (resp.ok) {
    return { ok: true, value: [resp.content.entries, resp.content.count ?? 0] };
  }
  return { ok: false, error: getErrorDetails(`Could not load feed logs`, resp) };
};

const isSameFeed = (a: Feed, b: Feed) => a.url === b.url && a.rolie === b.rolie;

const calculateMissingFeeds = (pmdFeeds: Feed[], feeds: Feed[]): Feed[] => {
  return pmdFeeds.filter((a) => !feeds.some((b) => isSameFeed(a, b)));
};

const parseHeaders = (source: Source): string[][] => {
  const headers = [];
  for (const header of source.headers) {
    const h = header.split(":");
    headers.push([h[0], h[1]]);
  }
  if (headers.length === 0) {
    headers.push(["", ""]);
  }
  return headers;
};

export {
  type Source,
  LogLevel,
  type Feed,
  saveSource,
  fetchSource,
  fetchPMD,
  getSourceName,
  calculateMissingFeeds,
  fetchFeedLogs,
  parseHeaders,
  parseFeeds,
  deleteFeed,
  fetchFeed,
  fetchFeeds,
  fetchSources,
  saveFeeds,
  logLevels
};
