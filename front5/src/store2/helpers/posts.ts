import {store} from "../store";
import {SetLikes, SetOnline, SetPhoto, SetPost, SetUser} from "../reducers";
import {PostsList} from "../../store/types";

export function parsePostsList(posts: PostsList) {
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
