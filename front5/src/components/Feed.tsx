import React, {useEffect} from "react";
import {ComponentPost} from "./Post";
import {Composer} from "./Composer";
import {api, Post} from "../store/types";
import {useCustomEventListener} from "react-custom-events";

export function Feed() {
    const [viewerID, setViewerID] = React.useState('');
    const [posts, setPosts] = React.useState<Post[]>([]);
    const [loaded, setIsLoaded] = React.useState(false);
    const [err, setIsError] = React.useState(false);

    const refreshData = () => {
        api("/feed", {}).then(data => {
            setViewerID(data[0]);
            setPosts(data[1]);
            setIsLoaded(true);
        }).catch(() => setIsError(true));
    }

    useEffect(refreshData, []);
    useCustomEventListener('postCreated', refreshData);

    if (!loaded) {
        return <>Загрузка...</>
    }
    if (err) {
        return <>Ошибка...</>
    }

    return <>
        {viewerID ? <Composer/> : <i>Авторизуйтесь, чтобы написать пост</i>}

        {posts.map(post => <ComponentPost post={post} key={post.id}/>)}
    </>
}
