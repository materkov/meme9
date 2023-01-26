import {store} from "../store";
import * as types from "../../api/types";

export function parsePostsList(posts: types.PostsList) {
    posts.items.forEach(post => {
        if (post.user) {
            store.dispatch({type: "users/set", user: post.user});

            if (post.user.online) {
                store.dispatch({
                    type: "online/set",
                    online: post.user.online,
                    userId: post.user.id
                })
                delete post.user.online;
            }
            delete post.user;
        }

        if (post.photo) {
            store.dispatch({type: "photos/set", photo: post.photo})
            delete post.photo;
        }

        if (post.likesConnection) {
            store.dispatch({
                type: "posts/setLikes",
                postId: post.id,
                isLiked: post.likesConnection.isViewerLiked || false,
                count: post.likesConnection.totalCount || 0,
            })

            delete post.likesConnection;
        }

        store.dispatch({type: "posts/set", post: post})
    })
}
