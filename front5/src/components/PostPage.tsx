import {ComponentPost} from "./Post";
import React, {useEffect} from "react";
import {api, Post} from "../store/types";

export function PostPage() {
    const [post, setPost] = React.useState<Post>();
    const [isLoaded, setIsLoaded] = React.useState(false);

    useEffect(() => {
        api("/postPage", {
            id: location.pathname.substring(7)
        }).then(r => {
            setPost(r);
            setIsLoaded(true);
        })
    }, [])

    if (!isLoaded || !post) {
        return <>Loading ...</>;
    }

    return <>
        <ComponentPost id={post.id}/>
    </>

}
