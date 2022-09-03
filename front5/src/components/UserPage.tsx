import React from "react";
import {PostUser} from "./PostUser";
import {Post, User, UserPage as UserPageRenderer} from "../store/types";
import {ComponentPost} from "./Post";

export function UserPage(props: { data: UserPageRenderer }) {
    const user = props.data.nodes?.users?.find(item => item.id == props.data.pageUser);

    const posts: Array<[Post, User]> = [];
    for (let postID of props.data.posts || []) {
        const post = props.data.nodes?.posts?.find(post => post.id == postID);
        if (!post) {
            continue;
        }

        const user = props.data.nodes?.users?.find(item => item.id == post.fromId);
        if (!user) {
            continue;
        }

        posts.push([post, user]);
    }

    return (
        <div>
            {user && <PostUser user={user}/>}
            <hr/>
            {posts.map(([post, user]) => <ComponentPost key={post.id} from={user} post={post}/>)}
        </div>
    )
}
