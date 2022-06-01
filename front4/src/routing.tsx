import React from "react";
import {PostPage} from "./postpage";
import {UserPage} from "./userpage";
import {FeedPage} from "./components";

export function ResolveRoute(props: { url: string }) {
    const url = props.url;

    if (url.match(/^\/$/)) {
        return <FeedPage/>
    } else if (url.match(/^\/vk-callback/)) {
        return <FeedPage/>
    } else if (url.match(/^\/posts\/(\w+)/)) {
        const postId = url.substr(7);
        return <PostPage id={postId}/>;
    } else if (url.match(/^\/users\/(\w+)/)) {
        const userId = url.substr(7);
        return <UserPage id={userId}/>;
    } else {
        return <>404</>;
    }
}
