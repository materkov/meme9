import React from "react";
import * as schema from "../../api/api2";
import {Post} from "../Feed/Post/Post";
import {Comment} from "../Feed/Comment/Comment";
import styles from "./PostPage.module.css";
import {CommentComposer} from "./CommentComposer";
import {GlobalStoreContext} from "../../Context";
import {Store} from "../../Store";
import {Link} from "../Link/Link";

export function PostPage(props: { data: schema.PostRenderer }) {
    if (!props.data.post) {
        return null;
    }

    const store = React.useContext(GlobalStoreContext) as Store;

    return <>
        <Post data={props.data.post}/>

        {props.data.composerPlaceholder &&
        <Link href={store.headerData.loginUrl} className={styles.ComposerPlaceholder}>
            {props.data.composerPlaceholder}
        </Link>
        }

        {props.data.composer &&
        <CommentComposer data={props.data.composer}/>
        }

        <Comments comments={props.data.comments}/>
    </>
}


const Comments = (props: { comments: schema.CommentRenderer[] }) => {
    if (props.comments.length == 0) {
        return null;
    }

    return <div className={styles.CommentsList}>
        {props.comments.map(comment => (
            <div key={comment.id} className={styles.CommentsItem}>
                <Comment data={comment}/>
            </div>
        ))}
    </div>
}


