import {Link} from "../Link";
import * as styles from "../Profile/Profile.module.css";
import React from "react";
import * as types from "../../types/types";

export function ArticleSnippet(props: { article: types.Article }) {
    const article = props.article;

    return <div>
        <Link className={styles.articleTitle} href={"/article/" + article.id}>{article.title}</Link>
        <div className={styles.articleSnippet}>{getTextSnippet(article)}</div>
    </div>
}

function getTextSnippet(article: types.Article): string {
    for (const p of article.paragraphs) {
        if (p.text) {
            return p.text.text;
        }
    }

    return "";
}
