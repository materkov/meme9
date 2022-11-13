import React from "react";
import {api, Edges, Post} from "../store/types";
import {PostUser} from "./PostUser";
import styles from "./Post.module.css";
import {Dots3} from "./icons/Dots3";
import classNames from "classnames";
import {PostLike} from "./PostLike";
import {queryClient, useQuery} from "../store/fetcher";

export type Props = {
    id: string;
}

export function ComponentPost(props: Props) {
    const {data: post} = useQuery<Post>("/posts/" + props.id);
    const [menuHidden, setMenuHidden] = React.useState(true);

    if (!post) return null;

    const onDelete = () => {
        const feedData = queryClient.getQueryData<Edges>(["/feed"]);
        if (feedData) {
            queryClient.setQueryData(["/feed"], {
                ...feedData,
                items: feedData.items.filter(item => item != props.id)
            })
        }

        api("/postDelete", {
            id: props.id,
        })
    }

    const onPhotoClick = () => {
        if (post && post.photoUrl) {
            window.open(post.photoUrl, '_blank');
        }
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

            {post.photoUrl &&
                <img src={post.photoUrl} onClick={onPhotoClick} className={styles.photoAttach}/>
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
