export type JsonDiffResultList = JsonDiffResult[];

export type JsonDiffResult = {
  op: string;
  path: string;
  value?: string | object | object[];
};
