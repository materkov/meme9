import React, {useEffect} from "react";
import {ComponentPost} from "./Post";
import {Composer} from "./Composer";
import {apiHost, Post} from "../store/types";

export function Feed() {
    const [viewerID, setViewerID] = React.useState('');
    const [posts, setPosts] = React.useState<Post[]>([]);
    const [loaded, setIsLoaded] = React.useState(false);

    useEffect(() => {
        fetch(apiHost + "/feed", {
            method: 'POST',
            headers: {
                'authorization': 'Bearer ' + localStorage.getItem('authToken'),
            },
        }).then(r => r.json()).then(r => {
            setViewerID(r[0]);
            setPosts(r[1]);
            setIsLoaded(true);
        })
    }, [])

    if (!loaded) {
        return <>Загрузка...</>
    }

    return <>
        {viewerID ? <Composer/> : <i>Авторизуйтесь, чтобы написать пост</i>}

        {posts.map(post => <ComponentPost post={post} key={post.id}/>)}
    </>
}
