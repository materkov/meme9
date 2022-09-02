import React from "react";
import {PostUser} from "./PostUser";
import {Post, UserPage as UserPageRenderer} from "../store2/types";
import {ComponentPost} from "./Post";

export function UserPage(props: { data: UserPageRenderer }) {
    const user = props.data.nodes?.users?.find(item => item.id == props.data.pageUser);

    const posts: Post[] = [];
    for (let postID of props.data.posts) {
        const post = props.data.nodes?.posts?.find(post => post.id == postID);
        if (post) {
            posts.push(post);
        }
    }

    return (
        <div>
            {user && <PostUser user={user}/>}
            <hr/>
            {posts.map(post => <ComponentPost post={post}/>)}
        </div>
    )
}
