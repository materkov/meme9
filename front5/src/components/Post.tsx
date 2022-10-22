import React from "react";
import {Post} from "../store/types";
import {PostUser} from "./PostUser";
import styles from "./Post.module.css";
import {Dots3} from "./icons/Dots3";
import classNames from "classnames";
import {PostLike} from "./PostLike";
import {feedStore} from "../store/Feed";
import {useQuery} from "../store/fetcher";

export type Props = {
    id: string;
}

export function ComponentPost(props: Props) {
    const {data: post} = useQuery<Post>("/posts/" + props.id);
    const [menuHidden, setMenuHidden] = React.useState(true);

    if (!post) return null;

    const onDelete = () => {
        //feedStore.delete(props.id);
    }

    // TODO post.canDelete -> true

    return (
        <div className={styles.post}>
            <div className={styles.topContainer}>
                <PostUser postId={props.id}/>

                {true && <>
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

            <PostLike id={props.id}/>
        </div>
    )
}
