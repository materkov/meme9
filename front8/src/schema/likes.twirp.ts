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
import { LikeRequest, LikeResponse } from "./likes";

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

export interface LikesClient {
  Like(request: LikeRequest): Promise<LikeResponse>;
  Unlike(request: LikeRequest): Promise<LikeResponse>;
}

export class LikesClientJSON implements LikesClient {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Like.bind(this);
    this.Unlike.bind(this);
  }
  Like(request: LikeRequest): Promise<LikeResponse> {
    const data = LikeRequest.toJson(request, {
      useProtoFieldName: true,
      emitDefaultValues: false,
    });
    const promise = this.rpc.request(
      "meme.likes.Likes",
      "Like",
      "application/json",
      data as object
    );
    return promise.then((data) =>
      LikeResponse.fromJson(data as any, { ignoreUnknownFields: true })
    );
  }

  Unlike(request: LikeRequest): Promise<LikeResponse> {
    const data = LikeRequest.toJson(request, {
      useProtoFieldName: true,
      emitDefaultValues: false,
    });
    const promise = this.rpc.request(
      "meme.likes.Likes",
      "Unlike",
      "application/json",
      data as object
    );
    return promise.then((data) =>
      LikeResponse.fromJson(data as any, { ignoreUnknownFields: true })
    );
  }
}

export class LikesClientProtobuf implements LikesClient {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Like.bind(this);
    this.Unlike.bind(this);
  }
  Like(request: LikeRequest): Promise<LikeResponse> {
    const data = LikeRequest.toBinary(request);
    const promise = this.rpc.request(
      "meme.likes.Likes",
      "Like",
      "application/protobuf",
      data
    );
    return promise.then((data) => LikeResponse.fromBinary(data as Uint8Array));
  }

  Unlike(request: LikeRequest): Promise<LikeResponse> {
    const data = LikeRequest.toBinary(request);
    const promise = this.rpc.request(
      "meme.likes.Likes",
      "Unlike",
      "application/protobuf",
      data
    );
    return promise.then((data) => LikeResponse.fromBinary(data as Uint8Array));
  }
}

//==================================//
//          Server Code             //
//==================================//

export interface LikesTwirp<T extends TwirpContext = TwirpContext> {
  Like(ctx: T, request: LikeRequest): Promise<LikeResponse>;
  Unlike(ctx: T, request: LikeRequest): Promise<LikeResponse>;
}

export enum LikesMethod {
  Like = "Like",
  Unlike = "Unlike",
}

export const LikesMethodList = [LikesMethod.Like, LikesMethod.Unlike];

export function createLikesServer<T extends TwirpContext = TwirpContext>(
  service: LikesTwirp<T>
) {
  return new TwirpServer<LikesTwirp, T>({
    service,
    packageName: "meme.likes",
    serviceName: "Likes",
    methodList: LikesMethodList,
    matchRoute: matchLikesRoute,
  });
}

function matchLikesRoute<T extends TwirpContext = TwirpContext>(
  method: string,
  events: RouterEvents<T>
) {
  switch (method) {
    case "Like":
      return async (
        ctx: T,
        service: LikesTwirp,
        data: Buffer,
        interceptors?: Interceptor<T, LikeRequest, LikeResponse>[]
      ) => {
        ctx = { ...ctx, methodName: "Like" };
        await events.onMatch(ctx);
        return handleLikesLikeRequest(ctx, service, data, interceptors);
      };
    case "Unlike":
      return async (
        ctx: T,
        service: LikesTwirp,
        data: Buffer,
        interceptors?: Interceptor<T, LikeRequest, LikeResponse>[]
      ) => {
        ctx = { ...ctx, methodName: "Unlike" };
        await events.onMatch(ctx);
        return handleLikesUnlikeRequest(ctx, service, data, interceptors);
      };
    default:
      events.onNotFound();
      const msg = `no handler found`;
      throw new TwirpError(TwirpErrorCode.BadRoute, msg);
  }
}

function handleLikesLikeRequest<T extends TwirpContext = TwirpContext>(
  ctx: T,
  service: LikesTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, LikeRequest, LikeResponse>[]
): Promise<string | Uint8Array> {
  switch (ctx.contentType) {
    case TwirpContentType.JSON:
      return handleLikesLikeJSON<T>(ctx, service, data, interceptors);
    case TwirpContentType.Protobuf:
      return handleLikesLikeProtobuf<T>(ctx, service, data, interceptors);
    default:
      const msg = "unexpected Content-Type";
      throw new TwirpError(TwirpErrorCode.BadRoute, msg);
  }
}

function handleLikesUnlikeRequest<T extends TwirpContext = TwirpContext>(
  ctx: T,
  service: LikesTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, LikeRequest, LikeResponse>[]
): Promise<string | Uint8Array> {
  switch (ctx.contentType) {
    case TwirpContentType.JSON:
      return handleLikesUnlikeJSON<T>(ctx, service, data, interceptors);
    case TwirpContentType.Protobuf:
      return handleLikesUnlikeProtobuf<T>(ctx, service, data, interceptors);
    default:
      const msg = "unexpected Content-Type";
      throw new TwirpError(TwirpErrorCode.BadRoute, msg);
  }
}
async function handleLikesLikeJSON<T extends TwirpContext = TwirpContext>(
  ctx: T,
  service: LikesTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, LikeRequest, LikeResponse>[]
) {
  let request: LikeRequest;
  let response: LikeResponse;

  try {
    const body = JSON.parse(data.toString() || "{}");
    request = LikeRequest.fromJson(body, { ignoreUnknownFields: true });
  } catch (e) {
    if (e instanceof Error) {
      const msg = "the json request could not be decoded";
      throw new TwirpError(TwirpErrorCode.Malformed, msg).withCause(e, true);
    }
  }

  if (interceptors && interceptors.length > 0) {
    const interceptor = chainInterceptors(...interceptors) as Interceptor<
      T,
      LikeRequest,
      LikeResponse
    >;
    response = await interceptor(ctx, request!, (ctx, inputReq) => {
      return service.Like(ctx, inputReq);
    });
  } else {
    response = await service.Like(ctx, request!);
  }

  return JSON.stringify(
    LikeResponse.toJson(response, {
      useProtoFieldName: true,
      emitDefaultValues: false,
    }) as string
  );
}

async function handleLikesUnlikeJSON<T extends TwirpContext = TwirpContext>(
  ctx: T,
  service: LikesTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, LikeRequest, LikeResponse>[]
) {
  let request: LikeRequest;
  let response: LikeResponse;

  try {
    const body = JSON.parse(data.toString() || "{}");
    request = LikeRequest.fromJson(body, { ignoreUnknownFields: true });
  } catch (e) {
    if (e instanceof Error) {
      const msg = "the json request could not be decoded";
      throw new TwirpError(TwirpErrorCode.Malformed, msg).withCause(e, true);
    }
  }

  if (interceptors && interceptors.length > 0) {
    const interceptor = chainInterceptors(...interceptors) as Interceptor<
      T,
      LikeRequest,
      LikeResponse
    >;
    response = await interceptor(ctx, request!, (ctx, inputReq) => {
      return service.Unlike(ctx, inputReq);
    });
  } else {
    response = await service.Unlike(ctx, request!);
  }

  return JSON.stringify(
    LikeResponse.toJson(response, {
      useProtoFieldName: true,
      emitDefaultValues: false,
    }) as string
  );
}
async function handleLikesLikeProtobuf<T extends TwirpContext = TwirpContext>(
  ctx: T,
  service: LikesTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, LikeRequest, LikeResponse>[]
) {
  let request: LikeRequest;
  let response: LikeResponse;

  try {
    request = LikeRequest.fromBinary(data);
  } catch (e) {
    if (e instanceof Error) {
      const msg = "the protobuf request could not be decoded";
      throw new TwirpError(TwirpErrorCode.Malformed, msg).withCause(e, true);
    }
  }

  if (interceptors && interceptors.length > 0) {
    const interceptor = chainInterceptors(...interceptors) as Interceptor<
      T,
      LikeRequest,
      LikeResponse
    >;
    response = await interceptor(ctx, request!, (ctx, inputReq) => {
      return service.Like(ctx, inputReq);
    });
  } else {
    response = await service.Like(ctx, request!);
  }

  return Buffer.from(LikeResponse.toBinary(response));
}

async function handleLikesUnlikeProtobuf<T extends TwirpContext = TwirpContext>(
  ctx: T,
  service: LikesTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, LikeRequest, LikeResponse>[]
) {
  let request: LikeRequest;
  let response: LikeResponse;

  try {
    request = LikeRequest.fromBinary(data);
  } catch (e) {
    if (e instanceof Error) {
      const msg = "the protobuf request could not be decoded";
      throw new TwirpError(TwirpErrorCode.Malformed, msg).withCause(e, true);
    }
  }

  if (interceptors && interceptors.length > 0) {
    const interceptor = chainInterceptors(...interceptors) as Interceptor<
      T,
      LikeRequest,
      LikeResponse
    >;
    response = await interceptor(ctx, request!, (ctx, inputReq) => {
      return service.Unlike(ctx, inputReq);
    });
  } else {
    response = await service.Unlike(ctx, request!);
  }

  return Buffer.from(LikeResponse.toBinary(response));
}
