/* eslint-disable */
// @ts-nocheck

/**
 * This file is a generated Typescript file for GRPC Gateway, DO NOT MODIFY
 */

import * as fm from "./fetch.pb";


export type VersionRequest = Record<string, never>;

export type VersionResponse = {
  version?: string;
};

export type Version = Record<string, never>;

export type ListEndpointsRequest = {
  serviceName?: string;
};

export type ListEndpointsResponse = {
  urls?: string[];
};

export type ListEndpoints = Record<string, never>;

export type Endpoint = {
  serviceName?: string;
  addrs?: string[];
};

export type UpsertEndpointsRequest = {
  endpoints?: Endpoint[];
};

export type UpsertEndpointsResponse = Record<string, never>;

export type UpsertEndpoints = Record<string, never>;

export class MakoshBeAPI {
  static Version(this:void, req: VersionRequest, initReq?: fm.InitReq): Promise<VersionResponse> {
    return fm.fetchRequest<VersionResponse>(`/api/version?${fm.renderURLSearchParams(req, [])}`, {...initReq, method: "GET"});
  }
  static ListEndpoints(this:void, req: ListEndpointsRequest, initReq?: fm.InitReq): Promise<ListEndpointsResponse> {
    return fm.fetchRequest<ListEndpointsResponse>(`/api/endpoints/${req.serviceName}?${fm.renderURLSearchParams(req, ["serviceName"])}`, {...initReq, method: "GET"});
  }
  static UpsertEndpoints(this:void, req: UpsertEndpointsRequest, initReq?: fm.InitReq): Promise<UpsertEndpointsResponse> {
    return fm.fetchRequest<UpsertEndpointsResponse>(`/api/endpoints`, {...initReq, method: "POST", body: JSON.stringify(req, fm.replacer)});
  }
}