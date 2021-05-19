import * as schema from "./api/posts";
import {AddCommentRequest, AddCommentResponse, PostsAddRequest, PostsAddResponse} from "./api/posts";
import {
    FeedGetHeaderRequest,
    FeedGetHeaderResponse,
    RelationsFollowRequest,
    RelationsFollowResponse,
    RelationsUnfollowRequest, RelationsUnfollowResponse
} from "./api/api2";
import {ResolveRouteRequest, UniversalRenderer} from "./api/renderer";

function api(method: string, args: any): Promise<any> {
    return new Promise((resolve, reject) => {
        fetch("/twirp/" + method, {
            method: 'POST',
            credentials: 'include',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(args),
        })
            .then(r => {
                if (r.status !== 200) {
                    reject();
                    return
                }

                return r.json()
            })
            .then(r => {
                resolve(r);
            })
            .catch(() => reject())
    })
}

export class API {
    static Posts_ToggleLike = (req: schema.ToggleLikeRequest): Promise<schema.ToggleLikeResponse> => {
        return new Promise(((resolve, reject) => {
            api("meme.Posts/ToggleLike", schema.ToggleLikeRequest.toJSON(req))
                .then(r => resolve(schema.ToggleLikeResponse.fromJSON(r)))
                .catch(e => reject(e));
        }))
    }

    static Feed_GetHeader = (req: FeedGetHeaderRequest): Promise<FeedGetHeaderResponse> => {
        return new Promise(((resolve, reject) => {
            api("meme.Feed/GetHeader", FeedGetHeaderRequest.toJSON(req))
                .then(r => resolve(FeedGetHeaderResponse.fromJSON(r)))
                .catch(e => reject(e));
        }))
    }

    static Utils_ResolveRoute = (req: ResolveRouteRequest): Promise<UniversalRenderer> => {
        return new Promise(((resolve, reject) => {
            api("meme.Utils/ResolveRoute", ResolveRouteRequest.toJSON(req))
                .then(r => resolve(UniversalRenderer.fromJSON(r)))
                .catch(e => reject(e));
        }))
    }

    static Posts_Add = (req: PostsAddRequest): Promise<PostsAddResponse> => {
        return new Promise(((resolve, reject) => {
            api("meme.Posts/Add", PostsAddRequest.toJSON(req))
                .then(r => resolve(PostsAddResponse.fromJSON(r)))
                .catch(e => reject(e));
        }))
    }

    static Posts_AddComment = (req: AddCommentRequest): Promise<AddCommentResponse> => {
        return new Promise(((resolve, reject) => {
            api("meme.Posts/AddComment", AddCommentRequest.toJSON(req))
                .then(r => resolve(AddCommentResponse.fromJSON(r)))
                .catch(e => reject(e));
        }))
    }

    static Relations_Follow = (req: RelationsFollowRequest): Promise<RelationsFollowResponse> => {
        return new Promise(((resolve, reject) => {
            api("meme.Relations/Follow", RelationsFollowRequest.toJSON(req))
                .then(r => resolve(RelationsFollowResponse.fromJSON(r)))
                .catch(e => reject(e));
        }))
    }

    static Relations_Unfollow = (req: RelationsUnfollowRequest): Promise<RelationsUnfollowResponse> => {
        return new Promise(((resolve, reject) => {
            api("meme.Relations/Unfollow", RelationsUnfollowRequest.toJSON(req))
                .then(r => resolve(RelationsUnfollowResponse.fromJSON(r)))
                .catch(e => reject(e));
        }))
    }
}
