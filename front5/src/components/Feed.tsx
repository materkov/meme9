import React from "react";
import {Composer} from "./Composer";
import {PostsList} from "./PostsList";
import {feedStore} from "../store/Feed";
import {useQuery} from "../store/fetcher";
import {Edges} from "../store/types";

export function Feed() {
    const {data} = useQuery<Edges>("/feed");
    if (!data) {
        return <>Загрузка...</>
    }

    return <>
        {feedStore.viewerId ? <Composer/> : <i>Авторизуйтесь, чтобы написать пост</i>}
        fee
        <PostsList posts={data.items} onShowMore={feedStore.fetch} showMore={Boolean(feedStore.nextCursor)}
                   showMoreDisabled={feedStore.isLoading}
        />
    </>
}
