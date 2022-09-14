import React from "react";
import {api, Post} from "../store/types";
import {PostUser} from "./PostUser";
import styles from "./Post.module.css";
import {Dots3} from "./icons/Dots3";
import classNames from "classnames";

export type Props = {
    post: Post;
    onDelete?: () => void;
}

export function ComponentPost(props: Props) {
    const post = props.post;
    const [menuHidden, setMenuHidden] = React.useState(true);

    const onDelete = () => {
        api("/postDelete", {id: post.id});
        props.onDelete && props.onDelete();
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
        </div>
    )
}
