import React from "react";

import {UserPage} from "./components/UserPage";
import {Page} from "./components/Page";
import {Feed} from "./components/Feed";
import {PostPage} from "./components/PostPage";

export function ResolveRoute(props: { url: string }) {
    const url = props.url;

    return <Page>{doResolveRoute(url)}</Page>;
}

function doResolveRoute(url: string) {
    if (url.match(/^\/$/)) {
        return <Feed/>
    } else if (url.match(/^\/vk-callback/)) {
        return <Feed/>
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