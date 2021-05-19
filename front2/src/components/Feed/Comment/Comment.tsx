import React from "react";
import * as schema from "../../../api/api2";
import {Link} from "../../Link/Link";
import styles from "./Comment.module.css";

export const Comment = (props: { data: schema.CommentRenderer }) => {
    return <div>
        <Link className={styles.Author} href={props.data.authorUrl}>
            {props.data.authorName}
        </Link>

        {props.data.text}
    </div>
}
