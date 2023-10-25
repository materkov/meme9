import React from "react";
import * as styles from "./Post.module.css";
import {Link} from "../Link/Link";
import {useGlobals} from "../../store/globals";
import {LikeAction, Post as ApiPost, postsDelete, postsLike} from "../../api/api";
import {useDiscoverPage} from "../../store/discoverPage";

const nl2br = (string: string) => {
    if (string) {
        return string.split("\n").map((item, key) => {
            return (
                <span key={key}>{item}<br/></span>
            );
        });
    }
};

export function Post(props: { post: ApiPost }) {
    const date = new Date(props.post.date).toLocaleString();
    const globals = useGlobals();
    const discoverPage = useDiscoverPage();

    const onDelete = (e: React.MouseEvent<HTMLAnchorElement>) => {
        e.preventDefault();
        postsDelete({postId: props.post.id}).then(() => {
            discoverPage.refetch();
            alert('Post deleted');
        });
    };

    const onLike = (e: React.MouseEvent<HTMLAnchorElement>) => {
        e.preventDefault();
        postsLike({
            postId: props.post.id,
            action: props.post.isLiked ? LikeAction.UNLIKE : LikeAction.LIKE,
        }).then(() => {
            // TODO: think about this
            discoverPage.refetch();
        })
    };

    return <div className={styles.post}>
        <Link href={"/users/" + props.post.user?.id} className={styles.name}>{props.post.user?.name}</Link>

        <Link href={"/posts/" + props.post.id} className={styles.date}>{date}</Link>

        <div>
            {nl2br(props.post.text)}
        </div>

        {globals.viewerId && props.post.user?.id == globals.viewerId &&
            <a onClick={onDelete} href="#" className={styles.deleteLink}>Delete post</a>
        }

        <div className={styles.likesLine}>
            {globals.viewerId &&
                <>
                    <a onClick={onLike} href="#">
                        {props.post.isLiked ? 'Unlike' : 'Like'}
                    </a>
                    {props.post.likesCount && <> | </>}
                </>
            }

            {props.post.likesCount && <>{props.post.likesCount} like(s)</>}
        </div>
    </div>
}
