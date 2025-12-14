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
import { FeedRequest, FeedResponse } from "./feed";

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

export interface FeedClient {
  GetFeed(request: FeedRequest): Promise<FeedResponse>;
}

export class FeedClientJSON implements FeedClient {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.GetFeed.bind(this);
  }
  GetFeed(request: FeedRequest): Promise<FeedResponse> {
    const data = FeedRequest.toJson(request, {
      useProtoFieldName: true,
      emitDefaultValues: false,
    });
    const promise = this.rpc.request(
      "meme.feed.Feed",
      "GetFeed",
      "application/json",
      data as object
    );
    return promise.then((data) =>
      FeedResponse.fromJson(data as any, { ignoreUnknownFields: true })
    );
  }
}

export class FeedClientProtobuf implements FeedClient {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.GetFeed.bind(this);
  }
  GetFeed(request: FeedRequest): Promise<FeedResponse> {
    const data = FeedRequest.toBinary(request);
    const promise = this.rpc.request(
      "meme.feed.Feed",
      "GetFeed",
      "application/protobuf",
      data
    );
    return promise.then((data) => FeedResponse.fromBinary(data as Uint8Array));
  }
}

//==================================//
//          Server Code             //
//==================================//

export interface FeedTwirp<T extends TwirpContext = TwirpContext> {
  GetFeed(ctx: T, request: FeedRequest): Promise<FeedResponse>;
}

export enum FeedMethod {
  GetFeed = "GetFeed",
}

export const FeedMethodList = [FeedMethod.GetFeed];

export function createFeedServer<T extends TwirpContext = TwirpContext>(
  service: FeedTwirp<T>
) {
  return new TwirpServer<FeedTwirp, T>({
    service,
    packageName: "meme.feed",
    serviceName: "Feed",
    methodList: FeedMethodList,
    matchRoute: matchFeedRoute,
  });
}

function matchFeedRoute<T extends TwirpContext = TwirpContext>(
  method: string,
  events: RouterEvents<T>
) {
  switch (method) {
    case "GetFeed":
      return async (
        ctx: T,
        service: FeedTwirp,
        data: Buffer,
        interceptors?: Interceptor<T, FeedRequest, FeedResponse>[]
      ) => {
        ctx = { ...ctx, methodName: "GetFeed" };
        await events.onMatch(ctx);
        return handleFeedGetFeedRequest(ctx, service, data, interceptors);
      };
    default:
      events.onNotFound();
      const msg = `no handler found`;
      throw new TwirpError(TwirpErrorCode.BadRoute, msg);
  }
}

function handleFeedGetFeedRequest<T extends TwirpContext = TwirpContext>(
  ctx: T,
  service: FeedTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, FeedRequest, FeedResponse>[]
): Promise<string | Uint8Array> {
  switch (ctx.contentType) {
    case TwirpContentType.JSON:
      return handleFeedGetFeedJSON<T>(ctx, service, data, interceptors);
    case TwirpContentType.Protobuf:
      return handleFeedGetFeedProtobuf<T>(ctx, service, data, interceptors);
    default:
      const msg = "unexpected Content-Type";
      throw new TwirpError(TwirpErrorCode.BadRoute, msg);
  }
}
async function handleFeedGetFeedJSON<T extends TwirpContext = TwirpContext>(
  ctx: T,
  service: FeedTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, FeedRequest, FeedResponse>[]
) {
  let request: FeedRequest;
  let response: FeedResponse;

  try {
    const body = JSON.parse(data.toString() || "{}");
    request = FeedRequest.fromJson(body, { ignoreUnknownFields: true });
  } catch (e) {
    if (e instanceof Error) {
      const msg = "the json request could not be decoded";
      throw new TwirpError(TwirpErrorCode.Malformed, msg).withCause(e, true);
    }
  }

  if (interceptors && interceptors.length > 0) {
    const interceptor = chainInterceptors(...interceptors) as Interceptor<
      T,
      FeedRequest,
      FeedResponse
    >;
    response = await interceptor(ctx, request!, (ctx, inputReq) => {
      return service.GetFeed(ctx, inputReq);
    });
  } else {
    response = await service.GetFeed(ctx, request!);
  }

  return JSON.stringify(
    FeedResponse.toJson(response, {
      useProtoFieldName: true,
      emitDefaultValues: false,
    }) as string
  );
}
async function handleFeedGetFeedProtobuf<T extends TwirpContext = TwirpContext>(
  ctx: T,
  service: FeedTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, FeedRequest, FeedResponse>[]
) {
  let request: FeedRequest;
  let response: FeedResponse;

  try {
    request = FeedRequest.fromBinary(data);
  } catch (e) {
    if (e instanceof Error) {
      const msg = "the protobuf request could not be decoded";
      throw new TwirpError(TwirpErrorCode.Malformed, msg).withCause(e, true);
    }
  }

  if (interceptors && interceptors.length > 0) {
    const interceptor = chainInterceptors(...interceptors) as Interceptor<
      T,
      FeedRequest,
      FeedResponse
    >;
    response = await interceptor(ctx, request!, (ctx, inputReq) => {
      return service.GetFeed(ctx, inputReq);
    });
  } else {
    response = await service.GetFeed(ctx, request!);
  }

  return Buffer.from(FeedResponse.toBinary(response));
}
