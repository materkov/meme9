import React, {useEffect} from "react";
import {Composer} from "./Composer";
import {PostsList} from "./PostsList";
import classNames from "classnames";
import styles from "./Feed.module.css";
import {CurrentFeed, Global, store} from "../store/store";
import {connect} from "react-redux";
import {fetchDiscover, fetchFeed} from "../store/actions/feed";

interface Props {
    currentFeed: CurrentFeed;
    viewerId: string;
    postIds: string[];
    isLoaded: boolean;
    hasMore: boolean;
}

function Component(props: Props) {
    function setFeedType(feed: CurrentFeed) {
        store.dispatch({type: "feed/setCurrentFeed", feed: feed});
    }

    if (props.currentFeed == CurrentFeed.FEED) {
        fetchFeed("firstPage");
    } else if (props.currentFeed == CurrentFeed.DISCOVER) {
        fetchDiscover("firstPage");
    }

    return <>
        {props.viewerId ? <Composer/> : <i>Авторизуйтесь, чтобы написать пост</i>}

        <div>Тип ленты:&nbsp;
            <a onClick={() => setFeedType(CurrentFeed.FEED)}
               className={classNames({
                   [styles.currentFeedType]: props.currentFeed == CurrentFeed.FEED
               })}>Мои подписки</a> |&nbsp;
            <a onClick={() => setFeedType(CurrentFeed.DISCOVER)}
               className={classNames({
                   [styles.currentFeedType]: props.currentFeed == CurrentFeed.DISCOVER
               })}>Дискавер</a>
        </div>

        {!props.isLoaded && <div>Loading...</div>}

        <PostsList posts={props.postIds} onShowMore={() => fetchFeed("nextPage")} showMore={props.hasMore}
                   showMoreDisabled={false}/>
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
        currentFeed: state.feed.currentFeed,
        viewerId: state.viewer.id,
        postIds: postIds,
        isLoaded: state.feed.pages.length > 0,
        hasMore: !!lastCursor,
    }
})(Component);

