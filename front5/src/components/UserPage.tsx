import React from "react";
import {PostUser} from "./PostUser";
import {Post, UserPage as UserPageRenderer} from "../store/types";
import {ComponentPost} from "./Post";

export function UserPage(props: { data: UserPageRenderer }) {
    const user = props.data[0];
    const posts = props.data[1];

    return (
        <div>
            {user.name}
            <hr/>
            {posts.map(post => <ComponentPost key={post.id} post={post}/>)}
        </div>
    )
}
