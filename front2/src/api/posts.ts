/* eslint-disable */
import { util, configure, Reader, Writer } from "protobufjs/minimal";
import * as Long from "long";

export const protobufPackage = "meme";

export interface PostsAddRequest {
  text: string;
}

export interface PostsAddResponse {
  postUrl: string;
}

export interface ToggleLikeRequest {
  action: ToggleLikeRequest_Action;
  postId: string;
}

export enum ToggleLikeRequest_Action {
  LIKE = "LIKE",
  UNLIKE = "UNLIKE",
  UNRECOGNIZED = "UNRECOGNIZED",
}

export function toggleLikeRequest_ActionFromJSON(
  object: any
): ToggleLikeRequest_Action {
  switch (object) {
    case 0:
    case "LIKE":
      return ToggleLikeRequest_Action.LIKE;
    case 1:
    case "UNLIKE":
      return ToggleLikeRequest_Action.UNLIKE;
    case -1:
    case "UNRECOGNIZED":
    default:
      return ToggleLikeRequest_Action.UNRECOGNIZED;
  }
}

export function toggleLikeRequest_ActionToJSON(
  object: ToggleLikeRequest_Action
): string {
  switch (object) {
    case ToggleLikeRequest_Action.LIKE:
      return "LIKE";
    case ToggleLikeRequest_Action.UNLIKE:
      return "UNLIKE";
    default:
      return "UNKNOWN";
  }
}

export function toggleLikeRequest_ActionToNumber(
  object: ToggleLikeRequest_Action
): number {
  switch (object) {
    case ToggleLikeRequest_Action.LIKE:
      return 0;
    case ToggleLikeRequest_Action.UNLIKE:
      return 1;
    default:
      return 0;
  }
}

export interface ToggleLikeResponse {
  likesCount: number;
}

export interface AddCommentRequest {
  text: string;
  postId: string;
}

export interface AddCommentResponse {}

export interface CommentComposerRenderer {
  postId: string;
  placeholder: string;
}

const basePostsAddRequest: object = { text: "" };

export const PostsAddRequest = {
  encode(message: PostsAddRequest, writer: Writer = Writer.create()): Writer {
    if (message.text !== "") {
      writer.uint32(10).string(message.text);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): PostsAddRequest {
    const reader = input instanceof Reader ? input : new Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...basePostsAddRequest } as PostsAddRequest;
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

  fromJSON(object: any): PostsAddRequest {
    const message = { ...basePostsAddRequest } as PostsAddRequest;
    if (object.text !== undefined && object.text !== null) {
      message.text = String(object.text);
    } else {
      message.text = "";
    }
    return message;
  },

  toJSON(message: PostsAddRequest): unknown {
    const obj: any = {};
    message.text !== undefined && (obj.text = message.text);
    return obj;
  },

  fromPartial(object: DeepPartial<PostsAddRequest>): PostsAddRequest {
    const message = { ...basePostsAddRequest } as PostsAddRequest;
    if (object.text !== undefined && object.text !== null) {
      message.text = object.text;
    } else {
      message.text = "";
    }
    return message;
  },
};

const basePostsAddResponse: object = { postUrl: "" };

export const PostsAddResponse = {
  encode(message: PostsAddResponse, writer: Writer = Writer.create()): Writer {
    if (message.postUrl !== "") {
      writer.uint32(10).string(message.postUrl);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): PostsAddResponse {
    const reader = input instanceof Reader ? input : new Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...basePostsAddResponse } as PostsAddResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.postUrl = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): PostsAddResponse {
    const message = { ...basePostsAddResponse } as PostsAddResponse;
    if (object.postUrl !== undefined && object.postUrl !== null) {
      message.postUrl = String(object.postUrl);
    } else {
      message.postUrl = "";
    }
    return message;
  },

  toJSON(message: PostsAddResponse): unknown {
    const obj: any = {};
    message.postUrl !== undefined && (obj.postUrl = message.postUrl);
    return obj;
  },

  fromPartial(object: DeepPartial<PostsAddResponse>): PostsAddResponse {
    const message = { ...basePostsAddResponse } as PostsAddResponse;
    if (object.postUrl !== undefined && object.postUrl !== null) {
      message.postUrl = object.postUrl;
    } else {
      message.postUrl = "";
    }
    return message;
  },
};

const baseToggleLikeRequest: object = {
  action: ToggleLikeRequest_Action.LIKE,
  postId: "",
};

export const ToggleLikeRequest = {
  encode(message: ToggleLikeRequest, writer: Writer = Writer.create()): Writer {
    if (message.action !== ToggleLikeRequest_Action.LIKE) {
      writer.uint32(8).int32(toggleLikeRequest_ActionToNumber(message.action));
    }
    if (message.postId !== "") {
      writer.uint32(18).string(message.postId);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): ToggleLikeRequest {
    const reader = input instanceof Reader ? input : new Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseToggleLikeRequest } as ToggleLikeRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.action = toggleLikeRequest_ActionFromJSON(reader.int32());
          break;
        case 2:
          message.postId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ToggleLikeRequest {
    const message = { ...baseToggleLikeRequest } as ToggleLikeRequest;
    if (object.action !== undefined && object.action !== null) {
      message.action = toggleLikeRequest_ActionFromJSON(object.action);
    } else {
      message.action = ToggleLikeRequest_Action.LIKE;
    }
    if (object.postId !== undefined && object.postId !== null) {
      message.postId = String(object.postId);
    } else {
      message.postId = "";
    }
    return message;
  },

  toJSON(message: ToggleLikeRequest): unknown {
    const obj: any = {};
    message.action !== undefined &&
      (obj.action = toggleLikeRequest_ActionToJSON(message.action));
    message.postId !== undefined && (obj.postId = message.postId);
    return obj;
  },

  fromPartial(object: DeepPartial<ToggleLikeRequest>): ToggleLikeRequest {
    const message = { ...baseToggleLikeRequest } as ToggleLikeRequest;
    if (object.action !== undefined && object.action !== null) {
      message.action = object.action;
    } else {
      message.action = ToggleLikeRequest_Action.LIKE;
    }
    if (object.postId !== undefined && object.postId !== null) {
      message.postId = object.postId;
    } else {
      message.postId = "";
    }
    return message;
  },
};

const baseToggleLikeResponse: object = { likesCount: 0 };

export const ToggleLikeResponse = {
  encode(
    message: ToggleLikeResponse,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.likesCount !== 0) {
      writer.uint32(8).int32(message.likesCount);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): ToggleLikeResponse {
    const reader = input instanceof Reader ? input : new Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseToggleLikeResponse } as ToggleLikeResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.likesCount = reader.int32();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): ToggleLikeResponse {
    const message = { ...baseToggleLikeResponse } as ToggleLikeResponse;
    if (object.likesCount !== undefined && object.likesCount !== null) {
      message.likesCount = Number(object.likesCount);
    } else {
      message.likesCount = 0;
    }
    return message;
  },

  toJSON(message: ToggleLikeResponse): unknown {
    const obj: any = {};
    message.likesCount !== undefined && (obj.likesCount = message.likesCount);
    return obj;
  },

  fromPartial(object: DeepPartial<ToggleLikeResponse>): ToggleLikeResponse {
    const message = { ...baseToggleLikeResponse } as ToggleLikeResponse;
    if (object.likesCount !== undefined && object.likesCount !== null) {
      message.likesCount = object.likesCount;
    } else {
      message.likesCount = 0;
    }
    return message;
  },
};

const baseAddCommentRequest: object = { text: "", postId: "" };

export const AddCommentRequest = {
  encode(message: AddCommentRequest, writer: Writer = Writer.create()): Writer {
    if (message.text !== "") {
      writer.uint32(10).string(message.text);
    }
    if (message.postId !== "") {
      writer.uint32(18).string(message.postId);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): AddCommentRequest {
    const reader = input instanceof Reader ? input : new Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseAddCommentRequest } as AddCommentRequest;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.text = reader.string();
          break;
        case 2:
          message.postId = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): AddCommentRequest {
    const message = { ...baseAddCommentRequest } as AddCommentRequest;
    if (object.text !== undefined && object.text !== null) {
      message.text = String(object.text);
    } else {
      message.text = "";
    }
    if (object.postId !== undefined && object.postId !== null) {
      message.postId = String(object.postId);
    } else {
      message.postId = "";
    }
    return message;
  },

  toJSON(message: AddCommentRequest): unknown {
    const obj: any = {};
    message.text !== undefined && (obj.text = message.text);
    message.postId !== undefined && (obj.postId = message.postId);
    return obj;
  },

  fromPartial(object: DeepPartial<AddCommentRequest>): AddCommentRequest {
    const message = { ...baseAddCommentRequest } as AddCommentRequest;
    if (object.text !== undefined && object.text !== null) {
      message.text = object.text;
    } else {
      message.text = "";
    }
    if (object.postId !== undefined && object.postId !== null) {
      message.postId = object.postId;
    } else {
      message.postId = "";
    }
    return message;
  },
};

const baseAddCommentResponse: object = {};

export const AddCommentResponse = {
  encode(_: AddCommentResponse, writer: Writer = Writer.create()): Writer {
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): AddCommentResponse {
    const reader = input instanceof Reader ? input : new Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseAddCommentResponse } as AddCommentResponse;
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

  fromJSON(_: any): AddCommentResponse {
    const message = { ...baseAddCommentResponse } as AddCommentResponse;
    return message;
  },

  toJSON(_: AddCommentResponse): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(_: DeepPartial<AddCommentResponse>): AddCommentResponse {
    const message = { ...baseAddCommentResponse } as AddCommentResponse;
    return message;
  },
};

const baseCommentComposerRenderer: object = { postId: "", placeholder: "" };

export const CommentComposerRenderer = {
  encode(
    message: CommentComposerRenderer,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.postId !== "") {
      writer.uint32(10).string(message.postId);
    }
    if (message.placeholder !== "") {
      writer.uint32(18).string(message.placeholder);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): CommentComposerRenderer {
    const reader = input instanceof Reader ? input : new Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = {
      ...baseCommentComposerRenderer,
    } as CommentComposerRenderer;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.postId = reader.string();
          break;
        case 2:
          message.placeholder = reader.string();
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): CommentComposerRenderer {
    const message = {
      ...baseCommentComposerRenderer,
    } as CommentComposerRenderer;
    if (object.postId !== undefined && object.postId !== null) {
      message.postId = String(object.postId);
    } else {
      message.postId = "";
    }
    if (object.placeholder !== undefined && object.placeholder !== null) {
      message.placeholder = String(object.placeholder);
    } else {
      message.placeholder = "";
    }
    return message;
  },

  toJSON(message: CommentComposerRenderer): unknown {
    const obj: any = {};
    message.postId !== undefined && (obj.postId = message.postId);
    message.placeholder !== undefined &&
      (obj.placeholder = message.placeholder);
    return obj;
  },

  fromPartial(
    object: DeepPartial<CommentComposerRenderer>
  ): CommentComposerRenderer {
    const message = {
      ...baseCommentComposerRenderer,
    } as CommentComposerRenderer;
    if (object.postId !== undefined && object.postId !== null) {
      message.postId = object.postId;
    } else {
      message.postId = "";
    }
    if (object.placeholder !== undefined && object.placeholder !== null) {
      message.placeholder = object.placeholder;
    } else {
      message.placeholder = "";
    }
    return message;
  },
};

export interface Posts {
  Add(request: PostsAddRequest): Promise<PostsAddResponse>;
  ToggleLike(request: ToggleLikeRequest): Promise<ToggleLikeResponse>;
  AddComment(request: AddCommentRequest): Promise<AddCommentResponse>;
}

export class PostsClientImpl implements Posts {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
  }
  Add(request: PostsAddRequest): Promise<PostsAddResponse> {
    const data = PostsAddRequest.encode(request).finish();
    const promise = this.rpc.request("meme.Posts", "Add", data);
    return promise.then((data) => PostsAddResponse.decode(new Reader(data)));
  }

  ToggleLike(request: ToggleLikeRequest): Promise<ToggleLikeResponse> {
    const data = ToggleLikeRequest.encode(request).finish();
    const promise = this.rpc.request("meme.Posts", "ToggleLike", data);
    return promise.then((data) => ToggleLikeResponse.decode(new Reader(data)));
  }

  AddComment(request: AddCommentRequest): Promise<AddCommentResponse> {
    const data = AddCommentRequest.encode(request).finish();
    const promise = this.rpc.request("meme.Posts", "AddComment", data);
    return promise.then((data) => AddCommentResponse.decode(new Reader(data)));
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
