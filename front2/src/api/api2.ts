/* eslint-disable */
import { util, configure, Reader, Writer } from "protobufjs/minimal";
import * as Long from "long";

export const protobufPackage = "meme";

export enum Renderers {
  UNKNOWN = "UNKNOWN",
  FEED = "FEED",
  PROFILE = "PROFILE",
  UNRECOGNIZED = "UNRECOGNIZED",
}

export function renderersFromJSON(object: any): Renderers {
  switch (object) {
    case 0:
    case "UNKNOWN":
      return Renderers.UNKNOWN;
    case 1:
    case "FEED":
      return Renderers.FEED;
    case 2:
    case "PROFILE":
      return Renderers.PROFILE;
    case -1:
    case "UNRECOGNIZED":
    default:
      return Renderers.UNRECOGNIZED;
  }
}

export function renderersToJSON(object: Renderers): string {
  switch (object) {
    case Renderers.UNKNOWN:
      return "UNKNOWN";
    case Renderers.FEED:
      return "FEED";
    case Renderers.PROFILE:
      return "PROFILE";
    default:
      return "UNKNOWN";
  }
}

export function renderersToNumber(object: Renderers): number {
  switch (object) {
    case Renderers.UNKNOWN:
      return 0;
    case Renderers.FEED:
      return 1;
    case Renderers.PROFILE:
      return 2;
    default:
      return 0;
  }
}

export interface ProfileGetRequest {
  id: string;
}

export interface ProfileGetResponse {
  renderer: ProfileRenderer | undefined;
}

export interface ProfileRenderer {
  id: string;
  name: string;
  avatar: string;
}

export interface FeedGetRequest {}

export interface FeedGetResponse {
  renderer: FeedRenderer | undefined;
}

export interface FeedRenderer {
  posts: FeedRenderer_Post[];
}

export interface FeedRenderer_Post {
  id: string;
  authorId: string;
  authorAvatar: string;
  authorName: string;
  authorUrl: string;
  dateDisplay: string;
  text: string;
  imageUrl: string;
}

export interface HeaderRenderer {
  mainUrl: string;
  userName: string;
  userAvatar: string;
}

const baseProfileGetRequest: object = { id: "" };

export const ProfileGetRequest = {
  encode(message: ProfileGetRequest, writer: Writer = Writer.create()): Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): ProfileGetRequest {
    const reader = input instanceof Reader ? input : new Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseProfileGetRequest } as ProfileGetRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ProfileGetRequest {
    const message = { ...baseProfileGetRequest } as ProfileGetRequest;
    if (object.id !== undefined && object.id !== null) {
      message.id = String(object.id);
    } else {
      message.id = "";
    }
    return message;
  },

  toJSON(message: ProfileGetRequest): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    return obj;
  },

  fromPartial(object: DeepPartial<ProfileGetRequest>): ProfileGetRequest {
    const message = { ...baseProfileGetRequest } as ProfileGetRequest;
    if (object.id !== undefined && object.id !== null) {
      message.id = object.id;
    } else {
      message.id = "";
    }
    return message;
  },
};

const baseProfileGetResponse: object = {};

export const ProfileGetResponse = {
  encode(
    message: ProfileGetResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.renderer !== undefined) {
      ProfileRenderer.encode(
        message.renderer,
        writer.uint32(10).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): ProfileGetResponse {
    const reader = input instanceof Reader ? input : new Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseProfileGetResponse } as ProfileGetResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.renderer = ProfileRenderer.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ProfileGetResponse {
    const message = { ...baseProfileGetResponse } as ProfileGetResponse;
    if (object.renderer !== undefined && object.renderer !== null) {
      message.renderer = ProfileRenderer.fromJSON(object.renderer);
    } else {
      message.renderer = undefined;
    }
    return message;
  },

  toJSON(message: ProfileGetResponse): unknown {
    const obj: any = {};
    message.renderer !== undefined &&
      (obj.renderer = message.renderer
        ? ProfileRenderer.toJSON(message.renderer)
        : undefined);
    return obj;
  },

  fromPartial(object: DeepPartial<ProfileGetResponse>): ProfileGetResponse {
    const message = { ...baseProfileGetResponse } as ProfileGetResponse;
    if (object.renderer !== undefined && object.renderer !== null) {
      message.renderer = ProfileRenderer.fromPartial(object.renderer);
    } else {
      message.renderer = undefined;
    }
    return message;
  },
};

const baseProfileRenderer: object = { id: "", name: "", avatar: "" };

export const ProfileRenderer = {
  encode(message: ProfileRenderer, writer: Writer = Writer.create()): Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.name !== "") {
      writer.uint32(18).string(message.name);
    }
    if (message.avatar !== "") {
      writer.uint32(26).string(message.avatar);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): ProfileRenderer {
    const reader = input instanceof Reader ? input : new Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseProfileRenderer } as ProfileRenderer;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.string();
          break;
        case 2:
          message.name = reader.string();
          break;
        case 3:
          message.avatar = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ProfileRenderer {
    const message = { ...baseProfileRenderer } as ProfileRenderer;
    if (object.id !== undefined && object.id !== null) {
      message.id = String(object.id);
    } else {
      message.id = "";
    }
    if (object.name !== undefined && object.name !== null) {
      message.name = String(object.name);
    } else {
      message.name = "";
    }
    if (object.avatar !== undefined && object.avatar !== null) {
      message.avatar = String(object.avatar);
    } else {
      message.avatar = "";
    }
    return message;
  },

  toJSON(message: ProfileRenderer): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.name !== undefined && (obj.name = message.name);
    message.avatar !== undefined && (obj.avatar = message.avatar);
    return obj;
  },

  fromPartial(object: DeepPartial<ProfileRenderer>): ProfileRenderer {
    const message = { ...baseProfileRenderer } as ProfileRenderer;
    if (object.id !== undefined && object.id !== null) {
      message.id = object.id;
    } else {
      message.id = "";
    }
    if (object.name !== undefined && object.name !== null) {
      message.name = object.name;
    } else {
      message.name = "";
    }
    if (object.avatar !== undefined && object.avatar !== null) {
      message.avatar = object.avatar;
    } else {
      message.avatar = "";
    }
    return message;
  },
};

const baseFeedGetRequest: object = {};

export const FeedGetRequest = {
  encode(_: FeedGetRequest, writer: Writer = Writer.create()): Writer {
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): FeedGetRequest {
    const reader = input instanceof Reader ? input : new Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseFeedGetRequest } as FeedGetRequest;
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

  fromJSON(_: any): FeedGetRequest {
    const message = { ...baseFeedGetRequest } as FeedGetRequest;
    return message;
  },

  toJSON(_: FeedGetRequest): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(_: DeepPartial<FeedGetRequest>): FeedGetRequest {
    const message = { ...baseFeedGetRequest } as FeedGetRequest;
    return message;
  },
};

const baseFeedGetResponse: object = {};

export const FeedGetResponse = {
  encode(message: FeedGetResponse, writer: Writer = Writer.create()): Writer {
    if (message.renderer !== undefined) {
      FeedRenderer.encode(message.renderer, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): FeedGetResponse {
    const reader = input instanceof Reader ? input : new Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseFeedGetResponse } as FeedGetResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.renderer = FeedRenderer.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): FeedGetResponse {
    const message = { ...baseFeedGetResponse } as FeedGetResponse;
    if (object.renderer !== undefined && object.renderer !== null) {
      message.renderer = FeedRenderer.fromJSON(object.renderer);
    } else {
      message.renderer = undefined;
    }
    return message;
  },

  toJSON(message: FeedGetResponse): unknown {
    const obj: any = {};
    message.renderer !== undefined &&
      (obj.renderer = message.renderer
        ? FeedRenderer.toJSON(message.renderer)
        : undefined);
    return obj;
  },

  fromPartial(object: DeepPartial<FeedGetResponse>): FeedGetResponse {
    const message = { ...baseFeedGetResponse } as FeedGetResponse;
    if (object.renderer !== undefined && object.renderer !== null) {
      message.renderer = FeedRenderer.fromPartial(object.renderer);
    } else {
      message.renderer = undefined;
    }
    return message;
  },
};

const baseFeedRenderer: object = {};

export const FeedRenderer = {
  encode(message: FeedRenderer, writer: Writer = Writer.create()): Writer {
    for (const v of message.posts) {
      FeedRenderer_Post.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): FeedRenderer {
    const reader = input instanceof Reader ? input : new Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseFeedRenderer } as FeedRenderer;
    message.posts = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.posts.push(FeedRenderer_Post.decode(reader, reader.uint32()));
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): FeedRenderer {
    const message = { ...baseFeedRenderer } as FeedRenderer;
    message.posts = [];
    if (object.posts !== undefined && object.posts !== null) {
      for (const e of object.posts) {
        message.posts.push(FeedRenderer_Post.fromJSON(e));
      }
    }
    return message;
  },

  toJSON(message: FeedRenderer): unknown {
    const obj: any = {};
    if (message.posts) {
      obj.posts = message.posts.map((e) =>
        e ? FeedRenderer_Post.toJSON(e) : undefined
      );
    } else {
      obj.posts = [];
    }
    return obj;
  },

  fromPartial(object: DeepPartial<FeedRenderer>): FeedRenderer {
    const message = { ...baseFeedRenderer } as FeedRenderer;
    message.posts = [];
    if (object.posts !== undefined && object.posts !== null) {
      for (const e of object.posts) {
        message.posts.push(FeedRenderer_Post.fromPartial(e));
      }
    }
    return message;
  },
};

const baseFeedRenderer_Post: object = {
  id: "",
  authorId: "",
  authorAvatar: "",
  authorName: "",
  authorUrl: "",
  dateDisplay: "",
  text: "",
  imageUrl: "",
};

export const FeedRenderer_Post = {
  encode(message: FeedRenderer_Post, writer: Writer = Writer.create()): Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.authorId !== "") {
      writer.uint32(18).string(message.authorId);
    }
    if (message.authorAvatar !== "") {
      writer.uint32(26).string(message.authorAvatar);
    }
    if (message.authorName !== "") {
      writer.uint32(34).string(message.authorName);
    }
    if (message.authorUrl !== "") {
      writer.uint32(66).string(message.authorUrl);
    }
    if (message.dateDisplay !== "") {
      writer.uint32(42).string(message.dateDisplay);
    }
    if (message.text !== "") {
      writer.uint32(50).string(message.text);
    }
    if (message.imageUrl !== "") {
      writer.uint32(58).string(message.imageUrl);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): FeedRenderer_Post {
    const reader = input instanceof Reader ? input : new Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseFeedRenderer_Post } as FeedRenderer_Post;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.id = reader.string();
          break;
        case 2:
          message.authorId = reader.string();
          break;
        case 3:
          message.authorAvatar = reader.string();
          break;
        case 4:
          message.authorName = reader.string();
          break;
        case 8:
          message.authorUrl = reader.string();
          break;
        case 5:
          message.dateDisplay = reader.string();
          break;
        case 6:
          message.text = reader.string();
          break;
        case 7:
          message.imageUrl = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): FeedRenderer_Post {
    const message = { ...baseFeedRenderer_Post } as FeedRenderer_Post;
    if (object.id !== undefined && object.id !== null) {
      message.id = String(object.id);
    } else {
      message.id = "";
    }
    if (object.authorId !== undefined && object.authorId !== null) {
      message.authorId = String(object.authorId);
    } else {
      message.authorId = "";
    }
    if (object.authorAvatar !== undefined && object.authorAvatar !== null) {
      message.authorAvatar = String(object.authorAvatar);
    } else {
      message.authorAvatar = "";
    }
    if (object.authorName !== undefined && object.authorName !== null) {
      message.authorName = String(object.authorName);
    } else {
      message.authorName = "";
    }
    if (object.authorUrl !== undefined && object.authorUrl !== null) {
      message.authorUrl = String(object.authorUrl);
    } else {
      message.authorUrl = "";
    }
    if (object.dateDisplay !== undefined && object.dateDisplay !== null) {
      message.dateDisplay = String(object.dateDisplay);
    } else {
      message.dateDisplay = "";
    }
    if (object.text !== undefined && object.text !== null) {
      message.text = String(object.text);
    } else {
      message.text = "";
    }
    if (object.imageUrl !== undefined && object.imageUrl !== null) {
      message.imageUrl = String(object.imageUrl);
    } else {
      message.imageUrl = "";
    }
    return message;
  },

  toJSON(message: FeedRenderer_Post): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.authorId !== undefined && (obj.authorId = message.authorId);
    message.authorAvatar !== undefined &&
      (obj.authorAvatar = message.authorAvatar);
    message.authorName !== undefined && (obj.authorName = message.authorName);
    message.authorUrl !== undefined && (obj.authorUrl = message.authorUrl);
    message.dateDisplay !== undefined &&
      (obj.dateDisplay = message.dateDisplay);
    message.text !== undefined && (obj.text = message.text);
    message.imageUrl !== undefined && (obj.imageUrl = message.imageUrl);
    return obj;
  },

  fromPartial(object: DeepPartial<FeedRenderer_Post>): FeedRenderer_Post {
    const message = { ...baseFeedRenderer_Post } as FeedRenderer_Post;
    if (object.id !== undefined && object.id !== null) {
      message.id = object.id;
    } else {
      message.id = "";
    }
    if (object.authorId !== undefined && object.authorId !== null) {
      message.authorId = object.authorId;
    } else {
      message.authorId = "";
    }
    if (object.authorAvatar !== undefined && object.authorAvatar !== null) {
      message.authorAvatar = object.authorAvatar;
    } else {
      message.authorAvatar = "";
    }
    if (object.authorName !== undefined && object.authorName !== null) {
      message.authorName = object.authorName;
    } else {
      message.authorName = "";
    }
    if (object.authorUrl !== undefined && object.authorUrl !== null) {
      message.authorUrl = object.authorUrl;
    } else {
      message.authorUrl = "";
    }
    if (object.dateDisplay !== undefined && object.dateDisplay !== null) {
      message.dateDisplay = object.dateDisplay;
    } else {
      message.dateDisplay = "";
    }
    if (object.text !== undefined && object.text !== null) {
      message.text = object.text;
    } else {
      message.text = "";
    }
    if (object.imageUrl !== undefined && object.imageUrl !== null) {
      message.imageUrl = object.imageUrl;
    } else {
      message.imageUrl = "";
    }
    return message;
  },
};

const baseHeaderRenderer: object = {
  mainUrl: "",
  userName: "",
  userAvatar: "",
};

export const HeaderRenderer = {
  encode(message: HeaderRenderer, writer: Writer = Writer.create()): Writer {
    if (message.mainUrl !== "") {
      writer.uint32(10).string(message.mainUrl);
    }
    if (message.userName !== "") {
      writer.uint32(18).string(message.userName);
    }
    if (message.userAvatar !== "") {
      writer.uint32(26).string(message.userAvatar);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): HeaderRenderer {
    const reader = input instanceof Reader ? input : new Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseHeaderRenderer } as HeaderRenderer;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.mainUrl = reader.string();
          break;
        case 2:
          message.userName = reader.string();
          break;
        case 3:
          message.userAvatar = reader.string();
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
    if (object.mainUrl !== undefined && object.mainUrl !== null) {
      message.mainUrl = String(object.mainUrl);
    } else {
      message.mainUrl = "";
    }
    if (object.userName !== undefined && object.userName !== null) {
      message.userName = String(object.userName);
    } else {
      message.userName = "";
    }
    if (object.userAvatar !== undefined && object.userAvatar !== null) {
      message.userAvatar = String(object.userAvatar);
    } else {
      message.userAvatar = "";
    }
    return message;
  },

  toJSON(message: HeaderRenderer): unknown {
    const obj: any = {};
    message.mainUrl !== undefined && (obj.mainUrl = message.mainUrl);
    message.userName !== undefined && (obj.userName = message.userName);
    message.userAvatar !== undefined && (obj.userAvatar = message.userAvatar);
    return obj;
  },

  fromPartial(object: DeepPartial<HeaderRenderer>): HeaderRenderer {
    const message = { ...baseHeaderRenderer } as HeaderRenderer;
    if (object.mainUrl !== undefined && object.mainUrl !== null) {
      message.mainUrl = object.mainUrl;
    } else {
      message.mainUrl = "";
    }
    if (object.userName !== undefined && object.userName !== null) {
      message.userName = object.userName;
    } else {
      message.userName = "";
    }
    if (object.userAvatar !== undefined && object.userAvatar !== null) {
      message.userAvatar = object.userAvatar;
    } else {
      message.userAvatar = "";
    }
    return message;
  },
};

export interface Feed {
  Get(request: FeedGetRequest): Promise<FeedGetResponse>;
}

export class FeedClientImpl implements Feed {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
  }
  Get(request: FeedGetRequest): Promise<FeedGetResponse> {
    const data = FeedGetRequest.encode(request).finish();
    const promise = this.rpc.request("meme.Feed", "Get", data);
    return promise.then((data) => FeedGetResponse.decode(new Reader(data)));
  }
}

export interface Profile {
  Get(request: ProfileGetRequest): Promise<ProfileGetResponse>;
}

export class ProfileClientImpl implements Profile {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
  }
  Get(request: ProfileGetRequest): Promise<ProfileGetResponse> {
    const data = ProfileGetRequest.encode(request).finish();
    const promise = this.rpc.request("meme.Profile", "Get", data);
    return promise.then((data) => ProfileGetResponse.decode(new Reader(data)));
  }
}

interface Rpc {
  request(
    service: string,
    method: string,
    data: Uint8Array
  ): Promise<Uint8Array>;
}

type Builtin = Date | Function | Uint8Array | string | number | undefined;
export type DeepPartial<T> = T extends Builtin
  ? T
  : T extends Array<infer U>
  ? Array<DeepPartial<U>>
  : T extends ReadonlyArray<infer U>
  ? ReadonlyArray<DeepPartial<U>>
  : T extends {}
  ? { [K in keyof T]?: DeepPartial<T[K]> }
  : Partial<T>;

// If you get a compile-error about 'Constructor<Long> and ... have no overlap',
// add '--ts_proto_opt=esModuleInterop=true' as a flag when calling 'protoc'.
if (util.Long !== Long) {
  util.Long = Long as any;
  configure();
}
