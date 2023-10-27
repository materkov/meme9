import React from "react";
import * as styles from "./LinkAttach.module.css";
import {Link} from "../../api/api";

export function LinkAttach(props: {
    link: Link
}) {
    const link = props.link;

    return <a href={link.url} target="_blank" className={styles.linkHref}>
        <div className={styles.link}>
            <img className={styles.picture} src={link.imageUrl}/>

            <div>
                <div className={styles.title}>{link.title}</div>
                <div className={styles.description}>{link.description}</div>
                <div className={styles.domain}>{link.domain}</div>
            </div>
        </div>
    </a>;
}