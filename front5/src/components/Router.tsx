import React from "react";
import {useCustomEventListener} from "react-custom-events";
import {Feed} from "./Feed";
import {UserPage} from "./UserPage";
import {VkCallback} from "./VkCallback";
import {PostPage} from "./PostPage";
import {LoginPage} from "./LoginPage";
import {RegisterPage} from "./RegisterPage";

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
    } else if (url.startsWith("/register")) {
        return <RegisterPage/>
    } else if (url.startsWith("/login")) {
        return <LoginPage/>
    } else if (url.startsWith("/vk-callback")) {
        return <VkCallback/>
    } else {
        return <>404 page</>;
    }
}
