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
import { GetUserRequest, GetUserResponse } from "./users";

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

export interface UsersClient {
  Get(request: GetUserRequest): Promise<GetUserResponse>;
}

export class UsersClientJSON implements UsersClient {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Get.bind(this);
  }
  Get(request: GetUserRequest): Promise<GetUserResponse> {
    const data = GetUserRequest.toJson(request, {
      useProtoFieldName: true,
      emitDefaultValues: false,
    });
    const promise = this.rpc.request(
      "meme.users.Users",
      "Get",
      "application/json",
      data as object
    );
    return promise.then((data) =>
      GetUserResponse.fromJson(data as any, { ignoreUnknownFields: true })
    );
  }
}

export class UsersClientProtobuf implements UsersClient {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Get.bind(this);
  }
  Get(request: GetUserRequest): Promise<GetUserResponse> {
    const data = GetUserRequest.toBinary(request);
    const promise = this.rpc.request(
      "meme.users.Users",
      "Get",
      "application/protobuf",
      data
    );
    return promise.then((data) =>
      GetUserResponse.fromBinary(data as Uint8Array)
    );
  }
}

//==================================//
//          Server Code             //
//==================================//

export interface UsersTwirp<T extends TwirpContext = TwirpContext> {
  Get(ctx: T, request: GetUserRequest): Promise<GetUserResponse>;
}

export enum UsersMethod {
  Get = "Get",
}

export const UsersMethodList = [UsersMethod.Get];

export function createUsersServer<T extends TwirpContext = TwirpContext>(
  service: UsersTwirp<T>
) {
  return new TwirpServer<UsersTwirp, T>({
    service,
    packageName: "meme.users",
    serviceName: "Users",
    methodList: UsersMethodList,
    matchRoute: matchUsersRoute,
  });
}

function matchUsersRoute<T extends TwirpContext = TwirpContext>(
  method: string,
  events: RouterEvents<T>
) {
  switch (method) {
    case "Get":
      return async (
        ctx: T,
        service: UsersTwirp,
        data: Buffer,
        interceptors?: Interceptor<T, GetUserRequest, GetUserResponse>[]
      ) => {
        ctx = { ...ctx, methodName: "Get" };
        await events.onMatch(ctx);
        return handleUsersGetRequest(ctx, service, data, interceptors);
      };
    default:
      events.onNotFound();
      const msg = `no handler found`;
      throw new TwirpError(TwirpErrorCode.BadRoute, msg);
  }
}

function handleUsersGetRequest<T extends TwirpContext = TwirpContext>(
  ctx: T,
  service: UsersTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, GetUserRequest, GetUserResponse>[]
): Promise<string | Uint8Array> {
  switch (ctx.contentType) {
    case TwirpContentType.JSON:
      return handleUsersGetJSON<T>(ctx, service, data, interceptors);
    case TwirpContentType.Protobuf:
      return handleUsersGetProtobuf<T>(ctx, service, data, interceptors);
    default:
      const msg = "unexpected Content-Type";
      throw new TwirpError(TwirpErrorCode.BadRoute, msg);
  }
}
async function handleUsersGetJSON<T extends TwirpContext = TwirpContext>(
  ctx: T,
  service: UsersTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, GetUserRequest, GetUserResponse>[]
) {
  let request: GetUserRequest;
  let response: GetUserResponse;

  try {
    const body = JSON.parse(data.toString() || "{}");
    request = GetUserRequest.fromJson(body, { ignoreUnknownFields: true });
  } catch (e) {
    if (e instanceof Error) {
      const msg = "the json request could not be decoded";
      throw new TwirpError(TwirpErrorCode.Malformed, msg).withCause(e, true);
    }
  }

  if (interceptors && interceptors.length > 0) {
    const interceptor = chainInterceptors(...interceptors) as Interceptor<
      T,
      GetUserRequest,
      GetUserResponse
    >;
    response = await interceptor(ctx, request!, (ctx, inputReq) => {
      return service.Get(ctx, inputReq);
    });
  } else {
    response = await service.Get(ctx, request!);
  }

  return JSON.stringify(
    GetUserResponse.toJson(response, {
      useProtoFieldName: true,
      emitDefaultValues: false,
    }) as string
  );
}
async function handleUsersGetProtobuf<T extends TwirpContext = TwirpContext>(
  ctx: T,
  service: UsersTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, GetUserRequest, GetUserResponse>[]
) {
  let request: GetUserRequest;
  let response: GetUserResponse;

  try {
    request = GetUserRequest.fromBinary(data);
  } catch (e) {
    if (e instanceof Error) {
      const msg = "the protobuf request could not be decoded";
      throw new TwirpError(TwirpErrorCode.Malformed, msg).withCause(e, true);
    }
  }

  if (interceptors && interceptors.length > 0) {
    const interceptor = chainInterceptors(...interceptors) as Interceptor<
      T,
      GetUserRequest,
      GetUserResponse
    >;
    response = await interceptor(ctx, request!, (ctx, inputReq) => {
      return service.Get(ctx, inputReq);
    });
  } else {
    response = await service.Get(ctx, request!);
  }

  return Buffer.from(GetUserResponse.toBinary(response));
}
