import React, {useEffect} from "react";
import {Post as PostTT} from "./Post";
import {PostComposer} from "./PostComposer";
import {Spinner} from "./Spinner";
import {getByID, getByType, storeSubscribe, storeUnsubscribe} from "../store/store";
import {storeFeedLoad} from "../store/feed";

export function Feed() {
    const [viewerId, setViewerId] = React.useState("");
    const [postIds, setPostIds] = React.useState<string[]>([]);
    const [isLoaded, setIsLoaded] = React.useState(false);

    useEffect(() => {
        storeFeedLoad();

        const callback = () => {
            const q = getByID("query");
            if (q && q.type == "Query") {
                setViewerId(q.viewer || "");

                const feedObj = getByType("Feed");
                if (feedObj && feedObj.type == "Feed") {
                    setPostIds(feedObj.feed);

                    setIsLoaded(true);
                }
            }

        }

        storeSubscribe(callback)

        return () => storeUnsubscribe(callback)
    })

    if (!isLoaded) {
        return <Spinner/>
    }

    return <>
        {viewerId && <>
            <PostComposer/>
            {postIds.map(postId => {
                return <PostTT postId={postId} key={postId}/>
            })}

            {!postIds && <div>Лента новостей пуста</div>}
        </>}

        {!viewerId && <div>Авторизуйтесь, чтобы посмотреть ленту</div>}
    </>;
}
