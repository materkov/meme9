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
  LoginRequest,
  LoginResponse,
  RegisterRequest,
  VerifyTokenRequest,
  VerifyTokenResponse,
} from "./auth";

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

export interface AuthClient {
  Login(request: LoginRequest): Promise<LoginResponse>;
  Register(request: RegisterRequest): Promise<LoginResponse>;
  VerifyToken(request: VerifyTokenRequest): Promise<VerifyTokenResponse>;
}

export class AuthClientJSON implements AuthClient {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Login.bind(this);
    this.Register.bind(this);
    this.VerifyToken.bind(this);
  }
  Login(request: LoginRequest): Promise<LoginResponse> {
    const data = LoginRequest.toJson(request, {
      useProtoFieldName: true,
      emitDefaultValues: false,
    });
    const promise = this.rpc.request(
      "meme.auth.Auth",
      "Login",
      "application/json",
      data as object
    );
    return promise.then((data) =>
      LoginResponse.fromJson(data as any, { ignoreUnknownFields: true })
    );
  }

  Register(request: RegisterRequest): Promise<LoginResponse> {
    const data = RegisterRequest.toJson(request, {
      useProtoFieldName: true,
      emitDefaultValues: false,
    });
    const promise = this.rpc.request(
      "meme.auth.Auth",
      "Register",
      "application/json",
      data as object
    );
    return promise.then((data) =>
      LoginResponse.fromJson(data as any, { ignoreUnknownFields: true })
    );
  }

  VerifyToken(request: VerifyTokenRequest): Promise<VerifyTokenResponse> {
    const data = VerifyTokenRequest.toJson(request, {
      useProtoFieldName: true,
      emitDefaultValues: false,
    });
    const promise = this.rpc.request(
      "meme.auth.Auth",
      "VerifyToken",
      "application/json",
      data as object
    );
    return promise.then((data) =>
      VerifyTokenResponse.fromJson(data as any, { ignoreUnknownFields: true })
    );
  }
}

export class AuthClientProtobuf implements AuthClient {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.Login.bind(this);
    this.Register.bind(this);
    this.VerifyToken.bind(this);
  }
  Login(request: LoginRequest): Promise<LoginResponse> {
    const data = LoginRequest.toBinary(request);
    const promise = this.rpc.request(
      "meme.auth.Auth",
      "Login",
      "application/protobuf",
      data
    );
    return promise.then((data) => LoginResponse.fromBinary(data as Uint8Array));
  }

  Register(request: RegisterRequest): Promise<LoginResponse> {
    const data = RegisterRequest.toBinary(request);
    const promise = this.rpc.request(
      "meme.auth.Auth",
      "Register",
      "application/protobuf",
      data
    );
    return promise.then((data) => LoginResponse.fromBinary(data as Uint8Array));
  }

  VerifyToken(request: VerifyTokenRequest): Promise<VerifyTokenResponse> {
    const data = VerifyTokenRequest.toBinary(request);
    const promise = this.rpc.request(
      "meme.auth.Auth",
      "VerifyToken",
      "application/protobuf",
      data
    );
    return promise.then((data) =>
      VerifyTokenResponse.fromBinary(data as Uint8Array)
    );
  }
}

//==================================//
//          Server Code             //
//==================================//

export interface AuthTwirp<T extends TwirpContext = TwirpContext> {
  Login(ctx: T, request: LoginRequest): Promise<LoginResponse>;
  Register(ctx: T, request: RegisterRequest): Promise<LoginResponse>;
  VerifyToken(
    ctx: T,
    request: VerifyTokenRequest
  ): Promise<VerifyTokenResponse>;
}

export enum AuthMethod {
  Login = "Login",
  Register = "Register",
  VerifyToken = "VerifyToken",
}

export const AuthMethodList = [
  AuthMethod.Login,
  AuthMethod.Register,
  AuthMethod.VerifyToken,
];

export function createAuthServer<T extends TwirpContext = TwirpContext>(
  service: AuthTwirp<T>
) {
  return new TwirpServer<AuthTwirp, T>({
    service,
    packageName: "meme.auth",
    serviceName: "Auth",
    methodList: AuthMethodList,
    matchRoute: matchAuthRoute,
  });
}

function matchAuthRoute<T extends TwirpContext = TwirpContext>(
  method: string,
  events: RouterEvents<T>
) {
  switch (method) {
    case "Login":
      return async (
        ctx: T,
        service: AuthTwirp,
        data: Buffer,
        interceptors?: Interceptor<T, LoginRequest, LoginResponse>[]
      ) => {
        ctx = { ...ctx, methodName: "Login" };
        await events.onMatch(ctx);
        return handleAuthLoginRequest(ctx, service, data, interceptors);
      };
    case "Register":
      return async (
        ctx: T,
        service: AuthTwirp,
        data: Buffer,
        interceptors?: Interceptor<T, RegisterRequest, LoginResponse>[]
      ) => {
        ctx = { ...ctx, methodName: "Register" };
        await events.onMatch(ctx);
        return handleAuthRegisterRequest(ctx, service, data, interceptors);
      };
    case "VerifyToken":
      return async (
        ctx: T,
        service: AuthTwirp,
        data: Buffer,
        interceptors?: Interceptor<T, VerifyTokenRequest, VerifyTokenResponse>[]
      ) => {
        ctx = { ...ctx, methodName: "VerifyToken" };
        await events.onMatch(ctx);
        return handleAuthVerifyTokenRequest(ctx, service, data, interceptors);
      };
    default:
      events.onNotFound();
      const msg = `no handler found`;
      throw new TwirpError(TwirpErrorCode.BadRoute, msg);
  }
}

function handleAuthLoginRequest<T extends TwirpContext = TwirpContext>(
  ctx: T,
  service: AuthTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, LoginRequest, LoginResponse>[]
): Promise<string | Uint8Array> {
  switch (ctx.contentType) {
    case TwirpContentType.JSON:
      return handleAuthLoginJSON<T>(ctx, service, data, interceptors);
    case TwirpContentType.Protobuf:
      return handleAuthLoginProtobuf<T>(ctx, service, data, interceptors);
    default:
      const msg = "unexpected Content-Type";
      throw new TwirpError(TwirpErrorCode.BadRoute, msg);
  }
}

function handleAuthRegisterRequest<T extends TwirpContext = TwirpContext>(
  ctx: T,
  service: AuthTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, RegisterRequest, LoginResponse>[]
): Promise<string | Uint8Array> {
  switch (ctx.contentType) {
    case TwirpContentType.JSON:
      return handleAuthRegisterJSON<T>(ctx, service, data, interceptors);
    case TwirpContentType.Protobuf:
      return handleAuthRegisterProtobuf<T>(ctx, service, data, interceptors);
    default:
      const msg = "unexpected Content-Type";
      throw new TwirpError(TwirpErrorCode.BadRoute, msg);
  }
}

function handleAuthVerifyTokenRequest<T extends TwirpContext = TwirpContext>(
  ctx: T,
  service: AuthTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, VerifyTokenRequest, VerifyTokenResponse>[]
): Promise<string | Uint8Array> {
  switch (ctx.contentType) {
    case TwirpContentType.JSON:
      return handleAuthVerifyTokenJSON<T>(ctx, service, data, interceptors);
    case TwirpContentType.Protobuf:
      return handleAuthVerifyTokenProtobuf<T>(ctx, service, data, interceptors);
    default:
      const msg = "unexpected Content-Type";
      throw new TwirpError(TwirpErrorCode.BadRoute, msg);
  }
}
async function handleAuthLoginJSON<T extends TwirpContext = TwirpContext>(
  ctx: T,
  service: AuthTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, LoginRequest, LoginResponse>[]
) {
  let request: LoginRequest;
  let response: LoginResponse;

  try {
    const body = JSON.parse(data.toString() || "{}");
    request = LoginRequest.fromJson(body, { ignoreUnknownFields: true });
  } catch (e) {
    if (e instanceof Error) {
      const msg = "the json request could not be decoded";
      throw new TwirpError(TwirpErrorCode.Malformed, msg).withCause(e, true);
    }
  }

  if (interceptors && interceptors.length > 0) {
    const interceptor = chainInterceptors(...interceptors) as Interceptor<
      T,
      LoginRequest,
      LoginResponse
    >;
    response = await interceptor(ctx, request!, (ctx, inputReq) => {
      return service.Login(ctx, inputReq);
    });
  } else {
    response = await service.Login(ctx, request!);
  }

  return JSON.stringify(
    LoginResponse.toJson(response, {
      useProtoFieldName: true,
      emitDefaultValues: false,
    }) as string
  );
}

async function handleAuthRegisterJSON<T extends TwirpContext = TwirpContext>(
  ctx: T,
  service: AuthTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, RegisterRequest, LoginResponse>[]
) {
  let request: RegisterRequest;
  let response: LoginResponse;

  try {
    const body = JSON.parse(data.toString() || "{}");
    request = RegisterRequest.fromJson(body, { ignoreUnknownFields: true });
  } catch (e) {
    if (e instanceof Error) {
      const msg = "the json request could not be decoded";
      throw new TwirpError(TwirpErrorCode.Malformed, msg).withCause(e, true);
    }
  }

  if (interceptors && interceptors.length > 0) {
    const interceptor = chainInterceptors(...interceptors) as Interceptor<
      T,
      RegisterRequest,
      LoginResponse
    >;
    response = await interceptor(ctx, request!, (ctx, inputReq) => {
      return service.Register(ctx, inputReq);
    });
  } else {
    response = await service.Register(ctx, request!);
  }

  return JSON.stringify(
    LoginResponse.toJson(response, {
      useProtoFieldName: true,
      emitDefaultValues: false,
    }) as string
  );
}

async function handleAuthVerifyTokenJSON<T extends TwirpContext = TwirpContext>(
  ctx: T,
  service: AuthTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, VerifyTokenRequest, VerifyTokenResponse>[]
) {
  let request: VerifyTokenRequest;
  let response: VerifyTokenResponse;

  try {
    const body = JSON.parse(data.toString() || "{}");
    request = VerifyTokenRequest.fromJson(body, { ignoreUnknownFields: true });
  } catch (e) {
    if (e instanceof Error) {
      const msg = "the json request could not be decoded";
      throw new TwirpError(TwirpErrorCode.Malformed, msg).withCause(e, true);
    }
  }

  if (interceptors && interceptors.length > 0) {
    const interceptor = chainInterceptors(...interceptors) as Interceptor<
      T,
      VerifyTokenRequest,
      VerifyTokenResponse
    >;
    response = await interceptor(ctx, request!, (ctx, inputReq) => {
      return service.VerifyToken(ctx, inputReq);
    });
  } else {
    response = await service.VerifyToken(ctx, request!);
  }

  return JSON.stringify(
    VerifyTokenResponse.toJson(response, {
      useProtoFieldName: true,
      emitDefaultValues: false,
    }) as string
  );
}
async function handleAuthLoginProtobuf<T extends TwirpContext = TwirpContext>(
  ctx: T,
  service: AuthTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, LoginRequest, LoginResponse>[]
) {
  let request: LoginRequest;
  let response: LoginResponse;

  try {
    request = LoginRequest.fromBinary(data);
  } catch (e) {
    if (e instanceof Error) {
      const msg = "the protobuf request could not be decoded";
      throw new TwirpError(TwirpErrorCode.Malformed, msg).withCause(e, true);
    }
  }

  if (interceptors && interceptors.length > 0) {
    const interceptor = chainInterceptors(...interceptors) as Interceptor<
      T,
      LoginRequest,
      LoginResponse
    >;
    response = await interceptor(ctx, request!, (ctx, inputReq) => {
      return service.Login(ctx, inputReq);
    });
  } else {
    response = await service.Login(ctx, request!);
  }

  return Buffer.from(LoginResponse.toBinary(response));
}

async function handleAuthRegisterProtobuf<
  T extends TwirpContext = TwirpContext
>(
  ctx: T,
  service: AuthTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, RegisterRequest, LoginResponse>[]
) {
  let request: RegisterRequest;
  let response: LoginResponse;

  try {
    request = RegisterRequest.fromBinary(data);
  } catch (e) {
    if (e instanceof Error) {
      const msg = "the protobuf request could not be decoded";
      throw new TwirpError(TwirpErrorCode.Malformed, msg).withCause(e, true);
    }
  }

  if (interceptors && interceptors.length > 0) {
    const interceptor = chainInterceptors(...interceptors) as Interceptor<
      T,
      RegisterRequest,
      LoginResponse
    >;
    response = await interceptor(ctx, request!, (ctx, inputReq) => {
      return service.Register(ctx, inputReq);
    });
  } else {
    response = await service.Register(ctx, request!);
  }

  return Buffer.from(LoginResponse.toBinary(response));
}

async function handleAuthVerifyTokenProtobuf<
  T extends TwirpContext = TwirpContext
>(
  ctx: T,
  service: AuthTwirp,
  data: Buffer,
  interceptors?: Interceptor<T, VerifyTokenRequest, VerifyTokenResponse>[]
) {
  let request: VerifyTokenRequest;
  let response: VerifyTokenResponse;

  try {
    request = VerifyTokenRequest.fromBinary(data);
  } catch (e) {
    if (e instanceof Error) {
      const msg = "the protobuf request could not be decoded";
      throw new TwirpError(TwirpErrorCode.Malformed, msg).withCause(e, true);
    }
  }

  if (interceptors && interceptors.length > 0) {
    const interceptor = chainInterceptors(...interceptors) as Interceptor<
      T,
      VerifyTokenRequest,
      VerifyTokenResponse
    >;
    response = await interceptor(ctx, request!, (ctx, inputReq) => {
      return service.VerifyToken(ctx, inputReq);
    });
  } else {
    response = await service.VerifyToken(ctx, request!);
  }

  return Buffer.from(VerifyTokenResponse.toBinary(response));
}
