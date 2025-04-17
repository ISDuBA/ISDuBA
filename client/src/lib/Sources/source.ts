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
import type { AggregatorMetadata } from "$lib/aggregatorTypes";

const dtClass: string = "ml-1 mt-1 text-gray-500 md:text-sm dark:text-gray-400";
const ddClass: string = "break-words font-semibold ml-2 mb-1";

type Source = {
  id?: number;
  name: string;
  url: string;
  active?: boolean;
  rate?: number;
  slots?: number;
  headers: string[];
  strict_mode?: boolean;
  secure?: boolean;
  signature_check?: boolean;
  age?: string;
  ignore_patterns: string[];
  client_cert_public?: string | null;
  client_cert_private?: string | null;
  client_cert_passphrase?: string | null;
  status?: string[];
  attention: boolean;
  stats?: {
    downloading: number;
    waiting: number;
  };
};

type SourceConfig = {
  slots: number;
  rate: number;
  log_level: LogLevel;
  strict_mode: boolean;
  secure: boolean;
  signature_check: boolean;
  age: string;
};

enum LogLevel {
  default = "default",
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
  { value: LogLevel.default, name: "Default" },
  { value: LogLevel.error, name: "Errors" },
  { value: LogLevel.warn, name: "Warning" },
  { value: LogLevel.info, name: "Info" },
  { value: LogLevel.debug, name: "Debug (verbose)" }
];

type Aggregator = {
  id?: number;
  name: string;
  url: string;
  attention?: boolean;
  active?: boolean;
  category?: string;
  contact_details?: string;
  issuing_authority?: string;
  namespace?: string;
  [key: string]: string | boolean | number | undefined;
};

type Attention = {
  id: number;
  name: string;
};

const fetchAggregatorAttentionList = async (): Promise<Result<Attention[], ErrorDetails>> => {
  const resp = await request(`/api/aggregators/attention`, "GET");
  if (resp.ok) {
    return {
      ok: true,
      value: resp.content
    };
  }
  return {
    ok: false,
    error: getErrorDetails(`Could not load attention list`, resp)
  };
};

const resetAggregatorAttention = async (
  aggregator: Aggregator
): Promise<Result<Attention[], ErrorDetails>> => {
  const formData = new FormData();
  formData.append("attention", "false");
  const resp = await request(`/api/aggregators/${aggregator.id}`, "PUT", formData);
  if (resp.ok) {
    return {
      ok: true,
      value: resp.content
    };
  }
  return {
    ok: false,
    error: getErrorDetails(`Could not update source attention`, resp)
  };
};

const fetchSourceAttentionList = async (): Promise<Result<Attention[], ErrorDetails>> => {
  const resp = await request(`/api/sources/attention`, "GET");
  if (resp.ok) {
    return {
      ok: true,
      value: resp.content
    };
  }
  return {
    ok: false,
    error: getErrorDetails(`Could not load attention list`, resp)
  };
};

const resetSourceAttention = async (source: Source): Promise<Result<Attention[], ErrorDetails>> => {
  const formData = new FormData();
  formData.append("attention", "false");
  const resp = await request(`/api/sources/${source.id}`, "PUT", formData);
  if (resp.ok) {
    return {
      ok: true,
      value: resp.content
    };
  }
  return {
    ok: false,
    error: getErrorDetails(`Could not update source attention`, resp)
  };
};

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
  if ((source.rate && source.rate < 0) || source.rate === undefined) {
    source.rate = 0;
  }
  formData.append("rate", source.rate.toString());
  if ((source.slots && source.slots < 0) || source.slots === undefined) {
    source.slots = 0;
  }
  formData.append("slots", source.slots.toString());
  if (source.strict_mode !== undefined) {
    formData.append("strict_mode", source.strict_mode.toString());
  }
  if (source.secure !== undefined) {
    formData.append("secure", source.secure.toString());
  }
  if (source.signature_check !== undefined) {
    formData.append("signature_check", source.signature_check.toString());
  }
  if (source.age != undefined) {
    formData.append("age", source.age.toString());
  }
  if (source.client_cert_public) {
    if (source.client_cert_public !== "***") {
      formData.append("client_cert_public", source.client_cert_public);
    }
  } else {
    formData.append("client_cert_public", "");
  }
  if (source.client_cert_private) {
    if (source.client_cert_private !== "***") {
      formData.append("client_cert_private", source.client_cert_private);
    }
  } else {
    formData.append("client_cert_private", "");
  }
  if (source.client_cert_passphrase && source.client_cert_passphrase !== "***") {
    formData.append("client_cert_passphrase", source.client_cert_passphrase);
  }
  if (source.headers) {
    for (const header of source.headers) {
      if (header != "") {
        formData.append("headers", header);
      }
    }
  }
  const patterns = source.ignore_patterns.filter((i) => i !== "");
  if (patterns.length > 0) {
    for (const pattern of patterns) {
      formData.append("ignore_patterns", pattern);
    }
  } else {
    formData.append("ignore_patterns", "");
  }
  formData.append("attention", "false");
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
          log_level: LogLevel.default,
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
        log_level: LogLevel.default,
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
    if (feed.log_level === LogLevel.default) {
      formData.append("log_level", "");
    } else {
      formData.append("log_level", feed.log_level);
    }

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
    let details = undefined;
    if (resp.content) {
      details = resp.content.join("\n");
    }
    const errorDetails: ErrorDetails = { message: `Could not load pmd`, details: details };
    return {
      ok: false,
      error: errorDetails
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

const fetchAggregators = async (): Promise<Result<Aggregator[], ErrorDetails>> => {
  const resp = await request(`/api/aggregators`, "GET");
  if (resp.ok) {
    return {
      ok: true,
      value: resp.content
    };
  }
  return {
    ok: false,
    error: getErrorDetails(`Could not load aggregators`, resp)
  };
};

const saveAggregator = async (aggregator: Aggregator): Promise<Result<number, ErrorDetails>> => {
  const formData = new FormData();
  formData.append("name", aggregator.name);
  formData.append("url", aggregator.url);
  const resp = await request(`/api/aggregators`, "POST", formData);
  if (resp.ok) {
    return {
      ok: true,
      value: resp.content.id
    };
  }
  return {
    ok: false,
    error: getErrorDetails(`Could not save aggregator`, resp)
  };
};

const updateAggregator = async (aggregator: Aggregator): Promise<Result<number, ErrorDetails>> => {
  const formData = new FormData();
  const keys = ["name", "url", "attention", "active"];
  keys.forEach((key) => {
    if (aggregator[key] !== undefined) {
      formData.append(key, `${aggregator[key]}`);
    }
  });
  const resp = await request(`/api/aggregators/${aggregator.id}`, "PUT", formData);
  if (resp.ok) {
    return {
      ok: true,
      value: resp.content.id
    };
  }
  return {
    ok: false,
    error: getErrorDetails(`Could not save aggregator`, resp)
  };
};

const deleteAggregator = async (id: number): Promise<Result<null, ErrorDetails>> => {
  const resp = await request(`/api/aggregators/${id}`, "DELETE");
  if (resp.error) {
    return {
      ok: false,
      error: getErrorDetails(`Could not delete aggregator`, resp)
    };
  } else {
    return {
      ok: true,
      value: null
    };
  }
};

const fetchAggregatorData = async (
  url: string
): Promise<Result<AggregatorMetadata, ErrorDetails>> => {
  const resp = await request(`/api/aggregator?url=${encodeURIComponent(url)}`, "GET");
  if (resp.ok) {
    return {
      ok: true,
      value: resp.content
    };
  }
  return {
    ok: false,
    error: getErrorDetails(`Could not fetch aggregator data`, resp)
  };
};

let defaultConfigCache: Promise<Result<SourceConfig, ErrorDetails>>;

const fetchSourceDefaultConfig = async (): Promise<Result<SourceConfig, ErrorDetails>> => {
  if (defaultConfigCache) {
    return defaultConfigCache;
  }
  const fetchRequest = async (): Promise<Result<SourceConfig, ErrorDetails>> => {
    const resp = await request(`/api/sources/default`, "GET");
    if (resp.ok) {
      return {
        ok: true,
        value: resp.content
      };
    }
    return {
      ok: false,
      error: getErrorDetails(`Could not load source default config`, resp)
    };
  };
  defaultConfigCache = fetchRequest();
  return defaultConfigCache;
};

const capitalize = (s: string) => {
  return s && s[0].toUpperCase() + s.slice(1);
};

const getLogLevels = async (
  enableDefault: boolean = false
): Promise<Result<{ value: LogLevel; name: string }[], ErrorDetails>> => {
  if (!enableDefault) {
    const levels = [...logLevels].filter((item) => item.value !== LogLevel.default);
    return {
      ok: true,
      value: levels
    };
  }
  const resp = await fetchSourceDefaultConfig();
  if (resp.ok) {
    const defaultLogLevel = {
      name: `Default (${capitalize(resp.value.log_level)})`,
      value: LogLevel.default
    };
    const levels = [...logLevels];
    const target = levels.find((i) => i.value === LogLevel.default) ?? defaultLogLevel;
    Object.assign(target, defaultLogLevel);

    levels.map((i) =>
      i.name === "Default" ? { name: `Default (${defaultLogLevel})`, value: i.value } : i
    );
    return {
      ok: true,
      value: levels
    };
  } else {
    return resp;
  }
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
      error: getErrorDetails(`Could not delete feed`, resp)
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
  from: Date | undefined = undefined,
  to: Date | undefined = undefined,
  search: string = "",
  logLevels: LogLevel[],
  count: boolean = false,
  abortController: AbortController | undefined = undefined
): Promise<Result<[any[], number], ErrorDetails>> => {
  const levels = `&levels=${logLevels.join(" ")}`;
  const fromParameter = from ? `&from=${from.toISOString()}` : "";
  const toParameter = to ? `&to=${to.toISOString()}` : "";
  const resp = await request(
    `/api/sources/feeds/${id}/log?limit=${limit}&offset=${offset}&count=${count}&search=${search}${fromParameter}${toParameter}${levels}`,
    "GET",
    undefined,
    abortController
  );
  if (resp.ok) {
    return { ok: true, value: [resp.content.entries, resp.content.count ?? 0] };
  }
  return { ok: false, error: getErrorDetails(`Could not load feed logs`, resp) };
};

const fetchAllFeedLogs = async (
  id: number,
  count: boolean = false,
  abortController: AbortController | undefined = undefined
): Promise<Result<[any[], number], ErrorDetails>> => {
  const resp = await request(
    `/api/sources/feeds/${id}/log?offset=0&count=${count}`,
    "GET",
    undefined,
    abortController
  );
  if (resp.ok) {
    return { ok: true, value: [resp.content.entries, resp.content.count ?? 0] };
  }
  return {
    ok: false,
    error:
      resp.error === "AbortError"
        ? { message: "AbortError" }
        : getErrorDetails(`Could not load feed logs`, resp)
  };
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
  type Aggregator,
  type Attention,
  LogLevel,
  type Feed,
  saveSource,
  saveAggregator,
  updateAggregator,
  fetchSource,
  fetchSourceDefaultConfig,
  getLogLevels,
  fetchPMD,
  getSourceName,
  calculateMissingFeeds,
  fetchFeedLogs,
  fetchAllFeedLogs,
  parseHeaders,
  parseFeeds,
  deleteFeed,
  deleteAggregator,
  fetchFeed,
  fetchFeeds,
  fetchSources,
  fetchAggregators,
  fetchAggregatorData,
  saveFeeds,
  fetchSourceAttentionList,
  resetSourceAttention,
  fetchAggregatorAttentionList,
  resetAggregatorAttention,
  dtClass,
  ddClass,
  logLevels
};
