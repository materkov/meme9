import React from "react";
import styles from "./Post.module.css";
import * as schema from "../../../api/api2";
import {Link} from "../../Link/Link";

export interface Props {
    data: schema.Post;
}

export class Post extends React.Component<Props> {
    render() {
        const data = this.props.data;

        return (
            <div className={styles.Post}>
                <div className={styles.Header}>
                    <img className={styles.AuthorAvatar} alt="" src={data.authorAvatar}/>
                    <div>
                        <Link className={styles.Author} href={data.authorUrl}>{data.authorName}</Link>
                        <div className={styles.Date}>{data.dateDisplay}</div>
                    </div>
                </div>

                <div className={styles.Text}>{data.text}</div>

                {data.imageUrl &&
                <img className={styles.Image} alt="" src={data.imageUrl}/>
                }
            </div>
        );
    }
}
