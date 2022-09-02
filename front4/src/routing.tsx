import React, {useEffect} from "react";

import {UserPage} from "./components/UserPage";
import {Page} from "./components/Page";
import {Feed} from "./components/Feed";
import {PostPage} from "./components/PostPage";
import {getByType, storeSubscribe, storeUnsubscribe} from "./store/store";
import {VKAuth} from "./components/VkAuth";

export function ResolveRoute() {
    const [url, setUrl] = React.useState("");

    useEffect(() => {
        const cb = () => {
            const route = getByType("CurrentRoute");
            if (route && route.type == "CurrentRoute") {
                setUrl(route.url);
            }
        };

        // @ts-ignore
        window.store["fake:id:currentRoute"] = {
            id: "fake:id:currentRoute",
            type: "CurrentRoute",
            url: location.pathname,
        };
        setUrl(location.pathname);

        storeSubscribe(cb);
        return () => storeUnsubscribe(cb)
    })

    return <Page>{doResolveRoute(url)}</Page>;
}

function doResolveRoute(url: string) {
    if (url.match(/^\/$/)) {
        return <Feed/>
    } else if (url.match(/^\/vk-callback/)) {
        return <VKAuth/>
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