import * as types from "../../types/types";
import React from "react";
import * as styles from "./Post.module.css";
import {Link} from "../Link";

export function Post(props: { post: types.Post }) {
    const date = new Date(props.post.date).toLocaleString();

    return <div className={styles.post}>
        <Link href={"/users/" + props.post.user?.id} className={styles.name}>{props.post.user?.name}</Link>

        <div className={styles.date}>{date}</div>

        {props.post.text}
    </div>
}
