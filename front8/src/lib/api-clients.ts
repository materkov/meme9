import { AuthClientJSON } from '@/schema/auth.twirp';
import { PostsClientJSON } from '@/schema/posts.twirp';
import { UsersClientJSON } from '@/schema/users.twirp';
import { SubscriptionsClientJSON } from '@/schema/subscriptions.twirp';
import { getAuthToken } from './authHelpers';

function getApiBaseUrl(): string {
  // Server-side
  if (typeof window === 'undefined') {
    return 'http://localhost:8080';
  }

  // Client-side
  if (window.location.hostname == "localhost") {
    return 'http://localhost:8080';
  } else {
    return 'https://meme2.mmaks.me';
  }
}

class TwirpRpcImpl {
  async request(
    service: string,
    method: string,
    contentType: string,
    data: any
  ): Promise<any> {
    const baseURL = getApiBaseUrl();
    const url = `${baseURL}/twirp/${service}/${method}`;
    
    const headers: Record<string, string> = {
      'Content-Type': 'application/json',
      'Accept': 'application/json',
    };

    const token = await getAuthToken();
    if (token) {
      headers['Authorization'] = token;
    }

    let response: Response;
    let responseBody: any = null;
    try {
      response = await fetch(url, {
        method: 'POST',
        headers,
        body: JSON.stringify(data),
      });
    } catch (error) {
      throw new Error(`Network error`);
    }

    try {
      responseBody = await response.json();
    } catch (error) {
      throw new Error('Response is not json');
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

export class ApiError extends Error {
  err: string;

  constructor(err: string) {
    super(err);
    this.name = 'ApiError';
    this.err = err;
  }
}
