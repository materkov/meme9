/* eslint-disable */
import { util, configure, Reader, Writer } from "protobufjs/minimal";
import * as Long from "long";

export const protobufPackage = "meme";

export enum Renderers {
  UNKNOWN = "UNKNOWN",
  FEED = "FEED",
  PROFILE = "PROFILE",
  LOGIN = "LOGIN",
  POST = "POST",
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
    case 3:
    case "LOGIN":
      return Renderers.LOGIN;
    case 4:
    case "POST":
      return Renderers.POST;
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
    case Renderers.LOGIN:
      return "LOGIN";
    case Renderers.POST:
      return "POST";
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
    case Renderers.LOGIN:
      return 3;
    case Renderers.POST:
      return 4;
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
  posts: Post[];
}

export interface FeedGetRequest {}

export interface FeedGetResponse {
  renderer: FeedRenderer | undefined;
}

export interface Post {
  id: string;
  authorId: string;
  authorAvatar: string;
  authorName: string;
  authorUrl: string;
  dateDisplay: string;
  text: string;
  imageUrl: string;
}

export interface FeedRenderer {
  posts: Post[];
}

export interface PostRenderer {
  post: Post | undefined;
}

export interface PostPageResponse {
  renderer: PostRenderer | undefined;
}

export interface FeedGetHeaderRequest {}

export interface FeedGetHeaderResponse {
  renderer: HeaderRenderer | undefined;
}

export interface HeaderRenderer {
  mainUrl: string;
  userName: string;
  userAvatar: string;
  isAuthorized: boolean;
}

export interface ResolveRouteResponse {
  renderer: Renderers;
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
    for (const v of message.posts) {
      Post.encode(v!, writer.uint32(34).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): ProfileRenderer {
    const reader = input instanceof Reader ? input : new Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseProfileRenderer } as ProfileRenderer;
    message.posts = [];
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
        case 4:
          message.posts.push(Post.decode(reader, reader.uint32()));
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
    message.posts = [];
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
    if (object.posts !== undefined && object.posts !== null) {
      for (const e of object.posts) {
        message.posts.push(Post.fromJSON(e));
      }
    }
    return message;
  },

  toJSON(message: ProfileRenderer): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.name !== undefined && (obj.name = message.name);
    message.avatar !== undefined && (obj.avatar = message.avatar);
    if (message.posts) {
      obj.posts = message.posts.map((e) => (e ? Post.toJSON(e) : undefined));
    } else {
      obj.posts = [];
    }
    return obj;
  },

  fromPartial(object: DeepPartial<ProfileRenderer>): ProfileRenderer {
    const message = { ...baseProfileRenderer } as ProfileRenderer;
    message.posts = [];
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
    if (object.posts !== undefined && object.posts !== null) {
      for (const e of object.posts) {
        message.posts.push(Post.fromPartial(e));
      }
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

const basePost: object = {
  id: "",
  authorId: "",
  authorAvatar: "",
  authorName: "",
  authorUrl: "",
  dateDisplay: "",
  text: "",
  imageUrl: "",
};

export const Post = {
  encode(message: Post, writer: Writer = Writer.create()): Writer {
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

  decode(input: Reader | Uint8Array, length?: number): Post {
    const reader = input instanceof Reader ? input : new Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...basePost } as Post;
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

  fromJSON(object: any): Post {
    const message = { ...basePost } as Post;
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

  toJSON(message: Post): unknown {
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

  fromPartial(object: DeepPartial<Post>): Post {
    const message = { ...basePost } as Post;
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

const baseFeedRenderer: object = {};

export const FeedRenderer = {
  encode(message: FeedRenderer, writer: Writer = Writer.create()): Writer {
    for (const v of message.posts) {
      Post.encode(v!, writer.uint32(10).fork()).ldelim();
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
          message.posts.push(Post.decode(reader, reader.uint32()));
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
        message.posts.push(Post.fromJSON(e));
      }
    }
    return message;
  },

  toJSON(message: FeedRenderer): unknown {
    const obj: any = {};
    if (message.posts) {
      obj.posts = message.posts.map((e) => (e ? Post.toJSON(e) : undefined));
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
        message.posts.push(Post.fromPartial(e));
      }
    }
    return message;
  },
};

const basePostRenderer: object = {};

export const PostRenderer = {
  encode(message: PostRenderer, writer: Writer = Writer.create()): Writer {
    if (message.post !== undefined) {
      Post.encode(message.post, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): PostRenderer {
    const reader = input instanceof Reader ? input : new Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...basePostRenderer } as PostRenderer;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.post = Post.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): PostRenderer {
    const message = { ...basePostRenderer } as PostRenderer;
    if (object.post !== undefined && object.post !== null) {
      message.post = Post.fromJSON(object.post);
    } else {
      message.post = undefined;
    }
    return message;
  },

  toJSON(message: PostRenderer): unknown {
    const obj: any = {};
    message.post !== undefined &&
      (obj.post = message.post ? Post.toJSON(message.post) : undefined);
    return obj;
  },

  fromPartial(object: DeepPartial<PostRenderer>): PostRenderer {
    const message = { ...basePostRenderer } as PostRenderer;
    if (object.post !== undefined && object.post !== null) {
      message.post = Post.fromPartial(object.post);
    } else {
      message.post = undefined;
    }
    return message;
  },
};

const basePostPageResponse: object = {};

export const PostPageResponse = {
  encode(message: PostPageResponse, writer: Writer = Writer.create()): Writer {
    if (message.renderer !== undefined) {
      PostRenderer.encode(message.renderer, writer.uint32(10).fork()).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): PostPageResponse {
    const reader = input instanceof Reader ? input : new Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...basePostPageResponse } as PostPageResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.renderer = PostRenderer.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): PostPageResponse {
    const message = { ...basePostPageResponse } as PostPageResponse;
    if (object.renderer !== undefined && object.renderer !== null) {
      message.renderer = PostRenderer.fromJSON(object.renderer);
    } else {
      message.renderer = undefined;
    }
    return message;
  },

  toJSON(message: PostPageResponse): unknown {
    const obj: any = {};
    message.renderer !== undefined &&
      (obj.renderer = message.renderer
        ? PostRenderer.toJSON(message.renderer)
        : undefined);
    return obj;
  },

  fromPartial(object: DeepPartial<PostPageResponse>): PostPageResponse {
    const message = { ...basePostPageResponse } as PostPageResponse;
    if (object.renderer !== undefined && object.renderer !== null) {
      message.renderer = PostRenderer.fromPartial(object.renderer);
    } else {
      message.renderer = undefined;
    }
    return message;
  },
};

const baseFeedGetHeaderRequest: object = {};

export const FeedGetHeaderRequest = {
  encode(_: FeedGetHeaderRequest, writer: Writer = Writer.create()): Writer {
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): FeedGetHeaderRequest {
    const reader = input instanceof Reader ? input : new Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseFeedGetHeaderRequest } as FeedGetHeaderRequest;
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

  fromJSON(_: any): FeedGetHeaderRequest {
    const message = { ...baseFeedGetHeaderRequest } as FeedGetHeaderRequest;
    return message;
  },

  toJSON(_: FeedGetHeaderRequest): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(_: DeepPartial<FeedGetHeaderRequest>): FeedGetHeaderRequest {
    const message = { ...baseFeedGetHeaderRequest } as FeedGetHeaderRequest;
    return message;
  },
};

const baseFeedGetHeaderResponse: object = {};

export const FeedGetHeaderResponse = {
  encode(
    message: FeedGetHeaderResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.renderer !== undefined) {
      HeaderRenderer.encode(
        message.renderer,
        writer.uint32(10).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): FeedGetHeaderResponse {
    const reader = input instanceof Reader ? input : new Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseFeedGetHeaderResponse } as FeedGetHeaderResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.renderer = HeaderRenderer.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): FeedGetHeaderResponse {
    const message = { ...baseFeedGetHeaderResponse } as FeedGetHeaderResponse;
    if (object.renderer !== undefined && object.renderer !== null) {
      message.renderer = HeaderRenderer.fromJSON(object.renderer);
    } else {
      message.renderer = undefined;
    }
    return message;
  },

  toJSON(message: FeedGetHeaderResponse): unknown {
    const obj: any = {};
    message.renderer !== undefined &&
      (obj.renderer = message.renderer
        ? HeaderRenderer.toJSON(message.renderer)
        : undefined);
    return obj;
  },

  fromPartial(
    object: DeepPartial<FeedGetHeaderResponse>
  ): FeedGetHeaderResponse {
    const message = { ...baseFeedGetHeaderResponse } as FeedGetHeaderResponse;
    if (object.renderer !== undefined && object.renderer !== null) {
      message.renderer = HeaderRenderer.fromPartial(object.renderer);
    } else {
      message.renderer = undefined;
    }
    return message;
  },
};

const baseHeaderRenderer: object = {
  mainUrl: "",
  userName: "",
  userAvatar: "",
  isAuthorized: false,
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
    if (message.isAuthorized === true) {
      writer.uint32(32).bool(message.isAuthorized);
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
        case 4:
          message.isAuthorized = reader.bool();
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
    if (object.isAuthorized !== undefined && object.isAuthorized !== null) {
      message.isAuthorized = Boolean(object.isAuthorized);
    } else {
      message.isAuthorized = false;
    }
    return message;
  },

  toJSON(message: HeaderRenderer): unknown {
    const obj: any = {};
    message.mainUrl !== undefined && (obj.mainUrl = message.mainUrl);
    message.userName !== undefined && (obj.userName = message.userName);
    message.userAvatar !== undefined && (obj.userAvatar = message.userAvatar);
    message.isAuthorized !== undefined &&
      (obj.isAuthorized = message.isAuthorized);
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
    if (object.isAuthorized !== undefined && object.isAuthorized !== null) {
      message.isAuthorized = object.isAuthorized;
    } else {
      message.isAuthorized = false;
    }
    return message;
  },
};

const baseResolveRouteResponse: object = { renderer: Renderers.UNKNOWN };

export const ResolveRouteResponse = {
  encode(
    message: ResolveRouteResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.renderer !== Renderers.UNKNOWN) {
      writer.uint32(8).int32(renderersToNumber(message.renderer));
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): ResolveRouteResponse {
    const reader = input instanceof Reader ? input : new Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseResolveRouteResponse } as ResolveRouteResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.renderer = renderersFromJSON(reader.int32());
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
    if (object.renderer !== undefined && object.renderer !== null) {
      message.renderer = renderersFromJSON(object.renderer);
    } else {
      message.renderer = Renderers.UNKNOWN;
    }
    return message;
  },

  toJSON(message: ResolveRouteResponse): unknown {
    const obj: any = {};
    message.renderer !== undefined &&
      (obj.renderer = renderersToJSON(message.renderer));
    return obj;
  },

  fromPartial(object: DeepPartial<ResolveRouteResponse>): ResolveRouteResponse {
    const message = { ...baseResolveRouteResponse } as ResolveRouteResponse;
    if (object.renderer !== undefined && object.renderer !== null) {
      message.renderer = object.renderer;
    } else {
      message.renderer = Renderers.UNKNOWN;
    }
    return message;
  },
};

export interface Feed {
  Get(request: FeedGetRequest): Promise<FeedGetResponse>;
  GetHeader(request: FeedGetHeaderRequest): Promise<FeedGetHeaderResponse>;
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

  GetHeader(request: FeedGetHeaderRequest): Promise<FeedGetHeaderResponse> {
    const data = FeedGetHeaderRequest.encode(request).finish();
    const promise = this.rpc.request("meme.Feed", "GetHeader", data);
    return promise.then((data) =>
      FeedGetHeaderResponse.decode(new Reader(data))
    );
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
