import {store} from "../store";
import * as types from "../../api/types";
import {SetUser} from "../reducers/users";
import {SetOnline} from "../reducers/online";
import {SetPhoto} from "../reducers/photos";
import {SetLikes, SetPost} from "../reducers/posts";

export function parsePostsList(posts: types.PostsList) {
    posts.items.forEach(post => {
        if (post.user) {
            store.dispatch({type: "users/set", user: post.user} as SetUser);

            if (post.user.online) {
                store.dispatch({
                    type: "online/set",
                    online: post.user.online,
                    userId: post.user.id
                } as SetOnline)
                delete post.user.online;
            }
            delete post.user;
        }

        if (post.photo) {
            store.dispatch({type: "photos/set", photo: post.photo} as SetPhoto)
            delete post.photo;
        }

        if (post.likesConnection) {
            store.dispatch({
                type: "posts/setLikes",
                postId: post.id,
                isLiked: post.likesConnection.isViewerLiked || false,
                count: post.likesConnection.totalCount || 0,
            } as SetLikes)

            delete post.likesConnection;
        }

        store.dispatch({type: "posts/set", post: post} as SetPost)
    })
}
