import React, {useEffect} from "react";
import {useDiscoverPage} from "../../store/discoverPage";
import * as styles from "./Discover.module.css";
import {useGlobals} from "../../store/globals";
import {FeedType, postsAdd} from "../../api/api";
import {useResources} from "../../store/resources";
import {PostsList} from "../Post/PostsList";
import {Composer} from "./Composer";

export function Discover() {
    const discoverState = useDiscoverPage();
    const globalState = useGlobals();
    const resources = useResources();

    useEffect(() => {
        discoverState.refetch();
    }, []);

    const postIds = discoverState.posts;
    const posts = postIds.map(postId => resources.posts[postId]);

    const switchType = (e: React.MouseEvent<HTMLAnchorElement>) => {
        discoverState.setType(discoverState.type === FeedType.FEED ? FeedType.DISCOVER : FeedType.FEED);
        discoverState.refetch();
        e.preventDefault();
    };

    const loadMore = () => {
        discoverState.fetch();
    }

    return <div>
        <h1>Discover</h1>

        {globalState.viewerId && <Composer/>}

        {globalState.viewerId && <>
            This is {discoverState.type == FeedType.DISCOVER ? 'discover' : 'feed'}. <a href="#" onClick={switchType}>
            Switch to {discoverState.type == FeedType.DISCOVER ? 'feed' : 'discover'}
        </a>
        </>}

        <PostsList posts={posts}/>

        {discoverState.postsPageToken && <button onClick={loadMore}>Load more</button>}
    </div>
}
