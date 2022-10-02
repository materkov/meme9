import React, {useEffect} from "react";
import {Composer} from "./Composer";
import {api, Post} from "../store/types";
import {useCustomEventListener} from "react-custom-events";
import {PostsList} from "./PostsList";

export function Feed() {
    const [viewerID, setViewerID] = React.useState('');
    const [posts, setPosts] = React.useState<Post[]>([]);
    const [loaded, setIsLoaded] = React.useState(false);
    const [err, setIsError] = React.useState(false);
    const [cursor, setCursor] = React.useState('');
    const [showMoreLocked, setShowMoreLocked] = React.useState(false);

    const refreshData = () => {
        setShowMoreLocked(true);

        api("/feed", {cursor: cursor}).then(data => {
            setViewerID(data.viewerId);
            setPosts([...posts, ...data.posts]);
            setIsLoaded(true);
            setShowMoreLocked(false);
            setCursor(data.nextCursor);
        }).catch(() => setIsError(true));
    }

    const onPostDelete = (postId: string) => {
        setPosts(posts.filter(post => post.id !== postId));
    }

    const onPostLike = (postId: string) => {
        const postsCopy = [...posts];
        for (let post of postsCopy) {
            if (post.id == postId) {
                post.isLiked = true;
                post.likesCount = (post.likesCount || 0) + 1;
            }
        }

        setPosts(postsCopy);
    }

    const onPostUnlike = (postId: string) => {
        const postsCopy = [...posts];
        for (let post of postsCopy) {
            if (post.id == postId) {
                post.isLiked = false;
                post.likesCount = (post.likesCount || 0) - 1;
            }
        }

        setPosts(postsCopy);
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

        <PostsList posts={posts} onPostDelete={onPostDelete} onShowMore={refreshData} showMore={Boolean(cursor)}
                   showMoreDisabled={showMoreLocked} onLike={onPostLike} onUnlike={onPostUnlike}
        />
    </>
}
