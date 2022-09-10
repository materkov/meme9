/* eslint-disable */
import { util, configure, Writer, Reader } from "protobufjs/minimal";
import * as Long from "long";

export const protobufPackage = "meme";

export interface LoginPageResponse {
  renderer: LoginPageRenderer | undefined;
}

export interface LoginPageRenderer {
  authUrl: string;
  text: string;
}

const baseLoginPageResponse: object = {};

export const LoginPageResponse = {
  encode(message: LoginPageResponse, writer: Writer = Writer.create()): Writer {
    if (message.renderer !== undefined) {
      LoginPageRenderer.encode(
        message.renderer,
        writer.uint32(10).fork()
      ).ldelim();
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): LoginPageResponse {
    const reader = input instanceof Reader ? input : new Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseLoginPageResponse } as LoginPageResponse;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.renderer = LoginPageRenderer.decode(reader, reader.uint32());
          break;
        default:
          reader.skipType(tag & 7);
          break;
      }
    }
    return message;
  },

  fromJSON(object: any): LoginPageResponse {
    const message = { ...baseLoginPageResponse } as LoginPageResponse;
    if (object.renderer !== undefined && object.renderer !== null) {
      message.renderer = LoginPageRenderer.fromJSON(object.renderer);
    } else {
      message.renderer = undefined;
    }
    return message;
  },

  toJSON(message: LoginPageResponse): unknown {
    const obj: any = {};
    message.renderer !== undefined &&
      (obj.renderer = message.renderer
        ? LoginPageRenderer.toJSON(message.renderer)
        : undefined);
    return obj;
  },

  fromPartial(object: DeepPartial<LoginPageResponse>): LoginPageResponse {
    const message = { ...baseLoginPageResponse } as LoginPageResponse;
    if (object.renderer !== undefined && object.renderer !== null) {
      message.renderer = LoginPageRenderer.fromPartial(object.renderer);
    } else {
      message.renderer = undefined;
    }
    return message;
  },
};

const baseLoginPageRenderer: object = { authUrl: "", text: "" };

export const LoginPageRenderer = {
  encode(message: LoginPageRenderer, writer: Writer = Writer.create()): Writer {
    if (message.authUrl !== "") {
      writer.uint32(10).string(message.authUrl);
    }
    if (message.text !== "") {
      writer.uint32(18).string(message.text);
    }
    return writer;
  },

  decode(input: Reader | Uint8Array, length?: number): LoginPageRenderer {
    const reader = input instanceof Reader ? input : new Reader(input);
    let end = length === undefined ? reader.len : reader.pos + length;
    const message = { ...baseLoginPageRenderer } as LoginPageRenderer;
    while (reader.pos < end) {
      const tag = reader.uint32();
      switch (tag >>> 3) {
        case 1:
          message.authUrl = reader.string();
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

  fromJSON(object: any): LoginPageRenderer {
    const message = { ...baseLoginPageRenderer } as LoginPageRenderer;
    if (object.authUrl !== undefined && object.authUrl !== null) {
      message.authUrl = String(object.authUrl);
    } else {
      message.authUrl = "";
    }
    if (object.text !== undefined && object.text !== null) {
      message.text = String(object.text);
    } else {
      message.text = "";
    }
    return message;
  },

  toJSON(message: LoginPageRenderer): unknown {
    const obj: any = {};
    message.authUrl !== undefined && (obj.authUrl = message.authUrl);
    message.text !== undefined && (obj.text = message.text);
    return obj;
  },

  fromPartial(object: DeepPartial<LoginPageRenderer>): LoginPageRenderer {
    const message = { ...baseLoginPageRenderer } as LoginPageRenderer;
    if (object.authUrl !== undefined && object.authUrl !== null) {
      message.authUrl = object.authUrl;
    } else {
      message.authUrl = "";
    }
    if (object.text !== undefined && object.text !== null) {
      message.text = object.text;
    } else {
      message.text = "";
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
