import {
  TwirpContext,
  TwirpServer,
  RouterEvents,
  TwirpError,
  TwirpErrorCode,
  Interceptor,
  TwirpContentType,
  chainInterceptors,
} from "twirp-ts";
import { GenerateUploadUrlRequest, GenerateUploadUrlResponse } from "./photos";

//==================================//
//          Client Code             //
//==================================//

interface Rpc {
  request(
    service: string,
    method: string,
    contentType: "application/json" | "application/protobuf",
    data: object | Uint8Array
  ): Promise<object | Uint8Array>;
}

export interface PhotosClient {
  GenerateUploadUrl(
    request: GenerateUploadUrlRequest
  ): Promise<GenerateUploadUrlResponse>;
}

export class PhotosClientJSON implements PhotosClient {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.GenerateUploadUrl.bind(this);
  }
  GenerateUploadUrl(
    request: GenerateUploadUrlRequest
  ): Promise<GenerateUploadUrlResponse> {
    const data = GenerateUploadUrlRequest.toJson(request, {
      useProtoFieldName: true,
      emitDefaultValues: false,
    });
    const promise = this.rpc.request(
      "meme.photos.Photos",
      "GenerateUploadUrl",
      "application/json",
      data as object
    );
    return promise.then((data) =>
      GenerateUploadUrlResponse.fromJson(data as any, {
        ignoreUnknownFields: true,
      })
    );
  }
}

export class PhotosClientProtobuf implements PhotosClient {
  private readonly rpc: Rpc;
  constructor(rpc: Rpc) {
    this.rpc = rpc;
    this.GenerateUploadUrl.bind(this);
  }
  GenerateUploadUrl(
    request: GenerateUploadUrlRequest
  ): Promise<GenerateUploadUrlResponse> {
    const data = GenerateUploadUrlRequest.toBinary(request);
    const promise = this.rpc.request(
      "meme.photos.Photos",
      "GenerateUploadUrl",
      "application/protobuf",
      data
    );
    return promise.then((data) =>
      GenerateUploadUrlResponse.fromBinary(data as Uint8Array)
    );
  }
}

//==================================//
//          Server Code             //
//==================================//

export interface PhotosTwirp<T extends TwirpContext = TwirpContext> {
  GenerateUploadUrl(
    ctx: T,
    request: GenerateUploadUrlRequest
  ): Promise<GenerateUploadUrlResponse>;
}

export enum PhotosMethod {
  GenerateUploadUrl = "GenerateUploadUrl",
}

export const PhotosMethodList = [PhotosMethod.GenerateUploadUrl];

export function createPhotosServer<T extends TwirpContext = TwirpContext>(
  service: PhotosTwirp<T>
) {
  return new TwirpServer<PhotosTwirp, T>({
    service,
    packageName: "meme.photos",
    serviceName: "Photos",
    methodList: PhotosMethodList,
    matchRoute: matchPhotosRoute,
  });
}

function matchPhotosRoute<T extends TwirpContext = TwirpContext>(
  method: string,
  events: RouterEvents<T>
) {
  switch (method) {
    case "GenerateUploadUrl":
      return async (
        ctx: T,
        service: PhotosTwirp,
        data: Buffer,
        interceptors?: Interceptor<
          T,
          GenerateUploadUrlRequest,
          GenerateUploadUrlResponse
        >[]
      ) => {
        ctx = { ...ctx, methodName: "GenerateUploadUrl" };
        await events.onMatch(ctx);
        return handlePhotosGenerateUploadUrlRequest(
          ctx,
          service,
          data,
          interceptors
        );
      };
    default:
      events.onNotFound();
      const msg = `no handler found`;
      throw new TwirpError(TwirpErrorCode.BadRoute, msg);
  }
}

function handlePhotosGenerateUploadUrlRequest<
  T extends TwirpContext = TwirpContext
>(
  ctx: T,
  service: PhotosTwirp,
  data: Buffer,
  interceptors?: Interceptor<
    T,
    GenerateUploadUrlRequest,
    GenerateUploadUrlResponse
  >[]
): Promise<string | Uint8Array> {
  switch (ctx.contentType) {
    case TwirpContentType.JSON:
      return handlePhotosGenerateUploadUrlJSON<T>(
        ctx,
        service,
        data,
        interceptors
      );
    case TwirpContentType.Protobuf:
      return handlePhotosGenerateUploadUrlProtobuf<T>(
        ctx,
        service,
        data,
        interceptors
      );
    default:
      const msg = "unexpected Content-Type";
      throw new TwirpError(TwirpErrorCode.BadRoute, msg);
  }
}
async function handlePhotosGenerateUploadUrlJSON<
  T extends TwirpContext = TwirpContext
>(
  ctx: T,
  service: PhotosTwirp,
  data: Buffer,
  interceptors?: Interceptor<
    T,
    GenerateUploadUrlRequest,
    GenerateUploadUrlResponse
  >[]
) {
  let request: GenerateUploadUrlRequest;
  let response: GenerateUploadUrlResponse;

  try {
    const body = JSON.parse(data.toString() || "{}");
    request = GenerateUploadUrlRequest.fromJson(body, {
      ignoreUnknownFields: true,
    });
  } catch (e) {
    if (e instanceof Error) {
      const msg = "the json request could not be decoded";
      throw new TwirpError(TwirpErrorCode.Malformed, msg).withCause(e, true);
    }
  }

  if (interceptors && interceptors.length > 0) {
    const interceptor = chainInterceptors(...interceptors) as Interceptor<
      T,
      GenerateUploadUrlRequest,
      GenerateUploadUrlResponse
    >;
    response = await interceptor(ctx, request!, (ctx, inputReq) => {
      return service.GenerateUploadUrl(ctx, inputReq);
    });
  } else {
    response = await service.GenerateUploadUrl(ctx, request!);
  }

  return JSON.stringify(
    GenerateUploadUrlResponse.toJson(response, {
      useProtoFieldName: true,
      emitDefaultValues: false,
    }) as string
  );
}
async function handlePhotosGenerateUploadUrlProtobuf<
  T extends TwirpContext = TwirpContext
>(
  ctx: T,
  service: PhotosTwirp,
  data: Buffer,
  interceptors?: Interceptor<
    T,
    GenerateUploadUrlRequest,
    GenerateUploadUrlResponse
  >[]
) {
  let request: GenerateUploadUrlRequest;
  let response: GenerateUploadUrlResponse;

  try {
    request = GenerateUploadUrlRequest.fromBinary(data);
  } catch (e) {
    if (e instanceof Error) {
      const msg = "the protobuf request could not be decoded";
      throw new TwirpError(TwirpErrorCode.Malformed, msg).withCause(e, true);
    }
  }

  if (interceptors && interceptors.length > 0) {
    const interceptor = chainInterceptors(...interceptors) as Interceptor<
      T,
      GenerateUploadUrlRequest,
      GenerateUploadUrlResponse
    >;
    response = await interceptor(ctx, request!, (ctx, inputReq) => {
      return service.GenerateUploadUrl(ctx, inputReq);
    });
  } else {
    response = await service.GenerateUploadUrl(ctx, request!);
  }

  return Buffer.from(GenerateUploadUrlResponse.toBinary(response));
}
