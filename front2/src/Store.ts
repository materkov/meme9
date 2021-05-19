import {AddCommentRequest, AddCommentResponse, ToggleLikeRequest_Action} from "./api/posts";
import {API} from "./Api";
import * as schema from "./api/renderer";
import {HeaderRenderer} from "./api/api2";

export class Store {
    data: schema.UniversalRenderer;
    headerData: HeaderRenderer;
    error?: any;
    onChange: (data: schema.UniversalRenderer, data2: HeaderRenderer) => void;

    constructor(onChange: (data: schema.UniversalRenderer, data2: HeaderRenderer) => void) {
        this.data = schema.UniversalRenderer.fromJSON({});
        this.headerData = HeaderRenderer.fromJSON({});
        this.onChange = onChange;

        this.refreshHeader();
        setInterval(this.refreshHeader, 60 * 1000);
    }

    changed() {
        this.onChange(this.data, this.headerData);
    }

    togglePostLike(postId: string) {
        if (this.data?.feedRenderer) {
            for (let post of this.data.feedRenderer.posts) {
                if (post.id == postId && post.canLike) {
                    let action: ToggleLikeRequest_Action;

                    if (post.isLiked) {
                        post.isLiked = false;
                        post.likesCount--;
                        action = ToggleLikeRequest_Action.UNLIKE;
                    } else {
                        post.isLiked = true;
                        post.likesCount++;
                        action = ToggleLikeRequest_Action.LIKE;
                    }

                    API.Posts_ToggleLike({
                        action: action,
                        postId: post.id,
                    }).then(r => {
                        post.likesCount = r.likesCount;
                        this.changed();
                    }).catch(e => {
                        console.error(e);
                    });
                    break;
                }
            }
        }
    }

    followUser(userId: string) {
        if (this.data.profileRenderer && this.data.profileRenderer.id == userId) {
            this.data.profileRenderer.isFollowing = true;
        }
        this.changed();

        API.Relations_Follow({userId: userId})
            .catch(console.error)
    }

    unfollowUser(userId: string) {
        if (this.data?.profileRenderer && this.data.profileRenderer.id == userId) {
            this.data.profileRenderer.isFollowing = false;
        }

        this.changed();

        API.Relations_Unfollow({userId: userId})
            .catch(console.error)
    }

    navigate(route: string) {
        window.history.pushState(null, "meme", route);

        API.Utils_ResolveRoute({url: route})
            .then(data => {
                this.data = schema.UniversalRenderer.fromJSON(data);
                this.error = undefined;
                this.changed();
            })
            .catch(() => {
                this.data = schema.UniversalRenderer.fromJSON({});
                this.error = true;
                this.changed();
            })
    }

    addPost(text: string): Promise<string> {
        return new Promise(((resolve, reject) => {
            API.Posts_Add({text: text})
                .then(r => resolve(r.postUrl))
                .catch(reject)
        }))
    }

    addComment(r: AddCommentRequest): Promise<AddCommentResponse> {
        return new Promise(((resolve, reject) => {
            API.Posts_AddComment(r)
                .then(result => {
                    this.changed();
                    this.navigate('/posts/' + r.postId);
                    resolve(result)
                })
                .catch(reject)
        }))
    }

    refreshHeader() {
        API.Feed_GetHeader({})
            .then(r => {
                this.headerData = HeaderRenderer.fromJSON(r.renderer);
                this.changed();
            })
            .catch(() => console.error('Failed updating header'))
    }
}
