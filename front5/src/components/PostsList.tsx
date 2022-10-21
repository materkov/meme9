import React from "react";
import {ComponentPost} from "./Post";
import styles from "./PostsList.module.css";

export type Props = {
    posts: string[];

    onShowMore?: () => void;
    showMore?: boolean;
    showMoreDisabled?: boolean;
}

export function PostsList(props: Props) {
    return <>
        {props.posts.map(post => <ComponentPost
            id={post}
            key={post}
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
