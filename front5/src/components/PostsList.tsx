import React from "react";
import {ComponentPost} from "./Post";
import {Post} from "../store/types";
import styles from "./PostsList.module.css";

export type Props = {
    posts: Post[];

    onShowMore?: () => void;
    showMore?: boolean;
    showMoreDisabled?: boolean;
}

export function PostsList(props: Props) {
    return <>
        {props.posts.map(post => <ComponentPost
            post={post}
            key={post.id}
        />)}

        {props.showMore && <button
            disabled={props.showMoreDisabled}
            className={styles.showMore}
            onClick={props.onShowMore}
        >
            Показать еще
        </button>}
    </>;
}
