import React, {useEffect} from "react";
import {Composer} from "./Composer";
import {PostsList} from "./PostsList";
import classNames from "classnames";
import styles from "./Feed.module.css";
import {Global} from "../store/store";
import {connect} from "react-redux";
import {loadFeed} from "../store/actions/feed";

enum FeedType {
    FEED = "FEED",
    DISCOVER = "DISCOVER",
}

interface Props {
    viewerId: string;
    feed: string[];
    isLoaded: boolean;
}

function Component(props: Props) {
    const [feedType, setFeedType] = React.useState<FeedType>(FeedType.DISCOVER);

    useEffect(() => {
        loadFeed();
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
        <PostsList posts={props.feed} onShowMore={() => {
        }} showMore={false}
                   showMoreDisabled={true}
        />
    </>
}

export const Feed = connect((state: Global) => {
    return {
        viewerId: state.viewer.id,
        feed: state.feed.items,
        isLoaded: state.feed.isLoaded,
    } as Props
})(Component);

