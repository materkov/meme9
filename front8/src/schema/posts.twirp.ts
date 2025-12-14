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
import {
  PublishRequest,
  PublishResponse,
  GetByUsersRequest,
  GetByUsersResponse,
  GetPostRequest,
  GetPostResponse,
} from "./posts";

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

export interface PostsClient {
  Publish(request: PublishRequest): Promise<PublishResponse>;
  GetByUsers(request: GetByUsersRequest): Promise<GetByUsersResponse>;
  Get(request: GetPostRequest): Promise<GetPostResponse>;
}

export class PostsClientJSON implements PostsClient {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Publish.bind(this);
    this.GetByUsers.bind(this);
    this.Get.bind(this);
  }
  Publish(request: PublishRequest): Promise<PublishResponse> {
    const data = PublishRequest.toJson(request, {
      useProtoFieldName: true,
      emitDefaultValues: false,
    });
    const promise = this.rpc.request(
      "meme.posts.Posts",
      "Publish",
      "application/json",
      data as object
    );
    return promise.then((data) =>
      PublishResponse.fromJson(data as any, { ignoreUnknownFields: true })
    );
  }

  GetByUsers(request: GetByUsersRequest): Promise<GetByUsersResponse> {
    const data = GetByUsersRequest.toJson(request, {
      useProtoFieldName: true,
      emitDefaultValues: false,
    });
    const promise = this.rpc.request(
      "meme.posts.Posts",
      "GetByUsers",
      "application/json",
      data as object
    );
    return promise.then((data) =>
      GetByUsersResponse.fromJson(data as any, { ignoreUnknownFields: true })
    );
  }

  Get(request: GetPostRequest): Promise<GetPostResponse> {
    const data = GetPostRequest.toJson(request, {
      useProtoFieldName: true,
      emitDefaultValues: false,
    });
    const promise = this.rpc.request(
      "meme.posts.Posts",
      "Get",
      "application/json",
      data as object
    );
    return promise.then((data) =>
      GetPostResponse.fromJson(data as any, { ignoreUnknownFields: true })
    );
  }
}

export class PostsClientProtobuf implements PostsClient {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Publish.bind(this);
    this.GetByUsers.bind(this);
    this.Get.bind(this);
  }
  Publish(request: PublishRequest): Promise<PublishResponse> {
    const data = PublishRequest.toBinary(request);
    const promise = this.rpc.request(
      "meme.posts.Posts",
      "Publish",
      "application/protobuf",
      data
    );
    return promise.then((data) =>
      PublishResponse.fromBinary(data as Uint8Array)
    );
  }

  GetByUsers(request: GetByUsersRequest): Promise<GetByUsersResponse> {
    const data = GetByUsersRequest.toBinary(request);
    const promise = this.rpc.request(
      "meme.posts.Posts",
      "GetByUsers",
      "application/protobuf",
      data
    );
    return promise.then((data) =>
      GetByUsersResponse.fromBinary(data as Uint8Array)
    );
  }

  Get(request: GetPostRequest): Promise<GetPostResponse> {
    const data = GetPostRequest.toBinary(request);
    const promise = this.rpc.request(
      "meme.posts.Posts",
      "Get",
      "application/protobuf",
      data
    );
    return promise.then((data) =>
      GetPostResponse.fromBinary(data as Uint8Array)
    );
  }
}

//==================================//
//          Server Code             //
//==================================//

export interface PostsTwirp<T extends TwirpContext = TwirpContext> {
  Publish(ctx: T, request: PublishRequest): Promise<PublishResponse>;
  GetByUsers(ctx: T, request: GetByUsersRequest): Promise<GetByUsersResponse>;
  Get(ctx: T, request: GetPostRequest): Promise<GetPostResponse>;
}

export enum PostsMethod {
  Publish = "Publish",
  GetByUsers = "GetByUsers",
  Get = "Get",
}

export const PostsMethodList = [
  PostsMethod.Publish,
  PostsMethod.GetByUsers,
  PostsMethod.Get,
];

export function createPostsServer<T extends TwirpContext = TwirpContext>(
  service: PostsTwirp<T>
) {
  return new TwirpServer<PostsTwirp, T>({
    service,
    packageName: "meme.posts",
    serviceName: "Posts",
    methodList: PostsMethodList,
    matchRoute: matchPostsRoute,
  });
}

function matchPostsRoute<T extends TwirpContext = TwirpContext>(
  method: string,
  events: RouterEvents<T>
) {
  switch (method) {
    case "Publish":
      return async (
        ctx: T,
        service: PostsTwirp,
        data: Buffer,
        interceptors?: Interceptor<T, PublishRequest, PublishResponse>[]
      ) => {
        ctx = { ...ctx, methodName: "Publish" };
        await events.onMatch(ctx);
        return handlePostsPublishRequest(ctx, service, data, interceptors);
      };
    case "GetByUsers":
      return async (
        ctx: T,
        service: PostsTwirp,
        data: Buffer,
        interceptors?: Interceptor<T, GetByUsersRequest, GetByUsersResponse>[]
      ) => {
        ctx = { ...ctx, methodName: "GetByUsers" };
        await events.onMatch(ctx);
        return handlePostsGetByUsersRequest(ctx, service, data, interceptors);
      };
    case "Get":
      return async (
        ctx: T,
        service: PostsTwirp,
        data: Buffer,
        interceptors?: Interceptor<T, GetPostRequest, GetPostResponse>[]
      ) => {
        ctx = { ...ctx, methodName: "Get" };
        await events.onMatch(ctx);
        return handlePostsGetRequest(ctx, service, data, interceptors);
      };
    default:
      events.onNotFound();
      const msg = `no handler found`;
      throw new TwirpError(TwirpErrorCode.BadRoute, msg);
  }
}

function handlePostsPublishRequest<T extends TwirpContext = TwirpContext>(
  ctx: T,
  service: PostsTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, PublishRequest, PublishResponse>[]
): Promise<string | Uint8Array> {
  switch (ctx.contentType) {
    case TwirpContentType.JSON:
      return handlePostsPublishJSON<T>(ctx, service, data, interceptors);
    case TwirpContentType.Protobuf:
      return handlePostsPublishProtobuf<T>(ctx, service, data, interceptors);
    default:
      const msg = "unexpected Content-Type";
      throw new TwirpError(TwirpErrorCode.BadRoute, msg);
  }
}

function handlePostsGetByUsersRequest<T extends TwirpContext = TwirpContext>(
  ctx: T,
  service: PostsTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, GetByUsersRequest, GetByUsersResponse>[]
): Promise<string | Uint8Array> {
  switch (ctx.contentType) {
    case TwirpContentType.JSON:
      return handlePostsGetByUsersJSON<T>(ctx, service, data, interceptors);
    case TwirpContentType.Protobuf:
      return handlePostsGetByUsersProtobuf<T>(ctx, service, data, interceptors);
    default:
      const msg = "unexpected Content-Type";
      throw new TwirpError(TwirpErrorCode.BadRoute, msg);
  }
}

function handlePostsGetRequest<T extends TwirpContext = TwirpContext>(
  ctx: T,
  service: PostsTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, GetPostRequest, GetPostResponse>[]
): Promise<string | Uint8Array> {
  switch (ctx.contentType) {
    case TwirpContentType.JSON:
      return handlePostsGetJSON<T>(ctx, service, data, interceptors);
    case TwirpContentType.Protobuf:
      return handlePostsGetProtobuf<T>(ctx, service, data, interceptors);
    default:
      const msg = "unexpected Content-Type";
      throw new TwirpError(TwirpErrorCode.BadRoute, msg);
  }
}
async function handlePostsPublishJSON<T extends TwirpContext = TwirpContext>(
  ctx: T,
  service: PostsTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, PublishRequest, PublishResponse>[]
) {
  let request: PublishRequest;
  let response: PublishResponse;

  try {
    const body = JSON.parse(data.toString() || "{}");
    request = PublishRequest.fromJson(body, { ignoreUnknownFields: true });
  } catch (e) {
    if (e instanceof Error) {
      const msg = "the json request could not be decoded";
      throw new TwirpError(TwirpErrorCode.Malformed, msg).withCause(e, true);
    }
  }

  if (interceptors && interceptors.length > 0) {
    const interceptor = chainInterceptors(...interceptors) as Interceptor<
      T,
      PublishRequest,
      PublishResponse
    >;
    response = await interceptor(ctx, request!, (ctx, inputReq) => {
      return service.Publish(ctx, inputReq);
    });
  } else {
    response = await service.Publish(ctx, request!);
  }

  return JSON.stringify(
    PublishResponse.toJson(response, {
      useProtoFieldName: true,
      emitDefaultValues: false,
    }) as string
  );
}

async function handlePostsGetByUsersJSON<T extends TwirpContext = TwirpContext>(
  ctx: T,
  service: PostsTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, GetByUsersRequest, GetByUsersResponse>[]
) {
  let request: GetByUsersRequest;
  let response: GetByUsersResponse;

  try {
    const body = JSON.parse(data.toString() || "{}");
    request = GetByUsersRequest.fromJson(body, { ignoreUnknownFields: true });
  } catch (e) {
    if (e instanceof Error) {
      const msg = "the json request could not be decoded";
      throw new TwirpError(TwirpErrorCode.Malformed, msg).withCause(e, true);
    }
  }

  if (interceptors && interceptors.length > 0) {
    const interceptor = chainInterceptors(...interceptors) as Interceptor<
      T,
      GetByUsersRequest,
      GetByUsersResponse
    >;
    response = await interceptor(ctx, request!, (ctx, inputReq) => {
      return service.GetByUsers(ctx, inputReq);
    });
  } else {
    response = await service.GetByUsers(ctx, request!);
  }

  return JSON.stringify(
    GetByUsersResponse.toJson(response, {
      useProtoFieldName: true,
      emitDefaultValues: false,
    }) as string
  );
}

async function handlePostsGetJSON<T extends TwirpContext = TwirpContext>(
  ctx: T,
  service: PostsTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, GetPostRequest, GetPostResponse>[]
) {
  let request: GetPostRequest;
  let response: GetPostResponse;

  try {
    const body = JSON.parse(data.toString() || "{}");
    request = GetPostRequest.fromJson(body, { ignoreUnknownFields: true });
  } catch (e) {
    if (e instanceof Error) {
      const msg = "the json request could not be decoded";
      throw new TwirpError(TwirpErrorCode.Malformed, msg).withCause(e, true);
    }
  }

  if (interceptors && interceptors.length > 0) {
    const interceptor = chainInterceptors(...interceptors) as Interceptor<
      T,
      GetPostRequest,
      GetPostResponse
    >;
    response = await interceptor(ctx, request!, (ctx, inputReq) => {
      return service.Get(ctx, inputReq);
    });
  } else {
    response = await service.Get(ctx, request!);
  }

  return JSON.stringify(
    GetPostResponse.toJson(response, {
      useProtoFieldName: true,
      emitDefaultValues: false,
    }) as string
  );
}
async function handlePostsPublishProtobuf<
  T extends TwirpContext = TwirpContext
>(
  ctx: T,
  service: PostsTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, PublishRequest, PublishResponse>[]
) {
  let request: PublishRequest;
  let response: PublishResponse;

  try {
    request = PublishRequest.fromBinary(data);
  } catch (e) {
    if (e instanceof Error) {
      const msg = "the protobuf request could not be decoded";
      throw new TwirpError(TwirpErrorCode.Malformed, msg).withCause(e, true);
    }
  }

  if (interceptors && interceptors.length > 0) {
    const interceptor = chainInterceptors(...interceptors) as Interceptor<
      T,
      PublishRequest,
      PublishResponse
    >;
    response = await interceptor(ctx, request!, (ctx, inputReq) => {
      return service.Publish(ctx, inputReq);
    });
  } else {
    response = await service.Publish(ctx, request!);
  }

  return Buffer.from(PublishResponse.toBinary(response));
}

async function handlePostsGetByUsersProtobuf<
  T extends TwirpContext = TwirpContext
>(
  ctx: T,
  service: PostsTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, GetByUsersRequest, GetByUsersResponse>[]
) {
  let request: GetByUsersRequest;
  let response: GetByUsersResponse;

  try {
    request = GetByUsersRequest.fromBinary(data);
  } catch (e) {
    if (e instanceof Error) {
      const msg = "the protobuf request could not be decoded";
      throw new TwirpError(TwirpErrorCode.Malformed, msg).withCause(e, true);
    }
  }

  if (interceptors && interceptors.length > 0) {
    const interceptor = chainInterceptors(...interceptors) as Interceptor<
      T,
      GetByUsersRequest,
      GetByUsersResponse
    >;
    response = await interceptor(ctx, request!, (ctx, inputReq) => {
      return service.GetByUsers(ctx, inputReq);
    });
  } else {
    response = await service.GetByUsers(ctx, request!);
  }

  return Buffer.from(GetByUsersResponse.toBinary(response));
}

async function handlePostsGetProtobuf<T extends TwirpContext = TwirpContext>(
  ctx: T,
  service: PostsTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, GetPostRequest, GetPostResponse>[]
) {
  let request: GetPostRequest;
  let response: GetPostResponse;

  try {
    request = GetPostRequest.fromBinary(data);
  } catch (e) {
    if (e instanceof Error) {
      const msg = "the protobuf request could not be decoded";
      throw new TwirpError(TwirpErrorCode.Malformed, msg).withCause(e, true);
    }
  }

  if (interceptors && interceptors.length > 0) {
    const interceptor = chainInterceptors(...interceptors) as Interceptor<
      T,
      GetPostRequest,
      GetPostResponse
    >;
    response = await interceptor(ctx, request!, (ctx, inputReq) => {
      return service.Get(ctx, inputReq);
    });
  } else {
    response = await service.Get(ctx, request!);
  }

  return Buffer.from(GetPostResponse.toBinary(response));
}
