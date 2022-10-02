import React from "react";
import {api, Post} from "../store/types";
import {PostUser} from "./PostUser";
import styles from "./Post.module.css";
import {Dots3} from "./icons/Dots3";
import classNames from "classnames";
import {PostLike} from "./PostLike";
import {feedStore} from "../store/Feed";

export type Props = {
    post: Post;
}

export function ComponentPost(props: Props) {
    const post = props.post;
    const [menuHidden, setMenuHidden] = React.useState(true);

    const onDelete = () => {
        feedStore.delete(post.id);
    }

    const onLikeToggle = () => {
        if (post.isLiked) {
            feedStore.unlike(post.id);
        } else {
            feedStore.like(post.id);
        }
    }

    return (
        <div className={styles.post}>
            <div className={styles.topContainer}>
                <PostUser post={post}/>

                {post.canDelete && <>
                    <Dots3 className={styles.menuIcon} onClick={() => setMenuHidden(!menuHidden)}/>

                    <div className={classNames({
                        [styles.menu]: true,
                        [styles.menuHidden]: menuHidden,
                    })}>
                        <div className={styles.menuItem} onClick={onDelete}>Удалить</div>
                    </div>
                </>
                }
            </div>

            <div className={styles.text}>{post.text}</div>

            <PostLike count={post.likesCount} isLiked={post.isLiked} onToggle={onLikeToggle}/>
        </div>
    )
}
