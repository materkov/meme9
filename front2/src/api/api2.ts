/* eslint-disable */
import { util, configure, Reader, Writer } from "protobufjs/minimal";
import * as Long from "long";
import { CommentComposerRenderer } from "./posts";

export const protobufPackage = "meme";

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
  isFollowing: boolean;
}

export interface FeedGetRequest {}

export interface FeedGetResponse {
  renderer: FeedRenderer | undefined;
}

export interface Post {
  id: string;
  url: string;
  authorId: string;
  authorAvatar: string;
  authorName: string;
  authorUrl: string;
  dateDisplay: string;
  text: string;
  imageUrl: string;
  /** Likes */
  isLiked: boolean;
  likesCount: number;
  canLike: boolean;
  /** Comments */
  commentsCount: number;
  topComment: CommentRenderer | undefined;
}

export interface CommentRenderer {
  id: string;
  text: string;
  authorId: string;
  authorName: string;
  authorUrl: string;
}

export interface FeedRenderer {
  posts: Post[];
  placeholderText: string;
}

export interface PostRenderer {
  post: Post | undefined;
  comments: CommentRenderer[];
  composer: CommentComposerRenderer | undefined;
  composerPlaceholder: string;
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
  logoutUrl: string;
  loginUrl: string;
}

/** Renderers renderer = 1; */
export interface ResolveRouteResponse {}

export interface RelationsFollowRequest {
  userId: string;
}

export interface RelationsFollowResponse {}

export interface RelationsUnfollowRequest {
  userId: string;
}

export interface RelationsUnfollowResponse {}

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

const baseProfileRenderer: object = {
  id: "",
  name: "",
  avatar: "",
  isFollowing: false,
};

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
    if (message.isFollowing === true) {
      writer.uint32(40).bool(message.isFollowing);
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
        case 5:
          message.isFollowing = reader.bool();
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
    if (object.isFollowing !== undefined && object.isFollowing !== null) {
      message.isFollowing = Boolean(object.isFollowing);
    } else {
      message.isFollowing = false;
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
    message.isFollowing !== undefined &&
      (obj.isFollowing = message.isFollowing);
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
    if (object.isFollowing !== undefined && object.isFollowing !== null) {
      message.isFollowing = object.isFollowing;
    } else {
      message.isFollowing = false;
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
  url: "",
  authorId: "",
  authorAvatar: "",
  authorName: "",
  authorUrl: "",
  dateDisplay: "",
  text: "",
  imageUrl: "",
  isLiked: false,
  likesCount: 0,
  canLike: false,
  commentsCount: 0,
};

export const Post = {
  encode(message: Post, writer: Writer = Writer.create()): Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.url !== "") {
      writer.uint32(74).string(message.url);
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
    if (message.isLiked === true) {
      writer.uint32(80).bool(message.isLiked);
    }
    if (message.likesCount !== 0) {
      writer.uint32(88).int32(message.likesCount);
    }
    if (message.canLike === true) {
      writer.uint32(96).bool(message.canLike);
    }
    if (message.commentsCount !== 0) {
      writer.uint32(104).int32(message.commentsCount);
    }
    if (message.topComment !== undefined) {
      CommentRenderer.encode(
        message.topComment,
        writer.uint32(114).fork()
      ).ldelim();
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
        case 9:
          message.url = reader.string();
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
        case 10:
          message.isLiked = reader.bool();
          break;
        case 11:
          message.likesCount = reader.int32();
          break;
        case 12:
          message.canLike = reader.bool();
          break;
        case 13:
          message.commentsCount = reader.int32();
          break;
        case 14:
          message.topComment = CommentRenderer.decode(reader, reader.uint32());
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
    if (object.url !== undefined && object.url !== null) {
      message.url = String(object.url);
    } else {
      message.url = "";
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
    if (object.isLiked !== undefined && object.isLiked !== null) {
      message.isLiked = Boolean(object.isLiked);
    } else {
      message.isLiked = false;
    }
    if (object.likesCount !== undefined && object.likesCount !== null) {
      message.likesCount = Number(object.likesCount);
    } else {
      message.likesCount = 0;
    }
    if (object.canLike !== undefined && object.canLike !== null) {
      message.canLike = Boolean(object.canLike);
    } else {
      message.canLike = false;
    }
    if (object.commentsCount !== undefined && object.commentsCount !== null) {
      message.commentsCount = Number(object.commentsCount);
    } else {
      message.commentsCount = 0;
    }
    if (object.topComment !== undefined && object.topComment !== null) {
      message.topComment = CommentRenderer.fromJSON(object.topComment);
    } else {
      message.topComment = undefined;
    }
    return message;
  },

  toJSON(message: Post): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.url !== undefined && (obj.url = message.url);
    message.authorId !== undefined && (obj.authorId = message.authorId);
    message.authorAvatar !== undefined &&
      (obj.authorAvatar = message.authorAvatar);
    message.authorName !== undefined && (obj.authorName = message.authorName);
    message.authorUrl !== undefined && (obj.authorUrl = message.authorUrl);
    message.dateDisplay !== undefined &&
      (obj.dateDisplay = message.dateDisplay);
    message.text !== undefined && (obj.text = message.text);
    message.imageUrl !== undefined && (obj.imageUrl = message.imageUrl);
    message.isLiked !== undefined && (obj.isLiked = message.isLiked);
    message.likesCount !== undefined && (obj.likesCount = message.likesCount);
    message.canLike !== undefined && (obj.canLike = message.canLike);
    message.commentsCount !== undefined &&
      (obj.commentsCount = message.commentsCount);
    message.topComment !== undefined &&
      (obj.topComment = message.topComment
        ? CommentRenderer.toJSON(message.topComment)
        : undefined);
    return obj;
  },

  fromPartial(object: DeepPartial<Post>): Post {
    const message = { ...basePost } as Post;
    if (object.id !== undefined && object.id !== null) {
      message.id = object.id;
    } else {
      message.id = "";
    }
    if (object.url !== undefined && object.url !== null) {
      message.url = object.url;
    } else {
      message.url = "";
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
    if (object.isLiked !== undefined && object.isLiked !== null) {
      message.isLiked = object.isLiked;
    } else {
      message.isLiked = false;
    }
    if (object.likesCount !== undefined && object.likesCount !== null) {
      message.likesCount = object.likesCount;
    } else {
      message.likesCount = 0;
    }
    if (object.canLike !== undefined && object.canLike !== null) {
      message.canLike = object.canLike;
    } else {
      message.canLike = false;
    }
    if (object.commentsCount !== undefined && object.commentsCount !== null) {
      message.commentsCount = object.commentsCount;
    } else {
      message.commentsCount = 0;
    }
    if (object.topComment !== undefined && object.topComment !== null) {
      message.topComment = CommentRenderer.fromPartial(object.topComment);
    } else {
      message.topComment = undefined;
    }
    return message;
  },
};

const baseCommentRenderer: object = {
  id: "",
  text: "",
  authorId: "",
  authorName: "",
  authorUrl: "",
};

export const CommentRenderer = {
  encode(message: CommentRenderer, writer: Writer = Writer.create()): Writer {
    if (message.id !== "") {
      writer.uint32(10).string(message.id);
    }
    if (message.text !== "") {
      writer.uint32(18).string(message.text);
    }
    if (message.authorId !== "") {
      writer.uint32(26).string(message.authorId);
    }
    if (message.authorName !== "") {
      writer.uint32(34).string(message.authorName);
    }
    if (message.authorUrl !== "") {
      writer.uint32(42).string(message.authorUrl);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): CommentRenderer {
    const reader = input instanceof Reader ? input : new Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseCommentRenderer } as CommentRenderer;
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
          message.authorId = reader.string();
          break;
        case 4:
          message.authorName = reader.string();
          break;
        case 5:
          message.authorUrl = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CommentRenderer {
    const message = { ...baseCommentRenderer } as CommentRenderer;
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
    if (object.authorId !== undefined && object.authorId !== null) {
      message.authorId = String(object.authorId);
    } else {
      message.authorId = "";
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
    return message;
  },

  toJSON(message: CommentRenderer): unknown {
    const obj: any = {};
    message.id !== undefined && (obj.id = message.id);
    message.text !== undefined && (obj.text = message.text);
    message.authorId !== undefined && (obj.authorId = message.authorId);
    message.authorName !== undefined && (obj.authorName = message.authorName);
    message.authorUrl !== undefined && (obj.authorUrl = message.authorUrl);
    return obj;
  },

  fromPartial(object: DeepPartial<CommentRenderer>): CommentRenderer {
    const message = { ...baseCommentRenderer } as CommentRenderer;
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
    if (object.authorId !== undefined && object.authorId !== null) {
      message.authorId = object.authorId;
    } else {
      message.authorId = "";
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
    return message;
  },
};

const baseFeedRenderer: object = { placeholderText: "" };

export const FeedRenderer = {
  encode(message: FeedRenderer, writer: Writer = Writer.create()): Writer {
    for (const v of message.posts) {
      Post.encode(v!, writer.uint32(10).fork()).ldelim();
    }
    if (message.placeholderText !== "") {
      writer.uint32(18).string(message.placeholderText);
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
        case 2:
          message.placeholderText = reader.string();
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
    if (
      object.placeholderText !== undefined &&
      object.placeholderText !== null
    ) {
      message.placeholderText = String(object.placeholderText);
    } else {
      message.placeholderText = "";
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
    message.placeholderText !== undefined &&
      (obj.placeholderText = message.placeholderText);
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
    if (
      object.placeholderText !== undefined &&
      object.placeholderText !== null
    ) {
      message.placeholderText = object.placeholderText;
    } else {
      message.placeholderText = "";
    }
    return message;
  },
};

const basePostRenderer: object = { composerPlaceholder: "" };

export const PostRenderer = {
  encode(message: PostRenderer, writer: Writer = Writer.create()): Writer {
    if (message.post !== undefined) {
      Post.encode(message.post, writer.uint32(10).fork()).ldelim();
    }
    for (const v of message.comments) {
      CommentRenderer.encode(v!, writer.uint32(18).fork()).ldelim();
    }
    if (message.composer !== undefined) {
      CommentComposerRenderer.encode(
        message.composer,
        writer.uint32(26).fork()
      ).ldelim();
    }
    if (message.composerPlaceholder !== "") {
      writer.uint32(34).string(message.composerPlaceholder);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): PostRenderer {
    const reader = input instanceof Reader ? input : new Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...basePostRenderer } as PostRenderer;
    message.comments = [];
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.post = Post.decode(reader, reader.uint32());
          break;
        case 2:
          message.comments.push(
            CommentRenderer.decode(reader, reader.uint32())
          );
          break;
        case 3:
          message.composer = CommentComposerRenderer.decode(
            reader,
            reader.uint32()
          );
          break;
        case 4:
          message.composerPlaceholder = reader.string();
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
    message.comments = [];
    if (object.post !== undefined && object.post !== null) {
      message.post = Post.fromJSON(object.post);
    } else {
      message.post = undefined;
    }
    if (object.comments !== undefined && object.comments !== null) {
      for (const e of object.comments) {
        message.comments.push(CommentRenderer.fromJSON(e));
      }
    }
    if (object.composer !== undefined && object.composer !== null) {
      message.composer = CommentComposerRenderer.fromJSON(object.composer);
    } else {
      message.composer = undefined;
    }
    if (
      object.composerPlaceholder !== undefined &&
      object.composerPlaceholder !== null
    ) {
      message.composerPlaceholder = String(object.composerPlaceholder);
    } else {
      message.composerPlaceholder = "";
    }
    return message;
  },

  toJSON(message: PostRenderer): unknown {
    const obj: any = {};
    message.post !== undefined &&
      (obj.post = message.post ? Post.toJSON(message.post) : undefined);
    if (message.comments) {
      obj.comments = message.comments.map((e) =>
        e ? CommentRenderer.toJSON(e) : undefined
      );
    } else {
      obj.comments = [];
    }
    message.composer !== undefined &&
      (obj.composer = message.composer
        ? CommentComposerRenderer.toJSON(message.composer)
        : undefined);
    message.composerPlaceholder !== undefined &&
      (obj.composerPlaceholder = message.composerPlaceholder);
    return obj;
  },

  fromPartial(object: DeepPartial<PostRenderer>): PostRenderer {
    const message = { ...basePostRenderer } as PostRenderer;
    message.comments = [];
    if (object.post !== undefined && object.post !== null) {
      message.post = Post.fromPartial(object.post);
    } else {
      message.post = undefined;
    }
    if (object.comments !== undefined && object.comments !== null) {
      for (const e of object.comments) {
        message.comments.push(CommentRenderer.fromPartial(e));
      }
    }
    if (object.composer !== undefined && object.composer !== null) {
      message.composer = CommentComposerRenderer.fromPartial(object.composer);
    } else {
      message.composer = undefined;
    }
    if (
      object.composerPlaceholder !== undefined &&
      object.composerPlaceholder !== null
    ) {
      message.composerPlaceholder = object.composerPlaceholder;
    } else {
      message.composerPlaceholder = "";
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
  logoutUrl: "",
  loginUrl: "",
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
    if (message.logoutUrl !== "") {
      writer.uint32(42).string(message.logoutUrl);
    }
    if (message.loginUrl !== "") {
      writer.uint32(50).string(message.loginUrl);
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
        case 5:
          message.logoutUrl = reader.string();
          break;
        case 6:
          message.loginUrl = reader.string();
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
    if (object.logoutUrl !== undefined && object.logoutUrl !== null) {
      message.logoutUrl = String(object.logoutUrl);
    } else {
      message.logoutUrl = "";
    }
    if (object.loginUrl !== undefined && object.loginUrl !== null) {
      message.loginUrl = String(object.loginUrl);
    } else {
      message.loginUrl = "";
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
    message.logoutUrl !== undefined && (obj.logoutUrl = message.logoutUrl);
    message.loginUrl !== undefined && (obj.loginUrl = message.loginUrl);
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
    if (object.logoutUrl !== undefined && object.logoutUrl !== null) {
      message.logoutUrl = object.logoutUrl;
    } else {
      message.logoutUrl = "";
    }
    if (object.loginUrl !== undefined && object.loginUrl !== null) {
      message.loginUrl = object.loginUrl;
    } else {
      message.loginUrl = "";
    }
    return message;
  },
};

const baseResolveRouteResponse: object = {};

export const ResolveRouteResponse = {
  encode(_: ResolveRouteResponse, writer: Writer = Writer.create()): Writer {
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): ResolveRouteResponse {
    const reader = input instanceof Reader ? input : new Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseResolveRouteResponse } as ResolveRouteResponse;
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

  fromJSON(_: any): ResolveRouteResponse {
    const message = { ...baseResolveRouteResponse } as ResolveRouteResponse;
    return message;
  },

  toJSON(_: ResolveRouteResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(_: DeepPartial<ResolveRouteResponse>): ResolveRouteResponse {
    const message = { ...baseResolveRouteResponse } as ResolveRouteResponse;
    return message;
  },
};

const baseRelationsFollowRequest: object = { userId: "" };

export const RelationsFollowRequest = {
  encode(
    message: RelationsFollowRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): RelationsFollowRequest {
    const reader = input instanceof Reader ? input : new Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseRelationsFollowRequest } as RelationsFollowRequest;
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

  fromJSON(object: any): RelationsFollowRequest {
    const message = { ...baseRelationsFollowRequest } as RelationsFollowRequest;
    if (object.userId !== undefined && object.userId !== null) {
      message.userId = String(object.userId);
    } else {
      message.userId = "";
    }
    return message;
  },

  toJSON(message: RelationsFollowRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    return obj;
  },

  fromPartial(
    object: DeepPartial<RelationsFollowRequest>
  ): RelationsFollowRequest {
    const message = { ...baseRelationsFollowRequest } as RelationsFollowRequest;
    if (object.userId !== undefined && object.userId !== null) {
      message.userId = object.userId;
    } else {
      message.userId = "";
    }
    return message;
  },
};

const baseRelationsFollowResponse: object = {};

export const RelationsFollowResponse = {
  encode(_: RelationsFollowResponse, writer: Writer = Writer.create()): Writer {
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): RelationsFollowResponse {
    const reader = input instanceof Reader ? input : new Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseRelationsFollowResponse,
    } as RelationsFollowResponse;
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

  fromJSON(_: any): RelationsFollowResponse {
    const message = {
      ...baseRelationsFollowResponse,
    } as RelationsFollowResponse;
    return message;
  },

  toJSON(_: RelationsFollowResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(
    _: DeepPartial<RelationsFollowResponse>
  ): RelationsFollowResponse {
    const message = {
      ...baseRelationsFollowResponse,
    } as RelationsFollowResponse;
    return message;
  },
};

const baseRelationsUnfollowRequest: object = { userId: "" };

export const RelationsUnfollowRequest = {
  encode(
    message: RelationsUnfollowRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.userId !== "") {
      writer.uint32(10).string(message.userId);
    }
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): RelationsUnfollowRequest {
    const reader = input instanceof Reader ? input : new Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseRelationsUnfollowRequest,
    } as RelationsUnfollowRequest;
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

  fromJSON(object: any): RelationsUnfollowRequest {
    const message = {
      ...baseRelationsUnfollowRequest,
    } as RelationsUnfollowRequest;
    if (object.userId !== undefined && object.userId !== null) {
      message.userId = String(object.userId);
    } else {
      message.userId = "";
    }
    return message;
  },

  toJSON(message: RelationsUnfollowRequest): unknown {
    const obj: any = {};
    message.userId !== undefined && (obj.userId = message.userId);
    return obj;
  },

  fromPartial(
    object: DeepPartial<RelationsUnfollowRequest>
  ): RelationsUnfollowRequest {
    const message = {
      ...baseRelationsUnfollowRequest,
    } as RelationsUnfollowRequest;
    if (object.userId !== undefined && object.userId !== null) {
      message.userId = object.userId;
    } else {
      message.userId = "";
    }
    return message;
  },
};

const baseRelationsUnfollowResponse: object = {};

export const RelationsUnfollowResponse = {
  encode(
    _: RelationsUnfollowResponse,
    writer: Writer = Writer.create()
  ): Writer {
    return writer;
  },

  decode(
    input: Reader | Uint8Array,
    length?: number
  ): RelationsUnfollowResponse {
    const reader = input instanceof Reader ? input : new Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseRelationsUnfollowResponse,
    } as RelationsUnfollowResponse;
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

  fromJSON(_: any): RelationsUnfollowResponse {
    const message = {
      ...baseRelationsUnfollowResponse,
    } as RelationsUnfollowResponse;
    return message;
  },

  toJSON(_: RelationsUnfollowResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(
    _: DeepPartial<RelationsUnfollowResponse>
  ): RelationsUnfollowResponse {
    const message = {
      ...baseRelationsUnfollowResponse,
    } as RelationsUnfollowResponse;
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

export interface Relations {
  Follow(request: RelationsFollowRequest): Promise<RelationsFollowResponse>;
  Unfollow(
    request: RelationsUnfollowRequest
  ): Promise<RelationsUnfollowResponse>;
}

export class RelationsClientImpl implements Relations {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
  }
  Follow(request: RelationsFollowRequest): Promise<RelationsFollowResponse> {
    const data = RelationsFollowRequest.encode(request).finish();
    const promise = this.rpc.request("meme.Relations", "Follow", data);
    return promise.then((data) =>
      RelationsFollowResponse.decode(new Reader(data))
    );
  }

  Unfollow(
    request: RelationsUnfollowRequest
  ): Promise<RelationsUnfollowResponse> {
    const data = RelationsUnfollowRequest.encode(request).finish();
    const promise = this.rpc.request("meme.Relations", "Unfollow", data);
    return promise.then((data) =>
      RelationsUnfollowResponse.decode(new Reader(data))
    );
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
