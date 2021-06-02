/* eslint-disable */
import { util, configure, Reader, Writer } from "protobufjs/minimal";
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
  sandboxRenderer: SandboxRenderer | undefined;
}

export interface ResolveRouteRequest {
  url: string;
}

export interface SandboxRenderer {}

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
    if (message.sandboxRenderer !== undefined) {
      SandboxRenderer.encode(
        message.sandboxRenderer,
        writer.uint32(50).fork()
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
        case 6:
          message.sandboxRenderer = SandboxRenderer.decode(
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
    if (
      object.sandboxRenderer !== undefined &&
      object.sandboxRenderer !== null
    ) {
      message.sandboxRenderer = SandboxRenderer.fromJSON(
        object.sandboxRenderer
      );
    } else {
      message.sandboxRenderer = undefined;
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
    message.sandboxRenderer !== undefined &&
      (obj.sandboxRenderer = message.sandboxRenderer
        ? SandboxRenderer.toJSON(message.sandboxRenderer)
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
    if (
      object.sandboxRenderer !== undefined &&
      object.sandboxRenderer !== null
    ) {
      message.sandboxRenderer = SandboxRenderer.fromPartial(
        object.sandboxRenderer
      );
    } else {
      message.sandboxRenderer = undefined;
    }
    return message;
  },
};

const baseResolveRouteRequest: object = { url: "" };

export const ResolveRouteRequest = {
  encode(
    message: ResolveRouteRequest,
    writer: Writer = Writer.create()
  ): Writer {
    if (message.url !== "") {
      writer.uint32(10).string(message.url);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): ResolveRouteRequest {
    const reader = input instanceof Reader ? input : new Reader(input);
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

  toJSON(message: ResolveRouteRequest): unknown {
    const obj: any = {};
    message.url !== undefined && (obj.url = message.url);
    return obj;
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
};

const baseSandboxRenderer: object = {};

export const SandboxRenderer = {
  encode(_: SandboxRenderer, writer: Writer = Writer.create()): Writer {
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): SandboxRenderer {
    const reader = input instanceof Reader ? input : new Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseSandboxRenderer } as SandboxRenderer;
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

  fromJSON(_: any): SandboxRenderer {
    const message = { ...baseSandboxRenderer } as SandboxRenderer;
    return message;
  },

  toJSON(_: SandboxRenderer): unknown {
    const obj: any = {};
    return obj;
  },

  fromPartial(_: DeepPartial<SandboxRenderer>): SandboxRenderer {
    const message = { ...baseSandboxRenderer } as SandboxRenderer;
    return message;
  },
};

export interface Utils {
  ResolveRoute(request: ResolveRouteRequest): Promise<UniversalRenderer>;
}

export class UtilsClientImpl implements Utils {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
  }
  ResolveRoute(request: ResolveRouteRequest): Promise<UniversalRenderer> {
    const data = ResolveRouteRequest.encode(request).finish();
    const promise = this.rpc.request("meme.Utils", "ResolveRoute", data);
    return promise.then((data) => UniversalRenderer.decode(new Reader(data)));
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
