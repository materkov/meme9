import {
  TwirpContext,
  TwirpServer,
  RouterEvents,
  TwirpError,
  TwirpErrorCode,
  Interceptor,
  TwirpContentType,
  chainInterceptors,
} from "twirp-ts";
import { SubscribeRequest, SubscribeResponse } from "./subscriptions";

//==================================//
//          Client Code             //
//==================================//

interface Rpc {
  request(
    service: string,
    method: string,
    contentType: "application/json" | "application/protobuf",
    data: object | Uint8Array
  ): Promise<object | Uint8Array>;
}

export interface SubscriptionsClient {
  Subscribe(request: SubscribeRequest): Promise<SubscribeResponse>;
  Unsubscribe(request: SubscribeRequest): Promise<SubscribeResponse>;
  GetStatus(request: SubscribeRequest): Promise<SubscribeResponse>;
}

export class SubscriptionsClientJSON implements SubscriptionsClient {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Subscribe.bind(this);
    this.Unsubscribe.bind(this);
    this.GetStatus.bind(this);
  }
  Subscribe(request: SubscribeRequest): Promise<SubscribeResponse> {
    const data = SubscribeRequest.toJson(request, {
      useProtoFieldName: true,
      emitDefaultValues: false,
    });
    const promise = this.rpc.request(
      "meme.subscriptions.Subscriptions",
      "Subscribe",
      "application/json",
      data as object
    );
    return promise.then((data) =>
      SubscribeResponse.fromJson(data as any, { ignoreUnknownFields: true })
    );
  }

  Unsubscribe(request: SubscribeRequest): Promise<SubscribeResponse> {
    const data = SubscribeRequest.toJson(request, {
      useProtoFieldName: true,
      emitDefaultValues: false,
    });
    const promise = this.rpc.request(
      "meme.subscriptions.Subscriptions",
      "Unsubscribe",
      "application/json",
      data as object
    );
    return promise.then((data) =>
      SubscribeResponse.fromJson(data as any, { ignoreUnknownFields: true })
    );
  }

  GetStatus(request: SubscribeRequest): Promise<SubscribeResponse> {
    const data = SubscribeRequest.toJson(request, {
      useProtoFieldName: true,
      emitDefaultValues: false,
    });
    const promise = this.rpc.request(
      "meme.subscriptions.Subscriptions",
      "GetStatus",
      "application/json",
      data as object
    );
    return promise.then((data) =>
      SubscribeResponse.fromJson(data as any, { ignoreUnknownFields: true })
    );
  }
}

export class SubscriptionsClientProtobuf implements SubscriptionsClient {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Subscribe.bind(this);
    this.Unsubscribe.bind(this);
    this.GetStatus.bind(this);
  }
  Subscribe(request: SubscribeRequest): Promise<SubscribeResponse> {
    const data = SubscribeRequest.toBinary(request);
    const promise = this.rpc.request(
      "meme.subscriptions.Subscriptions",
      "Subscribe",
      "application/protobuf",
      data
    );
    return promise.then((data) =>
      SubscribeResponse.fromBinary(data as Uint8Array)
    );
  }

  Unsubscribe(request: SubscribeRequest): Promise<SubscribeResponse> {
    const data = SubscribeRequest.toBinary(request);
    const promise = this.rpc.request(
      "meme.subscriptions.Subscriptions",
      "Unsubscribe",
      "application/protobuf",
      data
    );
    return promise.then((data) =>
      SubscribeResponse.fromBinary(data as Uint8Array)
    );
  }

  GetStatus(request: SubscribeRequest): Promise<SubscribeResponse> {
    const data = SubscribeRequest.toBinary(request);
    const promise = this.rpc.request(
      "meme.subscriptions.Subscriptions",
      "GetStatus",
      "application/protobuf",
      data
    );
    return promise.then((data) =>
      SubscribeResponse.fromBinary(data as Uint8Array)
    );
  }
}

//==================================//
//          Server Code             //
//==================================//

export interface SubscriptionsTwirp<T extends TwirpContext = TwirpContext> {
  Subscribe(ctx: T, request: SubscribeRequest): Promise<SubscribeResponse>;
  Unsubscribe(ctx: T, request: SubscribeRequest): Promise<SubscribeResponse>;
  GetStatus(ctx: T, request: SubscribeRequest): Promise<SubscribeResponse>;
}

export enum SubscriptionsMethod {
  Subscribe = "Subscribe",
  Unsubscribe = "Unsubscribe",
  GetStatus = "GetStatus",
}

export const SubscriptionsMethodList = [
  SubscriptionsMethod.Subscribe,
  SubscriptionsMethod.Unsubscribe,
  SubscriptionsMethod.GetStatus,
];

export function createSubscriptionsServer<
  T extends TwirpContext = TwirpContext
>(service: SubscriptionsTwirp<T>) {
  return new TwirpServer<SubscriptionsTwirp, T>({
    service,
    packageName: "meme.subscriptions",
    serviceName: "Subscriptions",
    methodList: SubscriptionsMethodList,
    matchRoute: matchSubscriptionsRoute,
  });
}

function matchSubscriptionsRoute<T extends TwirpContext = TwirpContext>(
  method: string,
  events: RouterEvents<T>
) {
  switch (method) {
    case "Subscribe":
      return async (
        ctx: T,
        service: SubscriptionsTwirp,
        data: Buffer,
        interceptors?: Interceptor<T, SubscribeRequest, SubscribeResponse>[]
      ) => {
        ctx = { ...ctx, methodName: "Subscribe" };
        await events.onMatch(ctx);
        return handleSubscriptionsSubscribeRequest(
          ctx,
          service,
          data,
          interceptors
        );
      };
    case "Unsubscribe":
      return async (
        ctx: T,
        service: SubscriptionsTwirp,
        data: Buffer,
        interceptors?: Interceptor<T, SubscribeRequest, SubscribeResponse>[]
      ) => {
        ctx = { ...ctx, methodName: "Unsubscribe" };
        await events.onMatch(ctx);
        return handleSubscriptionsUnsubscribeRequest(
          ctx,
          service,
          data,
          interceptors
        );
      };
    case "GetStatus":
      return async (
        ctx: T,
        service: SubscriptionsTwirp,
        data: Buffer,
        interceptors?: Interceptor<T, SubscribeRequest, SubscribeResponse>[]
      ) => {
        ctx = { ...ctx, methodName: "GetStatus" };
        await events.onMatch(ctx);
        return handleSubscriptionsGetStatusRequest(
          ctx,
          service,
          data,
          interceptors
        );
      };
    default:
      events.onNotFound();
      const msg = `no handler found`;
      throw new TwirpError(TwirpErrorCode.BadRoute, msg);
  }
}

function handleSubscriptionsSubscribeRequest<
  T extends TwirpContext = TwirpContext
>(
  ctx: T,
  service: SubscriptionsTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, SubscribeRequest, SubscribeResponse>[]
): Promise<string | Uint8Array> {
  switch (ctx.contentType) {
    case TwirpContentType.JSON:
      return handleSubscriptionsSubscribeJSON<T>(
        ctx,
        service,
        data,
        interceptors
      );
    case TwirpContentType.Protobuf:
      return handleSubscriptionsSubscribeProtobuf<T>(
        ctx,
        service,
        data,
        interceptors
      );
    default:
      const msg = "unexpected Content-Type";
      throw new TwirpError(TwirpErrorCode.BadRoute, msg);
  }
}

function handleSubscriptionsUnsubscribeRequest<
  T extends TwirpContext = TwirpContext
>(
  ctx: T,
  service: SubscriptionsTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, SubscribeRequest, SubscribeResponse>[]
): Promise<string | Uint8Array> {
  switch (ctx.contentType) {
    case TwirpContentType.JSON:
      return handleSubscriptionsUnsubscribeJSON<T>(
        ctx,
        service,
        data,
        interceptors
      );
    case TwirpContentType.Protobuf:
      return handleSubscriptionsUnsubscribeProtobuf<T>(
        ctx,
        service,
        data,
        interceptors
      );
    default:
      const msg = "unexpected Content-Type";
      throw new TwirpError(TwirpErrorCode.BadRoute, msg);
  }
}

function handleSubscriptionsGetStatusRequest<
  T extends TwirpContext = TwirpContext
>(
  ctx: T,
  service: SubscriptionsTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, SubscribeRequest, SubscribeResponse>[]
): Promise<string | Uint8Array> {
  switch (ctx.contentType) {
    case TwirpContentType.JSON:
      return handleSubscriptionsGetStatusJSON<T>(
        ctx,
        service,
        data,
        interceptors
      );
    case TwirpContentType.Protobuf:
      return handleSubscriptionsGetStatusProtobuf<T>(
        ctx,
        service,
        data,
        interceptors
      );
    default:
      const msg = "unexpected Content-Type";
      throw new TwirpError(TwirpErrorCode.BadRoute, msg);
  }
}
async function handleSubscriptionsSubscribeJSON<
  T extends TwirpContext = TwirpContext
>(
  ctx: T,
  service: SubscriptionsTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, SubscribeRequest, SubscribeResponse>[]
) {
  let request: SubscribeRequest;
  let response: SubscribeResponse;

  try {
    const body = JSON.parse(data.toString() || "{}");
    request = SubscribeRequest.fromJson(body, { ignoreUnknownFields: true });
  } catch (e) {
    if (e instanceof Error) {
      const msg = "the json request could not be decoded";
      throw new TwirpError(TwirpErrorCode.Malformed, msg).withCause(e, true);
    }
  }

  if (interceptors && interceptors.length > 0) {
    const interceptor = chainInterceptors(...interceptors) as Interceptor<
      T,
      SubscribeRequest,
      SubscribeResponse
    >;
    response = await interceptor(ctx, request!, (ctx, inputReq) => {
      return service.Subscribe(ctx, inputReq);
    });
  } else {
    response = await service.Subscribe(ctx, request!);
  }

  return JSON.stringify(
    SubscribeResponse.toJson(response, {
      useProtoFieldName: true,
      emitDefaultValues: false,
    }) as string
  );
}

async function handleSubscriptionsUnsubscribeJSON<
  T extends TwirpContext = TwirpContext
>(
  ctx: T,
  service: SubscriptionsTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, SubscribeRequest, SubscribeResponse>[]
) {
  let request: SubscribeRequest;
  let response: SubscribeResponse;

  try {
    const body = JSON.parse(data.toString() || "{}");
    request = SubscribeRequest.fromJson(body, { ignoreUnknownFields: true });
  } catch (e) {
    if (e instanceof Error) {
      const msg = "the json request could not be decoded";
      throw new TwirpError(TwirpErrorCode.Malformed, msg).withCause(e, true);
    }
  }

  if (interceptors && interceptors.length > 0) {
    const interceptor = chainInterceptors(...interceptors) as Interceptor<
      T,
      SubscribeRequest,
      SubscribeResponse
    >;
    response = await interceptor(ctx, request!, (ctx, inputReq) => {
      return service.Unsubscribe(ctx, inputReq);
    });
  } else {
    response = await service.Unsubscribe(ctx, request!);
  }

  return JSON.stringify(
    SubscribeResponse.toJson(response, {
      useProtoFieldName: true,
      emitDefaultValues: false,
    }) as string
  );
}

async function handleSubscriptionsGetStatusJSON<
  T extends TwirpContext = TwirpContext
>(
  ctx: T,
  service: SubscriptionsTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, SubscribeRequest, SubscribeResponse>[]
) {
  let request: SubscribeRequest;
  let response: SubscribeResponse;

  try {
    const body = JSON.parse(data.toString() || "{}");
    request = SubscribeRequest.fromJson(body, { ignoreUnknownFields: true });
  } catch (e) {
    if (e instanceof Error) {
      const msg = "the json request could not be decoded";
      throw new TwirpError(TwirpErrorCode.Malformed, msg).withCause(e, true);
    }
  }

  if (interceptors && interceptors.length > 0) {
    const interceptor = chainInterceptors(...interceptors) as Interceptor<
      T,
      SubscribeRequest,
      SubscribeResponse
    >;
    response = await interceptor(ctx, request!, (ctx, inputReq) => {
      return service.GetStatus(ctx, inputReq);
    });
  } else {
    response = await service.GetStatus(ctx, request!);
  }

  return JSON.stringify(
    SubscribeResponse.toJson(response, {
      useProtoFieldName: true,
      emitDefaultValues: false,
    }) as string
  );
}
async function handleSubscriptionsSubscribeProtobuf<
  T extends TwirpContext = TwirpContext
>(
  ctx: T,
  service: SubscriptionsTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, SubscribeRequest, SubscribeResponse>[]
) {
  let request: SubscribeRequest;
  let response: SubscribeResponse;

  try {
    request = SubscribeRequest.fromBinary(data);
  } catch (e) {
    if (e instanceof Error) {
      const msg = "the protobuf request could not be decoded";
      throw new TwirpError(TwirpErrorCode.Malformed, msg).withCause(e, true);
    }
  }

  if (interceptors && interceptors.length > 0) {
    const interceptor = chainInterceptors(...interceptors) as Interceptor<
      T,
      SubscribeRequest,
      SubscribeResponse
    >;
    response = await interceptor(ctx, request!, (ctx, inputReq) => {
      return service.Subscribe(ctx, inputReq);
    });
  } else {
    response = await service.Subscribe(ctx, request!);
  }

  return Buffer.from(SubscribeResponse.toBinary(response));
}

async function handleSubscriptionsUnsubscribeProtobuf<
  T extends TwirpContext = TwirpContext
>(
  ctx: T,
  service: SubscriptionsTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, SubscribeRequest, SubscribeResponse>[]
) {
  let request: SubscribeRequest;
  let response: SubscribeResponse;

  try {
    request = SubscribeRequest.fromBinary(data);
  } catch (e) {
    if (e instanceof Error) {
      const msg = "the protobuf request could not be decoded";
      throw new TwirpError(TwirpErrorCode.Malformed, msg).withCause(e, true);
    }
  }

  if (interceptors && interceptors.length > 0) {
    const interceptor = chainInterceptors(...interceptors) as Interceptor<
      T,
      SubscribeRequest,
      SubscribeResponse
    >;
    response = await interceptor(ctx, request!, (ctx, inputReq) => {
      return service.Unsubscribe(ctx, inputReq);
    });
  } else {
    response = await service.Unsubscribe(ctx, request!);
  }

  return Buffer.from(SubscribeResponse.toBinary(response));
}

async function handleSubscriptionsGetStatusProtobuf<
  T extends TwirpContext = TwirpContext
>(
  ctx: T,
  service: SubscriptionsTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, SubscribeRequest, SubscribeResponse>[]
) {
  let request: SubscribeRequest;
  let response: SubscribeResponse;

  try {
    request = SubscribeRequest.fromBinary(data);
  } catch (e) {
    if (e instanceof Error) {
      const msg = "the protobuf request could not be decoded";
      throw new TwirpError(TwirpErrorCode.Malformed, msg).withCause(e, true);
    }
  }

  if (interceptors && interceptors.length > 0) {
    const interceptor = chainInterceptors(...interceptors) as Interceptor<
      T,
      SubscribeRequest,
      SubscribeResponse
    >;
    response = await interceptor(ctx, request!, (ctx, inputReq) => {
      return service.GetStatus(ctx, inputReq);
    });
  } else {
    response = await service.GetStatus(ctx, request!);
  }

  return Buffer.from(SubscribeResponse.toBinary(response));
}
