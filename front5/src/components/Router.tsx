import React from "react";
import {useCustomEventListener} from "react-custom-events";
import {Feed} from "./Feed";
import {UserPage} from "./UserPage";
import {VkCallback} from "./VkCallback";
import {PostPage} from "./PostPage";

export function Router() {
    const [url, setUrl] = React.useState(location.pathname + location.search);

    useCustomEventListener('urlChanged', () => {
        setUrl(location.pathname + location.search);
    })

    if (url == "/") {
        return <Feed/>
    } else if (url.startsWith("/users/")) {
        return <UserPage/>
    } else if (url.startsWith("/posts/")) {
        return <PostPage/>
    } else if (url.startsWith("/vk-callback")) {
        return <VkCallback/>
    } else {
        return <>404 page</>;
    }
}
