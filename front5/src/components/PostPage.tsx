import {ComponentPost} from "./Post";
import React from "react";
import {PostPage as PostPageRenderer} from "../store2/types";

export function PostPage(props: { data: PostPageRenderer }) {
    const post = props.data.nodes?.posts?.find(item => item.id == props.data.pagePost);
    if (!post) {
        return null;
    }

    const user = props.data.nodes?.users?.find(item => item.id == post.fromId);

    return <>
        <ComponentPost post={post} from={user}/>
    </>

}
