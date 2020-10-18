/* eslint-disable */
import { Writer, Reader } from 'protobufjs/minimal';


export interface AnyRenderer {
  loginPageRenderer: LoginPageRenderer | undefined;
  postPageRenderer: PostPageRenderer | undefined;
  userPageRenderer: UserPageRenderer | undefined;
  errorRenderer: ErrorRenderer | undefined;
}

export interface AnyRequest {
  userPageRequest: UserPageRequest | undefined;
  postPageRequest: PostPageRequest | undefined;
}

export interface ErrorRenderer {
  displayText: string;
}

export interface LoginPageRenderer {
  submitUrl: string;
  welcomeText: string;
}

export interface PostPageRenderer {
  id: string;
  text: string;
  userId: string;
  currentUserId: string;
}

export interface UserPageRenderer {
  id: string;
  lastPostId: string;
  currentUserId: string;
  name: string;
}

export interface LoginRequest {
  login: string;
  password: string;
}

export interface LoginResponse {
}

export interface ResolveRouteRequest {
  url: string;
}

export interface ResolveRouteResponse {
  js: string[];
  request: AnyRequest | undefined;
}

export interface PostPageRequest {
  postId: string;
}

export interface UserPageRequest {
  userId: string;
}

export interface Error {
  message: string;
}

const baseAnyRenderer: object = {
};

const baseAnyRequest: object = {
};

const baseErrorRenderer: object = {
  displayText: "",
};

const baseLoginPageRenderer: object = {
  submitUrl: "",
  welcomeText: "",
};

const basePostPageRenderer: object = {
  id: "",
  text: "",
  userId: "",
  currentUserId: "",
};

const baseUserPageRenderer: object = {
  id: "",
  lastPostId: "",
  currentUserId: "",
  name: "",
};

const baseLoginRequest: object = {
  login: "",
  password: "",
};

const baseLoginResponse: object = {
};

const baseResolveRouteRequest: object = {
  url: "",
};

const baseResolveRouteResponse: object = {
  js: "",
};

const basePostPageRequest: object = {
  postId: "",
};

const baseUserPageRequest: object = {
  userId: "",
};

const baseError: object = {
  message: "",
};

export const protobufPackage = ''

export const AnyRenderer = {
  encode(message: AnyRenderer, writer: Writer = Writer.create()): Writer {
    if (message.loginPageRenderer !== undefined) {
      LoginPageRenderer.encode(message.loginPageRenderer, writer.uint32(10).fork()).ldelim();
    }
    if (message.postPageRenderer !== undefined) {
      PostPageRenderer.encode(message.postPageRenderer, writer.uint32(18).fork()).ldelim();
    }
    if (message.userPageRenderer !== undefined) {
      UserPageRenderer.encode(message.userPageRenderer, writer.uint32(26).fork()).ldelim();
    }
    if (message.errorRenderer !== undefined) {
      ErrorRenderer.encode(message.errorRenderer, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },
  decode(input: Uint8Array | Reader, length?: number): AnyRenderer {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseAnyRenderer } as AnyRenderer;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.loginPageRenderer = LoginPageRenderer.decode(reader, reader.uint32());
          break;
        case 2:
          message.postPageRenderer = PostPageRenderer.decode(reader, reader.uint32());
          break;
        case 3:
          message.userPageRenderer = UserPageRenderer.decode(reader, reader.uint32());
          break;
        case 4:
          message.errorRenderer = ErrorRenderer.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): AnyRenderer {
    const message = { ...baseAnyRenderer } as AnyRenderer;
    if (object.loginPageRenderer !== undefined && object.loginPageRenderer !== null) {
      message.loginPageRenderer = LoginPageRenderer.fromJSON(object.loginPageRenderer);
    } else {
      message.loginPageRenderer = undefined;
    }
    if (object.postPageRenderer !== undefined && object.postPageRenderer !== null) {
      message.postPageRenderer = PostPageRenderer.fromJSON(object.postPageRenderer);
    } else {
      message.postPageRenderer = undefined;
    }
    if (object.userPageRenderer !== undefined && object.userPageRenderer !== null) {
      message.userPageRenderer = UserPageRenderer.fromJSON(object.userPageRenderer);
    } else {
      message.userPageRenderer = undefined;
    }
    if (object.errorRenderer !== undefined && object.errorRenderer !== null) {
      message.errorRenderer = ErrorRenderer.fromJSON(object.errorRenderer);
    } else {
      message.errorRenderer = undefined;
    }
    return message;
  },
  fromPartial(object: DeepPartial<AnyRenderer>): AnyRenderer {
    const message = { ...baseAnyRenderer } as AnyRenderer;
    if (object.loginPageRenderer !== undefined && object.loginPageRenderer !== null) {
      message.loginPageRenderer = LoginPageRenderer.fromPartial(object.loginPageRenderer);
    } else {
      message.loginPageRenderer = undefined;
    }
    if (object.postPageRenderer !== undefined && object.postPageRenderer !== null) {
      message.postPageRenderer = PostPageRenderer.fromPartial(object.postPageRenderer);
    } else {
      message.postPageRenderer = undefined;
    }
    if (object.userPageRenderer !== undefined && object.userPageRenderer !== null) {
      message.userPageRenderer = UserPageRenderer.fromPartial(object.userPageRenderer);
    } else {
      message.userPageRenderer = undefined;
    }
    if (object.errorRenderer !== undefined && object.errorRenderer !== null) {
      message.errorRenderer = ErrorRenderer.fromPartial(object.errorRenderer);
    } else {
      message.errorRenderer = undefined;
    }
    return message;
  },
  toJSON(message: AnyRenderer): unknown {
    const obj: any = {};
    message.loginPageRenderer !== undefined && (obj.loginPageRenderer = message.loginPageRenderer ? LoginPageRenderer.toJSON(message.loginPageRenderer) : undefined);
    message.postPageRenderer !== undefined && (obj.postPageRenderer = message.postPageRenderer ? PostPageRenderer.toJSON(message.postPageRenderer) : undefined);
    message.userPageRenderer !== undefined && (obj.userPageRenderer = message.userPageRenderer ? UserPageRenderer.toJSON(message.userPageRenderer) : undefined);
    message.errorRenderer !== undefined && (obj.errorRenderer = message.errorRenderer ? ErrorRenderer.toJSON(message.errorRenderer) : undefined);
    return obj;
  },
};

export const AnyRequest = {
  encode(message: AnyRequest, writer: Writer = Writer.create()): Writer {
    if (message.userPageRequest !== undefined) {
      UserPageRequest.encode(message.userPageRequest, writer.uint32(10).fork()).ldelim();
    }
    if (message.postPageRequest !== undefined) {
      PostPageRequest.encode(message.postPageRequest, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: Uint8Array | Reader, length?: number): AnyRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseAnyRequest } as AnyRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.userPageRequest = UserPageRequest.decode(reader, reader.uint32());
          break;
        case 2:
          message.postPageRequest = PostPageRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): AnyRequest {
    const message = { ...baseAnyRequest } as AnyRequest;
    if (object.userPageRequest !== undefined && object.userPageRequest !== null) {
      message.userPageRequest = UserPageRequest.fromJSON(object.userPageRequest);
    } else {
      message.userPageRequest = undefined;
    }
    if (object.postPageRequest !== undefined && object.postPageRequest !== null) {
      message.postPageRequest = PostPageRequest.fromJSON(object.postPageRequest);
    } else {
      message.postPageRequest = undefined;
    }
    return message;
  },
  fromPartial(object: DeepPartial<AnyRequest>): AnyRequest {
    const message = { ...baseAnyRequest } as AnyRequest;
    if (object.userPageRequest !== undefined && object.userPageRequest !== null) {
      message.userPageRequest = UserPageRequest.fromPartial(object.userPageRequest);
    } else {
      message.userPageRequest = undefined;
    }
    if (object.postPageRequest !== undefined && object.postPageRequest !== null) {
      message.postPageRequest = PostPageRequest.fromPartial(object.postPageRequest);
    } else {
      message.postPageRequest = undefined;
    }
    return message;
  },
  toJSON(message: AnyRequest): unknown {
    const obj: any = {};
    message.userPageRequest !== undefined && (obj.userPageRequest = message.userPageRequest ? UserPageRequest.toJSON(message.userPageRequest) : undefined);
    message.postPageRequest !== undefined && (obj.postPageRequest = message.postPageRequest ? PostPageRequest.toJSON(message.postPageRequest) : undefined);
    return obj;
  },
};

export const ErrorRenderer = {
  encode(message: ErrorRenderer, writer: Writer = Writer.create()): Writer {
    writer.uint32(10).string(message.displayText);
    return writer;
  },
  decode(input: Uint8Array | Reader, length?: number): ErrorRenderer {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseErrorRenderer } as ErrorRenderer;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.displayText = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): ErrorRenderer {
    const message = { ...baseErrorRenderer } as ErrorRenderer;
    if (object.displayText !== undefined && object.displayText !== null) {
      message.displayText = String(object.displayText);
    } else {
      message.displayText = "";
    }
    return message;
  },
  fromPartial(object: DeepPartial<ErrorRenderer>): ErrorRenderer {
    const message = { ...baseErrorRenderer } as ErrorRenderer;
    if (object.displayText !== undefined && object.displayText !== null) {
      message.displayText = object.displayText;
    } else {
      message.displayText = "";
    }
    return message;
  },
  toJSON(message: ErrorRenderer): unknown {
    const obj: any = {};
    message.displayText !== undefined && (obj.displayText = message.displayText);
    return obj;
  },
};

export const LoginPageRenderer = {
  encode(message: LoginPageRenderer, writer: Writer = Writer.create()): Writer {
    writer.uint32(10).string(message.submitUrl);
    writer.uint32(18).string(message.welcomeText);
    return writer;
  },
  decode(input: Uint8Array | Reader, length?: number): LoginPageRenderer {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseLoginPageRenderer } as LoginPageRenderer;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.submitUrl = reader.string();
          break;
        case 2:
          message.welcomeText = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): LoginPageRenderer {
    const message = { ...baseLoginPageRenderer } as LoginPageRenderer;
    if (object.submitUrl !== undefined && object.submitUrl !== null) {
      message.submitUrl = String(object.submitUrl);
    } else {
      message.submitUrl = "";
    }
    if (object.welcomeText !== undefined && object.welcomeText !== null) {
      message.welcomeText = String(object.welcomeText);
    } else {
      message.welcomeText = "";
    }
    return message;
  },
  fromPartial(object: DeepPartial<LoginPageRenderer>): LoginPageRenderer {
    const message = { ...baseLoginPageRenderer } as LoginPageRenderer;
    if (object.submitUrl !== undefined && object.submitUrl !== null) {
      message.submitUrl = object.submitUrl;
    } else {
      message.submitUrl = "";
    }
    if (object.welcomeText !== undefined && object.welcomeText !== null) {
      message.welcomeText = object.welcomeText;
    } else {
      message.welcomeText = "";
    }
    return message;
  },
  toJSON(message: LoginPageRenderer): unknown {
    const obj: any = {};
    message.submitUrl !== undefined && (obj.submitUrl = message.submitUrl);
    message.welcomeText !== undefined && (obj.welcomeText = message.welcomeText);
    return obj;
  },
};

export const PostPageRenderer = {
  encode(message: PostPageRenderer, writer: Writer = Writer.create()): Writer {
    writer.uint32(10).string(message.id);
    writer.uint32(18).string(message.text);
    writer.uint32(26).string(message.userId);
    writer.uint32(34).string(message.currentUserId);
    return writer;
  },
  decode(input: Uint8Array | Reader, length?: number): PostPageRenderer {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...basePostPageRenderer } as PostPageRenderer;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.string();
          break;
        case 2:
          message.text = reader.string();
          break;
        case 3:
          message.userId = reader.string();
          break;
        case 4:
          message.currentUserId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): PostPageRenderer {
    const message = { ...basePostPageRenderer } as PostPageRenderer;
    if (object.id !== undefined && object.id !== null) {
      message.id = String(object.id);
    } else {
      message.id = "";
    }
    if (object.text !== undefined && object.text !== null) {
      message.text = String(object.text);
    } else {
      message.text = "";
    }
    if (object.userId !== undefined && object.userId !== null) {
      message.userId = String(object.userId);
    } else {
      message.userId = "";
    }
    if (object.currentUserId !== undefined && object.currentUserId !== null) {
      message.currentUserId = String(object.currentUserId);
    } else {
      message.currentUserId = "";
    }
    return message;
  },
  fromPartial(object: DeepPartial<PostPageRenderer>): PostPageRenderer {
    const message = { ...basePostPageRenderer } as PostPageRenderer;
    if (object.id !== undefined && object.id !== null) {
      message.id = object.id;
    } else {
      message.id = "";
    }
    if (object.text !== undefined && object.text !== null) {
      message.text = object.text;
    } else {
      message.text = "";
    }
    if (object.userId !== undefined && object.userId !== null) {
      message.userId = object.userId;
    } else {
      message.userId = "";
    }
    if (object.currentUserId !== undefined && object.currentUserId !== null) {
      message.currentUserId = object.currentUserId;
    } else {
      message.currentUserId = "";
    }
    return message;
  },
  toJSON(message: PostPageRenderer): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.text !== undefined && (obj.text = message.text);
    message.userId !== undefined && (obj.userId = message.userId);
    message.currentUserId !== undefined && (obj.currentUserId = message.currentUserId);
    return obj;
  },
};

export const UserPageRenderer = {
  encode(message: UserPageRenderer, writer: Writer = Writer.create()): Writer {
    writer.uint32(10).string(message.id);
    writer.uint32(18).string(message.lastPostId);
    writer.uint32(26).string(message.currentUserId);
    writer.uint32(34).string(message.name);
    return writer;
  },
  decode(input: Uint8Array | Reader, length?: number): UserPageRenderer {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseUserPageRenderer } as UserPageRenderer;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.string();
          break;
        case 2:
          message.lastPostId = reader.string();
          break;
        case 3:
          message.currentUserId = reader.string();
          break;
        case 4:
          message.name = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): UserPageRenderer {
    const message = { ...baseUserPageRenderer } as UserPageRenderer;
    if (object.id !== undefined && object.id !== null) {
      message.id = String(object.id);
    } else {
      message.id = "";
    }
    if (object.lastPostId !== undefined && object.lastPostId !== null) {
      message.lastPostId = String(object.lastPostId);
    } else {
      message.lastPostId = "";
    }
    if (object.currentUserId !== undefined && object.currentUserId !== null) {
      message.currentUserId = String(object.currentUserId);
    } else {
      message.currentUserId = "";
    }
    if (object.name !== undefined && object.name !== null) {
      message.name = String(object.name);
    } else {
      message.name = "";
    }
    return message;
  },
  fromPartial(object: DeepPartial<UserPageRenderer>): UserPageRenderer {
    const message = { ...baseUserPageRenderer } as UserPageRenderer;
    if (object.id !== undefined && object.id !== null) {
      message.id = object.id;
    } else {
      message.id = "";
    }
    if (object.lastPostId !== undefined && object.lastPostId !== null) {
      message.lastPostId = object.lastPostId;
    } else {
      message.lastPostId = "";
    }
    if (object.currentUserId !== undefined && object.currentUserId !== null) {
      message.currentUserId = object.currentUserId;
    } else {
      message.currentUserId = "";
    }
    if (object.name !== undefined && object.name !== null) {
      message.name = object.name;
    } else {
      message.name = "";
    }
    return message;
  },
  toJSON(message: UserPageRenderer): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.lastPostId !== undefined && (obj.lastPostId = message.lastPostId);
    message.currentUserId !== undefined && (obj.currentUserId = message.currentUserId);
    message.name !== undefined && (obj.name = message.name);
    return obj;
  },
};

export const LoginRequest = {
  encode(message: LoginRequest, writer: Writer = Writer.create()): Writer {
    writer.uint32(10).string(message.login);
    writer.uint32(18).string(message.password);
    return writer;
  },
  decode(input: Uint8Array | Reader, length?: number): LoginRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseLoginRequest } as LoginRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.login = reader.string();
          break;
        case 2:
          message.password = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): LoginRequest {
    const message = { ...baseLoginRequest } as LoginRequest;
    if (object.login !== undefined && object.login !== null) {
      message.login = String(object.login);
    } else {
      message.login = "";
    }
    if (object.password !== undefined && object.password !== null) {
      message.password = String(object.password);
    } else {
      message.password = "";
    }
    return message;
  },
  fromPartial(object: DeepPartial<LoginRequest>): LoginRequest {
    const message = { ...baseLoginRequest } as LoginRequest;
    if (object.login !== undefined && object.login !== null) {
      message.login = object.login;
    } else {
      message.login = "";
    }
    if (object.password !== undefined && object.password !== null) {
      message.password = object.password;
    } else {
      message.password = "";
    }
    return message;
  },
  toJSON(message: LoginRequest): unknown {
    const obj: any = {};
    message.login !== undefined && (obj.login = message.login);
    message.password !== undefined && (obj.password = message.password);
    return obj;
  },
};

export const LoginResponse = {
  encode(_: LoginResponse, writer: Writer = Writer.create()): Writer {
    return writer;
  },
  decode(input: Uint8Array | Reader, length?: number): LoginResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseLoginResponse } as LoginResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(_: any): LoginResponse {
    const message = { ...baseLoginResponse } as LoginResponse;
    return message;
  },
  fromPartial(_: DeepPartial<LoginResponse>): LoginResponse {
    const message = { ...baseLoginResponse } as LoginResponse;
    return message;
  },
  toJSON(_: LoginResponse): unknown {
    const obj: any = {};
    return obj;
  },
};

export const ResolveRouteRequest = {
  encode(message: ResolveRouteRequest, writer: Writer = Writer.create()): Writer {
    writer.uint32(10).string(message.url);
    return writer;
  },
  decode(input: Uint8Array | Reader, length?: number): ResolveRouteRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseResolveRouteRequest } as ResolveRouteRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.url = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): ResolveRouteRequest {
    const message = { ...baseResolveRouteRequest } as ResolveRouteRequest;
    if (object.url !== undefined && object.url !== null) {
      message.url = String(object.url);
    } else {
      message.url = "";
    }
    return message;
  },
  fromPartial(object: DeepPartial<ResolveRouteRequest>): ResolveRouteRequest {
    const message = { ...baseResolveRouteRequest } as ResolveRouteRequest;
    if (object.url !== undefined && object.url !== null) {
      message.url = object.url;
    } else {
      message.url = "";
    }
    return message;
  },
  toJSON(message: ResolveRouteRequest): unknown {
    const obj: any = {};
    message.url !== undefined && (obj.url = message.url);
    return obj;
  },
};

export const ResolveRouteResponse = {
  encode(message: ResolveRouteResponse, writer: Writer = Writer.create()): Writer {
    for (const v of message.js) {
      writer.uint32(10).string(v!);
    }
    if (message.request !== undefined && message.request !== undefined) {
      AnyRequest.encode(message.request, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: Uint8Array | Reader, length?: number): ResolveRouteResponse {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseResolveRouteResponse } as ResolveRouteResponse;
    message.js = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.js.push(reader.string());
          break;
        case 2:
          message.request = AnyRequest.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): ResolveRouteResponse {
    const message = { ...baseResolveRouteResponse } as ResolveRouteResponse;
    message.js = [];
    if (object.js !== undefined && object.js !== null) {
      for (const e of object.js) {
        message.js.push(String(e));
      }
    }
    if (object.request !== undefined && object.request !== null) {
      message.request = AnyRequest.fromJSON(object.request);
    } else {
      message.request = undefined;
    }
    return message;
  },
  fromPartial(object: DeepPartial<ResolveRouteResponse>): ResolveRouteResponse {
    const message = { ...baseResolveRouteResponse } as ResolveRouteResponse;
    message.js = [];
    if (object.js !== undefined && object.js !== null) {
      for (const e of object.js) {
        message.js.push(e);
      }
    }
    if (object.request !== undefined && object.request !== null) {
      message.request = AnyRequest.fromPartial(object.request);
    } else {
      message.request = undefined;
    }
    return message;
  },
  toJSON(message: ResolveRouteResponse): unknown {
    const obj: any = {};
    if (message.js) {
      obj.js = message.js.map(e => e);
    } else {
      obj.js = [];
    }
    message.request !== undefined && (obj.request = message.request ? AnyRequest.toJSON(message.request) : undefined);
    return obj;
  },
};

export const PostPageRequest = {
  encode(message: PostPageRequest, writer: Writer = Writer.create()): Writer {
    writer.uint32(10).string(message.postId);
    return writer;
  },
  decode(input: Uint8Array | Reader, length?: number): PostPageRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...basePostPageRequest } as PostPageRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.postId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): PostPageRequest {
    const message = { ...basePostPageRequest } as PostPageRequest;
    if (object.postId !== undefined && object.postId !== null) {
      message.postId = String(object.postId);
    } else {
      message.postId = "";
    }
    return message;
  },
  fromPartial(object: DeepPartial<PostPageRequest>): PostPageRequest {
    const message = { ...basePostPageRequest } as PostPageRequest;
    if (object.postId !== undefined && object.postId !== null) {
      message.postId = object.postId;
    } else {
      message.postId = "";
    }
    return message;
  },
  toJSON(message: PostPageRequest): unknown {
    const obj: any = {};
    message.postId !== undefined && (obj.postId = message.postId);
    return obj;
  },
};

export const UserPageRequest = {
  encode(message: UserPageRequest, writer: Writer = Writer.create()): Writer {
    writer.uint32(10).string(message.userId);
    return writer;
  },
  decode(input: Uint8Array | Reader, length?: number): UserPageRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseUserPageRequest } as UserPageRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.userId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): UserPageRequest {
    const message = { ...baseUserPageRequest } as UserPageRequest;
    if (object.userId !== undefined && object.userId !== null) {
      message.userId = String(object.userId);
    } else {
      message.userId = "";
    }
    return message;
  },
  fromPartial(object: DeepPartial<UserPageRequest>): UserPageRequest {
    const message = { ...baseUserPageRequest } as UserPageRequest;
    if (object.userId !== undefined && object.userId !== null) {
      message.userId = object.userId;
    } else {
      message.userId = "";
    }
    return message;
  },
  toJSON(message: UserPageRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    return obj;
  },
};

export const Error = {
  encode(message: Error, writer: Writer = Writer.create()): Writer {
    writer.uint32(10).string(message.message);
    return writer;
  },
  decode(input: Uint8Array | Reader, length?: number): Error {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseError } as Error;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.message = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): Error {
    const message = { ...baseError } as Error;
    if (object.message !== undefined && object.message !== null) {
      message.message = String(object.message);
    } else {
      message.message = "";
    }
    return message;
  },
  fromPartial(object: DeepPartial<Error>): Error {
    const message = { ...baseError } as Error;
    if (object.message !== undefined && object.message !== null) {
      message.message = object.message;
    } else {
      message.message = "";
    }
    return message;
  },
  toJSON(message: Error): unknown {
    const obj: any = {};
    message.message !== undefined && (obj.message = message.message);
    return obj;
  },
};

type Builtin = Date | Function | Uint8Array | string | number | undefined;
type DeepPartial<T> = T extends Builtin
  ? T
  : T extends Array<infer U>
  ? Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U>
  ? ReadonlyArray<DeepPartial<U>>
  : T extends {}
  ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;