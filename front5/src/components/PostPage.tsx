import {ComponentPost} from "./Post";
import React from "react";

export function PostPage() {
    const postId = location.pathname.substring(7);
    return <>
        <ComponentPost id={postId}/>
    </>

}
