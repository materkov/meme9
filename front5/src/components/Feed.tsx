import React, {useEffect} from "react";
import {Composer} from "./Composer";
import {PostsList} from "./PostsList";
import classNames from "classnames";
import styles from "./Feed.module.css";
import {Global} from "../store/store";
import {connect} from "react-redux";
import {fetchFeed} from "../store/actions/feed";

enum FeedType {
    FEED = "FEED",
    DISCOVER = "DISCOVER",
}

interface Props {
    viewerId: string;
    postIds: string[];
    isLoaded: boolean;
    hasMore: boolean;
}

function Component(props: Props) {
    const [feedType, setFeedType] = React.useState<FeedType>(FeedType.DISCOVER);

    useEffect(() => {
        fetchFeed();
    }, []);

    return <>
        {props.viewerId ? <Composer/> : <i>Авторизуйтесь, чтобы написать пост</i>}

        <div>Тип ленты:&nbsp;
            <a onClick={() => setFeedType(FeedType.FEED)} className={classNames({
                [styles.currentFeedType]: feedType == FeedType.FEED
            })}>Мои подписки</a> |&nbsp;
            <a onClick={() => setFeedType(FeedType.DISCOVER)} className={classNames({
                [styles.currentFeedType]: feedType == FeedType.DISCOVER
            })}>Дискавер</a>
        </div>

        {!props.isLoaded && <div>Loading...</div>}

        <PostsList posts={props.postIds} onShowMore={fetchFeed} showMore={props.hasMore} showMoreDisabled={false}/>
    </>
}

export const Feed = connect((state: Global): Props => {
    let postIds = [];
    let lastCursor = '';
    for (const page of state.feed.pages) {
        for (const id of page.items) {
            postIds.push(id);
        }
        lastCursor = page.nextCursor;
    }

    state.feed.pages.forEach(page => page.items)
    return {
        viewerId: state.viewer.id,
        postIds: postIds,
        isLoaded: state.feed.pages.length > 0,
        hasMore: !!lastCursor,
    }
})(Component);

