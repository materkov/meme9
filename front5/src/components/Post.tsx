import React from "react";
import {api, Post} from "../store/types";
import {PostUser} from "./PostUser";
import styles from "./Post.module.css";
import {Dots3} from "./icons/Dots3";
import classNames from "classnames";
import {PostLike} from "./PostLike";

export type Props = {
    post: Post;
    onDelete?: () => void;
    onLike?: () => void;
    onUnlike?: () => void;
}

export function ComponentPost(props: Props) {
    const post = props.post;
    const [menuHidden, setMenuHidden] = React.useState(true);

    const onDelete = () => {
        api("/postDelete", {id: post.id});
        props.onDelete && props.onDelete();
    }

    const onLikeToggle = () => {
        if (post.isLiked) {
            api("/postUnlike", {id: post.id});
            props.onUnlike && props.onUnlike();
        } else {
            api("/postLike", {id: post.id});
            props.onLike && props.onLike();
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
