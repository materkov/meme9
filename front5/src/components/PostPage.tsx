import {ComponentPost} from "./Post";
import React from "react";
import {PostPage as PostPageRenderer} from "../store/types";

export function PostPage(props: { data: PostPageRenderer }) {
    const post = props.data[1];
    if (!post) {
        return null;
    }

    return <>
        <ComponentPost post={post}/>
    </>

}
