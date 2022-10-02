import React, {useEffect} from "react";
import {Composer} from "./Composer";
import {useCustomEventListener} from "react-custom-events";
import {PostsList} from "./PostsList";
import {feedStore} from "../store/Feed";

export function Feed() {
    const [version, setVersion] = React.useState(0);

    useEffect(() => {
        const cb = () => {
            setVersion(feedStore.version);
        };

        feedStore.fetch();
        feedStore.subscribe(cb);
        return () => {
            feedStore.unsubscribe(cb);
        }
    })

    useCustomEventListener('postCreated', feedStore.reset);

    if (!feedStore.isFirstLoaded) {
        return <>Загрузка...</>
    }

    return <>
        {feedStore.viewerId ? <Composer/> : <i>Авторизуйтесь, чтобы написать пост</i>}

        <PostsList posts={feedStore.posts} onShowMore={feedStore.fetch} showMore={Boolean(feedStore.nextCursor)}
                   showMoreDisabled={feedStore.isLoading}
        />
    </>
}
