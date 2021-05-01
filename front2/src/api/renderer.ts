/* eslint-disable */
import { util, configure, Writer, Reader } from "protobufjs/minimal";
import * as Long from "long";
import {
  ProfileRenderer,
  FeedRenderer,
  PostRenderer,
  HeaderRenderer,
} from "./api2";
import { LoginPageRenderer } from "./login";

export const protobufPackage = "meme";

export interface UniversalRenderer {
  profileRenderer: ProfileRenderer | undefined;
  feedRenderer: FeedRenderer | undefined;
  postRenderer: PostRenderer | undefined;
  headerRenderer: HeaderRenderer | undefined;
  loginPageRenderer: LoginPageRenderer | undefined;
}

const baseUniversalRenderer: object = {};

export const UniversalRenderer = {
  encode(message: UniversalRenderer, writer: Writer = Writer.create()): Writer {
    if (message.profileRenderer !== undefined) {
      ProfileRenderer.encode(
        message.profileRenderer,
        writer.uint32(10).fork()
      ).ldelim();
    }
    if (message.feedRenderer !== undefined) {
      FeedRenderer.encode(
        message.feedRenderer,
        writer.uint32(18).fork()
      ).ldelim();
    }
    if (message.postRenderer !== undefined) {
      PostRenderer.encode(
        message.postRenderer,
        writer.uint32(26).fork()
      ).ldelim();
    }
    if (message.headerRenderer !== undefined) {
      HeaderRenderer.encode(
        message.headerRenderer,
        writer.uint32(34).fork()
      ).ldelim();
    }
    if (message.loginPageRenderer !== undefined) {
      LoginPageRenderer.encode(
        message.loginPageRenderer,
        writer.uint32(42).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): UniversalRenderer {
    const reader = input instanceof Reader ? input : new Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseUniversalRenderer } as UniversalRenderer;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.profileRenderer = ProfileRenderer.decode(
            reader,
            reader.uint32()
          );
          break;
        case 2:
          message.feedRenderer = FeedRenderer.decode(reader, reader.uint32());
          break;
        case 3:
          message.postRenderer = PostRenderer.decode(reader, reader.uint32());
          break;
        case 4:
          message.headerRenderer = HeaderRenderer.decode(
            reader,
            reader.uint32()
          );
          break;
        case 5:
          message.loginPageRenderer = LoginPageRenderer.decode(
            reader,
            reader.uint32()
          );
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): UniversalRenderer {
    const message = { ...baseUniversalRenderer } as UniversalRenderer;
    if (
      object.profileRenderer !== undefined &&
      object.profileRenderer !== null
    ) {
      message.profileRenderer = ProfileRenderer.fromJSON(
        object.profileRenderer
      );
    } else {
      message.profileRenderer = undefined;
    }
    if (object.feedRenderer !== undefined && object.feedRenderer !== null) {
      message.feedRenderer = FeedRenderer.fromJSON(object.feedRenderer);
    } else {
      message.feedRenderer = undefined;
    }
    if (object.postRenderer !== undefined && object.postRenderer !== null) {
      message.postRenderer = PostRenderer.fromJSON(object.postRenderer);
    } else {
      message.postRenderer = undefined;
    }
    if (object.headerRenderer !== undefined && object.headerRenderer !== null) {
      message.headerRenderer = HeaderRenderer.fromJSON(object.headerRenderer);
    } else {
      message.headerRenderer = undefined;
    }
    if (
      object.loginPageRenderer !== undefined &&
      object.loginPageRenderer !== null
    ) {
      message.loginPageRenderer = LoginPageRenderer.fromJSON(
        object.loginPageRenderer
      );
    } else {
      message.loginPageRenderer = undefined;
    }
    return message;
  },

  toJSON(message: UniversalRenderer): unknown {
    const obj: any = {};
    message.profileRenderer !== undefined &&
      (obj.profileRenderer = message.profileRenderer
        ? ProfileRenderer.toJSON(message.profileRenderer)
        : undefined);
    message.feedRenderer !== undefined &&
      (obj.feedRenderer = message.feedRenderer
        ? FeedRenderer.toJSON(message.feedRenderer)
        : undefined);
    message.postRenderer !== undefined &&
      (obj.postRenderer = message.postRenderer
        ? PostRenderer.toJSON(message.postRenderer)
        : undefined);
    message.headerRenderer !== undefined &&
      (obj.headerRenderer = message.headerRenderer
        ? HeaderRenderer.toJSON(message.headerRenderer)
        : undefined);
    message.loginPageRenderer !== undefined &&
      (obj.loginPageRenderer = message.loginPageRenderer
        ? LoginPageRenderer.toJSON(message.loginPageRenderer)
        : undefined);
    return obj;
  },

  fromPartial(object: DeepPartial<UniversalRenderer>): UniversalRenderer {
    const message = { ...baseUniversalRenderer } as UniversalRenderer;
    if (
      object.profileRenderer !== undefined &&
      object.profileRenderer !== null
    ) {
      message.profileRenderer = ProfileRenderer.fromPartial(
        object.profileRenderer
      );
    } else {
      message.profileRenderer = undefined;
    }
    if (object.feedRenderer !== undefined && object.feedRenderer !== null) {
      message.feedRenderer = FeedRenderer.fromPartial(object.feedRenderer);
    } else {
      message.feedRenderer = undefined;
    }
    if (object.postRenderer !== undefined && object.postRenderer !== null) {
      message.postRenderer = PostRenderer.fromPartial(object.postRenderer);
    } else {
      message.postRenderer = undefined;
    }
    if (object.headerRenderer !== undefined && object.headerRenderer !== null) {
      message.headerRenderer = HeaderRenderer.fromPartial(
        object.headerRenderer
      );
    } else {
      message.headerRenderer = undefined;
    }
    if (
      object.loginPageRenderer !== undefined &&
      object.loginPageRenderer !== null
    ) {
      message.loginPageRenderer = LoginPageRenderer.fromPartial(
        object.loginPageRenderer
      );
    } else {
      message.loginPageRenderer = undefined;
    }
    return message;
  },
};

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
