/* eslint-disable */
import { Reader, Writer } from 'protobufjs/minimal';


export interface ErrorRenderer {
  errorCode: string;
  displayText: string;
}

export interface LoginRequest {
  login: string;
  password: string;
}

export interface LoginRenderer {
  headerRenderer: HeaderRenderer | undefined;
}

export interface LogoutRequest {
}

export interface LogoutRenderer {
}

export interface LoginPageRequest {
}

export interface LoginPageRenderer {
  submitUrl: string;
  welcomeText: string;
  vkUrl: string;
  headerRenderer: HeaderRenderer | undefined;
}

export interface PostPageRenderer {
  id: string;
  text: string;
  userId: string;
  userUrl: string;
  currentUserId: string;
  postUrl: string;
  headerRenderer: HeaderRenderer | undefined;
}

export interface UserPageRenderer {
  id: string;
  lastPostId: string;
  lastPostUrl: string;
  name: string;
  headerRenderer: HeaderRenderer | undefined;
}

export interface ResolveRouteRequest {
  url: string;
}

export interface ResolveRouteResponse {
  js: string[];
  rootComponent: string;
  apiMethod: string;
  apiRequest: string;
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

export interface AddPostRequest {
  text: string;
}

export interface AddPostRenderer {
  id: string;
  text: string;
}

export interface GetFeedRequest {
}

export interface GetFeedRenderer {
  posts: PostPageRenderer[];
  headerRenderer: HeaderRenderer | undefined;
}

export interface ComposerRequest {
}

export interface ComposerRenderer {
  welcomeText: string;
  headerRenderer: HeaderRenderer | undefined;
  unathorizedText: string;
}

export interface IndexRequest {
}

export interface IndexRenderer {
  text: string;
  feedUrl: string;
  headerRenderer: HeaderRenderer | undefined;
}

export interface HeaderRenderer {
  currentUserId: string;
  currentUserName: string;
  links: HeaderRenderer_Link[];
}

export interface HeaderRenderer_Link {
  url: string;
  label: string;
}

export interface VKCallbackRequest {
  vkCode: string;
}

export interface VKCallbackRenderer {
  headerRenderer: HeaderRenderer | undefined;
}

const baseErrorRenderer: object = {
  errorCode: "",
  displayText: "",
};

const baseLoginRequest: object = {
  login: "",
  password: "",
};

const baseLoginRenderer: object = {
};

const baseLogoutRequest: object = {
};

const baseLogoutRenderer: object = {
};

const baseLoginPageRequest: object = {
};

const baseLoginPageRenderer: object = {
  submitUrl: "",
  welcomeText: "",
  vkUrl: "",
};

const basePostPageRenderer: object = {
  id: "",
  text: "",
  userId: "",
  userUrl: "",
  currentUserId: "",
  postUrl: "",
};

const baseUserPageRenderer: object = {
  id: "",
  lastPostId: "",
  lastPostUrl: "",
  name: "",
};

const baseResolveRouteRequest: object = {
  url: "",
};

const baseResolveRouteResponse: object = {
  js: "",
  rootComponent: "",
  apiMethod: "",
  apiRequest: "",
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

const baseAddPostRequest: object = {
  text: "",
};

const baseAddPostRenderer: object = {
  id: "",
  text: "",
};

const baseGetFeedRequest: object = {
};

const baseGetFeedRenderer: object = {
};

const baseComposerRequest: object = {
};

const baseComposerRenderer: object = {
  welcomeText: "",
  unathorizedText: "",
};

const baseIndexRequest: object = {
};

const baseIndexRenderer: object = {
  text: "",
  feedUrl: "",
};

const baseHeaderRenderer: object = {
  currentUserId: "",
  currentUserName: "",
};

const baseHeaderRenderer_Link: object = {
  url: "",
  label: "",
};

const baseVKCallbackRequest: object = {
  vkCode: "",
};

const baseVKCallbackRenderer: object = {
};

export interface API {

  LoginPage(request: LoginPageRequest): Promise<LoginPageRenderer>;

  PostPage(request: PostPageRequest): Promise<PostPageRenderer>;

  UserPage(request: UserPageRequest): Promise<UserPageRenderer>;

  Login(request: LoginRequest): Promise<LoginRenderer>;

  AddPost(request: AddPostRequest): Promise<AddPostRenderer>;

  GetFeed(request: GetFeedRequest): Promise<GetFeedRenderer>;

  Composer(request: ComposerRequest): Promise<ComposerRenderer>;

  Index(request: IndexRequest): Promise<IndexRenderer>;

  Logout(request: LogoutRequest): Promise<LogoutRenderer>;

  VKCallback(request: VKCallbackRenderer): Promise<VKCallbackRenderer>;

}

export class APIClientImpl implements API {

  private readonly rpc: Rpc;

  constructor(rpc: Rpc) {
    this.rpc = rpc;
  }

  LoginPage(request: LoginPageRequest): Promise<LoginPageRenderer> {
    const data = LoginPageRequest.encode(request).finish();
    const promise = this.rpc.request("meme.API", "LoginPage", data);
    return promise.then(data => LoginPageRenderer.decode(new Reader(data)));
  }

  PostPage(request: PostPageRequest): Promise<PostPageRenderer> {
    const data = PostPageRequest.encode(request).finish();
    const promise = this.rpc.request("meme.API", "PostPage", data);
    return promise.then(data => PostPageRenderer.decode(new Reader(data)));
  }

  UserPage(request: UserPageRequest): Promise<UserPageRenderer> {
    const data = UserPageRequest.encode(request).finish();
    const promise = this.rpc.request("meme.API", "UserPage", data);
    return promise.then(data => UserPageRenderer.decode(new Reader(data)));
  }

  Login(request: LoginRequest): Promise<LoginRenderer> {
    const data = LoginRequest.encode(request).finish();
    const promise = this.rpc.request("meme.API", "Login", data);
    return promise.then(data => LoginRenderer.decode(new Reader(data)));
  }

  AddPost(request: AddPostRequest): Promise<AddPostRenderer> {
    const data = AddPostRequest.encode(request).finish();
    const promise = this.rpc.request("meme.API", "AddPost", data);
    return promise.then(data => AddPostRenderer.decode(new Reader(data)));
  }

  GetFeed(request: GetFeedRequest): Promise<GetFeedRenderer> {
    const data = GetFeedRequest.encode(request).finish();
    const promise = this.rpc.request("meme.API", "GetFeed", data);
    return promise.then(data => GetFeedRenderer.decode(new Reader(data)));
  }

  Composer(request: ComposerRequest): Promise<ComposerRenderer> {
    const data = ComposerRequest.encode(request).finish();
    const promise = this.rpc.request("meme.API", "Composer", data);
    return promise.then(data => ComposerRenderer.decode(new Reader(data)));
  }

  Index(request: IndexRequest): Promise<IndexRenderer> {
    const data = IndexRequest.encode(request).finish();
    const promise = this.rpc.request("meme.API", "Index", data);
    return promise.then(data => IndexRenderer.decode(new Reader(data)));
  }

  Logout(request: LogoutRequest): Promise<LogoutRenderer> {
    const data = LogoutRequest.encode(request).finish();
    const promise = this.rpc.request("meme.API", "Logout", data);
    return promise.then(data => LogoutRenderer.decode(new Reader(data)));
  }

  VKCallback(request: VKCallbackRenderer): Promise<VKCallbackRenderer> {
    const data = VKCallbackRenderer.encode(request).finish();
    const promise = this.rpc.request("meme.API", "VKCallback", data);
    return promise.then(data => VKCallbackRenderer.decode(new Reader(data)));
  }

}

interface Rpc {

  request(service: string, method: string, data: Uint8Array): Promise<Uint8Array>;

}

export const protobufPackage = 'meme'

export const ErrorRenderer = {
  encode(message: ErrorRenderer, writer: Writer = Writer.create()): Writer {
    writer.uint32(10).string(message.errorCode);
    writer.uint32(18).string(message.displayText);
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
          message.errorCode = reader.string();
          break;
        case 2:
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
    if (object.errorCode !== undefined && object.errorCode !== null) {
      message.errorCode = String(object.errorCode);
    } else {
      message.errorCode = "";
    }
    if (object.displayText !== undefined && object.displayText !== null) {
      message.displayText = String(object.displayText);
    } else {
      message.displayText = "";
    }
    return message;
  },
  fromPartial(object: DeepPartial<ErrorRenderer>): ErrorRenderer {
    const message = { ...baseErrorRenderer } as ErrorRenderer;
    if (object.errorCode !== undefined && object.errorCode !== null) {
      message.errorCode = object.errorCode;
    } else {
      message.errorCode = "";
    }
    if (object.displayText !== undefined && object.displayText !== null) {
      message.displayText = object.displayText;
    } else {
      message.displayText = "";
    }
    return message;
  },
  toJSON(message: ErrorRenderer): unknown {
    const obj: any = {};
    message.errorCode !== undefined && (obj.errorCode = message.errorCode);
    message.displayText !== undefined && (obj.displayText = message.displayText);
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

export const LoginRenderer = {
  encode(message: LoginRenderer, writer: Writer = Writer.create()): Writer {
    if (message.headerRenderer !== undefined && message.headerRenderer !== undefined) {
      HeaderRenderer.encode(message.headerRenderer, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: Uint8Array | Reader, length?: number): LoginRenderer {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseLoginRenderer } as LoginRenderer;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.headerRenderer = HeaderRenderer.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): LoginRenderer {
    const message = { ...baseLoginRenderer } as LoginRenderer;
    if (object.headerRenderer !== undefined && object.headerRenderer !== null) {
      message.headerRenderer = HeaderRenderer.fromJSON(object.headerRenderer);
    } else {
      message.headerRenderer = undefined;
    }
    return message;
  },
  fromPartial(object: DeepPartial<LoginRenderer>): LoginRenderer {
    const message = { ...baseLoginRenderer } as LoginRenderer;
    if (object.headerRenderer !== undefined && object.headerRenderer !== null) {
      message.headerRenderer = HeaderRenderer.fromPartial(object.headerRenderer);
    } else {
      message.headerRenderer = undefined;
    }
    return message;
  },
  toJSON(message: LoginRenderer): unknown {
    const obj: any = {};
    message.headerRenderer !== undefined && (obj.headerRenderer = message.headerRenderer ? HeaderRenderer.toJSON(message.headerRenderer) : undefined);
    return obj;
  },
};

export const LogoutRequest = {
  encode(_: LogoutRequest, writer: Writer = Writer.create()): Writer {
    return writer;
  },
  decode(input: Uint8Array | Reader, length?: number): LogoutRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseLogoutRequest } as LogoutRequest;
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
  fromJSON(_: any): LogoutRequest {
    const message = { ...baseLogoutRequest } as LogoutRequest;
    return message;
  },
  fromPartial(_: DeepPartial<LogoutRequest>): LogoutRequest {
    const message = { ...baseLogoutRequest } as LogoutRequest;
    return message;
  },
  toJSON(_: LogoutRequest): unknown {
    const obj: any = {};
    return obj;
  },
};

export const LogoutRenderer = {
  encode(_: LogoutRenderer, writer: Writer = Writer.create()): Writer {
    return writer;
  },
  decode(input: Uint8Array | Reader, length?: number): LogoutRenderer {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseLogoutRenderer } as LogoutRenderer;
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
  fromJSON(_: any): LogoutRenderer {
    const message = { ...baseLogoutRenderer } as LogoutRenderer;
    return message;
  },
  fromPartial(_: DeepPartial<LogoutRenderer>): LogoutRenderer {
    const message = { ...baseLogoutRenderer } as LogoutRenderer;
    return message;
  },
  toJSON(_: LogoutRenderer): unknown {
    const obj: any = {};
    return obj;
  },
};

export const LoginPageRequest = {
  encode(_: LoginPageRequest, writer: Writer = Writer.create()): Writer {
    return writer;
  },
  decode(input: Uint8Array | Reader, length?: number): LoginPageRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseLoginPageRequest } as LoginPageRequest;
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
  fromJSON(_: any): LoginPageRequest {
    const message = { ...baseLoginPageRequest } as LoginPageRequest;
    return message;
  },
  fromPartial(_: DeepPartial<LoginPageRequest>): LoginPageRequest {
    const message = { ...baseLoginPageRequest } as LoginPageRequest;
    return message;
  },
  toJSON(_: LoginPageRequest): unknown {
    const obj: any = {};
    return obj;
  },
};

export const LoginPageRenderer = {
  encode(message: LoginPageRenderer, writer: Writer = Writer.create()): Writer {
    writer.uint32(10).string(message.submitUrl);
    writer.uint32(18).string(message.welcomeText);
    writer.uint32(34).string(message.vkUrl);
    if (message.headerRenderer !== undefined && message.headerRenderer !== undefined) {
      HeaderRenderer.encode(message.headerRenderer, writer.uint32(26).fork()).ldelim();
    }
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
        case 4:
          message.vkUrl = reader.string();
          break;
        case 3:
          message.headerRenderer = HeaderRenderer.decode(reader, reader.uint32());
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
    if (object.vkUrl !== undefined && object.vkUrl !== null) {
      message.vkUrl = String(object.vkUrl);
    } else {
      message.vkUrl = "";
    }
    if (object.headerRenderer !== undefined && object.headerRenderer !== null) {
      message.headerRenderer = HeaderRenderer.fromJSON(object.headerRenderer);
    } else {
      message.headerRenderer = undefined;
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
    if (object.vkUrl !== undefined && object.vkUrl !== null) {
      message.vkUrl = object.vkUrl;
    } else {
      message.vkUrl = "";
    }
    if (object.headerRenderer !== undefined && object.headerRenderer !== null) {
      message.headerRenderer = HeaderRenderer.fromPartial(object.headerRenderer);
    } else {
      message.headerRenderer = undefined;
    }
    return message;
  },
  toJSON(message: LoginPageRenderer): unknown {
    const obj: any = {};
    message.submitUrl !== undefined && (obj.submitUrl = message.submitUrl);
    message.welcomeText !== undefined && (obj.welcomeText = message.welcomeText);
    message.vkUrl !== undefined && (obj.vkUrl = message.vkUrl);
    message.headerRenderer !== undefined && (obj.headerRenderer = message.headerRenderer ? HeaderRenderer.toJSON(message.headerRenderer) : undefined);
    return obj;
  },
};

export const PostPageRenderer = {
  encode(message: PostPageRenderer, writer: Writer = Writer.create()): Writer {
    writer.uint32(10).string(message.id);
    writer.uint32(18).string(message.text);
    writer.uint32(26).string(message.userId);
    writer.uint32(50).string(message.userUrl);
    writer.uint32(34).string(message.currentUserId);
    writer.uint32(58).string(message.postUrl);
    if (message.headerRenderer !== undefined && message.headerRenderer !== undefined) {
      HeaderRenderer.encode(message.headerRenderer, writer.uint32(42).fork()).ldelim();
    }
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
        case 6:
          message.userUrl = reader.string();
          break;
        case 4:
          message.currentUserId = reader.string();
          break;
        case 7:
          message.postUrl = reader.string();
          break;
        case 5:
          message.headerRenderer = HeaderRenderer.decode(reader, reader.uint32());
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
    if (object.userUrl !== undefined && object.userUrl !== null) {
      message.userUrl = String(object.userUrl);
    } else {
      message.userUrl = "";
    }
    if (object.currentUserId !== undefined && object.currentUserId !== null) {
      message.currentUserId = String(object.currentUserId);
    } else {
      message.currentUserId = "";
    }
    if (object.postUrl !== undefined && object.postUrl !== null) {
      message.postUrl = String(object.postUrl);
    } else {
      message.postUrl = "";
    }
    if (object.headerRenderer !== undefined && object.headerRenderer !== null) {
      message.headerRenderer = HeaderRenderer.fromJSON(object.headerRenderer);
    } else {
      message.headerRenderer = undefined;
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
    if (object.userUrl !== undefined && object.userUrl !== null) {
      message.userUrl = object.userUrl;
    } else {
      message.userUrl = "";
    }
    if (object.currentUserId !== undefined && object.currentUserId !== null) {
      message.currentUserId = object.currentUserId;
    } else {
      message.currentUserId = "";
    }
    if (object.postUrl !== undefined && object.postUrl !== null) {
      message.postUrl = object.postUrl;
    } else {
      message.postUrl = "";
    }
    if (object.headerRenderer !== undefined && object.headerRenderer !== null) {
      message.headerRenderer = HeaderRenderer.fromPartial(object.headerRenderer);
    } else {
      message.headerRenderer = undefined;
    }
    return message;
  },
  toJSON(message: PostPageRenderer): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.text !== undefined && (obj.text = message.text);
    message.userId !== undefined && (obj.userId = message.userId);
    message.userUrl !== undefined && (obj.userUrl = message.userUrl);
    message.currentUserId !== undefined && (obj.currentUserId = message.currentUserId);
    message.postUrl !== undefined && (obj.postUrl = message.postUrl);
    message.headerRenderer !== undefined && (obj.headerRenderer = message.headerRenderer ? HeaderRenderer.toJSON(message.headerRenderer) : undefined);
    return obj;
  },
};

export const UserPageRenderer = {
  encode(message: UserPageRenderer, writer: Writer = Writer.create()): Writer {
    writer.uint32(10).string(message.id);
    writer.uint32(18).string(message.lastPostId);
    writer.uint32(50).string(message.lastPostUrl);
    writer.uint32(34).string(message.name);
    if (message.headerRenderer !== undefined && message.headerRenderer !== undefined) {
      HeaderRenderer.encode(message.headerRenderer, writer.uint32(42).fork()).ldelim();
    }
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
        case 6:
          message.lastPostUrl = reader.string();
          break;
        case 4:
          message.name = reader.string();
          break;
        case 5:
          message.headerRenderer = HeaderRenderer.decode(reader, reader.uint32());
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
    if (object.lastPostUrl !== undefined && object.lastPostUrl !== null) {
      message.lastPostUrl = String(object.lastPostUrl);
    } else {
      message.lastPostUrl = "";
    }
    if (object.name !== undefined && object.name !== null) {
      message.name = String(object.name);
    } else {
      message.name = "";
    }
    if (object.headerRenderer !== undefined && object.headerRenderer !== null) {
      message.headerRenderer = HeaderRenderer.fromJSON(object.headerRenderer);
    } else {
      message.headerRenderer = undefined;
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
    if (object.lastPostUrl !== undefined && object.lastPostUrl !== null) {
      message.lastPostUrl = object.lastPostUrl;
    } else {
      message.lastPostUrl = "";
    }
    if (object.name !== undefined && object.name !== null) {
      message.name = object.name;
    } else {
      message.name = "";
    }
    if (object.headerRenderer !== undefined && object.headerRenderer !== null) {
      message.headerRenderer = HeaderRenderer.fromPartial(object.headerRenderer);
    } else {
      message.headerRenderer = undefined;
    }
    return message;
  },
  toJSON(message: UserPageRenderer): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.lastPostId !== undefined && (obj.lastPostId = message.lastPostId);
    message.lastPostUrl !== undefined && (obj.lastPostUrl = message.lastPostUrl);
    message.name !== undefined && (obj.name = message.name);
    message.headerRenderer !== undefined && (obj.headerRenderer = message.headerRenderer ? HeaderRenderer.toJSON(message.headerRenderer) : undefined);
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
    writer.uint32(26).string(message.rootComponent);
    writer.uint32(34).string(message.apiMethod);
    writer.uint32(42).string(message.apiRequest);
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
        case 3:
          message.rootComponent = reader.string();
          break;
        case 4:
          message.apiMethod = reader.string();
          break;
        case 5:
          message.apiRequest = reader.string();
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
    if (object.rootComponent !== undefined && object.rootComponent !== null) {
      message.rootComponent = String(object.rootComponent);
    } else {
      message.rootComponent = "";
    }
    if (object.apiMethod !== undefined && object.apiMethod !== null) {
      message.apiMethod = String(object.apiMethod);
    } else {
      message.apiMethod = "";
    }
    if (object.apiRequest !== undefined && object.apiRequest !== null) {
      message.apiRequest = String(object.apiRequest);
    } else {
      message.apiRequest = "";
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
    if (object.rootComponent !== undefined && object.rootComponent !== null) {
      message.rootComponent = object.rootComponent;
    } else {
      message.rootComponent = "";
    }
    if (object.apiMethod !== undefined && object.apiMethod !== null) {
      message.apiMethod = object.apiMethod;
    } else {
      message.apiMethod = "";
    }
    if (object.apiRequest !== undefined && object.apiRequest !== null) {
      message.apiRequest = object.apiRequest;
    } else {
      message.apiRequest = "";
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
    message.rootComponent !== undefined && (obj.rootComponent = message.rootComponent);
    message.apiMethod !== undefined && (obj.apiMethod = message.apiMethod);
    message.apiRequest !== undefined && (obj.apiRequest = message.apiRequest);
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

export const AddPostRequest = {
  encode(message: AddPostRequest, writer: Writer = Writer.create()): Writer {
    writer.uint32(10).string(message.text);
    return writer;
  },
  decode(input: Uint8Array | Reader, length?: number): AddPostRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseAddPostRequest } as AddPostRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.text = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): AddPostRequest {
    const message = { ...baseAddPostRequest } as AddPostRequest;
    if (object.text !== undefined && object.text !== null) {
      message.text = String(object.text);
    } else {
      message.text = "";
    }
    return message;
  },
  fromPartial(object: DeepPartial<AddPostRequest>): AddPostRequest {
    const message = { ...baseAddPostRequest } as AddPostRequest;
    if (object.text !== undefined && object.text !== null) {
      message.text = object.text;
    } else {
      message.text = "";
    }
    return message;
  },
  toJSON(message: AddPostRequest): unknown {
    const obj: any = {};
    message.text !== undefined && (obj.text = message.text);
    return obj;
  },
};

export const AddPostRenderer = {
  encode(message: AddPostRenderer, writer: Writer = Writer.create()): Writer {
    writer.uint32(10).string(message.id);
    writer.uint32(18).string(message.text);
    return writer;
  },
  decode(input: Uint8Array | Reader, length?: number): AddPostRenderer {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseAddPostRenderer } as AddPostRenderer;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.string();
          break;
        case 2:
          message.text = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): AddPostRenderer {
    const message = { ...baseAddPostRenderer } as AddPostRenderer;
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
    return message;
  },
  fromPartial(object: DeepPartial<AddPostRenderer>): AddPostRenderer {
    const message = { ...baseAddPostRenderer } as AddPostRenderer;
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
    return message;
  },
  toJSON(message: AddPostRenderer): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.text !== undefined && (obj.text = message.text);
    return obj;
  },
};

export const GetFeedRequest = {
  encode(_: GetFeedRequest, writer: Writer = Writer.create()): Writer {
    return writer;
  },
  decode(input: Uint8Array | Reader, length?: number): GetFeedRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseGetFeedRequest } as GetFeedRequest;
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
  fromJSON(_: any): GetFeedRequest {
    const message = { ...baseGetFeedRequest } as GetFeedRequest;
    return message;
  },
  fromPartial(_: DeepPartial<GetFeedRequest>): GetFeedRequest {
    const message = { ...baseGetFeedRequest } as GetFeedRequest;
    return message;
  },
  toJSON(_: GetFeedRequest): unknown {
    const obj: any = {};
    return obj;
  },
};

export const GetFeedRenderer = {
  encode(message: GetFeedRenderer, writer: Writer = Writer.create()): Writer {
    for (const v of message.posts) {
      PostPageRenderer.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.headerRenderer !== undefined && message.headerRenderer !== undefined) {
      HeaderRenderer.encode(message.headerRenderer, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: Uint8Array | Reader, length?: number): GetFeedRenderer {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseGetFeedRenderer } as GetFeedRenderer;
    message.posts = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.posts.push(PostPageRenderer.decode(reader, reader.uint32()));
          break;
        case 2:
          message.headerRenderer = HeaderRenderer.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): GetFeedRenderer {
    const message = { ...baseGetFeedRenderer } as GetFeedRenderer;
    message.posts = [];
    if (object.posts !== undefined && object.posts !== null) {
      for (const e of object.posts) {
        message.posts.push(PostPageRenderer.fromJSON(e));
      }
    }
    if (object.headerRenderer !== undefined && object.headerRenderer !== null) {
      message.headerRenderer = HeaderRenderer.fromJSON(object.headerRenderer);
    } else {
      message.headerRenderer = undefined;
    }
    return message;
  },
  fromPartial(object: DeepPartial<GetFeedRenderer>): GetFeedRenderer {
    const message = { ...baseGetFeedRenderer } as GetFeedRenderer;
    message.posts = [];
    if (object.posts !== undefined && object.posts !== null) {
      for (const e of object.posts) {
        message.posts.push(PostPageRenderer.fromPartial(e));
      }
    }
    if (object.headerRenderer !== undefined && object.headerRenderer !== null) {
      message.headerRenderer = HeaderRenderer.fromPartial(object.headerRenderer);
    } else {
      message.headerRenderer = undefined;
    }
    return message;
  },
  toJSON(message: GetFeedRenderer): unknown {
    const obj: any = {};
    if (message.posts) {
      obj.posts = message.posts.map(e => e ? PostPageRenderer.toJSON(e) : undefined);
    } else {
      obj.posts = [];
    }
    message.headerRenderer !== undefined && (obj.headerRenderer = message.headerRenderer ? HeaderRenderer.toJSON(message.headerRenderer) : undefined);
    return obj;
  },
};

export const ComposerRequest = {
  encode(_: ComposerRequest, writer: Writer = Writer.create()): Writer {
    return writer;
  },
  decode(input: Uint8Array | Reader, length?: number): ComposerRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseComposerRequest } as ComposerRequest;
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
  fromJSON(_: any): ComposerRequest {
    const message = { ...baseComposerRequest } as ComposerRequest;
    return message;
  },
  fromPartial(_: DeepPartial<ComposerRequest>): ComposerRequest {
    const message = { ...baseComposerRequest } as ComposerRequest;
    return message;
  },
  toJSON(_: ComposerRequest): unknown {
    const obj: any = {};
    return obj;
  },
};

export const ComposerRenderer = {
  encode(message: ComposerRenderer, writer: Writer = Writer.create()): Writer {
    writer.uint32(10).string(message.welcomeText);
    if (message.headerRenderer !== undefined && message.headerRenderer !== undefined) {
      HeaderRenderer.encode(message.headerRenderer, writer.uint32(18).fork()).ldelim();
    }
    writer.uint32(26).string(message.unathorizedText);
    return writer;
  },
  decode(input: Uint8Array | Reader, length?: number): ComposerRenderer {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseComposerRenderer } as ComposerRenderer;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.welcomeText = reader.string();
          break;
        case 2:
          message.headerRenderer = HeaderRenderer.decode(reader, reader.uint32());
          break;
        case 3:
          message.unathorizedText = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): ComposerRenderer {
    const message = { ...baseComposerRenderer } as ComposerRenderer;
    if (object.welcomeText !== undefined && object.welcomeText !== null) {
      message.welcomeText = String(object.welcomeText);
    } else {
      message.welcomeText = "";
    }
    if (object.headerRenderer !== undefined && object.headerRenderer !== null) {
      message.headerRenderer = HeaderRenderer.fromJSON(object.headerRenderer);
    } else {
      message.headerRenderer = undefined;
    }
    if (object.unathorizedText !== undefined && object.unathorizedText !== null) {
      message.unathorizedText = String(object.unathorizedText);
    } else {
      message.unathorizedText = "";
    }
    return message;
  },
  fromPartial(object: DeepPartial<ComposerRenderer>): ComposerRenderer {
    const message = { ...baseComposerRenderer } as ComposerRenderer;
    if (object.welcomeText !== undefined && object.welcomeText !== null) {
      message.welcomeText = object.welcomeText;
    } else {
      message.welcomeText = "";
    }
    if (object.headerRenderer !== undefined && object.headerRenderer !== null) {
      message.headerRenderer = HeaderRenderer.fromPartial(object.headerRenderer);
    } else {
      message.headerRenderer = undefined;
    }
    if (object.unathorizedText !== undefined && object.unathorizedText !== null) {
      message.unathorizedText = object.unathorizedText;
    } else {
      message.unathorizedText = "";
    }
    return message;
  },
  toJSON(message: ComposerRenderer): unknown {
    const obj: any = {};
    message.welcomeText !== undefined && (obj.welcomeText = message.welcomeText);
    message.headerRenderer !== undefined && (obj.headerRenderer = message.headerRenderer ? HeaderRenderer.toJSON(message.headerRenderer) : undefined);
    message.unathorizedText !== undefined && (obj.unathorizedText = message.unathorizedText);
    return obj;
  },
};

export const IndexRequest = {
  encode(_: IndexRequest, writer: Writer = Writer.create()): Writer {
    return writer;
  },
  decode(input: Uint8Array | Reader, length?: number): IndexRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseIndexRequest } as IndexRequest;
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
  fromJSON(_: any): IndexRequest {
    const message = { ...baseIndexRequest } as IndexRequest;
    return message;
  },
  fromPartial(_: DeepPartial<IndexRequest>): IndexRequest {
    const message = { ...baseIndexRequest } as IndexRequest;
    return message;
  },
  toJSON(_: IndexRequest): unknown {
    const obj: any = {};
    return obj;
  },
};

export const IndexRenderer = {
  encode(message: IndexRenderer, writer: Writer = Writer.create()): Writer {
    writer.uint32(10).string(message.text);
    writer.uint32(26).string(message.feedUrl);
    if (message.headerRenderer !== undefined && message.headerRenderer !== undefined) {
      HeaderRenderer.encode(message.headerRenderer, writer.uint32(18).fork()).ldelim();
    }
    return writer;
  },
  decode(input: Uint8Array | Reader, length?: number): IndexRenderer {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseIndexRenderer } as IndexRenderer;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.text = reader.string();
          break;
        case 3:
          message.feedUrl = reader.string();
          break;
        case 2:
          message.headerRenderer = HeaderRenderer.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): IndexRenderer {
    const message = { ...baseIndexRenderer } as IndexRenderer;
    if (object.text !== undefined && object.text !== null) {
      message.text = String(object.text);
    } else {
      message.text = "";
    }
    if (object.feedUrl !== undefined && object.feedUrl !== null) {
      message.feedUrl = String(object.feedUrl);
    } else {
      message.feedUrl = "";
    }
    if (object.headerRenderer !== undefined && object.headerRenderer !== null) {
      message.headerRenderer = HeaderRenderer.fromJSON(object.headerRenderer);
    } else {
      message.headerRenderer = undefined;
    }
    return message;
  },
  fromPartial(object: DeepPartial<IndexRenderer>): IndexRenderer {
    const message = { ...baseIndexRenderer } as IndexRenderer;
    if (object.text !== undefined && object.text !== null) {
      message.text = object.text;
    } else {
      message.text = "";
    }
    if (object.feedUrl !== undefined && object.feedUrl !== null) {
      message.feedUrl = object.feedUrl;
    } else {
      message.feedUrl = "";
    }
    if (object.headerRenderer !== undefined && object.headerRenderer !== null) {
      message.headerRenderer = HeaderRenderer.fromPartial(object.headerRenderer);
    } else {
      message.headerRenderer = undefined;
    }
    return message;
  },
  toJSON(message: IndexRenderer): unknown {
    const obj: any = {};
    message.text !== undefined && (obj.text = message.text);
    message.feedUrl !== undefined && (obj.feedUrl = message.feedUrl);
    message.headerRenderer !== undefined && (obj.headerRenderer = message.headerRenderer ? HeaderRenderer.toJSON(message.headerRenderer) : undefined);
    return obj;
  },
};

export const HeaderRenderer = {
  encode(message: HeaderRenderer, writer: Writer = Writer.create()): Writer {
    writer.uint32(10).string(message.currentUserId);
    writer.uint32(18).string(message.currentUserName);
    for (const v of message.links) {
      HeaderRenderer_Link.encode(v!, writer.uint32(26).fork()).ldelim();
    }
    return writer;
  },
  decode(input: Uint8Array | Reader, length?: number): HeaderRenderer {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseHeaderRenderer } as HeaderRenderer;
    message.links = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.currentUserId = reader.string();
          break;
        case 2:
          message.currentUserName = reader.string();
          break;
        case 3:
          message.links.push(HeaderRenderer_Link.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): HeaderRenderer {
    const message = { ...baseHeaderRenderer } as HeaderRenderer;
    message.links = [];
    if (object.currentUserId !== undefined && object.currentUserId !== null) {
      message.currentUserId = String(object.currentUserId);
    } else {
      message.currentUserId = "";
    }
    if (object.currentUserName !== undefined && object.currentUserName !== null) {
      message.currentUserName = String(object.currentUserName);
    } else {
      message.currentUserName = "";
    }
    if (object.links !== undefined && object.links !== null) {
      for (const e of object.links) {
        message.links.push(HeaderRenderer_Link.fromJSON(e));
      }
    }
    return message;
  },
  fromPartial(object: DeepPartial<HeaderRenderer>): HeaderRenderer {
    const message = { ...baseHeaderRenderer } as HeaderRenderer;
    message.links = [];
    if (object.currentUserId !== undefined && object.currentUserId !== null) {
      message.currentUserId = object.currentUserId;
    } else {
      message.currentUserId = "";
    }
    if (object.currentUserName !== undefined && object.currentUserName !== null) {
      message.currentUserName = object.currentUserName;
    } else {
      message.currentUserName = "";
    }
    if (object.links !== undefined && object.links !== null) {
      for (const e of object.links) {
        message.links.push(HeaderRenderer_Link.fromPartial(e));
      }
    }
    return message;
  },
  toJSON(message: HeaderRenderer): unknown {
    const obj: any = {};
    message.currentUserId !== undefined && (obj.currentUserId = message.currentUserId);
    message.currentUserName !== undefined && (obj.currentUserName = message.currentUserName);
    if (message.links) {
      obj.links = message.links.map(e => e ? HeaderRenderer_Link.toJSON(e) : undefined);
    } else {
      obj.links = [];
    }
    return obj;
  },
};

export const HeaderRenderer_Link = {
  encode(message: HeaderRenderer_Link, writer: Writer = Writer.create()): Writer {
    writer.uint32(10).string(message.url);
    writer.uint32(18).string(message.label);
    return writer;
  },
  decode(input: Uint8Array | Reader, length?: number): HeaderRenderer_Link {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseHeaderRenderer_Link } as HeaderRenderer_Link;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.url = reader.string();
          break;
        case 2:
          message.label = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): HeaderRenderer_Link {
    const message = { ...baseHeaderRenderer_Link } as HeaderRenderer_Link;
    if (object.url !== undefined && object.url !== null) {
      message.url = String(object.url);
    } else {
      message.url = "";
    }
    if (object.label !== undefined && object.label !== null) {
      message.label = String(object.label);
    } else {
      message.label = "";
    }
    return message;
  },
  fromPartial(object: DeepPartial<HeaderRenderer_Link>): HeaderRenderer_Link {
    const message = { ...baseHeaderRenderer_Link } as HeaderRenderer_Link;
    if (object.url !== undefined && object.url !== null) {
      message.url = object.url;
    } else {
      message.url = "";
    }
    if (object.label !== undefined && object.label !== null) {
      message.label = object.label;
    } else {
      message.label = "";
    }
    return message;
  },
  toJSON(message: HeaderRenderer_Link): unknown {
    const obj: any = {};
    message.url !== undefined && (obj.url = message.url);
    message.label !== undefined && (obj.label = message.label);
    return obj;
  },
};

export const VKCallbackRequest = {
  encode(message: VKCallbackRequest, writer: Writer = Writer.create()): Writer {
    writer.uint32(10).string(message.vkCode);
    return writer;
  },
  decode(input: Uint8Array | Reader, length?: number): VKCallbackRequest {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseVKCallbackRequest } as VKCallbackRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.vkCode = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): VKCallbackRequest {
    const message = { ...baseVKCallbackRequest } as VKCallbackRequest;
    if (object.vkCode !== undefined && object.vkCode !== null) {
      message.vkCode = String(object.vkCode);
    } else {
      message.vkCode = "";
    }
    return message;
  },
  fromPartial(object: DeepPartial<VKCallbackRequest>): VKCallbackRequest {
    const message = { ...baseVKCallbackRequest } as VKCallbackRequest;
    if (object.vkCode !== undefined && object.vkCode !== null) {
      message.vkCode = object.vkCode;
    } else {
      message.vkCode = "";
    }
    return message;
  },
  toJSON(message: VKCallbackRequest): unknown {
    const obj: any = {};
    message.vkCode !== undefined && (obj.vkCode = message.vkCode);
    return obj;
  },
};

export const VKCallbackRenderer = {
  encode(message: VKCallbackRenderer, writer: Writer = Writer.create()): Writer {
    if (message.headerRenderer !== undefined && message.headerRenderer !== undefined) {
      HeaderRenderer.encode(message.headerRenderer, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },
  decode(input: Uint8Array | Reader, length?: number): VKCallbackRenderer {
    const reader = input instanceof Uint8Array ? new Reader(input) : input;
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseVKCallbackRenderer } as VKCallbackRenderer;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.headerRenderer = HeaderRenderer.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },
  fromJSON(object: any): VKCallbackRenderer {
    const message = { ...baseVKCallbackRenderer } as VKCallbackRenderer;
    if (object.headerRenderer !== undefined && object.headerRenderer !== null) {
      message.headerRenderer = HeaderRenderer.fromJSON(object.headerRenderer);
    } else {
      message.headerRenderer = undefined;
    }
    return message;
  },
  fromPartial(object: DeepPartial<VKCallbackRenderer>): VKCallbackRenderer {
    const message = { ...baseVKCallbackRenderer } as VKCallbackRenderer;
    if (object.headerRenderer !== undefined && object.headerRenderer !== null) {
      message.headerRenderer = HeaderRenderer.fromPartial(object.headerRenderer);
    } else {
      message.headerRenderer = undefined;
    }
    return message;
  },
  toJSON(message: VKCallbackRenderer): unknown {
    const obj: any = {};
    message.headerRenderer !== undefined && (obj.headerRenderer = message.headerRenderer ? HeaderRenderer.toJSON(message.headerRenderer) : undefined);
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