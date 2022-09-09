import {ComponentPost} from "./Post";
import React, {useEffect} from "react";
import {apiHost, Post} from "../store/types";

export function PostPage() {
    const [post, setPost] = React.useState<Post>();
    const [isLoaded, setIsLoaded] = React.useState(false);

    useEffect(() => {
        const f = new FormData();
        f.set("id", location.pathname.substring(7));

        fetch(apiHost + "/postPage", {
            method: 'POST',
            body: f,
            headers: {
                'authorization': 'Bearer ' + localStorage.getItem('authToken'),
            }
        })
            .then(r => r.json())
            .then(r => {
                setPost(r);
                setIsLoaded(true);
            })
    }, [])

    if (!isLoaded || !post) {
        return <>Loading ...</>;
    }

    return <>
        <ComponentPost post={post}/>
    </>

}
