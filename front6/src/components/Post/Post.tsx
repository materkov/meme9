import React from "react";
import * as styles from "./Post.module.css";
import {Link} from "../Link/Link";
import {useGlobals} from "../../store/globals";
import {Post as ApiPost, postsDelete} from "../../api/api";
import {useDiscoverPage} from "../../store/discoverPage";

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

    return <div className={styles.post}>
        <Link href={"/users/" + props.post.user?.id} className={styles.name}>{props.post.user?.name}</Link>

        <Link href={"/posts/" + props.post.id} className={styles.date}>{date}</Link>

        {props.post.text}

        {globals.viewerId && props.post.user?.id == globals.viewerId &&
            <a onClick={onDelete} href="#" className={styles.deleteLink}>Delete post</a>
        }
    </div>
}
