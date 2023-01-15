import React from "react";
import * as types from "../store/types";
import {PostUser} from "./PostUser";
import styles from "./Post.module.css";
import {Dots3} from "./icons/Dots3";
import classNames from "classnames";
import {PostLike} from "./PostLike";
import {PostPhoto} from "./PostPhoto";
import {Global} from "../store2/store";
import {connect} from "react-redux";
import {deletePost} from "../store2/actions/posts";

interface Props {
    post: types.Post;
}

function ComponentPostInner(props: Props) {
    const post = props.post;
    const [menuHidden, setMenuHidden] = React.useState(true);

    const onDelete = () => {
        deletePost(post.id);
    }

    if (post.isDeleted) return <DeletedStub/>

    return (
        <div className={styles.post}>
            <div className={styles.topContainer}>
                <PostUser postId={post.id}/>

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

            {post.photoId &&
                <PostPhoto id={post.photoId} className={styles.photoAttach}/>
            }

            <PostLike id={post.id}/>
        </div>
    )
}

export const ComponentPost = connect((state: Global, ownProps: { id: string }) => {
    return {
        post: state.posts.byId[ownProps.id],
    } as Props
})(ComponentPostInner);


function DeletedStub() {
    return <div className={styles.post}>
        <div className={styles.deletedStub}>
            Пост удален.
        </div>
    </div>
}
