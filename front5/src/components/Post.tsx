import React from "react";
import {api, Post} from "../store/types";
import {PostUser} from "./PostUser";
import styles from "./Post.module.css";
import {Dots3} from "./icons/Dots3";
import classNames from "classnames";
import {PostLike} from "./PostLike";
import {queryClient, useQuery} from "../store/fetcher";
import {PostPhoto} from "./PostPhoto";

export type Props = {
    id: string;
}

export function ComponentPost(props: Props) {
    const {data: post} = useQuery<Post>("/posts/" + props.id);
    const [menuHidden, setMenuHidden] = React.useState(true);

    if (!post) return null;

    const onDelete = () => {
        api("/postDelete", {
            id: props.id,
        }).then(() => {
            queryClient.invalidateQueries(["/feed"])
        })
    }

    if (post.isDeleted) return <DeletedStub/>

    return (
        <div className={styles.post}>
            <div className={styles.topContainer}>
                <PostUser postId={props.id}/>

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

            <PostLike id={props.id}/>
        </div>
    )
}

function DeletedStub() {
    return <div className={styles.post}>
        <div className={styles.deletedStub}>
            Пост удален.
        </div>
    </div>
}
