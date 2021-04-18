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

export interface Posts {
  Add(request: PostsAddRequest): Promise<PostsAddResponse>;
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
