import {ComponentPost} from "./Post";
import React from "react";
import {PostPage as PostPageRenderer} from "../store2/types";

export function PostPage(props: { data: PostPageRenderer }) {
    const post = props.data.nodes?.posts?.find(item => item.id == props.data.pagePost);

    return <>
        {post && <ComponentPost post={post}/>}
    </>

}
