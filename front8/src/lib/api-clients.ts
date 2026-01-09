import { AuthClientJSON } from "@/schema/auth.twirp";
import { PostsClientJSON } from "@/schema/posts.twirp";
import { UsersClientJSON } from "@/schema/users.twirp";
import { SubscriptionsClientJSON } from "@/schema/subscriptions.twirp";
import { LikesClientJSON } from "@/schema/likes.twirp";

import { getAuthToken } from "./authHelpers";

function getApiBaseUrl(service: string): string {
  const ports: Record<string, number> = {
    "meme.auth.Auth": 8081,
    "meme.users.Users": 8082,
    "meme.subscriptions.Subscriptions": 8083,
    "meme.likes.Likes": 8084,
    "meme.posts.Posts": 8085,
  };
  const port = ports[service] || 8080;

  // Server-side
  if (typeof window === "undefined") {
    return `http://localhost:${port}`;
  }

  // Client-side
  return "";
  if (window.location.hostname == "localhost") {
    return `http://localhost:${port}`;
  } else {
    return "https://meme2.mmaks.me";
  }
}

class TwirpRpcImpl {
  async request(
    service: string,
    method: string,
    contentType: string,
    data: any
  ): Promise<any> {
    const baseURL = getApiBaseUrl(service);
    const url = `${baseURL}/twirp/${service}/${method}`;

    const headers: Record<string, string> = {
      "Content-Type": "application/json",
      Accept: "application/json",
    };

    const token = await getAuthToken();
    if (token) {
      headers["Authorization"] = 'Bearer ' + token;
    }

    let response: Response;
    let responseBody: any = null;
    try {
      response = await fetch(url, {
        method: "POST",
        headers,
        body: JSON.stringify(data),
      });
    } catch (error) {
      throw new Error(`Network error`);
    }

    try {
      responseBody = await response.json();
    } catch (error) {
      throw new Error("Response is not json");
    }

    if (!response.ok) {
      throw new ApiError(responseBody.msg || `internal_error`);
    }

    return responseBody;
  }
}

const rpcImpl = new TwirpRpcImpl();

export const AuthClient = new AuthClientJSON(rpcImpl);
export const PostsClient = new PostsClientJSON(rpcImpl);
export const UsersClient = new UsersClientJSON(rpcImpl);
export const SubscriptionsClient = new SubscriptionsClientJSON(rpcImpl);
export const LikesClient = new LikesClientJSON(rpcImpl);

export class ApiError extends Error {
  err: string;

  constructor(err: string) {
    super(err);
    this.name = "ApiError";
    this.err = err;
  }
}
