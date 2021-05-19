import React from "react";
import * as schema from "../../api/api2";
import {Post} from "../Feed/Post/Post";
import {Comment} from "../Feed/Comment/Comment";
import styles from "./PostPage.module.css";
import {CommentComposer} from "./CommentComposer";

export function PostPage(props: { data: schema.PostRenderer }) {
    if (!props.data.post) {
        return null;
    }

    return <>
        <Post data={props.data.post}/>
        {props.data.composer && <CommentComposer data={props.data.composer}/>}
        <Comments comments={props.data.comments}/>
    </>
}


const Comments = (props: { comments: schema.CommentRenderer[] }) => (
    <div className={styles.CommentsList}>
        {props.comments.map(comment => (
            <div key={comment.id} className={styles.CommentsItem}>
                <Comment data={comment}/>
            </div>
        ))}
    </div>
)


